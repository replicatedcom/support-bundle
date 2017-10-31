package producers

import (
	"context"
	"io"
	"io/ioutil"
	"log"

	dockertypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// RunCommand returns stdout and stderr results.
func (d *Docker) RunCommand(opts types.DockerRunCommandOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		image := opts.ContainerCreateConfig.Config.Image
		if _, _, err := d.client.ImageInspectWithRaw(ctx, image); err != nil {
			if !opts.EnablePull || !docker.IsErrNotFound(err) {
				return nil, errors.Wrapf(err, "image inspect %s", image)
			}
			resp, err := d.client.ImagePull(ctx, image, dockertypes.ImagePullOptions{
			// TODO: registry auth
			})
			if err != nil {
				return nil, errors.Wrapf(err, "image pull %s", image)
			}
			if _, err := ioutil.ReadAll(resp); err != nil {
				return nil, err
			}
			resp.Close()
		}

		containerInstance, err := d.client.ContainerCreate(ctx, opts.ContainerCreateConfig.Config, opts.ContainerCreateConfig.HostConfig, opts.ContainerCreateConfig.NetworkingConfig, opts.ContainerCreateConfig.Name)
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
