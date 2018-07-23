package bundle

import (
	"context"

	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/core"
	dockerplugin "github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/journald"
	kubernetesplugin "github.com/replicatedcom/support-bundle/pkg/collect/plugins/kubernetes"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/retraced"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/supportbundle"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	dockerclient "github.com/replicatedcom/support-bundle/pkg/docker"
	kubernetesclient "github.com/replicatedcom/support-bundle/pkg/kubernetes"
	"github.com/replicatedcom/support-bundle/pkg/logger"
	jww "github.com/spf13/jwalterweatherman"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Planner struct {
	Plugins []types.Plugin
}

type PlannerOptions struct {
	EnableCore       bool
	EnableDocker     bool
	EnableJournald   bool
	EnableKubernetes bool
	EnableRetraced   bool

	RequireCore       bool
	RequireDocker     bool
	RequireJournald   bool
	RequireKubernetes bool
	RequireRetraced   bool
}

func NewPlanner(opts PlannerOptions, inContainer bool) (*Planner, error) {
	dockerClient, dockerErr := newDockerClient(opts, inContainer)
	kubernetesClient, kubernetesClientConfig, kubernetesErr := newKubernetesClientAndConfig(opts)

	var p Planner

	p.AddPlugin(supportbundle.New())

	if opts.EnableCore {
		if !inContainer || dockerErr == nil {
			p.AddPlugin(core.New(inContainer, dockerClient))
		} else if opts.RequireCore {
			return nil, errors.Wrap(dockerErr, "require core")
		}
	}

	if opts.EnableDocker {
		if dockerErr == nil {
			p.AddPlugin(dockerplugin.New(dockerClient))
		} else if opts.RequireDocker {
			return nil, errors.Wrap(dockerErr, "require docker")
		}
	}

	if opts.EnableJournald {
		if !inContainer || dockerErr == nil {
			p.AddPlugin(journald.New(inContainer, dockerClient))
		} else if opts.RequireJournald {
			return nil, errors.Wrap(dockerErr, "require journald")
		}
	}

	if opts.EnableKubernetes {
		if kubernetesErr == nil {
			p.AddPlugin(kubernetesplugin.New(kubernetesClient, kubernetesClientConfig))
		} else if opts.RequireKubernetes {
			return nil, errors.Wrap(kubernetesErr, "require kubernetes")
		}
	}

	if opts.EnableRetraced {
		p.AddPlugin(retraced.New())
	}

	return &p, nil
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

func newDockerClient(opts PlannerOptions, inContainer bool) (docker.CommonAPIClient, error) {
	if !opts.EnableDocker &&
		!(inContainer && (opts.EnableCore || opts.EnableJournald)) {
		return nil, errors.New("docker plugin not enabled")
	}
	kitLog := logger.New(
		logger.LevelFromJWWThreshold(jww.GetLogThreshold()),
	)
	client, err := dockerclient.NewEnvClient(context.Background(), kitLog)
	return client, errors.Wrap(err, "get docker client from environment")
}

func newKubernetesClientAndConfig(opts PlannerOptions) (kubernetes.Interface, *restclient.Config, error) {
	if !opts.EnableKubernetes {
		return nil, nil, errors.New("kubernetes plugin not enabled")
	}
	config, err := kubernetesclient.ClientConfig()
	if err != nil {
		return nil, config, errors.Wrap(err, "get kubernetes client config")
	}
	client, err := kubernetesclient.NewClient(config)
	return client, config, errors.Wrap(err, "get kubernetes client")
}
