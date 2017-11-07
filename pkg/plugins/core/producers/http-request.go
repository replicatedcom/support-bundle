package producers

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/docker/go-connections/sockets"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
	"golang.org/x/net/context/ctxhttp"
)

var insecureClient, secureClient *http.Client

func init() {
	trSecure := new(http.Transport)
	sockets.ConfigureTransport(trSecure, "tcp", "")
	// TODO: possibly use go-connections/tlsconfig
	trSecure.TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	secureClient = &http.Client{
		Transport: trSecure,
	}

	trInsecure := new(http.Transport)
	sockets.ConfigureTransport(trInsecure, "tcp", "")
	// TODO: possibly use go-connections/tlsconfig
	trInsecure.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	insecureClient = &http.Client{
		Transport: trInsecure,
	}
}

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

		var client *http.Client
		if opts.Insecure {
			client = insecureClient
		} else {
			client = secureClient
		}

		resp, err := ctxhttp.Do(ctx, client, req)
		if err != nil {
			return nil, errors.Wrap(err, "make request")
		}

		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("unexpected status %s", resp.Status)
		}
		return resp.Body, nil
	}
}
