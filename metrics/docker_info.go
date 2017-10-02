package metrics

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/docker/docker/client"
)

func DockerInfo(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	filename := "/docker/metrics/info"

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "dockerInfo",
			Description: "`docker info` command results",
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
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		jww.ERROR.Print(err)
	}

	infoJSON, err := json.Marshal(info)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		return err
	}

	// Send the raw
	dataCh <- types.Data{
		Filename: filepath.Join("/raw/", filename),
		Data:     infoJSON,
	}

	// Send the json
	dataCh <- types.Data{
		Filename: filepath.Join("/json/", filename+".json"),
		Data:     infoJSON,
	}

	infoIndentJSON, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		jww.ERROR.Print(err)
		humanError = err
		return err
	}

	// Human readable version
	dataCh <- types.Data{
		Filename: filepath.Join("/human/", filename+".json"),
		Data:     infoIndentJSON,
	}

	return nil
}
