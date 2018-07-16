package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func (d *Docker) ContainerCp(containerID string, path string) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		r, _, err := d.client.CopyFromContainer(ctx, containerID, path)
		return r, err
	}
}
