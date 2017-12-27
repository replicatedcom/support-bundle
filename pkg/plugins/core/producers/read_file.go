package producers

import (
	"context"
	"io"
	"os"

	"github.com/pkg/errors"
	dockerutil "github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func (c *Core) ReadFile(opts types.CoreReadFileOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		r, err := os.Open(opts.Filepath)
		if err != nil {
			return nil, err
		}
		return map[string]io.Reader{"contents": r}, nil
	}
}

func (c *Core) DockerReadFile(opts types.CoreReadFileOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		container, err := dockerutil.ThisContainer(ctx, c.dockerClient)
		if err != nil {
			return nil, errors.Wrap(err, "this container")
		}

		r, err := dockerutil.ReadFile(ctx, c.dockerClient, container.Image, opts.Filepath)
		if err != nil {
			return nil, errors.Wrap(err, "docker read file")
		}
		return map[string]io.Reader{"contents": r}, nil
	}
}
