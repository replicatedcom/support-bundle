package plans

import (
	"context"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/types"
)

// Prepared is a Task that returns preconfigured Results.
type Prepared struct {
	results []*types.Result
}

func (p Prepared) Exec(ctx context.Context, rootDir string) []*types.Result {
	for _, r := range p.results {
		if r.Path != "" {
			r.Path = filepath.Join(rootDir, r.Path)
		}
	}

	return p.results
}

// Prepare results for an incomplete spec
func PreparedError(err error, spec types.Spec) Prepared {
	results := []*types.Result{}

	if spec.Raw != "" {
		results = append(results, &types.Result{
			Error: err,
		})
	}
	if spec.JSON != "" {
		results = append(results, &types.Result{
			Error: err,
		})
	}
	if spec.Human != "" {
		results = append(results, &types.Result{
			Error: err,
		})
	}

	return Prepared{
		results: results,
	}
}
