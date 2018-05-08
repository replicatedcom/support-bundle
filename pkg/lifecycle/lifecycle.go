package lifecycle

import (
	"os"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/graphql"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Lifecycle struct {
	GenerateTimeout     int
	ConfirmUploadPrompt bool
	DenyUploadPrompt    bool
	Quiet               bool
	GenerateBundlePath  string
	UploadCustomerID    string
	GraphQLClient       *graphql.Client
	FileInfo            os.FileInfo
	// RealGeneratedBundlePath is the actual path the bundle was written to.
	// it will be different from GenerateBundlePath if a URL or `-` (stdout)
	// was passed in.
	RealGeneratedBundlePath string
	BundleTasks             []types.Task
	executors               []Executor
}

// Executor is a thing that can be executed.
//                                                    ...and I still think Rob is wrong
type Executor interface {
	Execute(lifecycle *Lifecycle) (bool, error)
}

func (l *Lifecycle) Build(tasks []types.LifecycleTask) error {

	for _, task := range tasks {
		t := task
		executor, err := resolveTask(&t)
		if err != nil {
			return errors.Wrap(err, "resolve event")
		}
		l.executors = append(l.executors, executor)
	}

	return nil
}

func resolveTask(t *types.LifecycleTask) (Executor, error) {
	switch {
	case t.Message != nil:
		return &MessageTask{*t.Message}, nil
	case t.Generate != nil:
		return &GenerateTask{*t.Generate}, nil
	case t.Upload != nil:
		return &UploadTask{*t.Upload}, nil
	}

	return nil, errors.New("no valid event found, requires one of: generate, message, boolean, upload")
}

func (l *Lifecycle) Run() error {
	for _, t := range l.executors {
		cont, err := t.Execute(l)

		if err != nil {
			return errors.Wrap(err, "running task")
		}

		if !cont {
			break
		}
	}

	return nil
}
