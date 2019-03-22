package util

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func ReadFile(ctx context.Context, client docker.CommonAPIClient, image, filename string, securityOpt []string) (io.ReadCloser, error) {
	if _, err := FileExists(ctx, client, image, filename, securityOpt); err != nil {
		return nil, err
	}

	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	config := dockertypes.ContainerCreateConfig{
		Config: &dockercontainertypes.Config{
			Image:      image,
			Entrypoint: []string{"tar"},
			Cmd:        []string{"-O", "-c", "-C", filepath.Join("/host", dir), base},
		},
		HostConfig: &dockercontainertypes.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", dir, filepath.Join("/host", dir)),
			},
			SecurityOpt: securityOpt,
		},
	}
	stdoutR, stderrR, cmdErrCh, err := ContainerRun(ctx, client, config, false)
	if err != nil {
		return nil, errors.Wrap(err, "container run")
	}

	// Error channel will not receive a message
	// until stdout buffer purged, so for large files
	// we ignore this error after 1s since stdout
	// not read until later
	var cmdErr ContainerCmdError
	select {
	case cmdErr = <-cmdErrCh:
	case <-time.After(time.Second):
	}

	stderrR.Close()
	if cmdErr.Error != nil {
		return stdoutR, cmdErr.Error
	} else if cmdErr.StatusCode == 2 {
		return stdoutR, errors.New("file not found")
	} else if cmdErr.StatusCode != 0 {
		return stdoutR, fmt.Errorf("error status code %d", cmdErr.StatusCode)
	}
	return stdoutR, nil
}
