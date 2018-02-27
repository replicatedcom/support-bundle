package producers

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	dockerutil "github.com/replicatedcom/support-bundle/pkg/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func recursiveReadFile(ctx context.Context, filepath string, prefix string) (map[string]io.Reader, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	readers := make(map[string]io.Reader)
	if !info.IsDir() {
		//get a reader for the file & add it to the map
		r, err := os.Open(filepath)
		if err != nil {
			return nil, err
		}
		readers[path.Join(prefix, info.Name())] = r
		return readers, nil
	}

	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		childReaders, err := recursiveReadFile(ctx, path.Join(filepath, file.Name()), path.Join(prefix, file.Name()))
		if err != nil {
			return readers, errors.Wrapf(err, "Reading directory %s: ", filepath)
		}
		if childReaders != nil {
			for k, v := range childReaders {
				if _, ok := readers[k]; ok {
					//name collisions like this shouldn't happen
					jww.DEBUG.Printf("Filename collision at %s when reading directory %s", k, filepath)
				} else {
					readers[k] = v
				}

			}
		}
	}
	return readers, nil
}

func (c *Core) ReadFile(opts types.CoreReadFileOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		return recursiveReadFile(ctx, opts.Filepath, "")
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
		return map[string]io.Reader{"": r}, nil
	}
}
