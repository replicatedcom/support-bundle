package bundle

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/core"
	dockerplugin "github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/journald"
	kubernetesplugin "github.com/replicatedcom/support-bundle/pkg/collect/plugins/kubernetes"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/retraced"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/supportbundle"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"

	docker "github.com/docker/docker/client"
	jww "github.com/spf13/jwalterweatherman"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Planner struct {
	Plugins []types.Plugin
}

type PlannerOptions struct {
	DockerClient           docker.CommonAPIClient
	KubernetesClient       kubernetes.Interface
	KubernetesClientConfig *restclient.Config

	InContainer bool

	EnableCore       bool
	EnableDocker     bool
	EnableJournald   bool
	EnableKubernetes bool
	EnableRetraced   bool
}

func NewPlanner(opts PlannerOptions) Planner {
	var p Planner

	p.AddPlugin(supportbundle.New())

	if opts.EnableCore {
		p.AddPlugin(core.New(opts.InContainer, opts.DockerClient))
	}

	if opts.EnableDocker {
		p.AddPlugin(dockerplugin.New(opts.DockerClient))
	}

	if opts.EnableJournald {
		p.AddPlugin(journald.New(opts.InContainer, opts.DockerClient))
	}

	if opts.EnableKubernetes {
		p.AddPlugin(kubernetesplugin.New(opts.KubernetesClient, opts.KubernetesClientConfig))
	}

	if opts.EnableRetraced {
		p.AddPlugin(retraced.New())
	}

	return p
}

func (p *Planner) AddPlugin(plugin types.Plugin) {
	p.Plugins = append(p.Plugins, plugin)
}

func (p *Planner) Plan(specs []types.Spec) []types.Task {
	var tasks []types.Task

Loop:
	for _, spec := range specs {
		for _, plugin := range p.Plugins {
			if planner := plugin.Plan(spec); planner != nil {
				tasks = append(tasks, planner(spec)...)
				continue Loop
			}
		}
		// don't show these to the customer, they're weird and scary
		jww.DEBUG.Printf("Producer not defined for spec %+v", spec)
	}

	return tasks
}
