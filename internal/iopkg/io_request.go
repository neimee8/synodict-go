package iopkg

type IORequest struct {
	Out bool
	In  bool

	Prompts []string
	InCh    chan string

	InValidationRegexes []string
	InErrorPrompts      []string
}
