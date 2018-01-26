package producers

import (
	"bytes"
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
)

func (d *Docker) StackServiceLogs(namespace string) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		services, err := d.client.ServiceList(ctx, dockertypes.ServiceListOptions{Filters: dockerNamespaceFilter(namespace)})
		if err != nil {
			return nil, err
		}

		serviceLogs := make(map[string]io.Reader)
		for _, service := range services {
			reader, err := d.client.ServiceLogs(ctx, service.ID,
				dockertypes.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true})
			if err != nil {
				gotError := bytes.NewReader([]byte("error: " + err.Error()))
				serviceLogs[service.Spec.Name+".error"] = gotError
				continue
			}

			demuxed, _ := util.DemuxLogs(ctx, reader, service.Spec.Name+".log")
			for k, v := range demuxed {
				serviceLogs[k] = v
			}
		}
		return serviceLogs, nil
	}
}
