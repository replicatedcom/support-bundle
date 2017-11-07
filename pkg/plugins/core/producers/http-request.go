package producers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
	"golang.org/x/net/context/ctxhttp"
)

func HTTPRequest(opts types.HTTPRequestCommandOptions) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {
		body := bytes.NewBuffer([]byte(opts.Body))
		if opts.Method != "" {
			jww.DEBUG.Printf("Making request to %s with method %s\n", opts.URL, opts.Method)
		} else {
			jww.DEBUG.Printf("Making request to %s\n", opts.URL)
		}
		req, err := http.NewRequest(strings.ToUpper(opts.Method), opts.URL, body)
		if err != nil {
			return nil, errors.Wrap(err, "new request")
		}

		resp, err := ctxhttp.Do(ctx, nil, req)
		if err != nil {
			return nil, errors.Wrap(err, "make request")
		}

		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("unexpected status %s", resp.Status)
		}
		return resp.Body, nil
	}
}
