package planners

import (
	"crypto/tls"
	"net/http"

	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/retraced/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	retraced "github.com/retracedhq/retraced-go"
)

type Retraced struct {
	producers *producers.Retraced
}

func New(retraced *producers.Retraced) *Retraced {
	return &Retraced{
		producers: retraced,
	}
}

func getClient(opts types.RetracedAPIClientOptions) (*retraced.Client, error) {
	client, err := retraced.NewClient("", opts.ProjectID, opts.APIToken)
	if err != nil {
		return nil, err
	}

	client.Endpoint = opts.APIEndpoint

	if opts.Insecure {
		client.HttpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	return client, nil
}
