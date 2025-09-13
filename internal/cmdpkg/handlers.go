package cmdpkg

import (
	"fmt"
	"synodict-go/internal/structpkg"
)

var cmdHandlers = map[string]func(d *structpkg.Dict, args []string) []string{
	"add": func(d *structpkg.Dict, args []string) []string {
		d.AddSynonyms(args...)

		return []string{}
	},
	"add-words": func(d *structpkg.Dict, args []string) []string {
		d.AddWords(args...)

		return []string{}
	},
	"remove": func(d *structpkg.Dict, args []string) []string {
		d.RemoveWords(args...)

		return []string{}
	},
	"unlink": func(d *structpkg.Dict, args []string) []string {
		d.UnlinkSynonyms(args[0], args[1])

		return []string{}
	},
	"unlink-clean": func(d *structpkg.Dict, args []string) []string {
		d.UnlinkSynonymsAndCleanup(args[0], args[1])

		return []string{}
	},
	"count": func(d *structpkg.Dict, args []string) []string {
		count := d.SynonymCount(args[0])
		response := []string{}

		if count > 1 {
			response = append(response, fmt.Sprintf("word \"%s\" has %d synonyms", args[0], count))
		} else if count == 1 {
			response = append(response, fmt.Sprintf("word \"%s\" has 1 synonym", args[0]))
		} else {
			response = append(response, fmt.Sprintf("word \"%s\" has no synonyms", args[0]))
		}

		return response
	},
	"synonyms": func(d *structpkg.Dict, args []string) []string {
		synonyms := d.GetSymomims(args[0])
		response := []string{}

		if len(synonyms) > 0 {
			response = append(
				response,
				fmt.Sprintf("word \"%s\" synonym list:", args[0]),
			)

			for i, s := range synonyms {
				response = append(
					response,
					fmt.Sprintf("%d) %s", i+1, s),
				)
			}
		} else {
			response = append(
				response,
				fmt.Sprintf("word \"%s\" has no synonyms", args[0]),
			)
		}

		return response
	},
	"groups": func(d *structpkg.Dict, args []string) []string {
		groups := d.GetSynonymGroups()
		response := []string{}

		for i, group := range groups {
			response = append(
				response,
				fmt.Sprintf("%d group", i+1),
			)

			response = append(response, "")

			for j, word := range group {
				response = append(
					response,
					fmt.Sprintf("%d) %s", j+1, word),
				)

				if j < len(group)-1 {
					response = append(response, "")
				}
			}
		}

		return response
	},
	"words": func(d *structpkg.Dict, args []string) []string {
		words := d.GetWords()
		response := []string{}

		for i, word := range words {
			response = append(
				response,
				fmt.Sprintf("%d) %s", i+1, word),
			)
		}

		return response
	},
	"cleanup": func(d *structpkg.Dict, args []string) []string {
		d.Cleanup()

		return []string{}
	},
	"clear": func(d *structpkg.Dict, args []string) []string {
		d.Clear()

		return []string{}
	},
	"help": func(d *structpkg.Dict, args []string) []string {
		return []string{
			"available commands:",
			"add \"word1\"... - adds each word to the dictionary if not already present, and links them as synonyms",
			"add-words \"word1\"... - adds each word to the dictionary if not already present (does not link them as synonyms)",
			"remove \"word1\"... - removes each word from the dictionary if already present",
			"unlink \"word1\" \"word2\" - removes synonym link between words (does not delete the words themselves)",
			"unlink-clean \"word1\" \"word2\" - removes synonym link between words (deletes words if they have no other synonyms)",
			"count \"word\" - prints the number of synonyms of the word",
			"synonys \"word\" - prints all synonyms of the word",
			"groups - prints all synonym groups",
			"words - prints all words",
			"cleanup - removes words that have no synonyms from the dictionary",
			"clear - clears the dictionary (warning: cannot be undone)",
			"help - prints this help message",
			"done - stops execution",
		}
	},
}
