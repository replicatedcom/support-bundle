package api

import (
	"github.com/ghodss/yaml"
)

// DeserializeDoc will deserialize a analyze spec.
// Possibly support more formats?
func DeserializeDoc(b []byte) (doc Doc, err error) {
	err = yaml.Unmarshal(b, &doc)
	return
}
