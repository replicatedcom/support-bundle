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
)

const (
	DefaultCustomerEndpoint       = "https://pg.replicated.com/graphql"
	DefaultGenerateTimeoutSeconds = 60
)

type GenerateOptions struct {
	CfgFiles            []string
	CfgDocs             []string
	BundlePath          string
	SkipDefault         bool
	TimeoutSeconds      int
	CustomerID          string
	CustomerEndpoint    string
	ConfirmUploadPrompt bool
	DenyUploadPrompt    bool
	Quiet               bool

	bundle.PlannerOptions
}

func (cli *Cli) Generate(opts GenerateOptions) error {
	planner, err := bundle.NewPlanner(opts.PlannerOptions, os.Getenv("IN_CONTAINER") != "")
	if err != nil {
		return errors.Wrap(err, "initialize planner")
	}

	customerEndpoint := opts.CustomerEndpoint
	if customerEndpoint == "" {
		customerEndpoint = DefaultCustomerEndpoint
	}

	graphQLClient := graphql.NewClient(customerEndpoint, http.DefaultClient)
	specs, err := resolveLocalSpecs(opts)
	if err != nil {
		return errors.Wrap(err, "resolve specs")
	}

	var customerDoc *types.Doc
	expectedDefaultTasks := 1 // there is always at least 1 for the version

	if opts.CustomerID != "" {
		customerDoc, err = getCustomerDoc(graphQLClient, opts.CustomerID)
		if err != nil {
			return errors.Wrap(err, "get customer specs")
		}
		specs = append(specs, customerDoc.Collect.V1...)
		specs = append(specs, bundle.CustomerJSONSpec(opts.CustomerID))

		if types.GetUseDefaults(customerDoc.Lifecycle) {
			specs = append(specs, bundle.DefaultSpecs()...)
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
		ConfirmUploadPrompt: opts.ConfirmUploadPrompt,
		DenyUploadPrompt:    opts.DenyUploadPrompt,
		Quiet:               opts.Quiet,
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
