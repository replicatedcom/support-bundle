package v1

type RawRequirement struct {
	Conditions []Condition `json:"conditions" yaml:"conditions" hcl:"conditions"`
}
