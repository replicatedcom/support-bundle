package producers

import (
	"context"
	"errors"

	"github.com/replicatedcom/support-bundle/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) Resource(rtype string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		ns := metav1.NamespaceAll
		opts := metav1.ListOptions{}
		c := k.client
		switch rtype {
		case "storageclasses":
			return c.Storage().StorageClasses().List(opts)
		case "statefulsets":
			return c.Apps().StatefulSets(ns).List(opts)
		case "services":
			return c.CoreV1().Services(ns).List(opts)
		case "secrets":
			return c.CoreV1().Secrets(ns).List(opts)
		case "resourcequotas":
			return c.CoreV1().ResourceQuotas(ns).List(opts)
		case "replicationcontrollers":
			return c.CoreV1().ReplicationControllers(ns).List(opts)
		case "replicasets":
			return c.Apps().ReplicaSets(ns).List(opts)
		case "podtemplates":
			return c.CoreV1().PodTemplates(ns).List(opts)
		case "podsecuritypolicies":
			return c.Extensions().PodSecurityPolicies(ns).List(opts)
		case "pods":
			return c.Core().Pods(ns).List(opts)
		case "persistentvolumes":
			return c.Core().PersistentVolumes(ns).List(opts)
		case "persistentvolumeclaims":
			return c.Core().PersistentVolumeClaims(ns).List(opts)
		case "nodes":
			return c.Core().Nodes(ns).List(opts)
		case "networkpolicies":
			return c.Networking().NetworkPolicies(ns).List(opts)
		case "limitranges":
			return c.Core().LimitRanges(ns).List(opts)
		case "jobs":
			return c.Batch().Jobs(ns).List(opts)
		case "ingresses":
			return c.Extensions().Ingresses(ns).List(opts)
		case "horizontalpodautoscalers":
			return c.Autoscaling().HorizontalPodAutoscalers(ns).List(opts)
		case "events":
			return c.Core().Events(ns).List(opts)
		case "endpoints":
			return c.Core().Endpoints(ns).List(opts)
		case "deployments":
			return c.Apps().Deployments(ns).List(opts)
		case "daemonsets":
			return c.Apps().DaemonSets(ns).List(opts)
		case "configmaps":
			return c.Core().ConfigMaps(ns).List(opts)
		case "configmaps":
			return c.Core().ComponentStatuses(ns).List(opts)
		default:
			return nil, errors.New("Unknown resource type, must be one of: storageclasses, statefulsets, services, secrets, resourcequotas, replicationcontrollers, replicasets, podtemplates, podsecuritypolicies, pods, persistentvolumes, persistentvolumeclaims, nodes, networkpolicies, limitranges, jobs, ingresses, horizontalpodautoscalers, events, endpoints, deployments, daemonsets, configmaps, componentstatuses")
		}
	}
}
