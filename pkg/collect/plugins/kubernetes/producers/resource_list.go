package producers

import (
	"context"
	"fmt"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type resourceListFn func(client kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error)

var resourceListFns = map[string]resourceListFn{
	"certificatesigningrequests": resourceCertificateSigningRequests, "csr": resourceCertificateSigningRequests,
	"clusterrolebindings": resourceListClusterRoleBindings,
	"clusterroles":        resourceListClusterRoles,
	"componentstatuses":   resourceListComponentStatuses, "cs": resourceListComponentStatuses,
	"configmaps": resourceListConfigMaps, "cm": resourceListConfigMaps,
	"controllerrevisions": resourceListControllerRevisions,
	"cronjobs":            resourceListCronJobs, "cj": resourceListCronJobs,
	"daemonsets": resourceListDaemonSets, "ds": resourceListDaemonSets,
	"deployments": resourceListDeployments, "deploy": resourceListDeployments,
	"endpoints": resourceListEndpoints, "ep": resourceListEndpoints,
	"events": resourceListEvents, "ev": resourceListEvents,
	"horizontalpodautoscalers": resourceListHorizontalPodAutoscalers, "hpa": resourceListHorizontalPodAutoscalers,
	"ingresses": resourceListIngresses, "ing": resourceListIngresses,
	"jobs":        resourceListJobs,
	"limitranges": resourceListLimitRanges, "limits": resourceListLimitRanges,
	"mutatingwebhookconfigurations": resourceListMutatingWebhookConfigurations,
	"namespaces":                    resourceListNamespaces, "ns": resourceListNamespaces,
	"networkpolicies": resourceListNetworkPolicies, "netpol": resourceListNetworkPolicies,
	"nodes": resourceListNodes, "no": resourceListNodes,
	"persistentvolumeclaims": resourceListPersistentVolumeClaims, "pvc": resourceListPersistentVolumeClaims,
	"persistentvolumes": resourceListPersistentVolumes, "pv": resourceListPersistentVolumes,
	"poddisruptionbudgets": resourceListPodDisruptionBudgets, "pdb": resourceListPodDisruptionBudgets,
	"pods": resourceListPods, "po": resourceListPods,
	"podsecuritypolicies": resourceListPodSecurityPolicies, "psp": resourceListPodSecurityPolicies,
	"podtemplates": resourceListPodTemplates,
	"replicasets":  resourceListReplicaSets, "rs": resourceListReplicaSets,
	"replicationcontrollers": resourceListReplicationControllers, "rc": resourceListReplicationControllers,
	"resourcequotas": resourceListResourceQuotas, "quota": resourceListResourceQuotas,
	"rolebindings":    resourceListRoleBindings,
	"roles":           resourceListRoles,
	"secrets":         resourceListSecrets,
	"serviceaccounts": resourceListServiceAccounts, "sa": resourceListServiceAccounts,
	"services": resourceListServices, "svc": resourceListServices,
	"statefulsets": resourceListStatefulSets, "sts": resourceListStatefulSets,
	"storageclasses":                  resourceListStorageClasses,
	"validatingwebhookconfigurations": resourceListValidatingWebhookConfigurations,
	"volumeattachments":               resourceListVolumeAttachments,
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

func resourceCertificateSigningRequests(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CertificatesV1beta1().CertificateSigningRequests().List(opts)
}

func resourceListClusterRoleBindings(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1().ClusterRoleBindings().List(opts)
	// return c.RbacV1beta1().ClusterRoleBindings().List(opts)
	// return c.RbacV1alpha1().ClusterRoleBindings().List(opts)
}

func resourceListClusterRoles(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1().ClusterRoles().List(opts)
	// return c.RbacV1beta1().ClusterRoles().List(opts)
	// return c.RbacV1alpha1().ClusterRoles().List(opts)
}

func resourceListComponentStatuses(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ComponentStatuses().List(opts)
}

func resourceListConfigMaps(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ConfigMaps(ns).List(opts)
}

func resourceListControllerRevisions(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().ControllerRevisions(ns).List(opts)
	// return c.AppsV1beta2().ControllerRevisions(ns).List(opts)
	// return c.AppsV1beta1().ControllerRevisions(ns).List(opts)
}

func resourceListCronJobs(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.BatchV1beta1().CronJobs(ns).List(opts)
	// return c.BatchV2alpha1().CronJobs(ns).List(opts)
}

func resourceListDaemonSets(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().DaemonSets(ns).List(opts)
	// return c.AppsV1beta2().DaemonSets(ns).List(opts)
	// return c.ExtensionsV1beta1().DaemonSets().List(opts)
}

func resourceListDeployments(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().Deployments(ns).List(opts)
	// return c.AppsV1beta2().Deployments(ns).List(opts)
	// return c.AppsV1beta1().Deployments(ns).List(opts)
}

func resourceListEndpoints(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Endpoints(ns).List(opts)
}

func resourceListEvents(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Events(ns).List(opts)
	// return c.EventsV1beta1().Events(ns).List(opts)
}

func resourceListHorizontalPodAutoscalers(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AutoscalingV1().HorizontalPodAutoscalers(ns).List(opts)
	// return c.AutoscalingV2beta1().HorizontalPodAutoscalers(ns).List(opts)
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

func resourceListMutatingWebhookConfigurations(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().List(opts)
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

func resourceListPodDisruptionBudgets(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.PolicyV1beta1().PodDisruptionBudgets(ns).List(opts)
}

func resourceListPods(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Pods(ns).List(opts)
}

func resourceListPodSecurityPolicies(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.PolicyV1beta1().PodSecurityPolicies().List(opts)
	// return c.ExtensionsV1beta1().PodSecurityPolicies().List(opts)
}

func resourceListPodTemplates(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().PodTemplates(ns).List(opts)
}

func resourceListReplicaSets(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().ReplicaSets(ns).List(opts)
	// return c.AppsV1beta2().ReplicaSets(ns).List(opts)
	// return c.ExtensionsV1beta1().ReplicaSets(ns).List(opts)
}

func resourceListReplicationControllers(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ReplicationControllers(ns).List(opts)
}

func resourceListResourceQuotas(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ResourceQuotas(ns).List(opts)
}

func resourceListRoleBindings(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1().RoleBindings(ns).List(opts)
	// return c.RbacV1beta1().RoleBindings().List(opts)
	// return c.RbacV1alpha1().RoleBindings().List(opts)
}

func resourceListRoles(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1().Roles(ns).List(opts)
	// return c.RbacV1beta1().Roles().List(opts)
	// return c.RbacV1alpha1().Roles().List(opts)
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
	return c.AppsV1().StatefulSets(ns).List(opts)
	// return c.AppsV1beta2().StatefulSets(ns).List(opts)
	// return c.AppsV1beta1().StatefulSets(ns).List(opts)
}

func resourceListStorageClasses(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.StorageV1().StorageClasses().List(opts)
	// return c.StorageV1beta1().StorageClasses().List(opts)
}

func resourceListValidatingWebhookConfigurations(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().List(opts)
}

func resourceListVolumeAttachments(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.StorageV1beta1().VolumeAttachments().List(opts)
	// return c.StorageV1alpha1().VolumeAttachments().List(opts)
}

func resourceListKinds() []string {
	var kinds []string
	for kind := range resourceListFns {
		kinds = append(kinds, kind)
	}
	return kinds
}
