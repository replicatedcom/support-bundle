package util

import (
	"context"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
)

// MultiErrorGroup is a collection of goroutines working on subtasks that are
// part of the same overall task.
type MultiErrorGroup struct {
	cancel func()

	wg sync.WaitGroup

	err   *multierror.Error
	errMu sync.Mutex
}

// MultiErrorGroupWithContext returns a new Group and an associated Context
// derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func MultiErrorGroupWithContext(ctx context.Context) (*MultiErrorGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &MultiErrorGroup{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the all non-nil errors (if any) as a MultiError.
func (g *MultiErrorGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}

	g.errMu.Lock()
	defer g.errMu.Unlock()
	return g.err.ErrorOrNil()
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will
// be returned by Wait.
func (g *MultiErrorGroup) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errMu.Lock()
			g.err = multierror.Append(g.err, err)
			g.errMu.Unlock()

			if g.cancel != nil {
				g.cancel()
			}
		}
	}()
}
