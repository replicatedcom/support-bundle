package kubernetes

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/kubernetes/planners"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/kubernetes/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Kubernetes struct {
	planner *planners.Kubernetes
}

func New(client kubernetes.Interface, config *restclient.Config) (*Kubernetes, error) {
	producers := producers.New(client, config)
	return &Kubernetes{
		planner: planners.New(producers),
	}, nil
}

func (p *Kubernetes) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.KubernetesAPIVersions != nil:
		return p.planner.APIVersions
	case spec.KubernetesClusterInfo != nil:
		return p.planner.ClusterInfo
	case spec.KubernetesContainerCp != nil:
		return p.planner.ContainerCp
	case spec.KubernetesLogs != nil:
		return p.planner.Logs
	case spec.KubernetesVersion != nil:
		return p.planner.Version
	case spec.KubernetesResourceList != nil:
		return p.planner.ResourceList
	default:
		return nil
	}
}
