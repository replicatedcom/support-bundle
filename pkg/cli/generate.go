package cli

import (
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/bundle"
	"github.com/replicatedcom/support-bundle/pkg/plugins/core"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker"
	"github.com/replicatedcom/support-bundle/pkg/plugins/supportbundle"
	"github.com/replicatedcom/support-bundle/pkg/spec"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func (cli *Cli) Generate(cfgFiles []string, cfgDocs []string, bundlePath string, skipDefault bool, timeoutSeconds int, enableKubernetes bool) error {
	jww.FEEDBACK.Println("Generating a new support bundle")
	var specs []types.Spec

	for _, cfgFile := range cfgFiles {
		yaml, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return errors.Wrap(err, "Failed to read spec file")
		}

		fileSpecs, err := spec.Parse(yaml)
		if err != nil {
			return errors.Wrap(err, "Failed to parse spec")
		}
		specs = append(specs, fileSpecs...)
	}

	for _, cfgDoc := range cfgDocs {
		argSpecs, err := spec.Parse([]byte(cfgDoc))
		if err != nil {
			return errors.Wrap(err, "Failed to parse spec")
		}
		specs = append(specs, argSpecs...)
	}

	if !skipDefault {
		defaultSpecs, err := bundle.DefaultSpecs()
		if err != nil {
			return errors.Wrap(err, "Failed to get default specs")
		}

		specs = append(defaultSpecs, specs...)
	}

	pluginSupportBundle, err := supportbundle.New()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize supportbundle plugin")
	}

	pluginCore, err := core.New()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize core plugin")
	}

	pluginDocker, err := docker.New()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize docker plugin")
	}

	// j, err := journald.New()
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to initialize journald plugin")
	// }

	planner := bundle.Planner{
		SupportBundle: pluginSupportBundle,
		Core:          pluginCore,
		Docker:        pluginDocker,
	}

	// if enableKubernetes {
	// 	k, err := kubernetes.New()
	// 	if err != nil {
	// 		return errors.Wrap(err, "Failed to initialize kubernetes plugin")
	// 	}

	// 	planner.Plugins["kubernetes"] = k
	// }

	var tasks = planner.Plan(specs)
	if len(tasks) == 0 {
		return errors.New("No tasks defined")
	}

	if err := bundle.Generate(tasks, time.Duration(time.Second*time.Duration(timeoutSeconds)), bundlePath); err != nil {
		return errors.Wrap(err, "Failed to generate bundle")
	}

	if !skipDefault {
		jww.FEEDBACK.Printf("Support bundle generated at %s and does contain the default files", bundlePath)
	} else {
		jww.FEEDBACK.Printf("Support bundle generated at %s and does not contain the default files", bundlePath)
	}

	return nil
}
