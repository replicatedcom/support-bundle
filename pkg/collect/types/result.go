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
	Path     string `json:"path"`
	Format   string `json:"format"`
	Size     int64  `json:"size"`
	Spec     Spec   `json:"spec"`
	Redacted bool   `json:"redacted"`
	Error    error  `json:"error,omitempty"`
}

type resultIntermediate struct {
	Path     string `json:"path"`
	Format   string `json:"format"`
	Size     int64  `json:"size"`
	Spec     Spec   `json:"spec"`
	Redacted bool   `json:"redacted"`
	Error    string `json:"error,omitempty"`
}

// MarshalJSON .Error will be {} if it has no exported fields, so replace it with a string.
func (r *Result) MarshalJSON() ([]byte, error) {
	intermediate := resultIntermediate{
		Path:     r.Path,
		Format:   r.Format,
		Spec:     r.Spec,
		Size:     r.Size,
		Redacted: r.Redacted,
	}
	if r.Error != nil {
		intermediate.Error = r.Error.Error()
	}
	return json.Marshal(intermediate)
}

// UnmarshalJSON will convert .Error from string to error
func (r *Result) UnmarshalJSON(raw []byte) error {
	var intermediate resultIntermediate
	if err := json.Unmarshal(raw, &intermediate); err != nil {
		return err
	}

	r.Path = intermediate.Path
	r.Format = intermediate.Format
	r.Spec = intermediate.Spec
	r.Size = intermediate.Size
	r.Redacted = intermediate.Redacted
	if intermediate.Error != "" {
		r.Error = errors.New(intermediate.Error)
	}
	return nil
}
