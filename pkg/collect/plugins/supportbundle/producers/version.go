package producers

import (
	"context"
	"sync"

	"github.com/replicatedcom/support-bundle/pkg/version"
)

var (
	versionInitOnce sync.Once
)

func (s *SupportBundle) Version(ctx context.Context) (interface{}, error) {
	versionInitOnce.Do(version.Init)
	return version.GetBuild(), nil
}
