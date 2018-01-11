package lifecycle

import "github.com/replicatedcom/support-bundle/pkg/types"

type Event func(*types.LifecycleTask) Task

type Task func() error

func Build(tasks []*types.LifecycleTask) ([]Task, error) {
	var t []Task

	for _, task := range tasks {
		eventFn := resolveEvent(task)
		t = append(t, eventFn(task))
	}

	return t, nil
}

func resolveEvent(t *types.LifecycleTask) Event {
	switch {
	case t.Message != nil:
		return MessageTask
	}

	return nil
}

// func Run([]*Event) error {

// }

// switch {
// case spec.KubernetesAPIVersions != nil:
// 	return p.planner.APIVersions
// case spec.KubernetesClusterInfo != nil:
// 	retu
