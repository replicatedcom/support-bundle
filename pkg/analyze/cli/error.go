package cli

import (
	"fmt"

	"github.com/mitchellh/cli"
)

// outputCmdError will display any unexpected errors to the end user.
func outputCmdError(err error, ui cli.Ui, logLevel string) {
	if logLevel == "debug" {
		ui.Error(fmt.Sprintf("There was an unexpected error! %+v", err))
	} else {
		ui.Error(fmt.Sprintf("There was an unexpected error! %v", err))
	}
	ui.Output("")

	if logLevel != "debug" {
		ui.Info("An unexpected error occured. Please re-run with --log-level=debug and include the output in any support inquiries.")
	}
}
