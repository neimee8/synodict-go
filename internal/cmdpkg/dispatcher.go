package cmdpkg

import (
	"strings"
	"synodict-go/internal/iopkg"
	"synodict-go/internal/structpkg"
)

var dict = structpkg.NewDict()

func Run(IORequestCh chan iopkg.IORequest, exitCh chan structpkg.Void) {
	for {
		request := iopkg.IORequest{
			Out:                 false,
			In:                  true,
			InCh:                make(chan string, 1),
			InValidationRegexes: cmdRegexes,
		}

		IORequestCh <- request

		select {
		case cmd := <-request.InCh:
			cmdParts := strings.Fields(cmd)
			op := cmdParts[0]
			args := []string{}

			for i := 1; i < len(cmdParts); i++ {
				args = append(args, strings.Trim(cmdParts[i], "\""))
			}

			output := cmdHandlers[op](dict, args)

			if len(output) > 0 {
				IORequestCh <- iopkg.IORequest{
					Out:     true,
					In:      false,
					Prompts: output,
				}
			}
		case <-exitCh:
			return
		}
	}
}
