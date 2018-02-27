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
)

func recursiveReadFile(ctx context.Context, filepath string, prefix string) (map[string]io.Reader, error) {
	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return nil, err
	}

	readers := make(map[string]io.Reader)
	for _, file := range files {
		if file.IsDir() {
			childReaders, err := recursiveReadFile(ctx, path.Join(filepath, file.Name()), path.Join(prefix, file.Name()))
			if err != nil {
				return readers, err
			}
			if childReaders != nil {
				for k, v := range childReaders {
					readers[k] = v
				}
			}
		} else {
			//get a reader for the file & add it to the map
			r, err := os.Open(path.Join(filepath, file.Name()))
			if err != nil {
				return nil, err
			}
			readers[path.Join(prefix, file.Name())] = r
		}
	}
	return readers, nil
}

func (c *Core) ReadFile(opts types.CoreReadFileOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		info, err := os.Stat(opts.Filepath)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			//this is a single file. We don't have to do anything special to handle it properly.
			r, err := os.Open(opts.Filepath)
			if err != nil {
				return nil, err
			}
			return map[string]io.Reader{info.Name(): r}, nil
		} else {
			//this is a directory. We need to loop through all files in the directory recursively and create io.Readers for each of them.
			return recursiveReadFile(ctx, opts.Filepath, info.Name())
		}
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
