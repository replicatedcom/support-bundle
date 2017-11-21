package producers

import (
	"context"
	"github.com/replicatedcom/support-bundle/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) APIVersions() types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		groups, _ := k.client.ServerGroups()
		apiVersions := metav1.ExtractGroupVersions(groups)

		return apiVersions, nil
	}
}
