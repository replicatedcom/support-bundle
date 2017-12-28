package producers

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
	coreutil "github.com/replicatedcom/support-bundle/pkg/plugins/core/util"
	dockerutil "github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (j *Journald) Logs(opts types.JournaldLogsOptions) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		opts := types.CoreRunCommandOptions{
			Name: "journalctl",
			Args: buildLogsCmdArgs(opts),
		}
		stdoutR, _, _, err := coreutil.RunCommand(ctx, opts)
		if err != nil {
			return nil, errors.Wrap(err, "run command")
		}
		return stdoutR, nil
	}
}

// DockerLogs gets journalctl logs by mounting the ephemeral and persistent
// journald data directories into the support bundle container, which has the
// journalctl command. It attempts to find the image and SELinux label to use to
// run the support bundle container by checking if the current process is in a
// container and, if so, copying the image and SELinux label of that container.
func (j *Journald) DockerLogs(opts types.JournaldLogsOptions) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		container, err := dockerutil.ThisContainer(ctx, j.dockerClient)
		if err != nil {
			return nil, errors.Wrap(err, "this container")
		}

		config := dockertypes.ContainerCreateConfig{
			Config: &dockercontainertypes.Config{
				Image:      container.Image,
				Entrypoint: []string{"journalctl"},
				Cmd:        buildLogsCmdArgs(opts),
			},
			HostConfig: &dockercontainertypes.HostConfig{
				Binds: []string{
					"/var/log/journal:/var/log/journal",
					"/run/log/journal:/run/log/journal",
				},
				SecurityOpt: container.HostConfig.SecurityOpt,
			},
		}

		stdoutR, _, _, err := dockerutil.ContainerRun(ctx, j.dockerClient, config, false)
		if err != nil {
			return nil, errors.Wrap(err, "container run")
		}
		return stdoutR, nil
	}
}

func buildLogsCmdArgs(opts types.JournaldLogsOptions) []string {
	var cmd []string
	if opts.Unit != "" {
		cmd = append(cmd, "--unit", opts.Unit)
	}
	if opts.Since != "" {
		cmd = append(cmd, "--since", opts.Since)
	}
	if opts.Reverse {
		cmd = append(cmd, "--reverse")
	}
	return cmd
}
