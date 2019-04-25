package variable

import (
	"io"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

// Interface is the Variable interface
// FIXME: This interface is really awkward.
// Variables like eval and whichCloud are only really Evaluators.
type Interface interface {
	ResultDistiller
}

type ResultDistiller interface {
	MatchResults(index []collecttypes.Result) ([]collecttypes.Result, error)
	DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error)
	ExtractValue(distilled interface{}) (interface{}, error)
}

type Evaluator interface {
	Evaluate(data map[string]interface{}) (interface{}, error)
}
