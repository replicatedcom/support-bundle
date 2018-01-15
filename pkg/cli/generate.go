package cli

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

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
	CustomerEndpoint string
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
		if err != nil {
			return errors.Wrap(err, "initialize journald plugin")
		}
		planner.AddPlugin(pluginJournald)
	}

	if opts.EnableKubernetes {
		pluginKubernetes, err := kubernetes.New()
		if err != nil {
			return errors.Wrap(err, "initialize kubernetes plugin")
		}
		planner.AddPlugin(pluginKubernetes)
	}

	if opts.EnableRetraced {
		pluginRetraced, err := retraced.New()
		if err != nil {
			return errors.Wrap(err, "initialize retraced plugin")
		}
		planner.AddPlugin(pluginRetraced)
	}

	lifecycleTasks := types.DefaultLifecycle

	builtTasks := lifecycle.Build(lifecycleTasks)

	if err = lifecycle.Run(builtTasks); err != nil {
		return errors.Wrap(err, "running tasks")
	}

	// graphQLClient := graphql.NewClient(opts.CustomerEndpoint, http.DefaultClient)
	// specs, err := resolveLocalSpecs(opts)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to resolve specs")
	// }

	// if opts.CustomerID != "" {
	// 	remoteSpecs, err := getCustomerSpecs(graphQLClient, opts.CustomerID)
	// 	if err != nil {
	// 		return errors.Wrap(err, "get customer specs")
	// 	}
	// 	specs = append(specs, remoteSpecs...)

	// 	// lifecycleTasks, err := lifecycle.Plan(remoteSpecs.Lifecycle)

	// 	return nil
	// }

	// var tasks = planner.Plan(specs)
	// if len(tasks) == 0 {
	// 	return errors.New("No tasks defined")
	// }

	// lifecycleTasks, err := lifecycle.Build()
	// if err != nil {
	// 	return errors.Wrap("build lifecycle tasks")
	// }

	// fileInfo, err := bundle.Generate(tasks, time.Duration(time.Second*time.Duration(opts.TimeoutSeconds)), opts.BundlePath)

	// if err != nil {
	// 	return errors.Wrap(err, "generate bundle")
	// }

	// if opts.CustomerID != "" {
	// 	bundleID, url, err := graphQLClient.GetSupportBundleUploadURI(opts.CustomerID, fileInfo.Size())

	// 	if err != nil {
	// 		return errors.Wrap(err, "get presigned URL")
	// 	}

	// 	err = putObject(fileInfo, url)
	// 	if err != nil {
	// 		return errors.Wrap(err, "uploading to presigned URL")
	// 	}

	// 	if err = graphQLClient.UpdateSupportBundleStatus(opts.CustomerID, bundleID, "uploaded"); err != nil {
	// 		return errors.Wrap(err, "updating bundle status")
	// 	}
	// }

	return nil
}

func putObject(fi os.FileInfo, url *url.URL) error {
	file, err := os.Open(fi.Name())
	if err != nil {
		return errors.Wrap(err, "opening file for upload")
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", url.String(), file)
	if err != nil {
		return errors.Wrap(err, "making request")
	}
	req.ContentLength = fi.Size()
	req.Header.Set("Content-Type", "application/tar+gzip")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "completing request")
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Error uploading support bundle, got %s", res.Status)
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
		defaultSpecs, err := bundle.DefaultSpecs()
		if err != nil {
			return nil, errors.Wrap(err, "get default specs")
		}

		specs = append(defaultSpecs, specs...)
	}

	return specs, nil
}

func getCustomerSpecs(gqlClient *graphql.Client, customerID string) ([]types.Spec, error) {
	remoteSpecBody, err := gqlClient.GetCustomerSpec(customerID)
	if err != nil {
		return nil, errors.Wrap(err, "get remote spec")
	}

	customerSpecs, err := spec.Parse(remoteSpecBody)
	if err != nil {
		return nil, errors.Wrap(err, "parse customer spec")
	}

	return customerSpecs, nil
}
