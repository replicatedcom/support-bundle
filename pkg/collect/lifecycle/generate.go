package lifecycle

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/bundle"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type GenerateTask struct {
	Options types.GenerateOptions
}

func (t *GenerateTask) Execute(l *Lifecycle) (bool, error) {

	fileInfo, finalPathname, err := bundle.Generate(l.BundleTasks, time.Duration(time.Second*time.Duration(l.GenerateTimeout)), l.GenerateBundlePath)
	if err != nil {
		return false, errors.Wrap(err, "generating bundle")
	}
	l.FileInfo = fileInfo
	l.RealGeneratedBundlePath = finalPathname

	return true, nil
}
