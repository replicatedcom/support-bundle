package producers

import (
	"context"
)

func (d *Docker) Version(ctx context.Context) (interface{}, error) {
	return d.client.ServerVersion(ctx)
}
