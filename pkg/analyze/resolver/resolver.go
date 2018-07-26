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

func New(logger log.Logger, fs afero.Fs) *Resolver {
	return &Resolver{
		Logger: logger,
		Fs:     fs,
	}
}

func (r *Resolver) ResolveSpec(
	ctx context.Context,
	files []string,
	inline []string,
	customerID string,
	customerEndpoint string,
) (resolved api.Doc, err error) {

	debug := level.Debug(log.With(r.Logger, "method", "Resolver.ResolveSpec"))

	cwd, _ := os.Getwd()
	for _, filename := range files {
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

	for i, inline := range inline {
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

	if customerID != "" {
		client := graphql.NewClient(customerEndpoint, http.DefaultClient)
		inline, errI := client.GetCustomerSpec(customerID)
		debug.Log(
			"phase", "doc.customer.retrieve",
			"error", err)
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
