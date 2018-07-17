package collector

import (
	"context"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/collect/bundle"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/core"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/journald"
	pluginskubernetes "github.com/replicatedcom/support-bundle/pkg/collect/plugins/kubernetes"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/retraced"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/supportbundle"
	"github.com/replicatedcom/support-bundle/pkg/spew"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Collector struct {
	Logger log.Logger

	KubernetesClient       kubernetes.Interface
	KubernetesClientConfig *restclient.Config
}

func New(
	logger log.Logger,
	kubernetesClient kubernetes.Interface,
	kubernetesClientConfig *restclient.Config,
) Interface {
	return &Collector{
		Logger: logger,

		KubernetesClient:       kubernetesClient,
		KubernetesClientConfig: kubernetesClientConfig,
	}
}

func (c *Collector) CollectBundle(ctx context.Context, collect api.Collect, timeout time.Duration, pathname string) (os.FileInfo, error) {
	debug := level.Debug(log.With(c.Logger, "method", "Collector.Collect"))

	planner, err := c.newPlanner()
	debug.Log(
		"phase", "planner.init",
		"error", err)
	if err != nil {
		return nil, errors.Wrap(err, "new planner")
	}

	tasks := planner.Plan(collect.V1)
	debug.Log(
		"phase", "planner.plan",
		"spec", spew.Sdump(collect))

	fi, _, err := bundle.Generate(tasks, timeout, pathname)
	debug.Log(
		"phase", "bundle.generate",
		"timeout", timeout,
		"pathname", pathname,
		"error", err)

	return fi, errors.Wrap(err, "generate")
}

// TODO: dependency injection
func (c *Collector) newPlanner() (*bundle.Planner, error) {
	planner := &bundle.Planner{}

	pluginSupportBundle, err := supportbundle.New()
	if err != nil {
		return planner, errors.Wrap(err, "initialize supportbundle plugin")
	}
	planner.AddPlugin(pluginSupportBundle)

	pluginCore, err := core.New()
	if err != nil {
		return planner, errors.Wrap(err, "initialize core plugin")
	}
	planner.AddPlugin(pluginCore)

	pluginDocker, err := docker.New()
	if err != nil {
		return planner, errors.Wrap(err, "initialize docker plugin")
	}
	planner.AddPlugin(pluginDocker)

	pluginJournald, err := journald.New()
	if err != nil {
		return planner, errors.Wrap(err, "initialize journald plugin")
	}
	planner.AddPlugin(pluginJournald)

	if c.KubernetesClient != nil {
		pluginKubernetes, err := pluginskubernetes.New(c.KubernetesClient, c.KubernetesClientConfig)
		if err != nil {
			return planner, errors.Wrap(err, "initialize kubernetes plugin")
		}
		planner.AddPlugin(pluginKubernetes)
	}

	pluginRetraced, err := retraced.New()
	if err != nil {
		return planner, errors.Wrap(err, "initialize retraced plugin")
	}
	planner.AddPlugin(pluginRetraced)

	return planner, nil
}
