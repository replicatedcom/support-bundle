package producers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	dockernetworktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/pkg/errors"
)

func (d *Docker) ReadHostFile(image string, path string) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		containerConfig := &dockercontainertypes.Config{
			Image:      image,
			Entrypoint: []string{"tar"},
			Cmd:        []string{"-O", "-c", "-C", "/host", filepath.Base(path)},
		}

		hostSettings := &dockercontainertypes.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", path, filepath.Join("/host", filepath.Base(path))),
			},
		}
		networkSettings := &dockernetworktypes.NetworkingConfig{}

		containerInstance, err := d.client.ContainerCreate(ctx, containerConfig, hostSettings, networkSettings, "")
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
		var stderrW bytes.Buffer

		go func() {
			// This will also be terminated when ctx.Done() fires and container is removed.
			var err error
			defer func() {
				// FIXME: error will be "io: read/write on closed pipe" even if there was no error
				// (see comment on StdCopy)
				stdoutW.CloseWithError(err)
				if err != nil {
					jww.ERROR.Printf("error copying %s: %v: %s\n", path, err, stderrW.Bytes())
				}
			}()

			// FIXME: This always returns error when stdoutW is connected to a tar reader
			// because tar reader is closed first and StdCopy is a magical multiplexer....
			_, err = stdcopy.StdCopy(stdoutW, &stderrW, resp.Reader)
		}()

		return stdoutR, nil
	}
}
