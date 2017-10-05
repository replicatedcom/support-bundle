package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerInspect - inspects a given docker container
// args[0] should be container id
func DockerInspect(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/docker/inspect/" + args[0]

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

		containerJSON, err := cli.ContainerInspect(ctx, args[0])
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		dataBytes, err := json.Marshal(containerJSON)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		// Send the raw (not really raw, since the source is json)
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename+".json"),
			Data:     dataBytes,
		})

		// Human readable version (not really raw, but )
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename+".json"),
			Data:     dataBytes,
		})

		type dockerLogsStruct struct {
			ContainerJSON dockertypes.ContainerJSON `json:"containerJSON"`
			Container     string                    `json:"container"`
		}
		u := dockerLogsStruct{
			ContainerJSON: containerJSON,
			Container:     args[0],
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
