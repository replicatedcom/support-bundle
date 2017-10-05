package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerLogs - reads the logs of a given docker container
// args[0] should be container id
func DockerLogs(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/docker/logs/" + args[0]

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

		options := dockertypes.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		}
		logsReader, err := cli.ContainerLogs(ctx, args[0], options)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		logs, err := ioutil.ReadAll(logsReader)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		// Send the raw
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename),
			Data:     logs,
		})

		// Human readable version (sorta, same as raw)
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename),
			Data:     logs,
		})

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
		err = types.TimeoutError{Message: fmt.Sprintf(`Getting logs from host:%s failed due to: %s`, args[0], ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Name:        "dockerLogs",
		Description: "`docker logs" + args[0] + "` command results",
		Filename:    filename,
		RawError:    rawError,
		JSONError:   jsonError,
		HumanError:  humanError,
	}

	return datas, result, err
}
