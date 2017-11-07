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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var cfgFile string
var cfgDoc string
var verbose bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "support-bundle",
	Short: "Generate and manage support bundles",
	Long: `A support bundle is an archive of files, output, metrics and state 
from a server that can be used to assist when troubleshooting a server.

The support-bundle utility can generate human readable archives, and can also 
be used to generate input to the support.io service.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		jww.SetLogOutput(os.Stderr)
		if verbose {
			jww.SetStdoutThreshold(jww.LevelTrace)
		} else {
			jww.SetStdoutThreshold(jww.LevelError)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "spec-file", "f", "", "config file (default is to run core tasks only)")
	RootCmd.PersistentFlags().StringVarP(&cfgDoc, "spec", "s", "", "config doc (default is to run core tasks only)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Turn on verbose logging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".support-bundle" (without extension).
		viper.SetConfigName(".support-bundle")
		home, err := homedir.Dir()
		if err != nil {
			jww.ERROR.Printf("Failed to find the user's home directory: %v\n", err)
		} else {
			viper.AddConfigPath(home)
		}
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		jww.DEBUG.Println("Using config file:", viper.ConfigFileUsed())
	}
}
