package v1alpha1

type RawAnalyzer struct {
	Conditions []Condition `json:"conditions" yaml:"conditions" hcl:"conditions"`
}

func (a *RawAnalyzer) Validate(spec AnalyzerSpec) error {
	return nil
}

func (a *RawAnalyzer) GetSpec() (AnalyzerSpec, error) {
	return AnalyzerSpec{Raw: a}, nil
}
