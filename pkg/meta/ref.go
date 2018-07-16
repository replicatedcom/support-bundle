package meta

type Meta struct {
	Name   string            `json:"name,omitempty" yaml:"name,omitempty" hcl:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty" hcl:"labels,omitempty"`
}

type Ref struct {
	Ref      string   `json:"ref" yaml:"ref" hcl:"ref"`
	Name     string   `json:"name,omitempty" yaml:"name,omitempty" hcl:"name,omitempty"`
	Selector Selector `json:"selector" yaml:"selector" hcl:"selector"`
}

func RefMatches(ref Ref, meta *Meta) bool {
	if meta == nil {
		return false
	}
	if len(ref.Name) > 0 && ref.Name != meta.Name {
		return false
	}
	if len(ref.Selector) > 0 && !SelectorMatches(ref.Selector, meta.Labels) {
		return false
	}
	return len(ref.Name) > 0 || len(ref.Selector) > 0
}
