package bundle

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/stretchr/testify/assert"
)

type taskStub struct {
	elapse   time.Duration
	results  []*types.Result
	deferred bool
}

func (t taskStub) Exec(ctx context.Context, rootDir string) []*types.Result {
	time.Sleep(t.elapse)
	return t.results
}

func (t taskStub) GetSpec() types.Spec {
	return types.Spec{
		SpecShared: types.SpecShared{Defer: t.deferred},
	}
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
	deferredResults := taskStub{
		elapse: time.Nanosecond,
		results: []*types.Result{
			{
				Spec: types.Spec{SpecShared: types.SpecShared{
					Description: "Support bundle goroutine dump",
					Defer:       true,
				}},
				Path: "/support-bundle/goroutines",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	results := Exec(ctx, "/dir", []types.Task{
		nilResults,
		noResults,
		singleResults,
		mixedResults,
		slowResults,
		deferredResults,
	}, func(opts *ExecOptions) {
		opts.DeferredTimeout = 100 * time.Millisecond
	})

	assert.Len(t, results, 6)
	assert.Contains(t, results, singleResults.results[0])
	assert.Contains(t, results, mixedResults.results[0])
	assert.Contains(t, results, mixedResults.results[1])
	assert.Contains(t, results, mixedResults.results[2])
	assert.NotContains(t, results, slowResults.results[0])
	assert.Contains(t, results, deferredResults.results[0])
}
