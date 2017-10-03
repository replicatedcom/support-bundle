package metrics

import (
	"context"
	"encoding/json"
	"fmt"
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

		containers, err := cli.ContainerList(context.Background(), dockertypes.ContainerListOptions{})
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			timeoutChan <- err
			return
		}

		containersJSON, err := json.Marshal(containers)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			timeoutChan <- err
			return
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
			timeoutChan <- err
			return
		}

		// Human readable version
		dataCh <- types.Data{
			Filename: filepath.Join("/human/", filename+".json"),
			Data:     containerIndentJSON,
		}
		timeoutChan <- nil
	}()

	select {
	case err := <-timeoutChan:
		//completed on time
		return err
	case <-time.After(timeout):
		//failed to complete on time
		err := types.TimeoutError{Message: fmt.Sprintf(`Docker ps timed out after %s`, timeout.String())}
		rawError = err
		jsonError = err
		humanError = err
		return err
	}
}
