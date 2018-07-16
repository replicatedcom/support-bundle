package types

import (
	"context"
)

type Task interface {
	Exec(ctx context.Context, rootDir string) []*Result
	GetSpec() Spec
}
