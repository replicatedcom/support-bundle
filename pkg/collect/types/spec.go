package types

import (
	dockertypes "github.com/docker/docker/api/types"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/retracedhq/retraced-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Doc struct {
	Collect   Collect         `json:"collect,omitempty"`
	Lifecycle []LifecycleTask `json:"lifecycle,omitempty"`
}

// Collect is the top level collect object
type Collect struct {
	V1 []Spec `json:"v1,omitempty"`
}

type LifecycleTask struct {
	Generate *GenerateOptions `json:"generate,omitempty"`
	Message  *MessageOptions  `json:"message,omitempty"`
	Notes    *NotesOptions    `json:"notes,omitempty"`
	Upload   *UploadOptions   `json:"upload,omitempty"`
}

type GenerateOptions struct {
	UseDefaults bool `json:"use_defaults,omitempty"`
}

type MessageOptions struct {
	Contents string `json:"contents,omitempty"`
}

type PromptOptions struct {
	Message        string `json:"message,omitempty"`
	Default        bool   `json:"default,omitempty"`
	AcceptMessage  string `json:"accept,omitempty"`
	DeclineMessage string `json:"decline,omitempty"`
}

type UploadOptions struct {
	Prompt *PromptOptions `json:"prompt,omitempty"`
}

type NotesOptions struct {
	Prompt string `json:"prompt,omitempty"`
}

type SpecShared struct {
	// Description describes this spec field
	Description    string     `json:"description,omitempty"`
	OutputDir      string     `json:"output_dir,omitempty"`
	TimeoutSeconds int        `json:"timeout_seconds,omitempty"`
	Scrub          *Scrub     `json:"scrub,omitempty"`
	IncludeEmpty   bool       `json:"include_empty,omitempty"`
	Meta           *meta.Meta `json:"meta,omitempty"`
}

type Spec struct {
	SpecShared `json:",inline"`

	// plan-specific config

	SupportBundleVersion *SupportBundleVersionOptions `json:"version,omitempty"`
	CustomerMeta         *CustomerMetaOptions         `json:"meta.customer,omitempty"`

	// CoreHostname gets the hostname of the machine on which we're running
	CoreHostname    *CoreHostnameOptions    `json:"os.hostname,omitempty"`
	CoreHTTPRequest *CoreHTTPRequestOptions `json:"os.http-request,omitempty"`
	CoreLoadavg     *CoreLoadavgOptions     `json:"os.loadavg,omitempty"`
	CoreReadFile    *CoreReadFileOptions    `json:"os.read-file,omitempty"`
	CoreRunCommand  *CoreRunCommandOptions  `json:"os.run-command,omitempty"`
	CoreUptime      *CoreUptimeOptions      `json:"os.uptime,omitempty"`

	JournaldLogs *JournaldLogsOptions `json:"journald.logs,omitempty"`

	// DockerContainerCp copies a file from a running docker container
	DockerContainerCp      *DockerContainerCpOptions      `json:"docker.container-cp,omitempty"`
	DockerContainerExec    *DockerContainerExecOptions    `json:"docker.container-exec,omitempty"`
	DockerContainerInspect *DockerContainerInspectOptions `json:"docker.container-inspect,omitempty"`
	DockerContainerLogs    *DockerContainerLogsOptions    `json:"docker.container-logs,omitempty"`
	DockerContainerLs      *DockerContainerLsOptions      `json:"docker.container-ls,omitempty"`
	DockerContainerRun     *DockerContainerRunOptions     `json:"docker.container-run,omitempty"`
	DockerExec             *DockerContainerExecOptions    `json:"docker.exec,omitempty"` // canonical: docker.container-exec
	DockerImageLs          *DockerImageLsOptions          `json:"docker.image-ls,omitempty"`
	DockerImages           *DockerImageLsOptions          `json:"docker.images,omitempty"` // canonical: docker.image-ls
	DockerInfo             *DockerInfoOptions             `json:"docker.info,omitempty"`
	// TODO: docker.inspect
	DockerLogs             *DockerContainerLogsOptions    `json:"docker.logs,omitempty"` // canonical: docker.container-logs
	DockerNodeLs           *DockerNodeLsOptions           `json:"docker.node-ls,omitempty"`
	DockerPs               *DockerContainerLsOptions      `json:"docker.ps,omitempty"`  // canonical: docker.container-ls
	DockerRun              *DockerContainerRunOptions     `json:"docker.run,omitempty"` // canonical: docker.container-run
	DockerServiceLogs      *DockerServiceLogsOptions      `json:"docker.service-logs,omitempty"`
	DockerServiceLs        *DockerServiceLsOptions        `json:"docker.service-ls,omitempty"`
	DockerServicePs        *DockerServicePsOptions        `json:"docker.service-ps,omitempty"`
	DockerStackServiceLogs *DockerStackServiceLogsOptions `json:"docker.stack-service-logs,omitempty"`
	DockerStackServiceLs   *DockerStackServiceLsOptions   `json:"docker.stack-service-ls,omitempty"`
	DockerStackServicePs   *DockerStackServicePsOptions   `json:"docker.stack-service-ps,omitempty"`
	DockerStackTaskLogs    *DockerStackTaskLogsOptions    `json:"docker.stack-task-logs,omitempty"`
	DockerTaskLogs         *DockerTaskLogsOptions         `json:"docker.task-logs,omitempty"`
	DockerTaskLs           *DockerTaskLsOptions           `json:"docker.task-ls,omitempty"`
	DockerVersion          *DockerVersionOptions          `json:"docker.version,omitempty"`

	KubernetesAPIVersions  *KubernetesAPIVersionsOptions  `json:"kubernetes.api-versions,omitempty"`
	KubernetesClusterInfo  *KubernetesClusterInfoOptions  `json:"kubernetes.cluster-info,omitempty"`
	KubernetesContainerCp  *KubernetesContainerCpOptions  `json:"kubernetes.container-cp,omitempty"`
	KubernetesLogs         *KubernetesLogsOptions         `json:"kubernetes.logs,omitempty"`
	KubernetesResourceList *KubernetesResourceListOptions `json:"kubernetes.resource-list,omitempty"`
	KubernetesVersion      *KubernetesVersionOptions      `json:"kubernetes.version,omitempty"`

	RetracedEvents *RetracedEventsOptions `json:"retraced.events,omitempty"`
}

// plugin.supportbundle options

type SupportBundleVersionOptions struct {
	SpecShared `json:",inline,omitempty"`
}

// meta.customer

type CustomerMetaOptions struct {
	SpecShared     `json:",inline,omitempty"`
	CustomerID     string `json:"customer_id,omitempty"`
	InstallationID string `json:"installation_id,omitempty"`
	CustomerName   string `json:"customer_name,omitempty"`
}

// plugin.core options

type CoreHostnameOptions struct {
	SpecShared `json:",inline,omitempty"`
}

type CoreHTTPRequestOptions struct {
	SpecShared `json:",inline,omitempty"`
	URL        string              `json:"url"`
	Method     string              `json:"method,omitempty"` // default "get"
	Body       string              `json:"body,omitempty"`
	Header     map[string][]string `json:"header,omitempty"`
	Insecure   bool                `json:"insecure,omitempty"`
}

type CoreLoadavgOptions struct {
	SpecShared `json:",inline,omitempty"`
	Template   string `json:"template"`
}

type CoreReadFileOptions struct {
	SpecShared `json:",inline,omitempty"`
	Filepath   string `json:"filepath"`
}

type CoreRunCommandOptions struct {
	SpecShared `json:",inline,omitempty"`
	Name       string   `json:"name"`
	Args       []string `json:"args,omitempty"`
	Env        []string `json:"env,omitempty"`
	Dir        string   `json:"dir,omitempty"`
}

type CoreUptimeOptions struct {
	SpecShared `json:",inline,omitempty"`
	Template   string `json:"template"`
}

// plugin.journald options

type JournaldLogsOptions struct {
	SpecShared `json:",inline,omitempty"`
	Unit       string `json:"unit,omitempty"`
	Since      string `json:"since,omitempty"`
	Reverse    bool   `json:"reverse,omitempty"`
}

// plugin.docker options

type DockerContainerCpOptions struct {
	SpecShared `json:",inline,omitempty"`
	Container  string   `json:"container,omitempty"`
	SrcPath    string   `json:"src_path,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

type DockerContainerExecOptions struct {
	SpecShared `json:",inline,omitempty"`
	Container  string                 `json:"container"`
	ExecConfig dockertypes.ExecConfig `json:"exec_config"`
}

type DockerContainerInspectOptions struct {
	SpecShared           `json:",inline,omitempty"`
	Container            string                `json:"container,omitempty"`
	ContainerListOptions *ContainerListOptions `json:"container_list_options,omitempty"`
}

type DockerContainerLogsOptions struct {
	SpecShared           `json:",inline,omitempty"`
	Container            string                            `json:"container,omitempty"`
	ContainerLogsOptions *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	ContainerListOptions *ContainerListOptions             `json:"container_list_options,omitempty"`
}

type DockerContainerLsOptions struct {
	SpecShared           `json:",inline,omitempty"`
	ContainerListOptions `json:",inline,omitempty"`
}

type DockerContainerRunOptions struct {
	SpecShared            `json:",inline,omitempty"`
	ContainerCreateConfig dockertypes.ContainerCreateConfig `json:"container_create_config"`
	EnablePull            bool                              `json:"enable_pull,omitempty"`
}

type DockerImageLsOptions struct {
	SpecShared       `json:",inline,omitempty"`
	ImageListOptions `json:",inline,omitempty"`
}

type DockerInfoOptions struct {
	SpecShared `json:",inline,omitempty"`
}

type DockerNodeLsOptions struct {
	SpecShared      `json:",inline,omitempty"`
	NodeListOptions `json:",inline,omitempty"`
}

type DockerServiceLogsOptions struct {
	SpecShared           `json:",inline,omitempty"`
	Service              string                            `json:"service,omitempty"`
	ContainerLogsOptions *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	ServiceListOptions   *ServiceListOptions               `json:"service_list_options,omitempty"`
}

type DockerServiceLsOptions struct {
	SpecShared         `json:",inline,omitempty"`
	ServiceListOptions `json:",inline,omitempty"`
}

type DockerServicePsOptions struct {
	SpecShared          `json:",inline,omitempty"`
	DockerTaskLsOptions `json:",inline,omitempty"`
}

type DockerStackServiceLogsOptions struct {
	SpecShared             `json:",inline,omitempty"`
	LabelsFilter           `json:",inline,omitempty"`
	Namespace              string                            `json:"namespace"`
	ContainerLogsOptions   *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	DockerServiceLsOptions *DockerServiceLsOptions           `json:"service_list_options,omitempty"`
}

type DockerStackServiceLsOptions struct {
	SpecShared             `json:",inline,omitempty"`
	Namespace              string                  `json:"namespace"`
	DockerServiceLsOptions *DockerServiceLsOptions `json:"service_list_options,omitempty"`
}

type DockerStackServicePsOptions struct {
	SpecShared             `json:",inline,omitempty"`
	Namespace              string                  `json:"namespace"`
	DockerServicePsOptions *DockerServicePsOptions `json:"task_list_options,omitempty"`
}

type DockerStackTaskLogsOptions struct {
	SpecShared             `json:",inline,omitempty"`
	LabelsFilter           `json:",inline,omitempty"`
	Namespace              string                            `json:"namespace"`
	ContainerLogsOptions   *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	DockerServicePsOptions *DockerServicePsOptions           `json:"task_list_options,omitempty"`
}

type DockerTaskLogsOptions struct {
	SpecShared           `json:",inline,omitempty"`
	ID                   string                            `json:"id,omitempty"`
	ContainerLogsOptions *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	TaskListOptions      *TaskListOptions                  `json:"task_list_options,omitempty"`
}

type DockerTaskLsOptions struct {
	SpecShared      `json:",inline,omitempty"`
	TaskListOptions `json:",inline,omitempty"`
}

type DockerVersionOptions struct {
	SpecShared `json:",inline,omitempty"`
}

// plugin.kubernetes options

type KubernetesAPIVersionsOptions struct {
	SpecShared `json:",inline,omitempty"`
}

type KubernetesClusterInfoOptions struct {
	SpecShared `json:",inline,omitempty"`
}

type KubernetesVersionOptions struct {
	SpecShared `json:",inline,omitempty"`
}

type KubernetesLogsOptions struct {
	SpecShared    `json:",inline,omitempty"`
	Pod           string              `json:"pod,omitempty"`
	Namespace     string              `json:"namespace,omitempty"`
	PodLogOptions *PodLogOptions      `json:"pod_log_options,omitempty"`
	ListOptions   *metav1.ListOptions `json:"list_options,omitempty"`
}

type KubernetesContainerCpOptions struct {
	SpecShared     `json:",inline,omitempty"`
	Pod            string              `json:"pod,omitempty"`
	PodListOptions *metav1.ListOptions `json:"pod_list_options,omitempty"`
	Container      string              `json:"container,omitempty"`
	Namespace      string              `json:"namespace,omitempty"`
	SrcPath        string              `json:"src_path,omitempty"`
}

type KubernetesResourceListOptions struct {
	SpecShared   `json:",inline,omitempty"`
	Kind         string              `json:"kind"`
	GroupVersion string              `json:"group_version,omitempty"`
	Namespace    string              `json:"namespace,omitempty"`
	ListOptions  *metav1.ListOptions `json:"resource_list_options,omitempty"`
}

type RetracedEventsOptions struct {
	SpecShared               `json:",inline,omitempty"`
	RetracedAPIClientOptions `json:",inline"`
	Mask                     *retraced.EventNodeMask   `json:"mask,omitempty"`
	Query                    *retraced.StructuredQuery `json:"query,omitempty"`
}

type RetracedAPIClientOptions struct {
	SpecShared  `json:",inline,omitempty"`
	APIEndpoint string `json:"api_endpoint"`
	ProjectID   string `json:"project_id,omitempty"`
	APIToken    string `json:"api_token,omitempty"`
	Insecure    bool   `json:"insecure,omitempty"`
}

type LabelsFilter struct {
	Labels map[string]string `json:"labels,omitempty"`
}

type Scrub struct {
	Regex   string `json:"regex"`
	Replace string `json:"replace"`
}

// Todo maybe cache this or something
// this gets the shared stuff for a spec
//
//  ideally we want to move to
//
//  - docker.container-inspect:
//      container_name: ...
//      output_dir: ...
//      description: ...
//
//  but for now its more like
//
//  - docker.container-inspect:
//      container_name: ...
//    output_dir: ...
//    description: ...
//
//  This method allows you to use either or both,
//  giving priority to the first, nested version
//
//  i.e. if you have
//
//  - docker.container-inspect:
//      output_dir: /foo
//    output_dir: /bar
//
//  then /foo will be used
func (s *Spec) Shared() SpecShared {
	shared := func() SpecShared {
		switch {
		case s.SupportBundleVersion != nil:
			return s.SupportBundleVersion.SpecShared
		case s.CustomerMeta != nil:
			return s.CustomerMeta.SpecShared

			// CoreHostname gets the hostname of the machine on which we're running
		case s.CoreHostname != nil:
			return s.CoreHostname.SpecShared
		case s.CoreHTTPRequest != nil:
			return s.CoreHTTPRequest.SpecShared
		case s.CoreLoadavg != nil:
			return s.CoreLoadavg.SpecShared
		case s.CoreReadFile != nil:
			return s.CoreReadFile.SpecShared
		case s.CoreRunCommand != nil:
			return s.CoreRunCommand.SpecShared
		case s.CoreUptime != nil:
			return s.CoreUptime.SpecShared

		case s.JournaldLogs != nil:
			return s.JournaldLogs.SpecShared

			// DockerContainerCp copies a file from a running docker container
		case s.DockerContainerCp != nil:
			return s.DockerContainerCp.SpecShared
		case s.DockerContainerExec != nil:
			return s.DockerContainerExec.SpecShared
		case s.DockerContainerInspect != nil:
			return s.DockerContainerInspect.SpecShared
		case s.DockerContainerLogs != nil:
			return s.DockerContainerLogs.SpecShared
		case s.DockerContainerLs != nil:
			return s.DockerContainerLs.SpecShared
		case s.DockerContainerRun != nil:
			return s.DockerContainerRun.SpecShared
		case s.DockerExec != nil:
			return s.DockerExec.SpecShared
		case s.DockerImageLs != nil:
			return s.DockerImageLs.SpecShared
		case s.DockerImages != nil:
			return s.DockerImages.SpecShared
		case s.DockerInfo != nil:
			return s.DockerInfo.SpecShared
			// TODO: docker.inspect
		case s.DockerLogs != nil:
			return s.DockerLogs.SpecShared
		case s.DockerNodeLs != nil:
			return s.DockerNodeLs.SpecShared
		case s.DockerPs != nil:
			return s.DockerPs.SpecShared
		case s.DockerRun != nil:
			return s.DockerRun.SpecShared
		case s.DockerServiceLogs != nil:
			return s.DockerServiceLogs.SpecShared
		case s.DockerServiceLs != nil:
			return s.DockerServiceLs.SpecShared
		case s.DockerServicePs != nil:
			return s.DockerServicePs.SpecShared
		case s.DockerStackServiceLogs != nil:
			return s.DockerStackServiceLogs.SpecShared
		case s.DockerStackServiceLs != nil:
			return s.DockerStackServiceLs.SpecShared
		case s.DockerStackServicePs != nil:
			return s.DockerStackServicePs.SpecShared
		case s.DockerStackTaskLogs != nil:
			return s.DockerStackTaskLogs.SpecShared
		case s.DockerTaskLogs != nil:
			return s.DockerTaskLogs.SpecShared
		case s.DockerTaskLs != nil:
			return s.DockerTaskLs.SpecShared
		case s.DockerVersion != nil:
			return s.DockerVersion.SpecShared

		case s.KubernetesAPIVersions != nil:
			return s.KubernetesAPIVersions.SpecShared
		case s.KubernetesClusterInfo != nil:
			return s.KubernetesClusterInfo.SpecShared
		case s.KubernetesContainerCp != nil:
			return s.KubernetesContainerCp.SpecShared
		case s.KubernetesLogs != nil:
			return s.KubernetesLogs.SpecShared
		case s.KubernetesResourceList != nil:
			return s.KubernetesResourceList.SpecShared
		case s.KubernetesVersion != nil:
			return s.KubernetesVersion.SpecShared

		case s.RetracedEvents != nil:
			return s.RetracedEvents.SpecShared
		default:
			return s.SpecShared
		}
	}()

	if shared.OutputDir == "" {
		shared.OutputDir = s.OutputDir
	}

	if shared.Meta == nil {
		shared.Meta = s.Meta
	}

	if shared.Description == "" {
		shared.Description = s.Description
	}

	if shared.Scrub == nil {
		shared.Scrub = s.Scrub
	}

	if shared.TimeoutSeconds == 0 {
		shared.TimeoutSeconds = s.TimeoutSeconds
	}

	return shared
}
