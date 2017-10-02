package metrics

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Dockerps(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	filename := "/docker/metrics/ps"

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "dockerps",
			Description: "`docker ps` command outputs",
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

	containers, err := cli.ContainerList(context.Background(), dockertypes.ContainerListOptions{})
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	containersJSON, err := json.Marshal(containers)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		return err
	}

	// Send the raw
	dataCh <- types.Data{
		Filename: filepath.Join("/raw/", filename),
		Data:     containersJSON,
	}

	// Send the json
	dataCh <- types.Data{
		Filename: filepath.Join("/json/", filename+".json"),
		Data:     containersJSON,
	}

	containerIndentJSON, err := json.MarshalIndent(containers, "", "  ")
	if err != nil {
		jww.ERROR.Print(err)
		humanError = err
		return err
	}

	// Human readable version
	dataCh <- types.Data{
		Filename: filepath.Join("/human/", filename+".json"),
		Data:     containerIndentJSON,
	}

	return nil
}
