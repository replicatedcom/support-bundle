package plans

import (
	"io"

	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

// add an error to every result, returning the results argument
// skips results that already have an error
func resultsWithErr(err error, results []*types.Result) []*types.Result {
	for _, r := range results {
		if r.Error == nil {
			r.Error = err
		}
	}

	return results
}

// closseLogErr
func closeLogErr(c io.Closer) {
	if err := c.Close(); err != nil {
		jww.ERROR.Print(err)
	}
}
