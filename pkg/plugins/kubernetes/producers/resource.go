package producers

import (
	"context"
	"errors"

	"github.com/replicatedcom/support-bundle/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) Resource(rtype string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		opts := metav1.ListOptions{}
		c := k.client
		switch rtype {
		case "storageclasses":
			return c.Storage().StorageClasses().List(opts)
		case "statefulsets":
			return c.Apps().StatefulSets(metav1.NamespaceAll).List(opts)
		default:
			return nil, errors.New("Unknown resource type, must be one of: storageclasses, statefulsets, services, secrets, resourcequotas, replicationcontrollers, replicasets, podtemplates, podsecuritypolicies, pods, persistentvolumes, persistentvolumeclaims, nodes, networkpolicies, limitranges, jobs, ingresses, horizontalpodautoscalers, events, endpoints, deployments, daemonsets, configmaps, componentstatuses")
		}
	}
}
