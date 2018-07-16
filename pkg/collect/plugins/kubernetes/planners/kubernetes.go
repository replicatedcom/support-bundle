package planners

import "github.com/replicatedcom/support-bundle/pkg/collect/plugins/kubernetes/producers"

type Kubernetes struct {
	producers *producers.Kubernetes
}

func New(kubernetes *producers.Kubernetes) *Kubernetes {
	return &Kubernetes{kubernetes}
}
