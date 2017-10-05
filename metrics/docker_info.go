package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/docker/docker/client"
)

func DockerInfo(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/docker/metrics/info"

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

		info, err := cli.Info(ctx)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		infoJSON, err := json.Marshal(info)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			completeChan <- err
			return
		}

		// Send the raw
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename),
			Data:     infoJSON,
		})

		// Send the json
		datas = append(datas, types.Data{
			Filename: filepath.Join("/json/", filename+".json"),
			Data:     infoJSON,
		})

		infoIndentJSON, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			jww.ERROR.Print(err)
			humanError = err
			completeChan <- err
			return
		}

		// Human readable version
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename+".json"),
			Data:     infoIndentJSON,
		})
		completeChan <- nil
	}()

	var err error

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf(`Docker info failed due to: %s`, ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Name:        "dockerInfo",
		Description: "`docker info` command results",
		Filename:    filename,
		RawError:    rawError,
		JSONError:   jsonError,
		HumanError:  humanError,
	}

	return datas, result, err
}
