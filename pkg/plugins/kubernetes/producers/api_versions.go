package producers

import (
	"context"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
)

func (k *Kubernetes) APIVersions(ctx context.Context) (interface{}, error) {
	discoveryClient, ok := k.client.(discovery.DiscoveryInterface)
	if !ok {
		return nil, errors.New("client is not a discovery client")
	}
	groups, err := discoveryClient.ServerGroups()
	if err != nil {
		return nil, errors.Wrap(err, "get server groups")
	}
	apiVersions := metav1.ExtractGroupVersions(groups)

	return apiVersions, nil
}
