package cmdpkg

import (
	"strings"
	"synodict-go/internal/common"
	"synodict-go/internal/iopkg"
	"synodict-go/internal/structpkg"
)

var dict = structpkg.NewDict()

func Run(IORequestCh chan iopkg.IORequest, exitCh chan common.Void) {
	for {
		request := iopkg.IORequest{
			Out:                 false,
			In:                  true,
			InCh:                make(chan string, 1),
			InValidationRegexes: cmdRegexes,
		}

		IORequestCh <- request

		select {
		case cmd, ok := <-request.InCh:
			if !ok {
				return
			}

			cmdParts := strings.Fields(cmd)
			op := cmdParts[0]
			args := []string{}

			for i := 1; i < len(cmdParts); i++ {
				args = append(args, strings.Trim(cmdParts[i], "\""))
			}

			output := cmdHandlers[op](dict, args, IORequestCh)

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
