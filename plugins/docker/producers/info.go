package producers

import (
	"context"
)

func (d *Docker) Info(ctx context.Context) (interface{}, error) {
	return d.client.Info(ctx)
}
