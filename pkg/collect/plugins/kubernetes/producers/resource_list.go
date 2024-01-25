package producers

import (
	"context"
	"fmt"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type resourceListFnMapping struct {
	aliases []string
	fn      resourceListFn
	gvfns   map[string]resourceListFn
}

type resourceListFn func(client kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error)

var resourceListFnMappings = []resourceListFnMapping{
	{
		aliases: []string{"certificatesigningrequests", "csr"},
		fn:      resourceCertificateSigningRequestsCertificatesV1beta1,
		gvfns: map[string]resourceListFn{
			"certificates.k8s.io/v1beta1": resourceCertificateSigningRequestsCertificatesV1beta1,
		},
	},
	{
		aliases: []string{"clusterrolebindings"},
		fn:      resourceListClusterRoleBindingsRbacV1,
		gvfns: map[string]resourceListFn{
			"rbac.authorization.k8s.io/v1":       resourceListClusterRoleBindingsRbacV1,
			"rbac.authorization.k8s.io/v1beta1":  resourceListClusterRoleBindingsRbacV1beta1,
			"rbac.authorization.k8s.io/v1alpha1": resourceListClusterRoleBindingsRbacV1alpha1,
		},
	},
	{
		aliases: []string{"clusterroles"},
		fn:      resourceListClusterRolesRbacV1,
		gvfns: map[string]resourceListFn{
			"rbac.authorization.k8s.io/v1":       resourceListClusterRolesRbacV1,
			"rbac.authorization.k8s.io/v1beta1":  resourceListClusterRolesRbacV1beta1,
			"rbac.authorization.k8s.io/v1alpha1": resourceListClusterRolesRbacV1alpha1,
		},
	},
	{
		aliases: []string{"componentstatuses", "cs"},
		fn:      resourceListComponentStatusesCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListComponentStatusesCoreV1,
		},
	},
	{
		aliases: []string{"configmaps", "cm"},
		fn:      resourceListConfigMapsCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListConfigMapsCoreV1,
		},
	},
	{
		aliases: []string{"controllerrevisions"},
		fn:      resourceListControllerRevisionsAppsV1,
		gvfns: map[string]resourceListFn{
			"apps/v1":      resourceListControllerRevisionsAppsV1,
			"apps/v1beta2": resourceListControllerRevisionsAppsV1beta2,
			"apps/v1beta1": resourceListControllerRevisionsAppsV1beta1,
		},
	},
	{
		aliases: []string{"cronjobs", "cj"},
		fn:      resourceListCronJobsBatchV1beta1,
		gvfns: map[string]resourceListFn{
			"batch/v1beta1": resourceListCronJobsBatchV1beta1,
		},
	},
	{
		aliases: []string{"daemonsets", "ds"},
		fn:      resourceListDaemonSetsAppsV1,
		gvfns: map[string]resourceListFn{
			"apps/v1":      resourceListDaemonSetsAppsV1,
			"apps/v1beta2": resourceListDaemonSetsAppsV1beta2,
		},
	},
	{
		aliases: []string{"deployments", "deploy"},
		fn:      resourceListDeploymentsAppsV1,
		gvfns: map[string]resourceListFn{
			"apps/v1":            resourceListDeploymentsAppsV1,
			"apps/v1beta2":       resourceListDeploymentsAppsV1beta2,
			"apps/v1beta1":       resourceListDeploymentsAppsV1beta1,
			"extensions/v1beta1": resourceListDeploymentsExtensionsV1beta1,
		},
	},
	{
		aliases: []string{"endpoints", "ep"},
		fn:      resourceListEndpointsCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListEndpointsCoreV1,
		},
	},
	{
		aliases: []string{"events", "ev"},
		fn:      resourceListEventsCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1":               resourceListEventsCoreV1,
			"events.k8s.io/v1beta1": resourceListEventsEventsV1beta1,
		},
	},
	{
		aliases: []string{"horizontalpodautoscalers", "hpa"},
		fn:      resourceListHorizontalPodAutoscalersAutoscalingV1,
		gvfns: map[string]resourceListFn{
			"autoscaling/v1":      resourceListHorizontalPodAutoscalersAutoscalingV1,
			"autoscaling/v1beta2": resourceListHorizontalPodAutoscalersAutoscalingV2beta1,
		},
	},
	{
		aliases: []string{"ingresses", "ing"},
		fn:      resourceListIngressesExtensionsV1beta1,
		gvfns: map[string]resourceListFn{
			"extensions/v1beta1": resourceListIngressesExtensionsV1beta1,
		},
	},
	{
		aliases: []string{"jobs"},
		fn:      resourceListJobsBatchV1,
		gvfns: map[string]resourceListFn{
			"batch/v1": resourceListJobsBatchV1,
		},
	},
	{
		aliases: []string{"limitranges", "limits"},
		fn:      resourceListLimitRangesCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListLimitRangesCoreV1,
		},
	},
	{
		aliases: []string{"mutatingwebhookconfigurations"},
		fn:      resourceListMutatingWebhookConfigurationsAdmissionregistrationV1beta1,
		gvfns: map[string]resourceListFn{
			"admissionregistration.k8s.io/v1beta1": resourceListMutatingWebhookConfigurationsAdmissionregistrationV1beta1,
		},
	},
	{
		aliases: []string{"namespaces", "ns"},
		fn:      resourceListNamespacesCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListNamespacesCoreV1,
		},
	},
	{
		aliases: []string{"networkpolicies", "netpol"},
		fn:      resourceListNetworkPoliciesNetworkingV1,
		gvfns: map[string]resourceListFn{
			"networking.k8s.io/v1": resourceListNetworkPoliciesNetworkingV1,
		},
	},
	{
		aliases: []string{"nodes", "no"},
		fn:      resourceListNodesCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListNodesCoreV1,
		},
	},
	{
		aliases: []string{"persistentvolumeclaims", "pvc"},
		fn:      resourceListPersistentVolumeClaimsCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListPersistentVolumeClaimsCoreV1,
		},
	},
	{
		aliases: []string{"persistentvolumes", "pv"},
		fn:      resourceListPersistentVolumesCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListPersistentVolumesCoreV1,
		},
	},
	{
		aliases: []string{"poddisruptionbudgets", "pdb"},
		fn:      resourceListPodDisruptionBudgetsPolicyV1beta1,
		gvfns: map[string]resourceListFn{
			"policy/v1beta1": resourceListPodDisruptionBudgetsPolicyV1beta1,
		},
	},
	{
		aliases: []string{"pods", "po"},
		fn:      resourceListPodsCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListPodsCoreV1,
		},
	},
	{
		aliases: []string{"podtemplates"},
		fn:      resourceListPodTemplatesCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListPodTemplatesCoreV1,
		},
	},
	{
		aliases: []string{"replicasets", "rs"},
		fn:      resourceListReplicaSetsAppsV1,
		gvfns: map[string]resourceListFn{
			"apps/v1":            resourceListReplicaSetsAppsV1,
			"apps/v1beta2":       resourceListReplicaSetsAppsV1beta2,
			"extensions/v1beta1": resourceListReplicaSetsExtensionsV1beta1,
		},
	},
	{
		aliases: []string{"replicationcontrollers", "rc"},
		fn:      resourceListReplicationControllersCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListReplicationControllersCoreV1,
		},
	},
	{
		aliases: []string{"resourcequotas", "quota"},
		fn:      resourceListResourceQuotasCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListResourceQuotasCoreV1,
		},
	},
	{
		aliases: []string{"rolebindings"},
		fn:      resourceListRoleBindingsRbacV1,
		gvfns: map[string]resourceListFn{
			"rbac.authorization.k8s.io/v1":       resourceListRoleBindingsRbacV1,
			"rbac.authorization.k8s.io/v1beta1":  resourceListRoleBindingsRbacV1beta1,
			"rbac.authorization.k8s.io/v1alpha1": resourceListRoleBindingsRbacV1alpha1,
		},
	},
	{
		aliases: []string{"roles"},
		fn:      resourceListRolesRbacV1,
		gvfns: map[string]resourceListFn{
			"rbac.authorization.k8s.io/v1":       resourceListRolesRbacV1,
			"rbac.authorization.k8s.io/v1beta1":  resourceListRolesRbacV1beta1,
			"rbac.authorization.k8s.io/v1alpha1": resourceListRolesRbacV1alpha1,
		},
	},
	{
		aliases: []string{"secrets"},
		fn:      resourceListSecretsCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListSecretsCoreV1,
		},
	},
	{
		aliases: []string{"serviceaccounts", "sa"},
		fn:      resourceListServiceAccountsCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListServiceAccountsCoreV1,
		},
	},
	{
		aliases: []string{"services", "svc"},
		fn:      resourceListServicesCoreV1,
		gvfns: map[string]resourceListFn{
			"core/v1": resourceListServicesCoreV1,
		},
	},
	{
		aliases: []string{"statefulsets", "sts"},
		fn:      resourceListStatefulSetsAppsV1,
		gvfns: map[string]resourceListFn{
			"apps/v1":      resourceListStatefulSetsAppsV1,
			"apps/v1beta2": resourceListStatefulSetsAppsV1beta2,
			"apps/v1beta1": resourceListStatefulSetsAppsV1beta1,
		},
	},
	{
		aliases: []string{"storageclasses", "sc"},
		fn:      resourceListStorageClassesStorageV1,
		gvfns: map[string]resourceListFn{
			"storage.k8s.io/v1":      resourceListStorageClassesStorageV1,
			"storage.k8s.io/v1beta1": resourceListStorageClassesStorageV1beta1,
		},
	},
	{
		aliases: []string{"validatingwebhookconfigurations"},
		fn:      resourceListValidatingWebhookConfigurationsAdmissionregistrationV1beta1,
		gvfns: map[string]resourceListFn{
			"admissionregistration.k8s.io/v1beta1": resourceListValidatingWebhookConfigurationsAdmissionregistrationV1beta1,
		},
	},
	{
		aliases: []string{"volumeattachments"},
		fn:      resourceListVolumeAttachmentsStorageV1beta1,
		gvfns: map[string]resourceListFn{
			"storage.k8s.io/v1beta1":  resourceListVolumeAttachmentsStorageV1beta1,
			"storage.k8s.io/v1alpha1": resourceListVolumeAttachmentsStorageV1alpha1,
		},
	},
}

func (k *Kubernetes) ResourceList(opts types.KubernetesResourceListOptions) types.StructuredProducer {
	ns := opts.Namespace

	var listOpts metav1.ListOptions
	if opts.ListOptions != nil {
		listOpts = *opts.ListOptions
	}
	return func(ctx context.Context) (interface{}, error) {
		fn, err := findResourceListFnFromAliasAndGv(opts.Kind, opts.GroupVersion)
		if err != nil {
			return nil, err
		}
		return fn(k.client, ns, listOpts)
	}
}

func resourceCertificateSigningRequestsCertificatesV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CertificatesV1beta1().CertificateSigningRequests().List(context.TODO(), opts)
}

func resourceListClusterRoleBindingsRbacV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1().ClusterRoleBindings().List(context.TODO(), opts)
}

func resourceListClusterRoleBindingsRbacV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1beta1().ClusterRoleBindings().List(context.TODO(), opts)
}

func resourceListClusterRoleBindingsRbacV1alpha1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1alpha1().ClusterRoleBindings().List(context.TODO(), opts)
}

func resourceListClusterRolesRbacV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1().ClusterRoles().List(context.TODO(), opts)
}

func resourceListClusterRolesRbacV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1beta1().ClusterRoles().List(context.TODO(), opts)
}

func resourceListClusterRolesRbacV1alpha1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1alpha1().ClusterRoles().List(context.TODO(), opts)
}

func resourceListComponentStatusesCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ComponentStatuses().List(context.TODO(), opts)
}

func resourceListConfigMapsCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ConfigMaps(ns).List(context.TODO(), opts)
}

func resourceListControllerRevisionsAppsV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().ControllerRevisions(ns).List(context.TODO(), opts)
}

func resourceListControllerRevisionsAppsV1beta2(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().ControllerRevisions(ns).List(context.TODO(), opts)
}

func resourceListControllerRevisionsAppsV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta1().ControllerRevisions(ns).List(context.TODO(), opts)
}

func resourceListCronJobsBatchV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.BatchV1beta1().CronJobs(ns).List(context.TODO(), opts)
}

func resourceListDaemonSetsAppsV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().DaemonSets(ns).List(context.TODO(), opts)
}

func resourceListDaemonSetsAppsV1beta2(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().DaemonSets(ns).List(context.TODO(), opts)
}

func resourceListDaemonSetsExtensionsV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.ExtensionsV1beta1().DaemonSets(ns).List(context.TODO(), opts)
}

func resourceListDeploymentsAppsV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().Deployments(ns).List(context.TODO(), opts)
}

func resourceListDeploymentsAppsV1beta2(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().Deployments(ns).List(context.TODO(), opts)
}

func resourceListDeploymentsAppsV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta1().Deployments(ns).List(context.TODO(), opts)
}

func resourceListDeploymentsExtensionsV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.ExtensionsV1beta1().Deployments(ns).List(context.TODO(), opts)
}

func resourceListEndpointsCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Endpoints(ns).List(context.TODO(), opts)
}

func resourceListEventsCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Events(ns).List(context.TODO(), opts)
}

func resourceListEventsEventsV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.EventsV1beta1().Events(ns).List(context.TODO(), opts)
}

func resourceListHorizontalPodAutoscalersAutoscalingV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AutoscalingV1().HorizontalPodAutoscalers(ns).List(context.TODO(), opts)
}

func resourceListHorizontalPodAutoscalersAutoscalingV2beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AutoscalingV2beta1().HorizontalPodAutoscalers(ns).List(context.TODO(), opts)
}

func resourceListIngressesExtensionsV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.ExtensionsV1beta1().Ingresses(ns).List(context.TODO(), opts)
}

func resourceListJobsBatchV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.BatchV1().Jobs(ns).List(context.TODO(), opts)
}

func resourceListLimitRangesCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().LimitRanges(ns).List(context.TODO(), opts)
}

func resourceListMutatingWebhookConfigurationsAdmissionregistrationV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().List(context.TODO(), opts)
}

func resourceListNamespacesCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Namespaces().List(context.TODO(), opts)
}

func resourceListNetworkPoliciesNetworkingV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.NetworkingV1().NetworkPolicies(ns).List(context.TODO(), opts)
}

func resourceListNodesCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Nodes().List(context.TODO(), opts)
}

func resourceListPersistentVolumeClaimsCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().PersistentVolumeClaims(ns).List(context.TODO(), opts)
}

func resourceListPersistentVolumesCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().PersistentVolumes().List(context.TODO(), opts)
}

func resourceListPodDisruptionBudgetsPolicyV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.PolicyV1beta1().PodDisruptionBudgets(ns).List(context.TODO(), opts)
}

func resourceListPodsCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Pods(ns).List(context.TODO(), opts)
}

func resourceListPodSecurityPoliciesPolicyV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.PolicyV1beta1().PodSecurityPolicies().List(context.TODO(), opts)
}

func resourceListPodTemplatesCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().PodTemplates(ns).List(context.TODO(), opts)
}

func resourceListReplicaSetsAppsV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().ReplicaSets(ns).List(context.TODO(), opts)
}

func resourceListReplicaSetsAppsV1beta2(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().ReplicaSets(ns).List(context.TODO(), opts)
}

func resourceListReplicaSetsExtensionsV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.ExtensionsV1beta1().ReplicaSets(ns).List(context.TODO(), opts)
}

func resourceListReplicationControllersCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ReplicationControllers(ns).List(context.TODO(), opts)
}

func resourceListResourceQuotasCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ResourceQuotas(ns).List(context.TODO(), opts)
}

func resourceListRoleBindingsRbacV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1().RoleBindings(ns).List(context.TODO(), opts)
}

func resourceListRoleBindingsRbacV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1beta1().RoleBindings(ns).List(context.TODO(), opts)
}

func resourceListRoleBindingsRbacV1alpha1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1alpha1().RoleBindings(ns).List(context.TODO(), opts)
}

func resourceListRolesRbacV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1().Roles(ns).List(context.TODO(), opts)
}

func resourceListRolesRbacV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1beta1().Roles(ns).List(context.TODO(), opts)
}

func resourceListRolesRbacV1alpha1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.RbacV1alpha1().Roles(ns).List(context.TODO(), opts)
}

func resourceListSecretsCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Secrets(ns).List(context.TODO(), opts)
}

func resourceListServiceAccountsCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().ServiceAccounts(ns).List(context.TODO(), opts)
}

func resourceListServicesCoreV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.CoreV1().Services(ns).List(context.TODO(), opts)
}

func resourceListStatefulSetsAppsV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1().StatefulSets(ns).List(context.TODO(), opts)
}

func resourceListStatefulSetsAppsV1beta2(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta2().StatefulSets(ns).List(context.TODO(), opts)
}

func resourceListStatefulSetsAppsV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AppsV1beta1().StatefulSets(ns).List(context.TODO(), opts)
}

func resourceListStorageClassesStorageV1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.StorageV1().StorageClasses().List(context.TODO(), opts)
}

func resourceListStorageClassesStorageV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.StorageV1beta1().StorageClasses().List(context.TODO(), opts)
}

func resourceListValidatingWebhookConfigurationsAdmissionregistrationV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().List(context.TODO(), opts)
}

func resourceListVolumeAttachmentsStorageV1beta1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.StorageV1beta1().VolumeAttachments().List(context.TODO(), opts)
}

func resourceListVolumeAttachmentsStorageV1alpha1(c kubernetes.Interface, ns string, opts metav1.ListOptions) (interface{}, error) {
	return c.StorageV1alpha1().VolumeAttachments().List(context.TODO(), opts)
}

func resourceListKinds() []string {
	var kinds []string
	for _, mapping := range resourceListFnMappings {
		kinds = append(kinds, mapping.aliases...)
	}
	return kinds
}

func resourceListGroupVersions(mapping resourceListFnMapping) []string {
	var gvs []string
	for gv := range mapping.gvfns {
		gvs = append(gvs, gv)
	}
	return gvs
}

var resourceListFnMappingsMap = map[string]resourceListFnMapping{}

func init() {
	for _, mapping := range resourceListFnMappings {
		for _, alias := range mapping.aliases {
			resourceListFnMappingsMap[alias] = mapping
		}
	}
}

func findResourceListFnFromAliasAndGv(alias, gv string) (resourceListFn, error) {
	mapping, ok := resourceListFnMappingsMap[alias]
	if !ok {
		return nil, fmt.Errorf("unknown resource type, must be one of %q", resourceListKinds())
	}
	if gv == "" {
		return mapping.fn, nil
	}
	fn, ok := mapping.gvfns[gv]
	if !ok {
		return nil, fmt.Errorf("unknown group version for resource type %s, must be one of %q", alias, resourceListGroupVersions(mapping))
	}
	return fn, nil
}
