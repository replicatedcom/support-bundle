package types

import (
	dockertypes "github.com/docker/docker/api/types"
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
	ContainerListOptions dockertypes.ContainerListOptions `json:"container_list_options"`
}

type DockerContainerInspectOptions struct {
	ContainerListOptions dockertypes.ContainerListOptions `json:"container_list_options"`
}

type Scrub struct {
	Regex   string `json:"regex"`
	Replace string `json:"replace"`
}
