package resolver

import (
	"strings"

	getter "github.com/hashicorp/go-getter"
	"github.com/spf13/afero"
)

var (
	fileDetector = new(getter.FileDetector)
)

type FileResolver struct {
}

func (g *FileResolver) Detect(name string) (string, bool, error) {
	forced, ok, err := fileDetector.Detect(name, "")
	return strings.TrimPrefix(forced, "file://"), ok, err
}

func (g *FileResolver) Resolve(name string) (afero.Fs, string, error) {
	return afero.NewOsFs(), name, nil
}
