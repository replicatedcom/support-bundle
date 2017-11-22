package planners

import "github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/producers"

type Kubernetes struct {
	p *producers.Kubernetes

}

func New(p *producers.Kubernetes) *Kubernetes {
	return &Kubernetes{p}
}
