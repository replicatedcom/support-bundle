package lifecycle

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/bundle"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func GenerateTask(task *types.LifecycleTask) Task {
	return func(l *Lifecycle) (bool, error) {

		fileInfo, err := bundle.Generate(l.BundleTasks, time.Duration(time.Second*time.Duration(l.GenerateTimeout)), l.GenerateBundlePath)
		if err != nil {
			return false, errors.Wrap(err, "generating bundle")
		}
		l.FileInfo = fileInfo
		return true, nil
	}
}
