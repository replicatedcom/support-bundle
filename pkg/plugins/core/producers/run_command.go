package producers

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func RunCommand(opts types.CoreRunCommandOptions) types.StreamsProducer {
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

/*
	stdoutDest := filepath.Join(rootDir, cr.StdoutDest)
	stdout, err := os.Create(stdoutDest)
	if err != nil {
		result.Error = err
	} else {
		cmd.Stdout = stdout
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		result.Error = err
		return []types.Result{result}
	}

	if err := cmd.Wait(); err != nil {
		result.Error = err
		return []types.Result{result}
	}
	result.Pathname = stdoutDest

	errOut := stderr.String()
	if errOut != "" {
		result.Error = errors.New(errOut)
	}

	return []types.Result{result}
}
*/
