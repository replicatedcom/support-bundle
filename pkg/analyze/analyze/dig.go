package analyze

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyzer"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/resolver"
	"github.com/replicatedcom/support-bundle/pkg/fs"
	kubernetesclient "github.com/replicatedcom/support-bundle/pkg/kubernetes"
	"github.com/replicatedcom/support-bundle/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

func RunE(ctx context.Context, bundle string) ([]api.Result, error) {
	a, err := Get()
	if err != nil {
		return nil, err
	}
	return a.Execute(ctx, bundle)
}

func Get() (*Analyze, error) {
	// who injects the injectors?
	logLevel := viper.GetViper().GetString("log-level")
	debug := log.With(level.Debug(logger.New(logLevel)), "component", "injector", "phase", "instance.get")

	debug.Log("event", "injector.build")
	injector, err := buildInjector()
	if err != nil {
		debug.Log("event", "injector.build.fail")
		return nil, errors.Wrap(err, "build injector")
	}

	var analyze *Analyze

	// we return nil below , so the error will only ever be a construction error
	debug.Log("event", "injector.invoke")
	if err := injector.Invoke(func(a *Analyze) {
		debug.Log("event", "injector.invoke.resolve")
		analyze = a
	}); err != nil {
		debug.Log("event", "injector.invoke.fail")
		return nil, errors.Wrap(err, "resolve dependencies")
	}

	return analyze, nil
}

func buildInjector() (*dig.Container, error) {
	providers := []interface{}{
		viper.GetViper,

		logger.FromViper,
		fs.FromViper,

		KubernetesClientConfigOptional,
		KubernetesNewClientOptional,

		resolver.New,
		analyzer.New,

		New,
	}

	container := dig.New()

	for _, provider := range providers {
		err := container.Provide(provider)
		if err != nil {
			return nil, errors.Wrap(err, "register providers")
		}
	}

	return container, nil
}

func KubernetesClientConfigOptional(logger log.Logger) *restclient.Config {
	config, err := kubernetesclient.ClientConfig()
	if err != nil {
		level.Debug(logger).Log(
			"method", "analyze.KubernetesClientConfigOptional",
			"err", err)
	}
	return config
}

func KubernetesNewClientOptional(config *restclient.Config, logger log.Logger) (kubernetes.Interface, error) {
	if config == nil {
		return nil, nil
	}
	return kubernetesclient.NewClient(config)
}
