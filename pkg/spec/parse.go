package spec

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func Parse(doc []byte) ([]types.Spec, error) {
	d := &types.Doc{}
	if err := yaml.Unmarshal(doc, d); err != nil {
		return nil, errors.Wrap(err, "parse yaml spec")
	}

	return d.Specs, nil
}
