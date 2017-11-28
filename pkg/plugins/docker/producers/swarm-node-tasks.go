package producers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/types"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (d *Docker) SwarmNodeTasks() types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		nodes, err := d.client.NodeList(ctx, dockertypes.NodeListOptions{})
		if err != nil {
			return nil, err
		}

		nodeTasks := make(map[string]io.Reader)
		for _, node := range nodes {
			filter := filters.NewArgs()
			filter.Add("node", node.ID)

			tasks, err := d.client.TaskList(ctx, dockertypes.TaskListOptions{Filters: filter})
			if err != nil {
				return nodeTasks, err
			}

			tasksJSON, err := json.Marshal(tasks)
			if err != nil {
				return nodeTasks, err
			}

			nodeTasks["swarm-node-"+node.ID+"-tasks.json"] = bytes.NewReader(tasksJSON)
		}
		return nodeTasks, nil
	}
}
