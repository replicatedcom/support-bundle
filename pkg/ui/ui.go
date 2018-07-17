package ui

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

func FromViper(v *viper.Viper) cli.Ui {
	return New(
		v.GetBool("force-color"),
		v.GetBool("no-color"))
}

func New(forceColor, noColor bool) cli.Ui {
	base := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	if !isInteractive() && !forceColor {
		return base
	}

	if noColor {
		return base
	}

	return &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorYellow,
		InfoColor:   cli.UiColorGreen,
		Ui:          base,
	}
}

// todo detect if this is an interactive session and/or if we have a tty
func isInteractive() bool {
	return true
}
