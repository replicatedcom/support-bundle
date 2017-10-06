package types

type Result struct {
	Task      string   `json:"task"`
	Args      []string `json:"arguments"`
	Filenames []string `json:"filenames"`
	Error     error    `json:"error,omitempty"`
}
