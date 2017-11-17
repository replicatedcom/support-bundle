package util

import (
	"bytes"
	"context"

	dockertypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
)

// RunCommand starts a container, runs a command, returns stdout as bytes and
// stderr as an error
func RunCommand(ctx context.Context, client *docker.Client, config dockertypes.ContainerCreateConfig) ([]byte, error) {
	containerInstance, err := client.ContainerCreate(ctx, config.Config, config.HostConfig, config.NetworkingConfig, config.Name)
	if err != nil {
		return nil, errors.Wrap(err, "container create")
	}

	localCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-localCtx.Done()
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
		return nil, errors.Wrap(err, "container attach")
	}

	if err := client.ContainerStart(ctx, containerInstance.ID, dockertypes.ContainerStartOptions{}); err != nil {
		resp.Close()
		return nil, errors.Wrap(err, "container start")
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// TODO context interruptible
	_, err = stdcopy.StdCopy(&stdout, &stderr, resp.Reader)
	resp.Close()
	if err != nil {
		return nil, err
	}

	var stderrErr error
	if stderr.Len() > 0 {
		stderrErr = errors.New(stderr.String())
	}

	return stdout.Bytes(), stderrErr
}
