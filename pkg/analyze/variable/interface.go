package variable

import (
	"io"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type Interface interface {
	ResultDistiller
	ValueExtractor
}

type ResultDistiller interface {
	MatchResults(index []collecttypes.Result) []collecttypes.Result
	DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error)
}

type ValueExtractor interface {
	ExtractValue(distilled interface{}, data interface{}) (interface{}, error)
}
