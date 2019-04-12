package condition

type Compare struct {
	Eq  string `json:"eq,omitempty" yaml:"eq,omitempty" hcl:"eq,omitempty"`
	Lt  string `json:"lt,omitempty" yaml:"lt,omitempty" hcl:"lt,omitempty"`
	Lte string `json:"lte,omitempty" yaml:"lte,omitempty" hcl:"lte,omitempty"`
	Gt  string `json:"gt,omitempty" yaml:"gt,omitempty" hcl:"gt,omitempty"`
	Gte string `json:"gte,omitempty" yaml:"gte,omitempty" hcl:"gte,omitempty"`
}
