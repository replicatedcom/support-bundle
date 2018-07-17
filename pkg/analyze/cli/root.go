package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
func RootCmd() *cobra.Command {
	version.Init()
	cmd := &cobra.Command{
		Use:     "analyze",
		Short:   "troubleshoot analysis tool",
		Version: getVersionString(),
	}
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /etc/replicated/analyze.yaml)")
	cmd.PersistentFlags().String("log-level", "off", "Log level")

	// required
	// cmd.PersistentFlags().String("customer-id", "", "customer id")

	// sub-commands
	cmd.AddCommand(RunCmd())

	viper.BindPFlags(cmd.Flags())
	viper.BindPFlags(cmd.PersistentFlags())

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd().Execute(); err != nil {
		os.Exit(2)
	}
}

// getVersionString formats the version string for --version flag cli output.
func getVersionString() string {
	return fmt.Sprintf(
		`%s (git="%s", built="%s")`,
		version.Version(), version.GitSHA(),
		version.BuildTime().Format(time.UnixDate))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/replicated")
		viper.AddConfigPath("/etc/sysconfig/")
		viper.SetConfigName("analyze")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
