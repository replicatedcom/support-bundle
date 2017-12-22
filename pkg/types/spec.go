package types

import (
	dockertypes "github.com/docker/docker/api/types"
	"github.com/retracedhq/retraced-go"
)

type Doc struct {
	Specs []Spec
}

type Spec struct {
	Description    string `json:"description,omitempty"`
	Meta           Meta   `json:"meta,omitempty"`
	OutputDir      string `json:"output_dir"`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty"`
	Scrub          *Scrub `json:"scrub,omitempty"`

	// plan-specific config

	CoreHostname    *CoreHostnameOptions    `json:"os.hostname,omitempty"`
	CoreHTTPRequest *CoreHTTPRequestOptions `json:"os.http-request,omitempty"`
	CoreLoadavg     *CoreLoadavgOptions     `json:"os.loadavg,omitempty"`
	CoreReadFile    *CoreReadFileOptions    `json:"os.read-file,omitempty"`
	CoreRunCommand  *CoreRunCommandOptions  `json:"os.run-command,omitempty"`
	CoreUptime      *CoreUptimeOptions      `json:"os.uptime,omitempty"`

	JournaldLogs *JournaldLogsOptions `json:"journald.logs,omitempty"`

	DockerContainerInspect *DockerContainerInspectOptions `json:"docker.container-inspect,omitempty"`
	DockerContainerLogs    *DockerContainerLogsOptions    `json:"docker.container-logs,omitempty"`
	DockerContainerLs      *DockerContainerLsOptions      `json:"docker.container-ls,omitempty"`
	DockerImageLs          *DockerImageLsOptions          `json:"docker.image-ls,omitempty"`
	DockerImages           *DockerImageLsOptions          `json:"docker.images,omitempty"` // canonical: docker.image-ls
	DockerInfo             *DockerInfoOptions             `json:"docker.info,omitempty"`
	DockerNodeLs           *DockerNodeLsOptions           `json:"docker.node-ls,omitempty"`
	DockerPs               *DockerContainerLsOptions      `json:"docker.ps,omitempty"` // canonical: docker.container-ls
	DockerServiceLs        *DockerServiceLsOptions        `json:"docker.service-ls,omitempty"`
	DockerServicePs        *DockerServicePsOptions        `json:"docker.service-ps,omitempty"` // TODO: is there a more canonical way to get service tasks?
	DockerStackServiceLs   *DockerStackServiceLsOptions   `json:"docker.stack-service-ls,omitempty"`
	DockerStackServicePs   *DockerStackServicePsOptions   `json:"docker.stack-service-ps,omitempty"`
	DockerTaskLs           *DockerTaskLsOptions           `json:"docker.task-ls,omitempty"`
	DockerVersion          *DockerVersionOptions          `json:"docker.version,omitempty"`

	DockerContainerCp   *DockerContainerCpOptions   `json:"docker.container-cp,omitempty"`
	DockerContainerExec *DockerContainerExecOptions `json:"docker.container-exec,omitempty"`
	DockerContainerRun  *DockerContainerRunOptions  `json:"docker.container-run,omitempty"`
	DockerServiceLogs   *DockerServiceLogsOptions   `json:"docker.service-logs,omitempty"`
	// DockerStackServiceLogs *DockerStackServiceLogsOptions `json:"docker.stack-service-logs,omitempty"`
	// DockerStackTasks *DockerStackTasksOptions `json:"docker.stack-tasks,omitempty"`
	// DockerStackTaskLogs *DockerStackTaskLogsOptions `json:"docker.stack-task-logs,omitempty"`

	KubernetesLogs     *KubernetesLogsOptions     `json:"kubernetes.logs,omitempty"`
	KubernetesResource *KubernetesResourceOptions `json:"kubernetes.resource,omitempty"`

	RetracedEvents *RetracedEventsOptions `json:"retraced.events,omitempty"`
}

type Meta struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
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
	Unit  string
	Since string
}

// plugin.docker options

type DockerContainerInspectOptions struct {
	ID                   string                `json:"id,omitempty"`
	Name                 string                `json:"name,omitempty"`
	ContainerListOptions *ContainerListOptions `json:"container_list_options,omitempty"`
}

type DockerContainerLogsOptions struct {
	ID                   string                            `json:"id,omitempty"`
	Name                 string                            `json:"name,omitempty"`
	ContainerLogsOptions *dockertypes.ContainerLogsOptions `json:"container_logs_options,omitempty"`
	ContainerListOptions *ContainerListOptions             `json:"container_list_options,omitempty"`
}

type DockerContainerLsOptions struct {
	ContainerListOptions `json:",inline,omitempty"`
}

type DockerImageLsOptions struct {
	ImageListOptions `json:",inline,omitempty"`
}

type DockerInfoOptions struct {
}

type DockerNodeLsOptions struct {
	NodeListOptions `json:",inline,omitempty"`
}

type DockerServiceLsOptions struct {
	ServiceListOptions `json:",inline,omitempty"`
}

type DockerServicePsOptions struct {
	DockerTaskLsOptions `json:",inline,omitempty"`
}

type DockerStackServiceLsOptions struct {
	Namespace string `json:"namespace"`
	// TODO: filters.Args will panic
	*DockerServiceLsOptions `json:"service_list_options,omitempty"`
}

type DockerStackServicePsOptions struct {
	Namespace string `json:"namespace"`
	// TODO: filters.Args will panic
	*DockerServicePsOptions `json:"task_list_options,omitempty"`
}

type DockerTaskLsOptions struct {
	TaskListOptions `json:",inline,omitempty"`
}

type DockerVersionOptions struct {
}

type DockerContainerCpOptions struct {
}

type DockerContainerExecOptions struct {
}

type DockerContainerRunOptions struct {
	ContainerCreateConfig dockertypes.ContainerCreateConfig
	EnablePull            bool
}

type DockerServiceLogsOptions struct {
}

type DockerStackServicesOptions struct {
}

type DockerStackServiceLogsOptions struct {
}

type DockerStackTasksOptions struct {
}

type DockerStackTaskLogsOptions struct {
}

// plugin.kubernetes options

type KubernetesLogsOptions struct {
	PodName       string `json:"pod_name"`
	ContainerName string `json:"container_name,omitempty"`
	Namespace     string `json:"namespace,omitempty"`
	Scrub         Scrub  `json:"scrub"`
}

type KubernetesResourceOptions struct {
	Type      string `json:"type"`
	Namespace string `json:"namespace,omitempty"`
}

type RetracedEventsOptions struct {
	APIEndpoint string                    `json:"api_endpoint"`
	ProjectID   string                    `json:"project_id,omitempty"`
	APIToken    string                    `json:"api_token,omitempty"`
	Insecure    bool                      `json:"insecure,omitempty"`
	Mask        *retraced.EventNodeMask   `json:"mask,omitempty"`
	Query       *retraced.StructuredQuery `json:"query,omitempty"`
}

type Scrub struct {
	Regex   string `json:"regex"`
	Replace string `json:"replace"`
}
