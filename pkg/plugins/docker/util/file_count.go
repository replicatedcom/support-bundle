package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func FileCount(ctx context.Context, client docker.CommonAPIClient, image string, dirname string, securityOpt []string) (int, error) {
	config := dockertypes.ContainerCreateConfig{
		Config: &dockercontainertypes.Config{
			Image:      image,
			Entrypoint: []string{"sh"},
			Cmd:        []string{"-c", "ls /host/ws 2>/dev/null | wc -l"},
		},
		HostConfig: &dockercontainertypes.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:/host/ws", dirname),
			},
			SecurityOpt: securityOpt,
		},
	}

	stdoutR, stderrR, cmdErrCh, err := ContainerRun(ctx, client, config, false)
	if err != nil {
		return 0, errors.Wrap(err, "container run")
	}
	cmdErr := <-cmdErrCh
	b, err := ioutil.ReadAll(stdoutR)
	stdoutR.Close()
	stderrR.Close()
	if err != nil {
		return 0, errors.Wrap(err, "read stdout")
	}
	if cmdErr.Error != nil {
		return 0, cmdErr.Error
	} else if cmdErr.StatusCode == 2 {
		return 0, errors.New("file not found")
	} else if cmdErr.StatusCode != 0 {
		return 0, fmt.Errorf("error status code %d", cmdErr.StatusCode)
	}
	return strconv.Atoi(strings.TrimSpace(string(b)))
}
