package producers

import (
	"context"
	"io"

	coreutil "github.com/replicatedcom/support-bundle/pkg/plugins/core/util"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (c *Core) RunCommand(opts types.CoreRunCommandOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		stdoutR, stderrR, _, err := coreutil.RunCommand(ctx, opts)
		if err != nil {
			return nil, err
		}
		return map[string]io.Reader{
			"stdout": stdoutR,
			"stderr": stderrR,
		}, nil
	}
}
