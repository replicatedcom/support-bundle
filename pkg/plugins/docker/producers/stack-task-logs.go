package producers

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
)

func (d *Docker) StackTaskLogs(namespace string) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		tasks, err := d.client.TaskList(ctx, dockertypes.TaskListOptions{Filters: dockerNamespaceFilter(namespace)})
		if err != nil {
			return nil, err
		}

		services, err := d.client.ServiceList(ctx, dockertypes.ServiceListOptions{Filters: dockerNamespaceFilter(namespace)})
		if err != nil {
			return nil, err
		}

		serviceNames := make(map[string]string)
		for _, service := range services {
			serviceNames[service.ID] = service.Spec.Name
		}

		taskLogs := make(map[string]io.Reader)
		for _, task := range tasks {
			serviceName, ok := serviceNames[task.ServiceID]
			if ok {
				if strings.HasPrefix(serviceName, namespace+"_") {
					serviceName = serviceName[len(namespace)+1:]
				}
			} else {
				serviceName = task.ServiceID
			}

			if task.NodeID == "" {
				notScheduled := bytes.NewReader([]byte("task " + task.ID + " is not yet scheduled"))
				taskLogs[serviceName+"-"+task.ID+".error"] = notScheduled
				continue
			}

			//TODO: this needs to be changed to TaskLogs, but that requires a docker library version >=17.05
			reader, err := d.client.ServiceLogs(ctx, task.ID,
				dockertypes.ContainerLogsOptions{ShowStderr: true, ShowStdout: true, Timestamps: true})
			if err != nil {
				gotError := bytes.NewReader([]byte("error: " + err.Error()))
				taskLogs[serviceName+"-"+task.ID+".error"] = gotError
				continue
			}

			taskLogs[serviceName+"-"+task.ID+".log"] = reader
		}
		return taskLogs, nil
	}
}
