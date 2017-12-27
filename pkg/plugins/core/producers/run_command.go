package producers

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (c *Core) RunCommand(opts types.CoreRunCommandOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		cmd := exec.CommandContext(ctx, opts.Name, opts.Args...)
		cmd.Env = append(os.Environ(), opts.Env...)
		cmd.Dir = opts.Dir

		prOut, pwOut := io.Pipe()
		prErr, pwErr := io.Pipe()

		cmd.Stdout = pwOut
		cmd.Stderr = pwErr

		streams := map[string]io.Reader{
			"stdout": prOut,
			"stderr": prErr,
		}

		go func() {
			err := cmd.Run()
			pwOut.CloseWithError(err)
			pwErr.CloseWithError(err)
		}()

		return streams, nil
	}
}
