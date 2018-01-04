package types

import (
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type ContainerListOptions struct {
	Quiet   bool
	Size    bool
	All     bool
	Latest  bool
	Since   string
	Before  string
	Limit   int
	Filters map[string][]string
}

func (opts *ContainerListOptions) ToDockerContainerListOptions() dockertypes.ContainerListOptions {
	return dockertypes.ContainerListOptions{
		Quiet:   opts.Quiet,
		Size:    opts.Size,
		All:     opts.All,
		Latest:  opts.Latest,
		Since:   opts.Since,
		Before:  opts.Before,
		Limit:   opts.Limit,
		Filters: FiltersToArgs(opts.Filters),
	}
}

type ImageListOptions struct {
	All     bool
	Filters map[string][]string
}

func (opts *ImageListOptions) ToDockerImageListOptions() dockertypes.ImageListOptions {
	return dockertypes.ImageListOptions{
		All:     opts.All,
		Filters: FiltersToArgs(opts.Filters),
	}
}

type NodeListOptions struct {
	Filters map[string][]string
}

func (opts *NodeListOptions) ToDockerNodeListOptions() dockertypes.NodeListOptions {
	return dockertypes.NodeListOptions{
		Filters: FiltersToArgs(opts.Filters),
	}
}

type ServiceListOptions struct {
	Filters map[string][]string
}

func (opts *ServiceListOptions) ToDockerServiceListOptions() dockertypes.ServiceListOptions {
	return dockertypes.ServiceListOptions{
		Filters: FiltersToArgs(opts.Filters),
	}
}

type TaskListOptions struct {
	Filters map[string][]string
}

func (opts *TaskListOptions) ToDockerTaskListOptions() dockertypes.TaskListOptions {
	return dockertypes.TaskListOptions{
		Filters: FiltersToArgs(opts.Filters),
	}
}

func FiltersToArgs(f map[string][]string) filters.Args {
	// Because filters.Args will not marshal
	args := filters.NewArgs()
	for k, vList := range f {
		for _, v := range vList {
			args.Add(k, v)
		}
	}
	return args
}
