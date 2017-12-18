package types

import (
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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

	CoreLoadavg     *CoreLoadavgOptions     `json:"os.loadavg,omitempty"`
	CoreHostname    *CoreHostnameOptions    `json:"os.hostname,omitempty"`
	CoreUptime      *CoreUptimeOptions      `json:"os.uptime,omitempty"`
	CoreReadFile    *CoreReadFileOptions    `json:"os.read-file,omitempty"`
	CoreRunCommand  *CoreRunCommandOptions  `json:"os.run-command,omitempty"`
	CoreHTTPRequest *CoreHTTPRequestOptions `json:"os.http-request,omitempty"`

	JournaldLogs *JournaldLogsOptions `json:"journald.logs,omitempty"`

	DockerInfo             *DockerInfoOptions             `json:"docker.info,omitempty"`
	DockerVersion          *DockerVersionOptions          `json:"docker.version,omitempty"`
	DockerPs               *DockerContainerLsOptions      `json:"docker.ps,omitempty"`     // canonical: docker.container-ls
	DockerImages           *DockerImageLsOptions          `json:"docker.images,omitempty"` // canonical: docker.image-ls
	DockerImageLs          *DockerImageLsOptions          `json:"docker.image-ls,omitempty"`
	DockerContainerLs      *DockerContainerLsOptions      `json:"docker.container-ls,omitempty"`
	DockerContainerLogs    *DockerContainerLogsOptions    `json:"docker.container-logs,omitempty"`
	DockerContainerInspect *DockerContainerInspectOptions `json:"docker.container-inspect,omitempty"`
	DockerContainerCp      *DockerContainerCpOptions      `json:"docker.container-cp,omitempty"`
	DockerContainerExec    *DockerContainerExecOptions    `json:"docker.container-exec,omitempty"`
	DockerContainerRun     *DockerContainerRunOptions     `json:"docker.container-run,omitempty"`
	DockerServiceLs        *DockerServiceLsOptions        `json:"docker.service-ls,omitempty"`
	DockerServiceLogs      *DockerServiceLogsOptions      `json:"docker.service-logs,omitempty"`
	DockerNodeLs           *DockerNodeLsOptions           `json:"docker.node-ls,omitempty"`
	DockerServicePs        *DockerServicePsOptions        `json:"docker.service-ps,omitempty"` // TODO: is there a more canonical way to get service tasks?
	// DockerStackServices *DockerStackServicesOptions `json:"docker.stack-services,omitempty"`
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

type Config struct {
	FilePath      string   `json:"file_path"`
	Args          []string `json:"args"`
	Image         string   `json:"image"`
	EnablePull    bool     `json:"enable_pull"`
	ContainerID   string   `json:"container_id"`
	ContainerName string   `json:"container_name"`
	Command       string   `json:"command"`
}

// plugin.core options

type CoreLoadavgOptions struct {
	Template string `json:"template"`
}

type CoreHostnameOptions struct {
}

type CoreUptimeOptions struct {
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

type CoreHTTPRequestOptions struct {
	URL      string              `json:"url"`
	Method   string              `json:"method,omitempty"` // default "get"
	Body     string              `json:"body,omitempty"`
	Header   map[string][]string `json:"header,omitempty"`
	Insecure bool                `json:"insecure,omitempty"`
}

// plugin.journald options

type JournaldLogsOptions struct {
	Unit  string
	Since string
}

// plugin.docker options

type DockerInfoOptions struct {
}

type DockerVersionOptions struct {
}

type DockerImageLsOptions struct {
}

type DockerContainerLsOptions struct {
}

type DockerContainerLogsOptions ContainerListOptions

type DockerContainerInspectOptions ContainerListOptions

type ContainerListOptions struct {
	dockertypes.ContainerListOptions
	Filters map[string][]string
}

func (opts ContainerListOptions) ToDockerContainerListOptions() dockertypes.ContainerListOptions {
	dockerOpts := opts.ContainerListOptions
	dockerOpts.Filters = FiltersToArgs(opts.Filters)
	return dockerOpts
}

func FiltersToArgs(f map[string][]string) filters.Args {
	args := filters.NewArgs()
	for k, vList := range f {
		for _, v := range vList {
			args.Add(k, v)
		}
	}
	return args
}

type DockerContainerCpOptions struct {
}

type DockerContainerExecOptions struct {
}

type DockerContainerRunOptions struct {
	ContainerCreateConfig dockertypes.ContainerCreateConfig
	EnablePull            bool
}

type DockerServiceLsOptions struct {
}

type DockerServiceLogsOptions struct {
}

type DockerNodeLsOptions struct {
}

type DockerServicePsOptions struct {
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
