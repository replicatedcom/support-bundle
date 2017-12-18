package producers

import (
	"context"
	"io"
	"os"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func ReadFile(opts types.CoreReadFileOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		r, err := os.Open(opts.Filepath)
		if err != nil {
			return nil, err
		}
		return map[string]io.Reader{"contents": r}, nil
	}
}
