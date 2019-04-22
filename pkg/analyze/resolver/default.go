package resolver

import (
	"context"
	"sort"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/v1/defaultspec"
)

func (r *Resolver) ResolveDefaultSpec(ctx context.Context) (resolved api.Doc, err error) {
	debug := level.Debug(log.With(r.Logger, "method", "Resolver.ResolveDefaultSpec"))

	for _, name := range sort.StringSlice(defaultspec.AssetNames()) {
		b, errI := defaultspec.Asset(name)
		debug.Log(
			"phase", "asset.load",
			"name", name,
			"error", errI)
		if err != nil {
			err = multierror.Append(err, errors.Wrapf(errI, "load asset %s", name))
			continue
		}

		doc, errI := api.DeserializeDoc(b)
		debug.Log(
			"phase", "asset.deserialize",
			"name", name,
			"error", errI)
		if errI != nil {
			err = multierror.Append(err, errors.Wrapf(errI, "deserialize asset %s", name))
			continue
		}

		resolved = mergeDocs(resolved, doc)
	}

	return
}
