package lifecycle

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type NotesTask struct {
	Options types.NotesOptions
}

func (task *NotesTask) Execute(l *Lifecycle) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(task.Options.Prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false, errors.Wrap(err, "reading input")
	}
	l.Notes = text
	return true, nil
}
