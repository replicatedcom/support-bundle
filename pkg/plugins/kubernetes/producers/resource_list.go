package producers

import (
	"context"
	"fmt"

	"github.com/replicatedcom/support-bundle/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type resourceListFn func(client kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error)

var resourceListFns = map[string]resourceListFn{
	"componentstatuses": resourceListComponentStatuses, "cs": resourceListComponentStatuses,
	"configmaps": resourceListConfigMaps, "cm": resourceListConfigMaps,
	"daemonsets": resourceListDaemonSets, "ds": resourceListDaemonSets,
	"deployments": resourceListDeployments, "deploy": resourceListDeployments,
	"endpoints": resourceListEndpoints, "ep": resourceListEndpoints,
	"events": resourceListEvents, "ev": resourceListEvents,
	"horizontalpodautoscalers": resourceListHorizontalPodAutoscalers, "hpa": resourceListHorizontalPodAutoscalers,
	"ingresses": resourceListIngresses, "ing": resourceListIngresses,
	"jobs":        resourceListJobs,
	"limitranges": resourceListLimitRanges, "limits": resourceListLimitRanges,
	"namespaces": resourceListNamespaces, "ns": resourceListNamespaces,
	"networkpolicies": resourceListNetworkPolicies, "netpol": resourceListNetworkPolicies,
	"nodes": resourceListNodes, "no": resourceListNodes,
	"persistentvolumeclaims": resourceListPersistentVolumeClaims, "pvc": resourceListPersistentVolumeClaims,
	"persistentvolumes": resourceListPersistentVolumes, "pv": resourceListPersistentVolumes,
	"pods": resourceListPods, "po": resourceListPods,
	"podsecuritypolicies": resourceListPodSecurityPolicies, "psp": resourceListPodSecurityPolicies,
	"podtemplates": resourceListPodTemplates,
	"replicasets":  resourceListReplicaSets, "rs": resourceListReplicaSets,
	"replicationcontrollers": resourceListReplicationControllers, "rc": resourceListReplicationControllers,
	"resourcequotas": resourceListResourceQuotas, "quota": resourceListResourceQuotas,
	"storageclasses":  resourceListStorageClasses,
	"statefulsets":    resourceListStatefulSets,
	"serviceaccounts": resourceListServiceAccounts, "sa": resourceListServiceAccounts,
	"services": resourceListServices, "svc": resourceListServices,
	"secrets": resourceListSecrets,
}

func (k *Kubernetes) ResourceList(opts types.KubernetesResourceListOptions) types.StructuredProducer {
	ns := opts.Namespace

	var listOpts metav1.ListOptions
	if opts.ListOptions != nil {
		listOpts = *opts.ListOptions
	}
	return func(ctx context.Context) (interface{}, error) {
		fn, ok := resourceListFns[opts.Kind]
		if !ok {
			return nil, fmt.Errorf("Unknown resource type, must be one of: %s", resourceListKinds())
		}
		return fn(k.client, ns, listOpts)
	}
}

func resourceListComponentStatuses(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ComponentStatuses().List(opts)
}

func resourceListConfigMaps(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ConfigMaps(ns).List(opts)
}

func resourceListDaemonSets(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().DaemonSets(ns).List(opts)
}

func resourceListDeployments(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().Deployments(ns).List(opts)
}

func resourceListEndpoints(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Endpoints(ns).List(opts)
}

func resourceListEvents(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Events(ns).List(opts)
}

func resourceListHorizontalPodAutoscalers(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AutoscalingV1().HorizontalPodAutoscalers(ns).List(opts)
}

func resourceListIngresses(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.ExtensionsV1beta1().Ingresses(ns).List(opts)
}

func resourceListJobs(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.BatchV1().Jobs(ns).List(opts)
}

func resourceListLimitRanges(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().LimitRanges(ns).List(opts)
}

func resourceListNamespaces(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Namespaces().List(opts)
}

func resourceListNetworkPolicies(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.NetworkingV1().NetworkPolicies(ns).List(opts)
}

func resourceListNodes(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Nodes().List(opts)
}

func resourceListPersistentVolumeClaims(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().PersistentVolumeClaims(ns).List(opts)
}

func resourceListPersistentVolumes(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().PersistentVolumes().List(opts)
}

func resourceListPods(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Pods(ns).List(opts)
}

func resourceListPodSecurityPolicies(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.ExtensionsV1beta1().PodSecurityPolicies().List(opts)
}

func resourceListPodTemplates(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().PodTemplates(ns).List(opts)
}

func resourceListReplicaSets(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().ReplicaSets(ns).List(opts)
}

func resourceListReplicationControllers(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ReplicationControllers(ns).List(opts)
}

func resourceListResourceQuotas(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ResourceQuotas(ns).List(opts)
}

func resourceListSecrets(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Secrets(ns).List(opts)
}

func resourceListServiceAccounts(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ServiceAccounts(ns).List(opts)
}

func resourceListServices(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Services(ns).List(opts)
}

func resourceListStatefulSets(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().StatefulSets(ns).List(opts)
}

func resourceListStorageClasses(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.StorageV1().StorageClasses().List(opts)
}

func resourceListKinds() []string {
	var kinds []string
	for kind := range resourceListFns {
		kinds = append(kinds, kind)
	}
	return kinds
}
