package collector

import (
	"context"
	"os"
	"time"

	docker "github.com/docker/docker/client"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/collect/bundle"
	"github.com/replicatedcom/support-bundle/pkg/spew"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Options struct {
	EnableCore       bool
	EnableDocker     bool
	EnableJournald   bool
	EnableKubernetes bool
	EnableRetraced   bool
}

type Collector struct {
	Logger log.Logger

	DockerClient           docker.CommonAPIClient
	KubernetesClient       kubernetes.Interface
	KubernetesClientConfig *restclient.Config
}

func New(
	logger log.Logger,
	dockerClient docker.CommonAPIClient,
	kubernetesClient kubernetes.Interface,
	kubernetesClientConfig *restclient.Config,
) Interface {
	return &Collector{
		Logger: logger,

		DockerClient:           dockerClient,
		KubernetesClient:       kubernetesClient,
		KubernetesClientConfig: kubernetesClientConfig,
	}
}

func (c *Collector) CollectBundle(
	ctx context.Context,
	collectSpec api.Collect,
	timeout time.Duration,
	dest string,
	opts Options,
) (os.FileInfo, error) {

	debug := level.Debug(log.With(c.Logger, "method", "Collector.Collect"))

	planner := bundle.NewPlanner(bundle.PlannerOptions{
		DockerClient:           c.DockerClient,
		KubernetesClient:       c.KubernetesClient,
		KubernetesClientConfig: c.KubernetesClientConfig,

		InContainer: os.Getenv("IN_CONTAINER") != "",

		EnableCore:       opts.EnableCore,
		EnableDocker:     opts.EnableDocker,
		EnableJournald:   opts.EnableJournald,
		EnableKubernetes: opts.EnableKubernetes,
		EnableRetraced:   opts.EnableRetraced,
	})

	tasks := planner.Plan(collectSpec.V1)
	debug.Log(
		"phase", "planner.plan",
		"spec", spew.Sdump(collectSpec))

	fi, _, err := bundle.Generate(tasks, timeout, dest)
	debug.Log(
		"phase", "bundle.generate",
		"timeout", timeout,
		"dest", dest,
		"error", err)

	return fi, errors.Wrap(err, "generate")
}
