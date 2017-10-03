package metrics

import (
	"context"
	"encoding/json"
	"fmt"
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
			Description: "`docker logs" + args[0] + "` command results",
			Filename:    filename,
			RawError:    rawError,
			JSONError:   jsonError,
			HumanError:  humanError,
		}
		completeCh <- true
	}()

	timeoutChan := make(chan error, 1)

	go func() {
		cli, err := client.NewEnvClient()
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			timeoutChan <- err
			return
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
			timeoutChan <- err
			return
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
			timeoutChan <- err
			return
		}

		// Send the json
		dataCh <- types.Data{
			Filename: filepath.Join("/json/", filename+".json"),
			Data:     j,
		}
		timeoutChan <- nil
	}()

	select {
	case err := <-timeoutChan:
		//completed on time
		return err
	case <-time.After(timeout):
		//failed to complete on time
		err := types.TimeoutError{Message: fmt.Sprintf(`Getting logs from host:%s timed out after %s`, args[0], timeout.String())}
		rawError = err
		jsonError = err
		humanError = err
		return err
	}
}
