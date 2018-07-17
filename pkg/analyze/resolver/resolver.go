package resolver

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
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

func (r *Resolver) ResolveSpec(ctx context.Context, files, inline []string) (resolved api.Spec, err error) {
	debug := level.Debug(log.With(r.Logger, "method", "Resolver.ResolveSpec"))

	cwd, _ := os.Getwd()
	for _, filename := range files {
		spec, errI := deserializeSpecFile(r.Fs, filename)
		debug.Log(
			"phase", "spec.file.deserialize",
			"cwd", cwd,
			"spec", filename,
			"error", errI)
		if errI != nil {
			err = multierror.Append(err, errors.Wrapf(errI, "deserialize spec file %s", filename))
		} else {
			resolved = mergeSpecs(resolved, spec)
		}
	}

	for i, inline := range inline {
		spec, errI := api.DeserializeSpec([]byte(inline))
		debug.Log(
			"phase", "spec.inline.deserialize",
			"spec", i,
			"error", errI)
		if errI != nil {
			err = multierror.Append(err, errors.Wrapf(errI, "deserialize inline spec %d", i))
		} else {
			resolved = mergeSpecs(resolved, spec)
		}
	}

	return
}

func deserializeSpecFile(fs afero.Fs, filename string) (api.Spec, error) {
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		return api.Spec{}, err
	}
	return api.DeserializeSpec(b)
}

func mergeSpecs(specs ...api.Spec) (merged api.Spec) {
	for _, spec := range specs {
		merged.Analyze.V1Alpha1 = append(merged.Analyze.V1Alpha1, spec.Analyze.V1Alpha1...)
		merged.Collect.V1 = append(merged.Collect.V1, spec.Collect.V1...)
	}
	return
}
