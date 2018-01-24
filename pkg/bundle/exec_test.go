package bundle

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
	"github.com/stretchr/testify/assert"
)

type taskStub struct {
	elapse  time.Duration
	results []*types.Result
}

func (t taskStub) Exec(ctx context.Context, rootDir string) []*types.Result {
	time.Sleep(t.elapse)
	return t.results
}

func TestExec(t *testing.T) {
	nilResults := taskStub{
		elapse:  time.Nanosecond,
		results: nil,
	}
	noResults := taskStub{
		elapse:  time.Nanosecond,
		results: []*types.Result{},
	}
	singleResults := taskStub{
		elapse: time.Nanosecond,
		results: []*types.Result{
			{
				Spec: types.Spec{SpecShared: types.SpecShared{Description: "Logs from db container"}},
				Path: "/docker/db.logs",
			},
		},
	}
	mixedResults := taskStub{
		elapse: time.Nanosecond,
		results: []*types.Result{
			{
				Spec: types.Spec{SpecShared: types.SpecShared{Description: "Stderr from app container"}},
				Path: "/docker/app/stderr.txt",
			},
			{
				Spec:  types.Spec{SpecShared: types.SpecShared{Description: "Stdout from app container"}},
				Error: errors.New("Docker API 500"),
			},
			{
				Spec:  types.Spec{SpecShared: types.SpecShared{Description: "Logs from app container"}},
				Path:  "/docker/app.logs",
				Error: errors.New("Timedout"),
			},
		},
	}
	slowResults := taskStub{
		elapse: time.Second,
		results: []*types.Result{
			{
				Spec: types.Spec{SpecShared: types.SpecShared{Description: "/usr/bin/free"}},
				Path: "/host/commands/free",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	results := Exec(ctx, "/dir", []types.Task{
		nilResults,
		noResults,
		singleResults,
		mixedResults,
		slowResults,
	})

	assert.Len(t, results, 4)
	assert.Contains(t, results, singleResults.results[0])
	assert.Contains(t, results, mixedResults.results[0])
	assert.Contains(t, results, mixedResults.results[1])
	assert.Contains(t, results, mixedResults.results[2])
	assert.NotContains(t, results, slowResults.results[0])
}
