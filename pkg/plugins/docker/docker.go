package docker

import (
	"github.com/docker/docker/client"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/docker/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Docker struct {
	planner *planners.Docker
}

func New() (*Docker, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	producers := producers.New(c)
	return &Docker{
		planner: planners.New(producers, c),
	}, nil
}

func (p *Docker) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.DockerContainerLs != nil:
		return p.planner.ContainerLs
	case spec.DockerImageLs != nil:
		return p.planner.ImageLs
	case spec.DockerInfo != nil:
		return p.planner.Info
	case spec.DockerVersion != nil:
		return p.planner.Version

	// DockerPs               *DockerContainerLsOptions      `json:"docker.ps,omitempty"`     // canonical: docker.container-ls
	// DockerImages           *DockerImageLsOptions          `json:"docker.images,omitempty"` // canonical: docker.image-ls
	// DockerImageLs          *DockerImageLsOptions          `json:"docker.image-ls,omitempty"`
	// DockerContainerLs      *DockerContainerLsOptions      `json:"docker.container-ls,omitempty"`

	default:
		return nil
	}
}
