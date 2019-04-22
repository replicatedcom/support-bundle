package variable

import (
	"io"
	"strings"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var (
	_ Interface = new(CoreReadFilePath)
)

type CoreReadFilePath struct {
	Paths []string `json:"paths" yaml:"path" hcl:"paths"`

	Distiller `json:",inline" yaml:",inline" hcl:",inline"`
}

func (v *CoreReadFilePath) MatchResults(index []collecttypes.Result) []collecttypes.Result {
	var results []collecttypes.Result
	for _, result := range index {
		if result.Size == 0 || result.Spec.CoreReadFile == nil {
			continue
		}
		for _, path := range v.Paths {
			if strings.TrimLeft(result.Spec.CoreReadFile.Filepath, "/") == strings.TrimLeft(path, "/") {
				results = append(results, result)
			}
		}
	}
	return results
}

func (v *CoreReadFilePath) DistillReader(r io.Reader, result collecttypes.Result) (interface{}, error) {
	return v.Distiller.Distill(r)
}

func (v *CoreReadFilePath) ExtractValue(distilled interface{}, data interface{}) (interface{}, error) {
	return distilled, nil
}
