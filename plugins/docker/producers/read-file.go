package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/types"
)

func (d *Docker) ReadFile(containerID string, path string) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		readcloser, _, err := d.client.CopyFromContainer(ctx, containerID, path)
		if err != nil {
			return nil, err
		}
		return readcloser, err
	}
}
