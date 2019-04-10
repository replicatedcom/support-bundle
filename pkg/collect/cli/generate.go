package cli

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/bundle"
	"github.com/replicatedcom/support-bundle/pkg/collect/graphql"
	"github.com/replicatedcom/support-bundle/pkg/collect/lifecycle"
	"github.com/replicatedcom/support-bundle/pkg/collect/spec"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/util"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	DefaultEndpoint               = "https://pg.replicated.com/graphql"
	DefaultGenerateTimeoutSeconds = 60
)

type GenerateOptions struct {
	CfgFiles            []string
	CfgDocs             []string
	BundlePath          string
	SkipDefault         bool
	TimeoutSeconds      int
	ConfirmUploadPrompt bool
	DenyUploadPrompt    bool
	Quiet               bool
	Endpoint            string
	ChannelID           string

	CustomerID string // Deprecated

	bundle.PlannerOptions
}

func (cli *Cli) Generate(opts GenerateOptions) error {
	logOutput := &util.Buffer{}
	jww.SetLogOutput(logOutput)
	jww.SetLogThreshold(jww.LevelDebug)

	planner, err := bundle.NewPlanner(opts.PlannerOptions, os.Getenv("IN_CONTAINER") != "", logOutput)
	if err != nil {
		return errors.Wrap(err, "initialize planner")
	}

	endpoint := opts.Endpoint
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	graphQLClient := graphql.NewClient(endpoint, http.DefaultClient)
	specs, err := resolveLocalSpecs(opts)
	if err != nil {
		return errors.Wrap(err, "resolve specs")
	}

	var customerDoc *types.Doc
	var channelDoc *types.Doc
	expectedDefaultTasks := 1 // there is always at least 1 for the version

	// this next if statement and included scope is deprecated
	if opts.CustomerID != "" {
		customerDoc, err = getCustomerDoc(graphQLClient, opts.CustomerID)
		if err != nil {
			return errors.Wrap(err, "get customer specs")
		}
		specs = append(specs, customerDoc.Collect.V1...)
		specs = append(specs, bundle.CustomerJSONSpec(opts.CustomerID))

		if types.GetUseDefaults(customerDoc.Lifecycle) {
			defaultSpecs, err := bundle.DefaultSpecs()
			if err != nil {
				return errors.Wrap(err, "get default spec")
			}
			specs = append(specs, defaultSpecs...)
		}

		expectedDefaultTasks++
	}

	if opts.ChannelID != "" {
		channelDoc, err = getChannelDoc(graphQLClient, opts.ChannelID)
		if err != nil {
			return errors.Wrap(err, "get channel spec")
		}
		specs = append(specs, channelDoc.Collect.V1...)
		specs = append(specs, bundle.ChannelJSONSpec(opts.ChannelID))

		if types.GetUseDefaults(channelDoc.Lifecycle) {
			defaultSpecs, err := bundle.DefaultSpecs()
			if err != nil {
				return errors.Wrap(err, "get default spec")
			}
			specs = append(specs, defaultSpecs...)
		}

		expectedDefaultTasks++
	}

	var tasks = planner.Plan(specs)
	if len(tasks) < expectedDefaultTasks {
		return errors.New("No tasks defined")
	}

	timeoutSeconds := opts.TimeoutSeconds
	if timeoutSeconds == 0 {
		timeoutSeconds = DefaultGenerateTimeoutSeconds
	}

	lf := lifecycle.Lifecycle{
		BundleTasks:         tasks,
		GenerateTimeout:     timeoutSeconds,
		GenerateBundlePath:  opts.BundlePath,
		GraphQLClient:       graphQLClient,
		UploadCustomerID:    opts.CustomerID,
		UploadChannelID:     opts.ChannelID,
		ConfirmUploadPrompt: opts.ConfirmUploadPrompt,
		DenyUploadPrompt:    opts.DenyUploadPrompt,
		Quiet:               opts.Quiet,
	}

	lifecycleTasks := types.DefaultLifecycleTasks
	if channelDoc != nil && channelDoc.Lifecycle != nil {
		lifecycleTasks = channelDoc.Lifecycle
	} else if customerDoc != nil && customerDoc.Lifecycle != nil {
		lifecycleTasks = customerDoc.Lifecycle
	}

	if opts.CustomerID == "" && opts.ChannelID == "" {
		lifecycleTasks = types.GenerateOnlyLifecycleTasks
	}

	if err = lf.Build(lifecycleTasks); err != nil {
		return errors.Wrap(err, "build lifecycle events")
	}

	if err = lf.Run(); err != nil {
		return errors.Wrap(err, "running tasks")
	}

	return nil
}

func resolveLocalSpecs(opts GenerateOptions) ([]types.Spec, error) {
	specs := []types.Spec{}

	specs = append(specs, bundle.SupportBundleSpecs()...)

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

	if opts.CustomerID == "" && opts.ChannelID == "" && !opts.SkipDefault {
		defaultSpecs, err := bundle.DefaultSpecs()
		if err != nil {
			return nil, errors.Wrap(err, "get default spec")
		}
		specs = append(defaultSpecs, specs...)
	}

	return specs, nil
}

// getCustomerDoc is deprecated
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

func getChannelDoc(gqlClient *graphql.Client, channelID string) (*types.Doc, error) {
	remoteSpecBody, err := gqlClient.GetChannelSpec(channelID)
	if err != nil {
		return nil, errors.Wrap(err, "get remote spec")
	}

	channelDoc, err := spec.Unmarshal(remoteSpecBody)
	if err != nil {
		return nil, errors.Wrap(err, "parse channel spec")
	}

	return channelDoc, nil
}
