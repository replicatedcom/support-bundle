package variable

import (
	"io"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

var (
	_ Interface = new(CollectRef)
)

type CollectRef struct {
	meta.Ref `json:",inline" yaml:",inline" hcl:",inline"`

	Distiller `json:",inline" yaml:",inline" hcl:",inline"`
}

func (v *CollectRef) MatchResults(index []collecttypes.Result) (results []collecttypes.Result) {
	for _, result := range index {
		if meta.RefMatches(v.Ref, result.Spec.Shared().Meta) {
			results = append(results, result)
		}
	}
	return
}

func (v *CollectRef) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return v.Distiller.Distill(r)
}

func (v *CollectRef) ExtractValue(distilled interface{}, data interface{}) (interface{}, error) {
	return distilled, nil
}
