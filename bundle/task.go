package bundle

import (
	"context"
	"time"

	"github.com/replicatedcom/support-bundle/types"
)

type Task struct {
	Description string
	ExecFunc    func(context.Context, []string) ([]types.Data, types.Result, error)
	Timeout     time.Duration
	Args        []string
}
