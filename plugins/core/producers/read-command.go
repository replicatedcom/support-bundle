package producers

import (
	"context"
	"os/exec"
	"strings"

	"github.com/replicatedcom/support-bundle/types"
)

func ReadCommand(command string, args ...string) types.BytesProducer {
	return func(ctx context.Context) ([]byte, error) {
		return exec.CommandContext(ctx, command, strings.Join(args, " ")).Output()
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
