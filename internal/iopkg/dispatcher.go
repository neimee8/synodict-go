package iopkg

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"regexp"
	"strings"
	"synodict-go/internal/structpkg"
)

const exitCmd = "done"

var reader = bufio.NewReader(os.Stdin)

var regexCache map[string]*regexp.Regexp = make(map[string]*regexp.Regexp)
var regexResultCache map[uint64]bool = make(map[uint64]bool)

func Request(requests chan IORequest, exitCh chan structpkg.Void) {
	for r := range requests {
		if r.Out {
			write(r.Prompts)
		}

		if r.In {
			input, exit := scan(r.InValidationRegexes)

			if exit {
				exitCh <- structpkg.Void{}
				close(exitCh)
			}

			r.InCh <- input
			close(r.InCh)

			cacheCleanup()
		}
	}
}

func write(prompts []string) {
	for _, p := range prompts {
		fmt.Fprintln(os.Stdout, "> "+p)
	}
}

func scan(regexes []string) (string, bool) {
	for {
		fmt.Print("type here> ")

		input, err := reader.ReadString('\n')

		if err != nil {
			return "", true
		}

		input = strings.TrimSpace(input)

		if input == exitCmd {
			return "", true
		}

		if validateByRegex(input, regexes) {
			return input, false
		} else {
			fmt.Fprintln(os.Stderr, "ERROR> incorrect syntax, type \"help\"")
		}
	}
}

func validateByRegex(input string, regexes []string) bool {
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
	if len(regexCache) > 5000 {
		regexCache = make(map[string]*regexp.Regexp)
	}

	if len(regexResultCache) > 5000 {
		regexResultCache = make(map[uint64]bool)
	}
}
