package types

type Result struct {
	Task       string   `json:"task"`
	Args       []string `json:"arguments"`
	Filenames  []string `json:"filenames"`
	RawError   error    `json:"error_raw,omitempty"`
	JSONError  error    `json:"error_json,omitempty"`
	HumanError error    `json:"error_human,omitempty"`
}
