package producers

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) ClusterInfo(ctx context.Context) (interface{}, error) {
	opts := metav1.ListOptions{}
	opts.LabelSelector = "kubernetes.io/cluster-service=true"
	return k.client.CoreV1().Services(metav1.NamespaceSystem).List(ctx, opts)
}
