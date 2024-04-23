package util

import (
	"context"
	"io"
	"io/ioutil"

	dockertypes "github.com/docker/docker/api/types"
	dockerbackendtypes "github.com/docker/docker/api/types/backend"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

func ContainerRun(ctx context.Context, client docker.CommonAPIClient, config dockerbackendtypes.ContainerCreateConfig, enablePull bool) (io.ReadCloser, io.ReadCloser, <-chan ContainerCmdError, error) {
	image := config.Config.Image
	if _, _, err := client.ImageInspectWithRaw(ctx, image); err != nil {
		if !enablePull || !docker.IsErrNotFound(err) {
			return nil, nil, nil, errors.Wrapf(err, "image inspect %s", image)
		}
		resp, err := client.ImagePull(ctx, image, dockertypes.ImagePullOptions{
			// TODO: registry auth
		})
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "image pull %s", image)
		}
		if _, err := ioutil.ReadAll(resp); err != nil {
			return nil, nil, nil, err
		}
		resp.Close()
	}

	config.Config.AttachStdout = true
	config.Config.AttachStderr = true
	containerInstance, err := client.ContainerCreate(ctx, config.Config, config.HostConfig, config.NetworkingConfig, nil, config.Name)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "container create")
	}
	go func() {
		<-ctx.Done()
		client.ContainerRemove(context.Background(), containerInstance.ID, dockertypes.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
	}()

	resp, err := client.ContainerAttach(ctx, containerInstance.ID, dockertypes.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "container attach")
	}

	okCh, errCh := client.ContainerWait(ctx, containerInstance.ID, dockercontainertypes.WaitConditionNextExit)

	if err := client.ContainerStart(ctx, containerInstance.ID, dockertypes.ContainerStartOptions{}); err != nil {
		resp.Close()
		return nil, nil, nil, errors.Wrap(err, "container start")
	}

	errorCh := make(chan ContainerCmdError, 1)
	go func() {
		select {
		case body := <-okCh:
			var err error
			if body.Error != nil {
				err = errors.New(body.Error.Message)
			}
			errorCh <- ContainerCmdError{int(body.StatusCode), err}

		case err := <-errCh:
			errorCh <- ContainerCmdError{0, err}
		}
	}()

	stdoutR, stdoutW := io.Pipe()
	stderrR, stderrW := io.Pipe()

	go func() {
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
