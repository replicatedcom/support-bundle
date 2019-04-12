package variable

import (
	"io"

	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type Interface interface {
	MatchResults(bundlereader.BundleReader) []collecttypes.Result
	ExtractValue(r io.Reader, result collecttypes.Result, data interface{}) (interface{}, error)
}
