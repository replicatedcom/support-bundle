package api

import (
	"github.com/ghodss/yaml"
)

// DeserializeSpec will deserialize a analyze spec.
// Possibly support more formats?
func DeserializeSpec(b []byte) (spec Spec, err error) {
	err = yaml.Unmarshal(b, &spec)
	return
}
