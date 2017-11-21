package kubernetes

import (
	"os"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func New() types.Plugin {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)

	clientset, _ := kubernetes.NewForConfig(config)

	p := producers.New(clientset)
	k := planners.New(p)
	return map[string]types.Planner{
		"version":      k.Version,
		"api-versions": k.APIVersions,
		"cluster-info": k.ClusterInfo,
		"resource":     k.Resource,
	}
}
