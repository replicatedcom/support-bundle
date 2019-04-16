package variable

import (
	"io"
	"strings"

	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	_ Interface = new(FileMatch)
)

type FileMatch struct {
	Path string `json:"path" yaml:"path" hcl:"path"`

	Distiller `json:",inline" yaml:",inline" hcl:",inline"`
}

func (v *FileMatch) MatchResults(bundleReader bundlereader.BundleReader) []collecttypes.Result {
	if v.Path == "" {
		return nil
	}
	var results []collecttypes.Result
	for _, result := range bundleReader.GetIndex() {
		if strings.HasPrefix(strings.TrimPrefix(result.Path, "/"), strings.TrimPrefix(v.Path, "/")) {
			results = append(results, result)
		}
	}
	return results
}

func (v *FileMatch) ExtractValue(r io.Reader, result collecttypes.Result, data interface{}) (interface{}, error) {
	return v.Distiller.Distill(r)
}
