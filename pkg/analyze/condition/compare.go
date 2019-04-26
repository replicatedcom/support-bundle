package condition

type Compare struct {
	Eq  interface{}   `json:"eq,omitempty" yaml:"eq,omitempty" hcl:"eq,omitempty"`
	Lt  interface{}   `json:"lt,omitempty" yaml:"lt,omitempty" hcl:"lt,omitempty"`
	Lte interface{}   `json:"lte,omitempty" yaml:"lte,omitempty" hcl:"lte,omitempty"`
	Gt  interface{}   `json:"gt,omitempty" yaml:"gt,omitempty" hcl:"gt,omitempty"`
	Gte interface{}   `json:"gte,omitempty" yaml:"gte,omitempty" hcl:"gte,omitempty"`
	In  []interface{} `json:"in,omitempty" yaml:"in,omitempty" hcl:"in,omitempty"`
}
