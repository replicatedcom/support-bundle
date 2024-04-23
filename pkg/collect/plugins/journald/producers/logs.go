package producers

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	dockerbackendtypes "github.com/docker/docker/api/types/backend"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
	coreutil "github.com/replicatedcom/support-bundle/pkg/collect/plugins/core/util"
	dockerutil "github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	jww "github.com/spf13/jwalterweatherman"
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

		dir := "/var/log/journal"
		count, err := dockerutil.FileCount(ctx, j.dockerClient, container.Image, dir, container.HostConfig.SecurityOpt)
		if err != nil {
			jww.DEBUG.Printf("Failed to get file count for directory %s: %v", dir, err)
			dir = "/run/log/journal"
		} else if count == 0 {
			jww.DEBUG.Printf("No files in directory %s, deferring to /run/log/journal for journald logs", dir)
			dir = "/run/log/journal"
		}

		config := dockerbackendtypes.ContainerCreateConfig{
			Config: &dockercontainertypes.Config{
				Image:      container.Image,
				Entrypoint: []string{"journalctl"},
				Cmd:        append(buildLogsCmdArgs(opts), "-D", "/run/log/journal"),
			},
			HostConfig: &dockercontainertypes.HostConfig{
				Binds: []string{
					fmt.Sprintf("%s:/run/log/journal", dir),
				},
				SecurityOpt: container.HostConfig.SecurityOpt,
			},
		}

		stdoutR, stderrR, _, err := dockerutil.ContainerRun(ctx, j.dockerClient, config, false)
		if err != nil {
			return nil, errors.Wrap(err, "container run")
		}
		go func() {
			io.Copy(ioutil.Discard, stderrR)
			stderrR.Close()
		}()
		// FIXME: stdoutR never closed
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
