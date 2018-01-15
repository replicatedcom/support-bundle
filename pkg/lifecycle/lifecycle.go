package lifecycle

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Event func(*types.LifecycleTask) Task

type Task func() error

func Build(tasks []*types.LifecycleTask) []Task {
	var t []Task

	for _, task := range tasks {
		eventFn := resolveEvent(task)
		t = append(t, eventFn(task))
	}

	return t
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

func Run(tasks []Task) error {
	for _, t := range tasks {
		if err := t(); err != nil {
			return errors.Wrap(err, "running task")
		}
	}

	return nil
}

// switch {
// case spec.KubernetesAPIVersions != nil:
// 	return p.planner.APIVersions
// case spec.KubernetesClusterInfo != nil:
// 	retu
