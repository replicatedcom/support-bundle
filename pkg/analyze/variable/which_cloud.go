package variable

import (
	"io"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	_ Interface = new(WhichCloud)
)

type WhichCloud struct {
	VariableRef string `json:"variableRef" yaml:"variableRef" hcl:"variableRef"`
}

func (v *WhichCloud) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return nil, nil
}

func (v *WhichCloud) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return nil, nil
}

func (v *WhichCloud) ExtractValue(distilled interface{}) (interface{}, error) {
	return nil, nil
}

func (v *WhichCloud) Evaluate(data map[string]interface{}) (interface{}, error) {
	address, ok := data[v.VariableRef]
	if !ok || address == "" {
		return nil, nil
	}
	return nil, nil // TODO
}
