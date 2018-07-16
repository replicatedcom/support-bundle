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
	config := dockertypes.ContainerCreateConfig{
		Config: &dockercontainertypes.Config{
			Image:      image,
			Entrypoint: []string{},
			Cmd:        []string{"test", "-e", "/host/ws"},
		},
		HostConfig: &dockercontainertypes.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:/host", filepath.Dir(filename)),
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
	}
	return cmdErr.Error != nil && cmdErr.StatusCode == 0, cmdErr.Error
}
