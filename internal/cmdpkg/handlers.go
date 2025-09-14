package cmdpkg

import (
	"fmt"
	"synodict-go/internal/common"
	"synodict-go/internal/iopkg"
	"synodict-go/internal/stgpkg"
	"synodict-go/internal/structpkg"
)

// helpers
func askUserChoice(IORequestCh chan iopkg.IORequest) bool {
	request := iopkg.IORequest{
		Out:                 true,
		In:                  true,
		Prompts:             []string{"are you sure? this action cannot be undone (y/n or done)"},
		InCh:                make(chan string),
		InValidationRegexes: []string{`^(y|n)$`},
		InErrorPrompts:      []string{"type \"y\", \"n\" or \"done\""},
	}

	IORequestCh <- request
	response := <-request.InCh

	return response != "n"
}

func collectErrors(errs []error, log *[]string) {
	if len(errs) > 0 {
		for _, err := range errs {
			*log = append(*log, err.Error())
		}
	}
}

// handlers
func add(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	err_log := []string{}
	errs := d.AddSynonyms(args...)
	collectErrors(errs, &err_log)

	return err_log
}

func addWords(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	err_log := []string{}
	errs := d.AddWords(args...)
	collectErrors(errs, &err_log)

	return err_log
}

func remove(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	err_log := []string{}
	errs := d.RemoveWords(args...)
	collectErrors(errs, &err_log)

	return err_log
}

func unlink(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	err_log := []string{}
	errs := d.UnlinkSynonyms(args[0], args[1])
	collectErrors(errs, &err_log)

	return err_log
}

func unlinkClean(d *structpkg.Dict, args []string, IORequestCh chan iopkg.IORequest) []string {
	if !askUserChoice(IORequestCh) {
		return []string{}
	}

	err_log := []string{}
	errs := d.UnlinkSynonymsAndCleanup(args[0], args[1])
	collectErrors(errs, &err_log)

	return err_log
}

func check(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	err_log := []string{}
	result, errs := d.AreSynonyms(args[0], args[1])
	collectErrors(errs, &err_log)

	if len(errs) > 0 {
		return err_log
	}

	if result {
		return []string{"synonyms: yes"}
	} else {
		return []string{"synonyms: no"}
	}
}

func checkDirect(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	err_log := []string{}
	result, errs := d.AreDirectSynonyms(args[0], args[1])
	collectErrors(errs, &err_log)

	if len(errs) > 0 {
		return err_log
	}

	if result {
		return []string{"direct-linked synonyms: yes"}
	} else {
		return []string{"direct-linked synonyms: no"}
	}
}

func exists(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	if d.WordExists(args[0]) {
		return []string{"exists: yes"}
	} else {
		return []string{"exists: no"}
	}
}

func count(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	result, err := d.SynonymCount(args[0])

	if err != nil {
		return []string{err.Error()}
	}

	response := []string{}

	if result > 1 {
		response = append(response, fmt.Sprintf("word \"%s\" has %d synonyms", args[0], result))
	} else if result == 1 {
		response = append(response, fmt.Sprintf("word \"%s\" has 1 synonym", args[0]))
	} else {
		response = append(response, fmt.Sprintf("word \"%s\" has no synonyms yet", args[0]))
	}

	return response
}

func synonyms(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	result, err := d.GetSynomyms(args[0])

	if err != nil {
		return []string{err.Error()}
	}

	response := []string{}

	if len(result) > 0 {
		response = append(
			response,
			fmt.Sprintf("word \"%s\" synonym list:", args[0]),
		)

		for i, s := range result {
			response = append(
				response,
				fmt.Sprintf("%d) %s", i+1, s),
			)
		}
	} else {
		response = append(
			response,
			fmt.Sprintf("word \"%s\" has no synonyms yet", args[0]),
		)
	}

	return response
}

func directSynonyms(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	result, err := d.GetDirectSynonyms(args[0])

	if err != nil {
		return []string{err.Error()}
	}

	response := []string{}

	if len(result) > 0 {
		response = append(
			response,
			fmt.Sprintf("word \"%s\" direct-linked synonym list:", args[0]),
		)

		for i, s := range result {
			response = append(
				response,
				fmt.Sprintf("%d) %s", i+1, s),
			)
		}
	} else {
		response = append(
			response,
			fmt.Sprintf("word \"%s\" has no direct-linked synonyms yet", args[0]),
		)
	}

	return response
}

func countGroups(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	count := len(d.GetSynonymGroups())

	switch count {
	case 0:
		return []string{"dictionary has no synonym groups yet"}

	case 1:
		return []string{"dictionary contains 1 synonym group"}

	default:
		return []string{fmt.Sprintf("dictionary contains %d synonym groups", count)}
	}
}

func groups(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	groups := d.GetSynonymGroups()
	response := []string{}

	for i, group := range groups {
		response = append(
			response,
			fmt.Sprintf("%d group", i+1),
		)

		for j, word := range group {
			response = append(
				response,
				fmt.Sprintf("%d) %s", j+1, word),
			)
		}

		if i < len(group)-1 {
			response = append(response, "")
		}
	}

	return response
}

func countWords(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	count := len(d.GetWords())

	switch count {
	case 0:
		return []string{"dictionary has no words yet"}

	case 1:
		return []string{"dictionary contains 1 word"}

	default:
		return []string{fmt.Sprintf("dictionary contains %d words", count)}
	}
}

func words(d *structpkg.Dict, args []string, _ chan iopkg.IORequest) []string {
	words := d.GetWords()
	response := []string{}

	for i, word := range words {
		response = append(
			response,
			fmt.Sprintf("%d) %s", i+1, word),
		)
	}

	return response
}

func cleanup(d *structpkg.Dict, args []string, IORequestCh chan iopkg.IORequest) []string {
	if !askUserChoice(IORequestCh) {
		return []string{}
	}

	d.Cleanup()

	return []string{}
}

func clear(d *structpkg.Dict, args []string, IORequestCh chan iopkg.IORequest) []string {
	if !askUserChoice(IORequestCh) {
		return []string{}
	}

	d.Clear()

	return []string{}
}

func importDict(d *structpkg.Dict, args []string, IORequestCh chan iopkg.IORequest) []string {
	stages := [][]string{
		{
			"please choose the import format:",
			"gob        - GOB",
			"csv        - CSV",
			"csvc       - CSV condensed",
			"c (cancel) - go back without import",
			"done       - stop execution",
		},
		{
			"please specify the file location:",
			"c (cancel) - go back to the previous step",
			"done       - stop execution",
		},
	}

	stage := 0
	format := ""
	path := ""

	for stage < len(stages) && stage >= 0 {
		errorPrompts := []string{}
		regexes := []string{}

		switch stage {
		case 0:
			errorPrompts = append(
				errorPrompts,
				"choose one of listed below:",
			)

			regexes = append(regexes, `^(gob|csv|csvc|c)$`)

		case 1:
			errorPrompts = append(
				errorPrompts,
				"type the path to import dictionary:",
			)
		}

		errorPrompts = append(errorPrompts, stages[stage][1:]...)

		request := iopkg.IORequest{
			Out:                 true,
			In:                  true,
			Prompts:             stages[stage],
			InCh:                make(chan string),
			InValidationRegexes: regexes,
			InErrorPrompts:      errorPrompts,
		}

		IORequestCh <- request
		response, ok := <-request.InCh

		if !ok {
			return []string{}
		}

		if response == "c" {
			stage--
			continue
		}

		switch stage {
		case 0:
			format = response
			stage++
			continue

		case 1:
			path = response

			for {
				_, err := stgpkg.Read(path)

				if err == nil {
					stage++
					break
				}

				request = iopkg.IORequest{
					Out:     true,
					In:      true,
					Prompts: []string{err.Error()},
					InCh:    make(chan string),
				}

				IORequestCh <- request
				path, ok = <-request.InCh

				if !ok {
					return []string{}
				}
			}
		}
	}

	if stage == -1 {
		return []string{"import canceled"}
	}

	if !d.IsEmpty() {
		stages = [][]string{
			{
				"current dictionary is not empty. choose the action:",
				"overwrite (o) - clear the current dictionary and replace it with the imported one",
				"merge (m)     - merge the imported dictionary with the current one",
				"cancel (c)    - go back without importing",
				"done          - stop execution",
			},
			{
				"do you want to save the dictionary before overwriting?",
				"yes (y)    - export and save the current dictionary first",
				"no (n)     - continue without saving",
				"cancel (c) - go back to the previous step",
				"done       - stop execution",
			},
		}

		stage = 0

		for stage < len(stages) && stage >= 0 {
			errorPrompts := []string{"choose one of listed below:"}
			errorPrompts = append(errorPrompts, stages[stage][1:]...)

			regexes := []string{}

			switch stage {
			case 0:
				regexes = append(regexes, `^(o|m|c)$`)

			case 1:
				regexes = append(regexes, `^(y|n|c)$`)
			}

			request := iopkg.IORequest{
				Out:                 true,
				In:                  true,
				Prompts:             stages[stage],
				InCh:                make(chan string),
				InValidationRegexes: regexes,
				InErrorPrompts:      errorPrompts,
			}

			IORequestCh <- request
			response, ok := <-request.InCh

			if !ok {
				return []string{}
			}

			if response == "c" {
				stage--
				continue
			}

			switch stage {
			case 0:
				switch response {
				case "o":
					stage++

				case "m":
					stage = len(stages)
				}

			case 1:
				switch response {
				case "y":
					exportDict(d, args, IORequestCh)

				case "n":
					d.Clear()
					stage++
				}
			}
		}

		if stage == -1 {
			return []string{"import canceled"}
		}
	}

	err := d.Import(path, format)

	if err != nil {
		return []string{err.Error()}
	}

	return []string{"imported successfully"}
}

func exportDict(d *structpkg.Dict, args []string, IORequestCh chan iopkg.IORequest) []string {
	stages := [][]string{
		{
			"please choose the export format:",
			"gob        - GOB",
			"csv        - CSV",
			"csvc       - CSV condensed",
			"c (cancel) - go back without export",
			"done       - stop execution",
		},
		{
			"please specify the desired file location:",
			"c (cancel) - go back to the previous step",
			"done       - stop execution",
		},
	}

	stage := 0
	format := ""
	path := ""
	fullPath := ""

	for stage < len(stages) && stage >= 0 {
		errorPrompts := []string{}
		regexes := []string{}

		switch stage {
		case 0:
			errorPrompts = append(
				errorPrompts,
				"choose one of listed below:",
			)

			regexes = append(regexes, `^(gob|csv|csvc|c)$`)

		case 1:
			errorPrompts = append(
				errorPrompts,
				"type the path to create file:",
			)
		}

		errorPrompts = append(errorPrompts, stages[stage][1:]...)

		request := iopkg.IORequest{
			Out:                 true,
			In:                  true,
			Prompts:             stages[stage],
			InCh:                make(chan string),
			InValidationRegexes: regexes,
			InErrorPrompts:      errorPrompts,
		}

		IORequestCh <- request
		response, ok := <-request.InCh

		if !ok {
			return []string{}
		}

		if response == "c" {
			stage--
			continue
		}

		switch stage {
		case 0:
			format = response
			stage++
			continue

		case 1:
			path = response

			for {
				fullPath = path + common.FormatFileExtensions[format]
				err := d.Export(fullPath, format)

				if err == nil {
					stage++
					break
				}

				request = iopkg.IORequest{
					Out:     true,
					In:      true,
					Prompts: []string{err.Error()},
					InCh:    make(chan string),
				}

				IORequestCh <- request
				path, ok = <-request.InCh

				if !ok {
					return []string{}
				}
			}
		}
	}

	if stage == -1 {
		return []string{"export canceled"}
	}

	return []string{fmt.Sprintf("exported successfully to %s", fullPath)}
}

func help(d *structpkg.Dict, args []string, IORequestCh chan iopkg.IORequest) []string {
	return []string{
		"available commands:",
		"add \"word1\"...                 - adds each word to the dictionary if not already present, and links them as synonyms",
		"add-words \"word1\"...           - adds each word to the dictionary if not already present (does not link them as synonyms)",
		"remove \"word1\"...              - removes each word from the dictionary if already present",
		"unlink \"word1\" \"word2\"       - removes synonym link between words (does not delete the words themselves)",
		"unlink-clean \"word1\" \"word2\" - removes synonym link between words (deletes words if they have no other synonyms)",
		"check \"word1\" \"word2\".       - checks if the words are synonyms (directly or transitively)",
		"check-direct \"word1\" \"word2\" - checks if the words are directly linked as synonyms",
		"exists \"word\"                  - checks if the word exists in the dictionary",
		"count \"word\"                   - prints the number of synonyms of the word",
		"synonyms \"word\"                - prints all synonyms of the word",
		"direct-synonyms \"word\"         - prints only directly linked synonyms (words that were explicitly connected)",
		"count-groups                     - prints the number of synonym groups",
		"groups                           - prints all synonym groups",
		"count-words                      - prints the total number of words in the dictionary",
		"words                            - prints all words",
		"cleanup                          - removes words that have no synonyms from the dictionary",
		"clear                            - clears the dictionary (warning: cannot be undone)",
		"import                           - import dictionary (supports gob/csv); if current dictionary is not empty, you will be prompted to save, merge, or overwrite",
		"export                           - export dictionary (supports gob/csv)",
		"help                             - prints this help message",
		"done                             - stops execution",
	}
}

// map
var cmdHandlers = map[string]func(d *structpkg.Dict, args []string, IORequestCh chan iopkg.IORequest) []string{
	"add":             add,
	"add-words":       addWords,
	"remove":          remove,
	"unlink":          unlink,
	"unlink-clean":    unlinkClean,
	"check":           check,
	"check-direct":    checkDirect,
	"exists":          exists,
	"count":           count,
	"synonyms":        synonyms,
	"direct-synonyms": directSynonyms,
	"count-groups":    countGroups,
	"groups":          groups,
	"count-words":     countWords,
	"words":           words,
	"cleanup":         cleanup,
	"clear":           clear,
	"import":          importDict,
	"export":          exportDict,
	"help":            help,
}
