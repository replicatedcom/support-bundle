// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io/ioutil"
	"time"

	"github.com/replicatedcom/support-bundle/bundle"
	"github.com/replicatedcom/support-bundle/plugins/core"
	"github.com/replicatedcom/support-bundle/plugins/docker"
	"github.com/replicatedcom/support-bundle/spec"
	"github.com/replicatedcom/support-bundle/types"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new support bundle",
	Long:  `Collect data and generate a new support bundle`,
	RunE:  generate,
}

var bundlePath string
var skipDefault bool
var timeoutSeconds int

func init() {
	RootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	generateCmd.Flags().StringVar(&bundlePath, "out", "supportbundle.tar.gz", "Path where the generated bundle should be stored")
	generateCmd.Flags().BoolVarP(&skipDefault, "skipdefault", "s", false, "If present, skip the default support bundle files")
	generateCmd.Flags().IntVar(&timeoutSeconds, "timeout", 60, "The overall support bundle generation timeout")
}

func generate(cmd *cobra.Command, args []string) error {
	jww.SetStdoutThreshold(jww.LevelTrace)

	jww.FEEDBACK.Println("Generating a new support bundle")

	var specs []types.Spec

	if cfgFile != "" {
		yaml, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			jww.ERROR.Fatal(err)
		}

		specs, err = spec.Parse(yaml)
		if err != nil {
			jww.ERROR.Fatal(err)
		}
	}
	if !skipDefault {
		defaultSpecs, err := bundle.DefaultSpecs()
		if err != nil {
			jww.ERROR.Fatal(err)
		}

		specs = append(defaultSpecs, specs...)
	}

	d, err := docker.New()
	if err != nil {
		jww.ERROR.Fatal(err)
	}
	planner := bundle.Planner{
		Plugins: map[string]types.Plugin{
			"core":   core.New(),
			"docker": d,
		},
	}

	var tasks = planner.Plan(specs)

	if err := bundle.Generate(tasks, time.Duration(time.Second*time.Duration(timeoutSeconds)), bundlePath); err != nil {
		jww.ERROR.Fatal(err)
	}

	if !skipDefault {
		jww.FEEDBACK.Printf("Support bundle generated at %s and does contain the default files\n", bundlePath)
	} else {
		jww.FEEDBACK.Printf("Support bundle generated at %s and does not contain the default files\n", bundlePath)
	}

	return nil
}
