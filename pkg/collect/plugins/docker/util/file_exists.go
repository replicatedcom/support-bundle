package util

import (
	"context"
	"fmt"
	"path/filepath"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func FileExists(ctx context.Context, client docker.CommonAPIClient, image string, filename string, securityOpt []string) (bool, error) {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	config := dockertypes.ContainerCreateConfig{
		Config: &dockercontainertypes.Config{
			Image:      image,
			Entrypoint: []string{},
			Cmd:        []string{"test", "-e", filepath.Join("/host", base)},
		},
		HostConfig: &dockercontainertypes.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:/host", dir),
			},
			SecurityOpt: securityOpt,
		},
	}

	stdoutR, stderrR, cmdErrCh, err := ContainerRun(ctx, client, config, false)
	if err != nil {
		return false, errors.Wrap(err, "container run")
	}
	cmdErr := <-cmdErrCh
	stdoutR.Close()
	stderrR.Close()
	if cmdErr.Error != nil {
		return false, cmdErr.Error
	} else if cmdErr.StatusCode == 1 {
		return false, errors.New("file not found")
	} else if cmdErr.StatusCode != 0 {
		return false, fmt.Errorf("error status code %d", cmdErr.StatusCode)
	}
	return true, nil
}
