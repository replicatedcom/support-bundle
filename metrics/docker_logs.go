package metrics

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerLogs - reads the logs of a given docker container
// args[0] should be container id
func DockerLogs(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	filename := "/docker/logs/" + args[0]

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "dockerLogs",
			Description: "`docker logs CONTAINER` command results",
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

	options := dockertypes.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	logsReader, err := cli.ContainerLogs(context.Background(), args[0], options)
	logs, err := ioutil.ReadAll(logsReader)

	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	// Send the raw
	dataCh <- types.Data{
		Filename: filepath.Join("/raw/", filename),
		Data:     logs,
	}

	// Human readable version (sorta, same as raw)
	dataCh <- types.Data{
		Filename: filepath.Join("/human/", filename),
		Data:     logs,
	}

	type dockerLogsStruct struct {
		Logs      string `json:"logs"`
		Container string `json:"container"`
	}
	u := dockerLogsStruct{
		Logs:      string(logs),
		Container: args[0],
	}
	j, err := json.Marshal(u)
	if err != nil {
		jww.ERROR.Print(err)
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
