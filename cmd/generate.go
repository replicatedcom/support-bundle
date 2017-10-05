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
	"time"

	"github.com/replicatedcom/support-bundle/bundle"
	"github.com/replicatedcom/support-bundle/metrics"
	"github.com/replicatedcom/support-bundle/systemutil"

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

func init() {
	RootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generate(cmd *cobra.Command, args []string) error {
	jww.SetStdoutThreshold(jww.LevelTrace)

	jww.FEEDBACK.Println("Generating a new support bundle")

	var tasks = []bundle.Task{
		bundle.Task{
			Description: "Get File",
			ExecFunc:    systemutil.ReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"C:/Go/VERSION"},
		},

		bundle.Task{
			Description: "Get Other File",
			ExecFunc:    systemutil.ReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"C:/Go/README.md"},
		},

		bundle.Task{
			Description: "Run Command",
			ExecFunc:    systemutil.RunCommand,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"docker", "ps"},
		},

		bundle.Task{
			Description: "Run Other Command",
			ExecFunc:    systemutil.RunCommand,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"docker", "info"},
		},

		bundle.Task{
			Description: "Docker run command in container",
			ExecFunc:    systemutil.DockerRunCommand,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057", "root", "ls", "-a"},
		},

		bundle.Task{
			Description: "Docker run command in container with large output",
			ExecFunc:    systemutil.DockerRunCommand,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057", "root", "ls", "-a", "-R"},
		},

		bundle.Task{
			Description: "Docker run command in container, timeout",
			ExecFunc:    systemutil.DockerRunCommand,
			Timeout:     time.Duration(time.Second * 1),
			Args:        []string{"7e47d28f0057", "root", "sleep", "1m"},
		},

		bundle.Task{
			Description: "Docker read file from container",
			ExecFunc:    systemutil.DockerReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057", "/usr/local/bin/docker-entrypoint.sh"},
		},

		bundle.Task{
			Description: "Docker read file from container",
			ExecFunc:    systemutil.DockerReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057", "/usr/local/bin/"},
		},

		bundle.Task{
			Description: "Docker ps",
			ExecFunc:    metrics.Dockerps,
			Timeout:     time.Duration(time.Second * 15),
		},

		bundle.Task{
			Description: "Docker info",
			ExecFunc:    metrics.DockerInfo,
			Timeout:     time.Duration(time.Second * 15),
		},

		bundle.Task{
			Description: "Docker logs",
			ExecFunc:    metrics.DockerLogs,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057"},
		},

		bundle.Task{
			Description: "Docker inspect",
			ExecFunc:    metrics.DockerInspect,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057"},
		},

		bundle.Task{
			Description: "System hostname",
			ExecFunc:    metrics.Hostname,
			Timeout:     time.Duration(time.Second * 1),
		},

		bundle.Task{
			Description: "System loadavg",
			ExecFunc:    metrics.LoadAvg,
			Timeout:     time.Duration(time.Second * 1),
		},

		bundle.Task{
			Description: "System uptime",
			ExecFunc:    metrics.Uptime,
			Timeout:     time.Duration(time.Second * 1),
		},
	}

	if _, err := bundle.Generate(tasks); err != nil {
		jww.ERROR.Fatal(err)
	}

	return nil
}
