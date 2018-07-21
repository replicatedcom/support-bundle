package cli

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/bundle"
	"github.com/replicatedcom/support-bundle/pkg/collect/graphql"
	"github.com/replicatedcom/support-bundle/pkg/collect/lifecycle"
	"github.com/replicatedcom/support-bundle/pkg/collect/spec"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/docker"
	kubernetesclient "github.com/replicatedcom/support-bundle/pkg/kubernetes"
	"github.com/replicatedcom/support-bundle/pkg/logger"
	jww "github.com/spf13/jwalterweatherman"
	"k8s.io/client-go/kubernetes"
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
	RequireJournald   bool
	RequireKubernetes bool
	RequireRetraced   bool
}

func (cli *Cli) Generate(opts GenerateOptions) error {
	planner, err := cli.newPlanner(opts)
	if err != nil {
		return errors.Wrap(err, "initialize planner")
	}

	graphQLClient := graphql.NewClient(opts.CustomerEndpoint, http.DefaultClient)
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

	lf := lifecycle.Lifecycle{
		BundleTasks:         tasks,
		GenerateTimeout:     opts.TimeoutSeconds,
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

func (cli *Cli) newPlanner(opts GenerateOptions) (*bundle.Planner, error) {
	plannerOpts := opts.PlannerOptions
	plannerOpts.InContainer = os.Getenv("IN_CONTAINER") != ""
	if opts.EnableDocker || (plannerOpts.InContainer && (opts.EnableCore || opts.EnableJournald)) {
		kitLog := logger.New(
			logger.LevelFromJWWThreshold(jww.GetLogThreshold()),
		)
		dockerClient, err := docker.NewEnvClient(context.Background(), kitLog)
		if err != nil {
			return nil, errors.Wrap(err, "get docker client from environment")
		}
		plannerOpts.DockerClient = dockerClient
	}
	if opts.EnableKubernetes {
		var client kubernetes.Interface
		clientConfig, err := kubernetesclient.ClientConfig()
		if err != nil {
			if opts.RequireKubernetes {
				return nil, errors.Wrap(err, "get kubernetes client config")
			}
			jww.DEBUG.Printf("get kubernetes client config: %s", err.Error())
		} else {
			client, err = kubernetesclient.NewClient(clientConfig)
			if err != nil {
				if opts.RequireKubernetes {
					return nil, errors.Wrap(err, "get kubernetes client")
				}
				jww.DEBUG.Printf("get kubernetes client: %s", err.Error())
			}
		}
		if client != nil {
			plannerOpts.KubernetesClientConfig = clientConfig
			plannerOpts.KubernetesClient = client
		}
	}

	planner := bundle.NewPlanner(plannerOpts)
	return &planner, nil
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
