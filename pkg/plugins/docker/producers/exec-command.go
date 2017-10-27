package producers

import (
	"context"
	"io"
	"log"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// ExecCommand returns stdout and stderr results.
func (d *Docker) ExecCommand(containerID string, cmd []string) types.StreamsProducer {
	return func(ctx context.Context) (io.Reader, io.Reader, error) {
		execOpts := dockertypes.ExecConfig{
			Cmd:          cmd,
			AttachStderr: true,
			AttachStdout: true,
			AttachStdin:  false,
		}

		execInstance, err := d.client.ContainerExecCreate(ctx, containerID, execOpts)
		if err != nil {
			return nil, nil, err
		}

		att, err := d.client.ContainerExecAttach(ctx, execInstance.ID, execOpts)
		if err != nil {
			return nil, nil, err
		}

		execStartOpts := dockertypes.ExecStartCheck{
			Detach: false,
			Tty:    false,
		}
		if err := d.client.ContainerExecStart(ctx, execInstance.ID, execStartOpts); err != nil {
			att.Close()
			return nil, nil, err
		}

		stdoutR, stdoutW := io.Pipe()
		stderrR, stderrW := io.Pipe()

		go func() {
			// TODO context interruptible
			_, err = stdcopy.StdCopy(stdoutW, stderrW, att.Reader)
			att.Close()
			if err != nil {
				stdoutW.CloseWithError(err)
				stderrW.CloseWithError(err)
				return
			}
			if err := stdoutW.Close(); err != nil {
				log.Print(err)
			}
			if err := stderrW.Close(); err != nil {
				log.Print(err)
			}
		}()

		return stdoutR, stderrR, nil
	}
}

func (d *Docker) ExecCommandByName(containerName string, cmd []string) types.StreamsProducer {
	return func(ctx context.Context) (io.Reader, io.Reader, error) {
		containerID, err := d.getContainerID(ctx, containerName)
		if err != nil {
			return nil, nil, err
		}
		return d.ExecCommand(containerID, cmd)(ctx)
	}
}
