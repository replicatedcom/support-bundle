package lifecycle

import (
	"fmt"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

type MessageTask struct {
	Options types.MessageOptions
}

func (t *MessageTask) Execute(l *Lifecycle) (bool, error) {
	fmt.Println(t.Options.Contents)
	return true, nil
}
