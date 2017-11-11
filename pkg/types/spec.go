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
	ID             string
	Builtin        string
	TimeoutSeconds int
	// paths
	Raw   string
	JSON  string
	Human string

	// Plan-specific config
	Config Config

	// New plan-specific config
	// TODO: new format for all tasks
	DockerRunCommand       *DockerRunCommandOptions       `json:"docker.run-command,omitempty"`
	DockerContainerLogs    *DockerContainerLogsOptions    `json:"docker.container-logs,omitempty"`
	DockerContainerInspect *DockerContainerInspectOptions `json:"docker.container-inspect,omitempty"`
	HTTPRequestCommand     *HTTPRequestCommandOptions     `json:"core.http-request,omitempty"`
	RetracedEventsCommand  *RetracedEventsOptions         `json:"retraced.events,omitempty"`
}

type Config struct {
	FilePath      string   `json:"file_path"`
	Args          []string `json:"args"`
	Image         string   `json:"image"`
	EnablePull    bool     `json:"enable_pull"`
	ContainerID   string   `json:"container_id"`
	ContainerName string   `json:"container_name"`
	Command       string   `json:"command"`
	Scrub         Scrub    `json:"scrub"` // TODO: should we pull scrub up one level into Spec?
}

type DockerRunCommandOptions struct {
	ContainerCreateConfig dockertypes.ContainerCreateConfig
	EnablePull            bool
}

type DockerContainerLogsOptions struct {
	ContainerListOptions ContainerListOptions `json:"container_list_options"`
}

type DockerContainerInspectOptions struct {
	ContainerListOptions ContainerListOptions `json:"container_list_options"`
}

type HTTPRequestCommandOptions struct {
	URL      string              `json:"url"`
	Method   string              `json:"method,omitempty"` // default "get"
	Body     string              `json:"body,omitempty"`
	Header   map[string][]string `json:"header,omitempty"`
	Insecure bool                `json:"insecure,omitempty"`
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
