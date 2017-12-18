package producers

import (
	"context"
	"errors"

	"github.com/replicatedcom/support-bundle/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) Resource(rtype string, ns string) types.StructuredProducer {
	return func(ctx context.Context) (interface{}, error) {
		opts := metav1.ListOptions{}
		c := k.client

		switch rtype {
		case "storageclasses":
			return c.StorageV1().StorageClasses().List(opts)
		case "statefulsets":
			return c.AppsV1beta2().StatefulSets(ns).List(opts)
		case "services":
			return c.CoreV1().Services(ns).List(opts)
		case "secrets":
			return c.CoreV1().Secrets(ns).List(opts)
		case "resourcequotas":
			return c.CoreV1().ResourceQuotas(ns).List(opts)
		case "replicationcontrollers":
			return c.CoreV1().ReplicationControllers(ns).List(opts)
		case "replicasets":
			return c.AppsV1beta2().ReplicaSets(ns).List(opts)
		case "podtemplates":
			return c.CoreV1().PodTemplates(ns).List(opts)
		case "podsecuritypolicies":
			return c.ExtensionsV1beta1().PodSecurityPolicies().List(opts)
		case "pods":
			return c.CoreV1().Pods(ns).List(opts)
		case "persistentvolumes":
			return c.CoreV1().PersistentVolumes().List(opts)
		case "persistentvolumeclaims":
			return c.CoreV1().PersistentVolumeClaims(ns).List(opts)
		case "nodes":
			return c.CoreV1().Nodes().List(opts)
		case "networkpolicies":
			return c.NetworkingV1().NetworkPolicies(ns).List(opts)
		case "limitranges":
			return c.CoreV1().LimitRanges(ns).List(opts)
		case "jobs":
			return c.BatchV1().Jobs(ns).List(opts)
		case "ingresses":
			return c.ExtensionsV1beta1().Ingresses(ns).List(opts)
		case "horizontalpodautoscalers":
			return c.AutoscalingV1().HorizontalPodAutoscalers(ns).List(opts)
		case "events":
			return c.CoreV1().Events(ns).List(opts)
		case "endpoints":
			return c.CoreV1().Endpoints(ns).List(opts)
		case "deployments":
			return c.AppsV1beta2().Deployments(ns).List(opts)
		case "daemonsets":
			return c.AppsV1beta2().DaemonSets(ns).List(opts)
		case "configmaps":
			return c.CoreV1().ConfigMaps(ns).List(opts)
		case "componentstatuses":
			return c.CoreV1().ComponentStatuses().List(opts)
		default:
			return nil, errors.New("Unknown resource type, must be one of: storageclasses, statefulsets, services, secrets, resourcequotas, replicationcontrollers, replicasets, podtemplates, podsecuritypolicies, pods, persistentvolumes, persistentvolumeclaims, nodes, networkpolicies, limitranges, jobs, ingresses, horizontalpodautoscalers, events, endpoints, deployments, daemonsets, configmaps, componentstatuses")
		}
	}
}
