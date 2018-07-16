package util

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type CmdError struct {
	StatusCode int
	Error      error
}

func RunCommand(ctx context.Context, opts types.CoreRunCommandOptions) (io.Reader, io.Reader, <-chan CmdError, error) {
	cmd := exec.CommandContext(ctx, opts.Name, opts.Args...)
	cmd.Env = append(os.Environ(), opts.Env...)
	cmd.Dir = opts.Dir

	prOut, pwOut := io.Pipe()
	prErr, pwErr := io.Pipe()

	cmd.Stdout = pwOut
	cmd.Stderr = pwErr

	errorCh := make(chan CmdError, 1)

	go func() {
		err := cmd.Run()
		pwOut.CloseWithError(err)
		pwErr.CloseWithError(err)

		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				errorCh <- CmdError{status.ExitStatus(), err}
				return
			}
		}
		errorCh <- CmdError{0, err}
	}()

	return prOut, prErr, errorCh, nil
}
