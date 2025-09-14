package iopkg

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"regexp"
	"strings"
	"synodict-go/internal/common"
)

const exitCmd = "done"

var reader = bufio.NewReader(os.Stdin)

var regexCache map[string]*regexp.Regexp = make(map[string]*regexp.Regexp)
var regexResultCache map[uint64]bool = make(map[uint64]bool)

func Request(requests chan IORequest, exitCh chan common.Void) {
	for r := range requests {
		if r.Out {
			write(r.Prompts)
		}

		if r.In {
			input, exit := scan(r.InValidationRegexes, r.InErrorPrompts)

			if exit {
				fmt.Println(" goodbye!")
				close(r.InCh)

				exitCh <- common.Void{}
				close(exitCh)

				return
			}

			r.InCh <- input
			close(r.InCh)

			cacheCleanup()
		}
	}
}

func write(prompts []string) {
	for _, p := range prompts {
		fmt.Fprintln(os.Stdout, ""+p)
	}
}

func scan(regexes, errPrompts []string) (string, bool) {
	for {
		fmt.Print(" > ")

		input, err := reader.ReadString('\n')

		if err != nil {
			return "", true
		}

		input = strings.ToLower(strings.TrimSpace(input))

		if input == exitCmd {
			return "", true
		}

		if validateByRegex(input, regexes) {
			return input, false
		}

		if len(errPrompts) == 0 {
			fmt.Fprintln(os.Stderr, "ERROR> incorrect syntax, type \"help\"")
		} else {
			for _, prompt := range errPrompts {
				fmt.Fprintf(os.Stderr, "ERROR> %s\n", prompt)
			}
		}
	}
}

func validateByRegex(input string, regexes []string) bool {
	if len(regexes) == 0 {
		return true
	}

	for _, regex := range regexes {
		resultHash := hashStrings([]string{input, regex})

		if _, ok := regexResultCache[resultHash]; ok && regexResultCache[resultHash] {
			return true
		}

		var regexCompiled *regexp.Regexp

		if _, ok := regexCache[regex]; ok {
			regexCompiled = regexCache[regex]
		} else {
			regexCompiled = regexp.MustCompile(regex)
			regexCache[regex] = regexCompiled
		}

		result := regexCompiled.MatchString(input)
		regexResultCache[resultHash] = result

		if result {
			return true
		}
	}

	return false
}

func hashStrings(parts []string) uint64 {
	h := fnv.New64a()

	for _, part := range parts {
		h.Write([]byte(part))
	}

	return h.Sum64()
}

func cacheCleanup() {
	if len(regexCache) > 2500 {
		regexCache = make(map[string]*regexp.Regexp)
	}

	if len(regexResultCache) > 2500 {
		regexResultCache = make(map[uint64]bool)
	}
}
