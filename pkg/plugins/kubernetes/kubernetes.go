package kubernetes

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	planner *planners.Kubernetes
}

func New() (*Kubernetes, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "get kubernetes config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "get kubernetes client")
	}

	producers := producers.New(clientset)
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
