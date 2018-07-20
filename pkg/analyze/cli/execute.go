package cli

import (
	"os"

	pkgerrors "github.com/replicatedcom/support-bundle/pkg/errors"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	c := RootCmd()
	if err := c.Execute(); err != nil {
		exitCode := 1
		if cmdErr, ok := err.(pkgerrors.CmdError); !ok {
			c.Println("Error:", err.Error())
			c.Printf("Run '%v --help' for usage.\n", c.CommandPath())
		} else if cmdErr.ExitCode != 0 {
			exitCode = cmdErr.ExitCode
		}
		os.Exit(exitCode)
	}
}
