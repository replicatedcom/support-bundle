package systemutil

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerRunCommand Run a command on a specified docker instance.
// args: ["id", "user", "cmd", "arg1", "arg2"...]
func DockerRunCommand(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	filename := "/docker/runcommand/"

	r, _ := regexp.Compile(`[^\w]`)

	for _, arg := range args {
		filename += r.ReplaceAllString(arg, "_")
	}

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "dockerRunCommand",
			Description: "Results of running a command within a docker container",
			Filename:    filename,
			RawError:    rawError,
			JSONError:   jsonError,
			HumanError:  humanError,
		}
		completeCh <- true
	}()

	cli, err := client.NewEnvClient()
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	execOpts := dockertypes.ExecConfig{
		User:         args[1],
		Cmd:          args[2:],
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
	}

	execInstance, err := cli.ContainerExecCreate(context.Background(), args[0], execOpts)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	att, err := cli.ContainerExecAttach(context.Background(), execInstance.ID, execOpts)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	execStartOpts := dockertypes.ExecStartCheck{
		Detach: false,
		Tty:    false,
	}
	err = cli.ContainerExecStart(context.Background(), execInstance.ID, execStartOpts)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	// read everything
	response, err := ioutil.ReadAll(att.Reader)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	//close connection
	att.Close()

	// Send the raw
	dataCh <- types.Data{
		Filename: filepath.Join("/raw/", filename),
		Data:     response,
	}

	// Human readable version
	dataCh <- types.Data{
		Filename: filepath.Join("/human/", filename+".txt"),
		Data:     response,
	}

	type runCommandStruct struct {
		Output string `json:"output"`
	}
	u := runCommandStruct{
		Output: string(response),
	}
	j, err := json.Marshal(u)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		return err
	}

	// Send the json
	dataCh <- types.Data{
		Filename: filepath.Join("/json/", filename+".json"),
		Data:     j,
	}

	return nil
}
