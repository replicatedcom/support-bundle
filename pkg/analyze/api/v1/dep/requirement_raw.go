package v1

type RawRequirement struct {
	Conditions []Condition `json:"conditions" yaml:"conditions" hcl:"conditions"`
}

func (r *RawRequirement) Validate(spec AnalyzerSpec) error {
	return nil
}

func (r *RawRequirement) GetRawSpec() (RawSpec, error) {
	return RawSpec{Raw: r}, nil
}
