package kubernetes

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func New() (types.Plugin, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting kubernetes config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "getting kubernetes client")
	}

	p := producers.New(clientset)
	k := planners.New(p)
	return map[string]types.Planner{
		"version":      k.Version,
		"api-versions": k.APIVersions,
		"cluster-info": k.ClusterInfo,
		"resource":     k.Resource,
		"logs":         k.Logs,
	}, nil
}
