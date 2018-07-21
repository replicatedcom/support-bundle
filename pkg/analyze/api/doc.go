package api

import (
	"github.com/ghodss/yaml"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

// Doc is the top level Analyze document that defines an set of collectors and analysis
type Doc struct {
	Collect Collect `json:"collect,omitempty" yaml:"collect,omitempty" hcl:"collect,omitempty"`
	Analyze Analyze `json:"analyze" yaml:"analyze" hcl:"analyze"`
}

func (d *Doc) UnmarshalJSON(data []byte) error {
	// TODO: hcl
	var intermediate struct {
		Analyze Analyze             `json:"analyze" yaml:"analyze" hcl:"analyze"`
		Collect Collect             `json:"collect" yaml:"collect" hcl:"collect"`
		Specs   []collecttypes.Spec `json:"specs" yaml:"specs" hcl:"specs"`
	}
	if err := yaml.Unmarshal(data, &intermediate); err != nil {
		return err
	}
	d.Collect.V1 = append(intermediate.Collect.V1, intermediate.Specs...)
	d.Analyze.V1Alpha1 = intermediate.Analyze.V1Alpha1
	return nil
}
