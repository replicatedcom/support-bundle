package collector

import (
	"context"
	"os"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
)

type Interface interface {
	CollectBundle(ctx context.Context, collect api.Collect, timeout time.Duration, dest string, opts Options) (os.FileInfo, error)
}
