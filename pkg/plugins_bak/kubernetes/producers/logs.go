package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/types"
	"k8s.io/api/core/v1"
)

func (k *Kubernetes) Logs(pod string, container string, ns string) types.StreamProducer {
	if ns == "" {
		ns = "default"
	}
	return func(ctx context.Context) (io.Reader, error) {
		req := k.client.CoreV1().Pods(ns).GetLogs(pod, &v1.PodLogOptions{Container: container})
		return req.Stream()
	}
}
