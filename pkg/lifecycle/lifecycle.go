package lifecycle

import (
	"os"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/graphql"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Lifecycle struct {
	BundleTasks        []types.Task
	GenerateTimeout    int
	GenerateBundlePath string
	tasks              []Task
	FileInfo           os.FileInfo
	UploadCustomerID   string
	GraphQLClient      *graphql.Client
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
}

func resolveEvent(t *types.LifecycleTask) (Event, error) {
	switch {
	case t.Message != nil:
		return MessageTask, nil
	case t.BooleanPrompt != nil:
		return PromptTask, nil
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

// switch {
// case spec.KubernetesAPIVersions != nil:
// 	return p.planner.APIVersions
// case spec.KubernetesClusterInfo != nil:
// 	retu
