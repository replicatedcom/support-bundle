package systemutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// DockerRunCommand Run a command on a specified docker instance.
// args: ["id", "user", "cmd", "arg1", "arg2"...]
func DockerRunCommand(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/docker/runcommand/"
	commandString := ""
	r, _ := regexp.Compile(`[^\w]`)

	for _, arg := range args {
		commandString += arg + "_"
	}
	commandString = commandString[:len(commandString)-1]

	filename += r.ReplaceAllString(commandString, "_")

	var rawError, jsonError, humanError error = nil, nil, nil

	var datas []types.Data

	completeChan := make(chan error, 1)

	go func() {
		cli, err := client.NewEnvClient()
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		execOpts := dockertypes.ExecConfig{
			User:         args[1],
			Cmd:          args[2:],
			AttachStderr: true,
			AttachStdout: true,
			AttachStdin:  true,
		}

		execInstance, err := cli.ContainerExecCreate(ctx, args[0], execOpts)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		att, err := cli.ContainerExecAttach(ctx, execInstance.ID, execOpts)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		execStartOpts := dockertypes.ExecStartCheck{
			Detach: false,
			Tty:    false,
		}
		err = cli.ContainerExecStart(ctx, execInstance.ID, execStartOpts)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		var dstdout, dstderr bytes.Buffer

		//read and demultiplex
		_, err = stdcopy.StdCopy(&dstdout, &dstderr, att.Reader)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		// close connection
		att.Close()

		// get stdout and stderr byte arrays
		stdoutResult := dstdout.Bytes()
		stderrResult := dstderr.Bytes()

		// Send the raw result
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename+".out.txt"),
			Data:     stdoutResult,
		})
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename+".err.txt"),
			Data:     stderrResult,
		})

		// Human readable version
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename+".out.txt"),
			Data:     stdoutResult,
		})
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename+".err.txt"),
			Data:     stderrResult,
		})

		type runCommandStruct struct {
			Out string `json:"stdout"`
			Err string `json:"stderr"`
		}
		u := runCommandStruct{
			Out: string(stdoutResult),
			Err: string(stderrResult),
		}
		j, err := json.Marshal(u)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			completeChan <- err
			return
		}

		// Send the json
		datas = append(datas, types.Data{
			Filename: filepath.Join("/json/", filename+".json"),
			Data:     j,
		})
		completeChan <- nil
	}()

	var err error
	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf(`Command "%s" failed due to: %s`, commandString, ctx.Err().Error())}
		// err := errors.Wrap(ctx.Err(), `Command "`+commandString+`" failed`) //would be nice to use but doesn't convert to json
		rawError = err
		jsonError = err
		humanError = err
	}

	results := types.Result{
		Name:        "dockerRunCommand",
		Description: "Results of running a command within a docker container",
		Filename:    filename,
		RawError:    rawError,
		JSONError:   jsonError,
		HumanError:  humanError,
	}

	return datas, results, err
}
