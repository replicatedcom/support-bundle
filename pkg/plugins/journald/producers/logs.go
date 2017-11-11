package producers

import (
	"context"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
	dockerutil "github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// Logs gets journalctl logs by mounting the ephemeral and persistent journald
// data directories into the support bundle container, which has the journalctl
// command. It attempts to find the image and SELinux label to use to run the
// support bundle container by checking if the current process is in a container
// and, if so, copying the image and SELinux label of that container.
func Logs(client *docker.Client, opts types.JournaldLogs) types.BytesProducer {
	return func(ctx context.Context) ([]byte, error) {
		// defaults if values can't be pulled from current container
		image := "registry.replicated.com/library/support-bundle"
		secOpts := []string{"label=type:spc_t"}
		c, err := dockerutil.ThisContainer(context.Background(), client)
		if err == nil {
			image = c.Image
			details, err := client.ContainerInspect(context.Background(), c.ID)
			if err == nil {
				secOpts = details.HostConfig.SecurityOpt
			}
		}

		since := "-7 days"
		if opts.Since != "" {
			since = opts.Since
		}

		config := dockertypes.ContainerCreateConfig{
			Config: &container.Config{
				Entrypoint: []string{"journalctl"},
				Cmd: []string{
					"--reverse",
					"--merge",
					"--no-pager",
					"--since", since,
					"-u",
					opts.Unit,
				},
				Image: image,
			},
			HostConfig: &container.HostConfig{
				Binds: []string{
					"/var/log/journal:/var/log/journal",
					"/run/log/journal:/run/log/journal",
				},
				SecurityOpt: secOpts,
			},
		}

		logs, err := dockerutil.RunCommand(ctx, client, config)
		if err != nil {
			return nil, errors.Wrapf(err, "journald.logs")
		}

		return logs, nil
	}
}
