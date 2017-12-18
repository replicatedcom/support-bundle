package producers

import "k8s.io/client-go/kubernetes"

type Kubernetes struct {
	client *kubernetes.Clientset
}

func New(client *kubernetes.Clientset) *Kubernetes {
	return &Kubernetes{client}
}
