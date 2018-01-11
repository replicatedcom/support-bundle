package types

import (
	dockertypes "github.com/docker/docker/api/types"
	"github.com/retracedhq/retraced-go"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Doc struct {
	Specs     []Spec
	Lifecycle []Task
}

type Task struct {
	Generate      *GenerateOptions      `json:"generate,omitempty"`
	Message       *MessageOptions       `json:"message,omitempty"`
	BooleanPrompt *BooleanPromptOptions `json:"message,omitempty"`
	UploadOptions *BooleanPromptOptions `json:"message,omitempty"`
}

type GenerateOptions struct {
}

type MessageOptions struct {
	Contents string `json:"contents,omitempty"`
}

type BooleanPromptOptions struct {
	Contents string `json:"contents,omitempty"`
	Default  rune   `json:"default,omitempty"`
}

type UploadOptions struct {
}

type Spec struct {
	Description    string `json:"description,omitempty"`
	Meta           Meta   `json:"meta,omitempty"`
	OutputDir      string `json:"output_dir"`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty"`
	Scrub          *Scrub `json:"scrub,omitempty"`

	// plan-specific config

	SupportBundleVersion *SupportBundleVersionOptions `json:"version,omitempty"`
	CustomerMeta         *CustomerMetaOptions         `json:"meta.customer,omitempty"`

	CoreHostname    *CoreHostnameOptions    `json:"os.hostname,omitempty"`
	CoreHTTPRequest *CoreHTTPRequestOptions `json:"os.http-request,omitempty"`
	CoreLoadavg     *CoreLoadavgOptions     `json:"os.loadavg,omitempty"`
	CoreReadFile    *CoreReadFileOptions    `json:"os.read-file,omitempty"`
	CoreRunCommand  *CoreRunCommandOptions  `json:"os.run-command,omitempty"`
	CoreUptime      *CoreUptimeOptions      `json:"os.uptime,omitempty"`

	JournaldLogs *JournaldLogsOptions `json:"journald.logs,omitempty"`

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
	KubernetesLogs         *KubernetesLogsOptions         `json:"kubernetes.logs,omitempty"`
	KubernetesResourceList *KubernetesResourceListOptions `json:"kubernetes.resource-list,omitempty"`
	KubernetesVersion      *KubernetesVersionOptions      `json:"kubernetes.version,omitempty"`

	RetracedEvents *RetracedEventsOptions `json:"retraced.events,omitempty"`
}

type Meta struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// plugin.supportbundle options

type SupportBundleVersionOptions struct {
}

// meta.customer

type CustomerMetaOptions struct {
	CustomerID     string `json:"customer_id,omitempty"`
	InstallationID string `json:"installation_id,omitempty"`
	CustomerName   string `json:"customer_name,omitempty"`
}

// plugin.core options

type CoreHostnameOptions struct {
}

type CoreHTTPRequestOptions struct {
	URL      string              `json:"url"`
	Method   string              `json:"method,omitempty"` // default "get"
	Body     string              `json:"body,omitempty"`
	Header   map[string][]string `json:"header,omitempty"`
	Insecure bool                `json:"insecure,omitempty"`
}

type CoreLoadavgOptions struct {
	Template string `json:"template"`
}

type CoreReadFileOptions struct {
	Filepath string `json:"filepath"`
}

type CoreRunCommandOptions struct {
	Name string   `json:"name"`
	Args []string `json:"args,omitempty"`
	Env  []string `json:"env,omitempty"`
	Dir  string   `json:"dir,omitempty"`
}

type CoreUptimeOptions struct {
	Template string `json:"template"`
}

// plugin.journald options

type JournaldLogsOptions struct {
	Unit    string `json:"unit,omitempty"`
	Since   string `json:"since,omitempty"`
	Reverse bool   `json:"reverse,omitempty"`
}

// plugin.docker options

type DockerContainerCpOptions struct {
	Container string `json:"container,omitempty"`
	SrcPath   string `json:"src_path,omitempty"`
}

type DockerContainerExecOptions struct {
	Container  string                 `json:"container"`
	ExecConfig dockertypes.ExecConfig `json:"exec_config"`
}

type DockerContainerInspectOptions struct {
	Container            string                `json:"container,omitempty"`
	ContainerListOptions *ContainerListOptions `json:"container_list_options,omitempty"`
}

type DockerContainerLogsOptions struct {
	Container            string                            `json:"container,omitempty"`
	ContainerLogsOptions *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	ContainerListOptions *ContainerListOptions             `json:"container_list_options,omitempty"`
}

type DockerContainerLsOptions struct {
	ContainerListOptions `json:",inline,omitempty"`
}

type DockerContainerRunOptions struct {
	ContainerCreateConfig dockertypes.ContainerCreateConfig `json:"container_create_config"`
	EnablePull            bool                              `json:"enable_pull,omitempty"`
}

type DockerImageLsOptions struct {
	ImageListOptions `json:",inline,omitempty"`
}

type DockerInfoOptions struct {
}

type DockerNodeLsOptions struct {
	NodeListOptions `json:",inline,omitempty"`
}

type DockerServiceLogsOptions struct {
	Service              string                            `json:"service,omitempty"`
	ContainerLogsOptions *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	ServiceListOptions   *ServiceListOptions               `json:"service_list_options,omitempty"`
}

type DockerServiceLsOptions struct {
	ServiceListOptions `json:",inline,omitempty"`
}

type DockerServicePsOptions struct {
	DockerTaskLsOptions `json:",inline,omitempty"`
}

type DockerStackServiceLogsOptions struct {
	Namespace              string                            `json:"namespace"`
	ContainerLogsOptions   *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	DockerServiceLsOptions *DockerServiceLsOptions           `json:"service_list_options,omitempty"`
}

type DockerStackServiceLsOptions struct {
	Namespace              string                  `json:"namespace"`
	DockerServiceLsOptions *DockerServiceLsOptions `json:"service_list_options,omitempty"`
}

type DockerStackServicePsOptions struct {
	Namespace              string                  `json:"namespace"`
	DockerServicePsOptions *DockerServicePsOptions `json:"task_list_options,omitempty"`
}

type DockerStackTaskLogsOptions struct {
	Namespace              string                            `json:"namespace"`
	ContainerLogsOptions   *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	DockerServicePsOptions *DockerServicePsOptions           `json:"task_list_options,omitempty"`
}

type DockerTaskLogsOptions struct {
	ID                   string                            `json:"id,omitempty"`
	ContainerLogsOptions *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	TaskListOptions      *TaskListOptions                  `json:"task_list_options,omitempty"`
}

type DockerTaskLsOptions struct {
	TaskListOptions `json:",inline,omitempty"`
}

type DockerVersionOptions struct {
}

// plugin.kubernetes options

type KubernetesAPIVersionsOptions struct {
}

type KubernetesClusterInfoOptions struct {
}

type KubernetesVersionOptions struct {
}

type KubernetesLogsOptions struct {
	Pod           string            `json:"pod"`
	Namespace     string            `json:"namespace,omitempty"`
	PodLogOptions *v1.PodLogOptions `json:"pod_log_options,omitempty"`
}

type KubernetesResourceListOptions struct {
	Kind        string              `json:"kind"`
	Namespace   string              `json:"namespace,omitempty"`
	ListOptions *metav1.ListOptions `json:"resource_list_options,omitempty"`
}

type RetracedEventsOptions struct {
	RetracedAPIClientOptions `json:",inline"`
	Mask                     *retraced.EventNodeMask   `json:"mask,omitempty"`
	Query                    *retraced.StructuredQuery `json:"query,omitempty"`
}

type RetracedAPIClientOptions struct {
	APIEndpoint string `json:"api_endpoint"`
	ProjectID   string `json:"project_id,omitempty"`
	APIToken    string `json:"api_token,omitempty"`
	Insecure    bool   `json:"insecure,omitempty"`
}

type Scrub struct {
	Regex   string `json:"regex"`
	Replace string `json:"replace"`
}
