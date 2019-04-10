package lifecycle

import (
	"fmt"
	"os"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type MessageTask struct {
	Options types.MessageOptions
}

func (t *MessageTask) Execute(l *Lifecycle) (bool, error) {
	if !l.Quiet {
		fmt.Fprintln(os.Stdout, t.Options.Contents)
	}
	return true, nil
}
