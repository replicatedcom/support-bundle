package producers

import (
	"context"

	"github.com/replicatedcom/support-bundle/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) ClusterInfo() types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		opts := metav1.ListOptions{}
		opts.LabelSelector = "kubernetes.io/cluster-service=true"

		return k.client.CoreV1().Services(metav1.NamespaceSystem).List(opts)
	}
}
