package main

import (
	"synodict-go/internal/cmdpkg"
	"synodict-go/internal/common"
	"synodict-go/internal/iopkg"
)

func main() {
	IORequestCh := make(chan iopkg.IORequest)
	exitCh := make(chan common.Void)

	go iopkg.Request(IORequestCh, exitCh)
	go cmdpkg.Run(IORequestCh, exitCh)

	IORequestCh <- iopkg.IORequest{
		Out:     true,
		In:      false,
		Prompts: []string{"type \"help\" for instructions"},
	}

	<-exitCh

	close(IORequestCh)
}
