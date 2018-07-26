package v1alpha1

type RawAnalyzer struct {
	Conditions []Condition `json:"conditions" yaml:"conditions" hcl:"conditions"`
}
