package types

import "encoding/json"

// Result represents a single file within a support bundle or the failure to
// collect the data for a single file within a support bundle. A Result may have
// both a Pathname and an Error if the file written was corrupted or incomplete.
type Result struct {
	Description string `json:"description"`
	// The subpath within the bundle
	Path  string `json:"path"`
	Error error  `json:"error,omitempty"`
}

// Result.Error will be {} if it has no exported fields, so replace it with a
// string.
func (r *Result) MarshalJSON() ([]byte, error) {
	intermediate := map[string]string{}

	if r.Description != "" {
		intermediate["description"] = r.Description
	}
	if r.Path != "" {
		intermediate["path"] = r.Path
	}
	if r.Error != nil {
		intermediate["error"] = r.Error.Error()
	}

	return json.Marshal(intermediate)
}
