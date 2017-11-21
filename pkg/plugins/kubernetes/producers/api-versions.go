package producers

import (
	"context"

	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) APIVersions() types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		groups, err := k.client.ServerGroups()
		if err != nil {
			return nil, errors.Wrap(err, "getting server groups")
		}
		apiVersions := metav1.ExtractGroupVersions(groups)

		return apiVersions, nil
	}
}
