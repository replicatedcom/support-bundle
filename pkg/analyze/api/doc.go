package api

// Doc is the top level Analyze document that defines an set of collectors and analysis
type Doc struct {
	Analyze Analyze `json:"analyze" yaml:"analyze" hcl:"analyze"`
}
