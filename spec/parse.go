package spec

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/types"
	yaml "gopkg.in/yaml.v2"
)

type Doc struct {
	Specs []types.Spec
}

func Parse(doc []byte) ([]types.Spec, error) {
	d := &Doc{}
	if err := yaml.UnmarshalStrict(doc, d); err != nil {
		return nil, errors.Wrap(err, "parse yaml spec")
	}

	return d.Specs, nil
}
