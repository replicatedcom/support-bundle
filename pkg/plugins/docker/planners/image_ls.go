package planners

import (
	"errors"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/pkg/plans"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func (d *Docker) ImageLs(spec types.Spec) []types.Task {
	if spec.DockerImageLs == nil {
		if spec.DockerImages != nil {
			jww.DEBUG.Println("spec for docker.image-ls aliased from docker.images")
			spec.DockerImageLs = spec.DockerImages
		} else {
			err := errors.New("spec for docker.image-ls options required")
			task := plans.PreparedError(err, spec)
			return []types.Task{task}
		}
	}

	task := plans.StructuredSource{
		Producer:  d.producers.ImageLs(*spec.DockerImageLs),
		RawPath:   filepath.Join(spec.OutputDir, "image_ls.raw"),
		JSONPath:  filepath.Join(spec.OutputDir, "image_ls.json"),
		HumanPath: filepath.Join(spec.OutputDir, "image_ls.human"),
	}
	var err error
	task, err = plans.SetCommonFieldsStructuredSource(task, spec)
	if err != nil {
		task := plans.PreparedError(err, spec)
		return []types.Task{task}
	}
	return []types.Task{&task}
}
