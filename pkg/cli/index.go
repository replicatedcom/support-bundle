package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/bundle"
	"github.com/replicatedcom/support-bundle/pkg/spec"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (cli *Cli) Index(cfgFiles []string, cfgDocs []string, skipDefault bool, format string) error {
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

	index := []types.Result{}

	// TODO: it would be nice if specs had descriptions rather than results
	for _, spec := range specs {
		if spec.Raw != "" {
			index = append(index, types.Result{Path: spec.Raw, Description: "Raw"})
		}
		if spec.Human != "" {
			index = append(index, types.Result{Path: spec.Human, Description: "Human"})
		}
		if spec.JSON != "" {
			index = append(index, types.Result{Path: spec.JSON, Description: "JSON"})
		}
	}

	switch format {
	case "yaml":
		b, err := yaml.Marshal(index)
		if err != nil {
			return errors.Wrap(err, "Failed to marshal index yaml")
		}
		fmt.Fprint(os.Stdout, string(b))

	default: // json
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(index); err != nil {
			return errors.Wrap(err, "Failed to encode index")
		}
	}

	return nil
}
