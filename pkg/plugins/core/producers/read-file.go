package producers

import (
	"context"
	"io"
	"os"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func ReadFile(src string) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		return os.Open(src)
	}
}
