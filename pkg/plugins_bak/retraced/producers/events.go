package producers

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
	"github.com/retracedhq/retraced-go"
	jww "github.com/spf13/jwalterweatherman"
)

var (
	defaultMask = &retraced.EventNodeMask{
		Action:        true,
		CRUD:          true,
		CanonicalTime: true,
	}
	defaultQuery = &retraced.StructuredQuery{
		CRUD: "c,u,d",
	}
)

func Events(spec types.Spec) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {

		client, err := getClient(spec)
		if err != nil {
			return nil, errors.Wrap(err, "create retraced client")
		}

		query := getQuery(spec)
		mask := getMask(spec)
		reader, writer := io.Pipe()

		go func() {
			var err error
			defer func() {
				writer.CloseWithError(err)
			}()
			err = client.ExportCSV(ctx, writer, query, mask)
			if err != nil {
				jww.ERROR.Printf("Failed to collect retraced events %v", err)
			}
		}()

		return reader, nil
	}
}
func getClient(spec types.Spec) (*retraced.Client, error) {
	client, err := retraced.NewClient(spec.RetracedEventsCommand.ProjectID, spec.RetracedEventsCommand.APIToken)
	if err != nil {
		return nil, err
	}

	client.Endpoint = spec.RetracedEventsCommand.APIEndpoint

	if spec.RetracedEventsCommand.Insecure {
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

func getQuery(spec types.Spec) *retraced.StructuredQuery {
	var query *retraced.StructuredQuery
	if spec.RetracedEventsCommand.Query != nil {
		query = spec.RetracedEventsCommand.Query
	} else {
		query = defaultQuery
	}
	return query
}

func getMask(spec types.Spec) *retraced.EventNodeMask {
	var mask *retraced.EventNodeMask
	if spec.RetracedEventsCommand.Mask != nil {
		mask = spec.RetracedEventsCommand.Mask
	} else {
		mask = defaultMask
	}
	return mask
}
