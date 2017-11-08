package spec

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// left here as a convenience method
func Parse(doc []byte) ([]types.Spec, error) {
	d, err := Unmarshal(doc)
	if err != nil {
		return nil, err
	}
	return d.Specs, nil
}

func Marshal(doc *types.Doc) ([]byte, error) {
	marshaled, err := yaml.Marshal(doc)
	if err != nil {
		return nil, errors.Wrap(err, "marshal yaml spec")
	}
	return marshaled, nil
}

func Unmarshal(doc []byte) (*types.Doc, error) {
	d := &types.Doc{}
	if err := yaml.Unmarshal(doc, d); err != nil {
		return nil, errors.Wrap(err, "parse yaml spec")
	}

	return d, nil
}
