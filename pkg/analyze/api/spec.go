package api

// Spec is the top level Analyze document that defines an set of collectors and analysis
type Spec struct {
	Collect Collect `json:"collect" yaml:"collect" hcl:"collect"`
	Analyze Analyze `json:"analyze" yaml:"analyze" hcl:"analyze"`
}
