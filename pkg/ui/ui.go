package ui

import (
	"io"
	"os"

	"github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func New(
	outW io.Writer, errW io.Writer,
	forceColor, noColor bool,
) cli.Ui {
	base := &cli.BasicUi{
		Writer:      outW,
		ErrorWriter: errW,
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

func isInteractive() bool {
	return terminal.IsTerminal(int(os.Stdin.Fd()))
}
