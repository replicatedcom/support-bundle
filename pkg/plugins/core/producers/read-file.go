package producers

import (
	"context"
	"io/ioutil"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

func ReadFile(src string) types.BytesProducer {
	return func(ctx context.Context) ([]byte, error) {
		return ioutil.ReadFile(src)
	}
}
