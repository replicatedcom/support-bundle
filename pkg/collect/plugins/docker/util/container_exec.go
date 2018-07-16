package util

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

func ContainerExec(ctx context.Context, client docker.CommonAPIClient, container string, config dockertypes.ExecConfig) (io.ReadCloser, io.ReadCloser, <-chan ContainerCmdError, error) {
	config.AttachStdout = true
	config.AttachStderr = true
	exec, err := client.ContainerExecCreate(ctx, container, config)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "exec create")
	}

	startConfig := dockertypes.ExecStartCheck{
		Detach: config.Detach,
		Tty:    config.Tty,
	}
	resp, err := client.ContainerExecAttach(ctx, exec.ID, startConfig)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "exec attach")
	}

	if err := client.ContainerExecStart(ctx, exec.ID, startConfig); err != nil {
		resp.Close()
		return nil, nil, nil, errors.Wrap(err, "exec start")
	}

	doneCh := make(chan struct{})

	errorCh := make(chan ContainerCmdError, 1)
	go func() {
		<-doneCh
		resp, err := client.ContainerExecInspect(ctx, exec.ID)
		if err != nil {
			errorCh <- ContainerCmdError{0, err}
		} else {
			errorCh <- ContainerCmdError{resp.ExitCode, nil}
		}
	}()

	stdoutR, stdoutW := io.Pipe()
	stderrR, stderrW := io.Pipe()

	go func() {
		defer close(doneCh)

		// TODO context interruptible
		_, err := stdcopy.StdCopy(stdoutW, stderrW, resp.Reader)
		resp.Close()
		if err := stdoutW.CloseWithError(err); err != nil {
			jww.ERROR.Printf("Failed to close stdout writer: %v", err)
		}
		if err := stderrW.CloseWithError(err); err != nil {
			jww.ERROR.Printf("Failed to close stderr writer: %v", err)
		}
	}()

	return stdoutR, stderrR, errorCh, nil
}
