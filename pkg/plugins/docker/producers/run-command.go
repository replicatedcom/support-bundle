package producers

import (
	"context"
	"io"
	"log"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// RunCommand returns stdout and stderr results.
func (d *Docker) RunCommand(image string, cmd []string, env []string, binds []string) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		if _, _, err := d.client.ImageInspectWithRaw(ctx, image); err != nil {
			if docker.IsErrNotFound(err) {
				// TODO: support for image pull
			}
			return nil, errors.Wrapf(err, "image inspect %s", image)
		}

		createOpts := dockertypes.ContainerCreateConfig{
			Config: &container.Config{
				AttachStderr: true,
				AttachStdout: true,
				Cmd:          strslice.StrSlice(cmd),
				Image:        image,
				Env:          env,
			},
			HostConfig: &container.HostConfig{
				Binds: binds,
			},
		}
		containerInstance, err := d.client.ContainerCreate(ctx, createOpts.Config, createOpts.HostConfig, nil, "")
		if err != nil {
			return nil, errors.Wrap(err, "container create")
		}
		go func() {
			<-ctx.Done()
			d.client.ContainerRemove(context.Background(), containerInstance.ID, dockertypes.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			})
		}()

		resp, err := d.client.ContainerAttach(ctx, containerInstance.ID, dockertypes.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		})
		if err != nil {
			return nil, errors.Wrap(err, "container attach")
		}

		if err := d.client.ContainerStart(ctx, containerInstance.ID, dockertypes.ContainerStartOptions{}); err != nil {
			resp.Close()
			return nil, errors.Wrap(err, "container start")
		}

		stdoutR, stdoutW := io.Pipe()
		stderrR, stderrW := io.Pipe()

		go func() {
			// TODO context interruptible
			_, err = stdcopy.StdCopy(stdoutW, stderrW, resp.Reader)
			resp.Close()
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

		readers := make(map[string]io.Reader)
		readers["stdout"] = stdoutR
		readers["stderr"] = stderrR
		return readers, nil
	}
}
