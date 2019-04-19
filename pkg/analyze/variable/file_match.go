package variable

import (
	"io"
	"strings"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	_ Interface = new(FileMatch)
)

type FileMatch struct {
	Path string `json:"path" yaml:"path" hcl:"path"`

	Distiller `json:",inline" yaml:",inline" hcl:",inline"`
}

func (v *FileMatch) MatchResults(index []collecttypes.Result) []collecttypes.Result {
	if v.Path == "" {
		return nil
	}
	var results []collecttypes.Result
	for _, result := range index {
		if strings.HasPrefix(strings.TrimPrefix(result.Path, "/"), strings.TrimPrefix(v.Path, "/")) {
			results = append(results, result)
		}
	}
	return results
}

func (v *FileMatch) DistillReader(r io.Reader, result collecttypes.Result) (string, error) {
	return v.Distiller.Distill(r)
}

func (v *FileMatch) ExtractValue(distilled interface{}, data interface{}) (interface{}, error) {
	return distilled, nil
}
