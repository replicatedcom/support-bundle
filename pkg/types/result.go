package types

import (
	"encoding/json"
	"errors"
)

// Result represents a single file within a support bundle or the failure to
// collect the data for a single file within a support bundle. A Result may have
// both a Pathname and an Error if the file written was corrupted or incomplete.
type Result struct {
	// The subpath within the bundle
	Path   string `json:"path"`
	Format string `json:"format"`
	Spec   Spec   `json:"spec"`
	Error  error  `json:"error,omitempty"`
}

// Result.Error will be {} if it has no exported fields, so replace it with a
// string.
func (r *Result) MarshalJSON() ([]byte, error) {
	intermediate := map[string]interface{}{
		"path":   r.Path,
		"format": r.Format,
		"spec":   r.Spec,
	}
	if r.Error != nil {
		intermediate["error"] = r.Error.Error()
	}
	return json.Marshal(intermediate)
}

// convert Error from string to error
func (r *Result) UnmarshalJSON(raw []byte) error {
	intermediate := map[string]interface{}{}
	if err := json.Unmarshal(raw, &intermediate); err != nil {
		return err
	}

	r.Path, _ = intermediate["path"].(string)
	r.Format, _ = intermediate["format"].(string)
	r.Spec, _ = intermediate["spec"].(Spec)
	if errMsg, ok := intermediate["error"].(string); ok {
		r.Error = errors.New(errMsg)
	}
	return nil
}
