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

	var tasks = []bundle.Task{}

	if _, err := bundle.Generate(tasks, time.Duration(time.Second*15)); err != nil {
		jww.ERROR.Fatal(err)
	}

	return nil
}
