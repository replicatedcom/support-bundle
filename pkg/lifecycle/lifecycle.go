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

func (l *Lifecycle) Build(tasks []*types.LifecycleTask) {

	for _, task := range tasks {
		eventFn := resolveEvent(task)
		l.tasks = append(l.tasks, eventFn(task))
	}
}

func resolveEvent(t *types.LifecycleTask) Event {
	switch {
	case t.Message != nil:
		return MessageTask
	case t.BooleanPrompt != nil:
		return PromptTask
	case t.Generate != nil:
		return GenerateTask
	case t.Upload != nil:
		return UploadTask
	}

	return nil
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
