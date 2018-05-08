package lifecycle

import (
	"fmt"
	"os"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

type MessageTask struct {
	Options types.MessageOptions
}

func (t *MessageTask) Execute(l *Lifecycle) (bool, error) {
	if !l.Quiet {
		fmt.Fprintln(os.Stderr, t.Options.Contents)
	}
	return true, nil
}
