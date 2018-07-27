package meta

type Meta struct {
	Name   string            `json:"name,omitempty" yaml:"name,omitempty" hcl:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty" hcl:"labels,omitempty"`
}

// A Ref refers to a collected file
type Ref struct {
	// Ref is how this referenced file can be described in a template function
	// This can be used to support multiple file in a single analyzer
	Ref string `json:"ref,omitempty" yaml:"ref,omitempty" hcl:"ref,omitempty"`
	// Name is a way to reference a collected file by name
	Name string `json:"name,omitempty" yaml:"name,omitempty" hcl:"name,omitempty"`
	// Selector is a way to reference a collected file by labels
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
