package producers

import (
	"context"
	"io"
	"io/ioutil"
	"log"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type RunCommandOptions struct {
	Image      string
	Cmd        []string
	Env        []string
	Binds      []string
	EnablePull bool
}

// RunCommand returns stdout and stderr results.
func (d *Docker) RunCommand(opts RunCommandOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		if _, _, err := d.client.ImageInspectWithRaw(ctx, opts.Image); err != nil {
			if !opts.EnablePull || !docker.IsErrNotFound(err) {
				return nil, errors.Wrapf(err, "image inspect %s", opts.Image)
			}
			resp, err := d.client.ImagePull(ctx, opts.Image, dockertypes.ImagePullOptions{
			// TODO: registry auth
			})
			if err != nil {
				return nil, errors.Wrapf(err, "image pull %s", opts.Image)
			}
			if _, err := ioutil.ReadAll(resp); err != nil {
				return nil, err
			}
			resp.Close()
		}

		createOpts := dockertypes.ContainerCreateConfig{
			Config: &container.Config{
				AttachStderr: true,
				AttachStdout: true,
				Cmd:          strslice.StrSlice(opts.Cmd),
				Image:        opts.Image,
				Env:          opts.Env,
			},
			HostConfig: &container.HostConfig{
				Binds: opts.Binds,
			},
		}
		containerInstance, err := d.client.ContainerCreate(ctx, createOpts.Config, createOpts.HostConfig, nil, "")
		if err != nil {
			return nil, errors.Wrap(err, "container create")
		}
		go func() {
			<-ctx.Done()
			d.client.ContainerRemove(context.Background(), containerInstance.ID, dockertypes.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			})
		}()

		resp, err := d.client.ContainerAttach(ctx, containerInstance.ID, dockertypes.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		})
		if err != nil {
			return nil, errors.Wrap(err, "container attach")
		}

		if err := d.client.ContainerStart(ctx, containerInstance.ID, dockertypes.ContainerStartOptions{}); err != nil {
			resp.Close()
			return nil, errors.Wrap(err, "container start")
		}

		stdoutR, stdoutW := io.Pipe()
		stderrR, stderrW := io.Pipe()

		go func() {
			// TODO context interruptible
			_, err = stdcopy.StdCopy(stdoutW, stderrW, resp.Reader)
			resp.Close()
			if err != nil {
				stdoutW.CloseWithError(err)
				stderrW.CloseWithError(err)
				return
			}
			if err := stdoutW.Close(); err != nil {
				log.Print(err)
			}
			if err := stderrW.Close(); err != nil {
				log.Print(err)
			}
		}()

		readers := make(map[string]io.Reader)
		readers["stdout"] = stdoutR
		readers["stderr"] = stderrR
		return readers, nil
	}
}

func RunCommandOptionsFromSpec(config types.Config) (RunCommandOptions, error) {
	if config.Image == "" {
		return RunCommandOptions{}, errors.New("spec for docker.run-command requires an image within config")
	}

	var fullCommand []string
	if config.Command != "" {
		fullCommand = append(fullCommand, config.Command)
	}
	fullCommand = append(fullCommand, config.Args...)

	return RunCommandOptions{
		Image:      config.Image,
		Cmd:        fullCommand,
		Env:        config.Env,
		Binds:      config.Binds,
		EnablePull: config.EnablePull,
	}, nil
}
