package main

import (
	"bufio"
	"os"
	"strings"
	"synodict-go/internal/utils"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	input, _ := reader.ReadString('\n')
	utils.CmdDispatch(strings.TrimSpace(input))
}
