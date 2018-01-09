package cli

import (
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/bundle"
	"github.com/replicatedcom/support-bundle/pkg/graphql"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker"
	"github.com/replicatedcom/support-bundle/pkg/plugins/journald"
	"github.com/replicatedcom/support-bundle/pkg/plugins/kubernetes"
	"github.com/replicatedcom/support-bundle/pkg/plugins/retraced"
	"github.com/replicatedcom/support-bundle/pkg/plugins/supportbundle"
	"github.com/replicatedcom/support-bundle/pkg/spec"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

type GenerateOptions struct {
	CfgFiles         []string
	CfgDocs          []string
	BundlePath       string
	SkipDefault      bool
	TimeoutSeconds   int
	EnableCore       bool
	EnableDocker     bool
	EnableJournald   bool
	EnableKubernetes bool
	EnableRetraced   bool
	CustomerID       string
}

func (cli *Cli) Generate(opts GenerateOptions) error {
	jww.FEEDBACK.Println("Generating a new support bundle")

	specs, err := resolveSpecs(opts)
	if err != nil {
		return errors.Wrap(err, "failed to resolve specs")
	}

	var planner bundle.Planner

	pluginSupportBundle, err := supportbundle.New()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize supportbundle plugin")
	}
	planner.AddPlugin(pluginSupportBundle)

	if opts.EnableCore {
		pluginCore, err := core.New()
		if err != nil {
			return errors.Wrap(err, "Failed to initialize core plugin")
		}
		planner.AddPlugin(pluginCore)
	}

	if opts.EnableDocker {
		pluginDocker, err := docker.New()
		if err != nil {
			return errors.Wrap(err, "Failed to initialize docker plugin")
		}
		planner.AddPlugin(pluginDocker)
	}

	if opts.EnableJournald {
		pluginJournald, err := journald.New()
		if err != nil {
			return errors.Wrap(err, "Failed to initialize journald plugin")
		}
		planner.AddPlugin(pluginJournald)
	}

	if opts.EnableKubernetes {
		pluginKubernetes, err := kubernetes.New()
		if err != nil {
			return errors.Wrap(err, "Failed to initialize kubernetes plugin")
		}
		planner.AddPlugin(pluginKubernetes)
	}

	if opts.EnableRetraced {
		pluginRetraced, err := retraced.New()
		if err != nil {
			return errors.Wrap(err, "Failed to initialize retraced plugin")
		}
		planner.AddPlugin(pluginRetraced)
	}

	var tasks = planner.Plan(specs)
	if len(tasks) == 0 {
		return errors.New("No tasks defined")
	}

	if err := bundle.Generate(tasks, time.Duration(time.Second*time.Duration(opts.TimeoutSeconds)), opts.BundlePath); err != nil {
		return errors.Wrap(err, "Failed to generate bundle")
	}

	jww.FEEDBACK.Printf("Support bundle generated at %s", opts.BundlePath)

	return nil
}

func resolveSpecs(opts GenerateOptions) ([]types.Spec, error) {
	specs := []types.Spec{bundle.SupportBundleVersionSpec()}

	if opts.CustomerID != "" {
		sbs := &graphql.SupportBundleSpec{
			CustomerID: opts.CustomerID,
		}

		remoteSpecBody, err := sbs.Get()
		if err != nil {
			return nil, errors.Wrap(err, "getting remote spec")
		}

		customerSpecs, err := spec.Parse(remoteSpecBody)
		if err != nil {
			return nil, errors.Wrap(err, "parsing customer spec")
		}

		specs = append(specs, customerSpecs...)
	}

	for _, cfgFile := range opts.CfgFiles {
		yaml, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read spec file")
		}

		fileSpecs, err := spec.Parse(yaml)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse spec")
		}
		specs = append(specs, fileSpecs...)
	}

	for _, cfgDoc := range opts.CfgDocs {
		argSpecs, err := spec.Parse([]byte(cfgDoc))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse spec")
		}
		specs = append(specs, argSpecs...)
	}

	if !opts.SkipDefault {
		defaultSpecs, err := bundle.DefaultSpecs()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get default specs")
		}

		specs = append(defaultSpecs, specs...)
	}

	return specs, nil
}
