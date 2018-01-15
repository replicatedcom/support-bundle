package lifecycle

import (
	"os"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/graphql"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Lifecycle struct {
	GenerateTimeout    int
	SkipPrompts        bool
	GenerateBundlePath string
	UploadCustomerID   string
	GraphQLClient      *graphql.Client
	FileInfo           os.FileInfo
	BundleTasks        []types.Task
	tasks              []Task
}

type Event func(*types.LifecycleTask) Task

type Task func(*Lifecycle) (bool, error)

func (l *Lifecycle) Build(tasks []*types.LifecycleTask) error {

	for _, task := range tasks {
		eventFn, err := resolveEvent(task)
		if err != nil {
			return errors.Wrap(err, "resolve event")
		}
		l.tasks = append(l.tasks, eventFn(task))
	}

	return nil
}

func resolveEvent(t *types.LifecycleTask) (Event, error) {
	switch {
	case t.Message != nil:
		return MessageTask, nil
	case t.Generate != nil:
		return GenerateTask, nil
	case t.Upload != nil:
		return UploadTask, nil
	}

	return nil, errors.New("no valid event found, requires one of: generate, message, boolean, upload")
}

func (l *Lifecycle) Run() error {
	for _, t := range l.tasks {
		cont, err := t(l)

		if err != nil {
			return errors.Wrap(err, "running task")
		}

		if !cont {
			break
		}
	}

	return nil
}
