package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/types"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) Logs() types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		req := k.client.Core().Pods(metav1.NamespaceAll).GetLogs("", &v1.PodLogOptions{})
		return req.Stream()
	}
}
