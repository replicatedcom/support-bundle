package producers

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/client-go/discovery"
)

func (k *Kubernetes) Version(ctx context.Context) (interface{}, error) {
	discoveryClient, ok := k.client.(discovery.DiscoveryInterface)
	if !ok {
		return nil, errors.New("client is not a discovery client")
	}
	return discoveryClient.ServerVersion()
}
