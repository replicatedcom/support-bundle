package commands

import (
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/replicatedcom/support-bundle/pkg/collect/cli"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type supportBundleOptions struct {
	cfgFile string
	verbose bool
	quiet   bool
}

func NewSupportBundleCommand(cli *cli.Cli) *cobra.Command {
	opts := supportBundleOptions{}

	cmd := NewRootCommand(cli)
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if opts.verbose {
			jww.SetStdoutThreshold(jww.LevelTrace)
		} else if opts.quiet {
			jww.SetStdoutThreshold(jww.LevelError)
		} else {
			//get log level from env var
			logLevel := os.Getenv("LOG_LEVEL")
			logLevel = strings.ToLower(logLevel)
			switch logLevel {
			case "error":
				jww.SetStdoutThreshold(jww.LevelError)
			case "warn":
				jww.SetStdoutThreshold(jww.LevelWarn)
			case "warning":
				jww.SetStdoutThreshold(jww.LevelWarn)
			case "info":
				jww.SetStdoutThreshold(jww.LevelInfo)
			case "debug":
				jww.SetStdoutThreshold(jww.LevelDebug)
			case "trace":
				jww.SetStdoutThreshold(jww.LevelTrace)
			default:
				jww.SetStdoutThreshold(jww.LevelError)
			}
		}
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cmd.PersistentFlags().StringVarP(&opts.cfgFile, "config-file", "c", "", "Config file to use")
	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "Enable verbose logging. Overrides any LOG_LEVEL from the environment.")
	cmd.PersistentFlags().BoolVarP(&opts.quiet, "quiet", "q", false, "If set, supress all non-error output and messages. Overrides any LOG_LEVEL from the environment.")

	cobra.OnInitialize(func() {
		initConfig(opts)
	})

	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig(opts supportBundleOptions) {
	if opts.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(opts.cfgFile)
	} else {
		// Search config in home directory with name ".support-bundle" (without extension).
		viper.SetConfigName(".support-bundle")
		home, err := homedir.Dir()
		if err != nil {
			jww.ERROR.Printf("Failed to find the user's home directory: %v", err)
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
