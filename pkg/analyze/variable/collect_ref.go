package variable

import (
	"io"

	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
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

func (v *CollectRef) MatchResults(bundleReader bundlereader.BundleReader) []collecttypes.Result {
	return bundleReader.ResultsFromRef(v.Ref)
}

func (v *CollectRef) ExtractValue(r io.Reader, result collecttypes.Result, data interface{}) (interface{}, error) {
	return v.Distiller.Distill(r)
}
