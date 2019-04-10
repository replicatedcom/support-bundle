package resolver

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/collect/graphql"
	"github.com/spf13/afero"
)

type Resolver struct {
	Logger log.Logger
	Fs     afero.Fs
}

type Input struct {
	Files      []string
	Inline     []string
	CustomerID string
	ChannelID  string
	Endpoint   string
}

func New(logger log.Logger, fs afero.Fs) *Resolver {
	return &Resolver{
		Logger: logger,
		Fs:     fs,
	}
}

func (r *Resolver) ResolveSpec(ctx context.Context, input Input) (resolved api.Doc, err error) {
	debug := level.Debug(log.With(r.Logger, "method", "Resolver.ResolveSpec"))

	cwd, _ := os.Getwd()
	for _, filename := range input.Files {
		doc, errI := deserializeDocFile(r.Fs, filename)
		debug.Log(
			"phase", "doc.file.deserialize",
			"cwd", cwd,
			"doc", filename,
			"error", errI)
		if errI != nil {
			err = multierror.Append(err, errors.Wrapf(errI, "deserialize doc file %s", filename))
		} else {
			resolved = mergeDocs(resolved, doc)
		}
	}

	for i, inline := range input.Inline {
		doc, errI := api.DeserializeDoc([]byte(inline))
		debug.Log(
			"phase", "doc.inline.deserialize",
			"doc", i,
			"error", errI)
		if errI != nil {
			err = multierror.Append(err, errors.Wrapf(errI, "deserialize inline doc %d", i))
		} else {
			resolved = mergeDocs(resolved, doc)
		}
	}

	if input.CustomerID != "" {
		client := graphql.NewClient(input.Endpoint, http.DefaultClient)
		inline, errI := client.GetCustomerSpec(input.CustomerID)
		debug.Log(
			"phase", "doc.customer.retrieve",
			"error", errI)
		if errI != nil {
			err = multierror.Append(err, errors.Wrap(errI, "retrieve customer doc"))
		} else {
			doc, errI := api.DeserializeDoc([]byte(inline))
			debug.Log(
				"phase", "doc.customer.deserialize",
				"error", errI)
			if errI != nil {
				err = multierror.Append(err, errors.Wrap(errI, "deserialize customer doc"))
			} else {
				resolved = mergeDocs(resolved, doc)
			}
		}
	}

	if input.ChannelID != "" {
		client := graphql.NewClient(input.Endpoint, http.DefaultClient)
		inline, errI := client.GetChannelSpec(input.ChannelID)
		debug.Log(
			"phase", "doc.channel.retrieve",
			"error", errI)
		if errI != nil {
			err = multierror.Append(err, errors.Wrap(errI, "retrieve channel doc"))
		} else {
			doc, errI := api.DeserializeDoc([]byte(inline))
			debug.Log(
				"phase", "doc.channel.deserialize",
				"error", errI)
			if errI != nil {
				err = multierror.Append(err, errors.Wrap(errI, "deserialize channel doc"))
			} else {
				resolved = mergeDocs(resolved, doc)
			}
		}
	}

	return
}

func deserializeDocFile(fs afero.Fs, filename string) (api.Doc, error) {
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		return api.Doc{}, err
	}
	return api.DeserializeDoc(b)
}

func mergeDocs(docs ...api.Doc) (merged api.Doc) {
	for _, doc := range docs {
		merged.Analyze.V1 = append(merged.Analyze.V1, doc.Analyze.V1...)
	}
	return
}
