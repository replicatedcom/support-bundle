package producers

import (
	"context"
	"io"

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

func (r *Retraced) Events(client *retraced.Client, opts types.RetracedEventsOptions) types.StreamProducer {
	return func(ctx context.Context) (io.Reader, error) {

		query := getQuery(opts)
		mask := getMask(opts)
		reader, writer := io.Pipe()

		go func() {
			err := client.ExportCSV(ctx, writer, query, mask)
			if err != nil {
				jww.ERROR.Printf("Failed to collect retraced events: %v", err)
			}
			writer.CloseWithError(err)
		}()

		return reader, nil
	}
}

func getQuery(opts types.RetracedEventsOptions) *retraced.StructuredQuery {
	var query *retraced.StructuredQuery
	if opts.Query != nil {
		query = opts.Query
	} else {
		query = defaultQuery
	}
	return query
}

func getMask(opts types.RetracedEventsOptions) *retraced.EventNodeMask {
	var mask *retraced.EventNodeMask
	if opts.Mask != nil {
		mask = opts.Mask
	} else {
		mask = defaultMask
	}
	return mask
}
