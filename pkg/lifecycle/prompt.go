package lifecycle

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func PromptTask(task *types.LifecycleTask) Task {
	return func(l *Lifecycle) (bool, error) {
		reader := bufio.NewReader(os.Stdin)

		for {
			def := "[y/N]"
			if task.BooleanPrompt.Default {
				def = "[Y/n]"
			}
			fmt.Printf("%s %s: ", task.BooleanPrompt.Contents, def)

			response, err := reader.ReadString('\n')
			if err != nil {
				return false, errors.Wrap(err, "prompt user")
			}

			response = strings.ToLower(strings.TrimSpace(response))

			if response == "" {
				return task.BooleanPrompt.Default, nil
			}

			if response == "y" {
				return true, nil
			}

			if response == "n" {
				return false, nil
			}
		}
	}
}
