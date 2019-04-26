package variable

import (
	"io"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/templates"
)

var (
	_ Interface = new(Eval)
)

type Eval string

func (v *Eval) MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error) {
	return nil, nil
}

func (v *Eval) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return nil, nil
}

func (v *Eval) ExtractValue(distilled interface{}) (interface{}, error) {
	return nil, nil
}

func (v *Eval) Evaluate(data map[string]interface{}) (interface{}, error) {
	return templates.String(string(*v), data)
}
