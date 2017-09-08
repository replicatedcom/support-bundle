package types

type Result struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Filename    string `json:"filename"`
	RawError    error  `json:"error_raw,omitempty"`
	JSONError   error  `json:"error_json,omitempty"`
	HumanError  error  `json:"error_human,omitempty"`
}
