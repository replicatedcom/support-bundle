package producers

import (
	"context"
	"io"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) Logs(opts types.KubernetesLogsOptions) types.StreamProducer {
	ns := opts.Namespace
	if ns == "" {
		ns = metav1.NamespaceDefault
	}

	return func(ctx context.Context) (io.Reader, error) {
		req := k.client.CoreV1().
			Pods(ns).
			GetLogs(opts.Pod, opts.PodLogOptions)
		return req.Stream(ctx)
	}
}
