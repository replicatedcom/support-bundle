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
	Paths []string `json:"paths" yaml:"path" hcl:"paths"`

	Distiller `json:",inline" yaml:",inline" hcl:",inline"`
}

func (v *FileMatch) MatchResults(index []collecttypes.Result) []collecttypes.Result {
	var results []collecttypes.Result
	for _, result := range index {
		for _, path := range v.Paths {
			if strings.HasPrefix(strings.TrimPrefix(result.Path, "/"), strings.TrimPrefix(path, "/")) {
				results = append(results, result)
			}
		}
	}
	return results
}

func (v *FileMatch) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return v.Distiller.Distill(r)
}

func (v *FileMatch) ExtractValue(distilled interface{}, data interface{}) (interface{}, error) {
	return distilled, nil
}
