package ui

import (
	"io"
	"os"

	"github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func New(inR io.Reader, outW io.Writer, errW io.Writer) cli.Ui {
	return &cli.BasicUi{
		Reader:      inR,
		Writer:      outW,
		ErrorWriter: errW,
	}
}

func Colored(base cli.Ui, force bool) cli.Ui {
	if !(isInteractive() || force) {
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
