package cli

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/bundle"
	"github.com/replicatedcom/support-bundle/pkg/graphql"
	"github.com/replicatedcom/support-bundle/pkg/lifecycle"
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
	CfgFiles          []string
	CfgDocs           []string
	BundlePath        string
	SkipDefault       bool
	TimeoutSeconds    int
	EnableCore        bool
	EnableDocker      bool
	EnableJournald    bool
	RequireJournald   bool
	EnableKubernetes  bool
	RequireKubernetes bool
	EnableRetraced    bool
	RequireRetraced   bool
	CustomerID        string
	CustomerEndpoint  string
	SkipPrompts       bool
}

func (cli *Cli) Generate(opts GenerateOptions) error {
	var planner bundle.Planner

	pluginSupportBundle, err := supportbundle.New()
	if err != nil {
		return errors.Wrap(err, "initialize supportbundle plugin")
	}
	planner.AddPlugin(pluginSupportBundle)

	if opts.EnableCore {
		pluginCore, err := core.New()
		if err != nil {
			return errors.Wrap(err, "initialize core plugin")
		}
		planner.AddPlugin(pluginCore)
	}

	if opts.EnableDocker {
		pluginDocker, err := docker.New()
		if err != nil {
			return errors.Wrap(err, "initialize docker plugin")
		}
		planner.AddPlugin(pluginDocker)
	}

	if opts.EnableJournald {
		pluginJournald, err := journald.New()
		if err != nil && opts.RequireJournald {
			return errors.Wrap(err, "initialize journald plugin")
		} else if err != nil {
			jww.DEBUG.Printf("initialize journald plugin: %s", err.Error())
		}

		planner.AddPlugin(pluginJournald)
	}

	if opts.EnableKubernetes {
		pluginKubernetes, err := kubernetes.New()

		if err != nil && opts.RequireKubernetes {
			return errors.Wrap(err, "initialize kubernetes plugin")
		} else if err != nil {
			jww.DEBUG.Printf("initialize kubernetes plugin: %s", err.Error())
		}

		planner.AddPlugin(pluginKubernetes)
	}

	if opts.EnableRetraced {
		pluginRetraced, err := retraced.New()
		if err != nil && opts.RequireRetraced {
			return errors.Wrap(err, "initialize retraced plugin")
		} else if err != nil {
			jww.DEBUG.Printf("initialize retraced plugin: %s", err.Error())
		}
		planner.AddPlugin(pluginRetraced)
	}

	graphQLClient := graphql.NewClient(opts.CustomerEndpoint, http.DefaultClient)
	specs, err := resolveLocalSpecs(opts)
	if err != nil {
		return errors.Wrap(err, "failed to resolve specs")
	}

	var customerDoc *types.Doc
	expectedDefaultTasks := 1 // there is always at least 1 for the version

	if opts.CustomerID != "" {
		customerDoc, err = getCustomerDoc(graphQLClient, opts.CustomerID)
		if err != nil {
			return errors.Wrap(err, "get customer specs")
		}
		specs = append(specs, customerDoc.Specs...)
		specs = append(specs, bundle.CustomerJsonSpec(opts.CustomerID))
		expectedDefaultTasks += 1
	}

	var tasks = planner.Plan(specs)
	if len(tasks) < expectedDefaultTasks {
		return errors.New("No tasks defined")
	}

	lf := lifecycle.Lifecycle{
		BundleTasks:        tasks,
		GenerateTimeout:    opts.TimeoutSeconds,
		GenerateBundlePath: opts.BundlePath,
		GraphQLClient:      graphQLClient,
		UploadCustomerID:   opts.CustomerID,
		SkipPrompts:        opts.SkipPrompts,
	}

	lt := types.DefaultLifecycleTasks
	if customerDoc != nil && customerDoc.Lifecycle != nil {
		lt = customerDoc.Lifecycle
	}

	if opts.CustomerID == "" {
		lt = types.GenerateOnlyLifecycleTasks
	}

	if err = lf.Build(lt); err != nil {
		return errors.Wrap(err, "build lifecycle events")
	}

	if err = lf.Run(); err != nil {
		return errors.Wrap(err, "running tasks")
	}

	return nil
}

func resolveLocalSpecs(opts GenerateOptions) ([]types.Spec, error) {
	specs := []types.Spec{bundle.SupportBundleVersionSpec()}

	for _, cfgFile := range opts.CfgFiles {
		yaml, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return nil, errors.Wrap(err, "read spec file")
		}

		fileSpecs, err := spec.Parse(yaml)
		if err != nil {
			return nil, errors.Wrap(err, "parse config file spec")
		}
		specs = append(specs, fileSpecs...)
	}

	for _, cfgDoc := range opts.CfgDocs {
		argSpecs, err := spec.Parse([]byte(cfgDoc))
		if err != nil {
			return nil, errors.Wrap(err, "parse config doc spec")
		}
		specs = append(specs, argSpecs...)
	}

	if opts.CustomerID == "" && !opts.SkipDefault {
		defaultSpecs := bundle.DefaultSpecs()

		specs = append(defaultSpecs, specs...)
	}

	return specs, nil
}

func getCustomerDoc(gqlClient *graphql.Client, customerID string) (*types.Doc, error) {
	remoteSpecBody, err := gqlClient.GetCustomerSpec(customerID)
	if err != nil {
		return nil, errors.Wrap(err, "get remote spec")
	}

	customerDoc, err := spec.Unmarshal(remoteSpecBody)
	if err != nil {
		return nil, errors.Wrap(err, "parse customer spec")
	}

	return customerDoc, nil
}
