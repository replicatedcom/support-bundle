package producers

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Kubernetes struct {
	client kubernetes.Interface
	config *restclient.Config
}

func New(client kubernetes.Interface, config *restclient.Config) *Kubernetes {
	return &Kubernetes{client, config}
}
