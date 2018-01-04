package producers

import "k8s.io/client-go/kubernetes"

type Kubernetes struct {
	client kubernetes.Interface
}

func New(client kubernetes.Interface) *Kubernetes {
	return &Kubernetes{client}
}
