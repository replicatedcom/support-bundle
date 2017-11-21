package kubernetes

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func New() (types.Plugin, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
