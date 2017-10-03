package metrics

import (
	"context"
	"encoding/json"
	"fmt"
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

		info, err := cli.Info(context.Background())
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			timeoutChan <- err
			return
		}

		infoJSON, err := json.Marshal(info)
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
			timeoutChan <- err
			return
		}

		// Human readable version
		dataCh <- types.Data{
			Filename: filepath.Join("/human/", filename+".json"),
			Data:     infoIndentJSON,
		}
		timeoutChan <- nil
	}()

	select {
	case err := <-timeoutChan:
		//completed on time
		return err
	case <-time.After(timeout):
		//failed to complete on time
		err := types.TimeoutError{Message: fmt.Sprintf(`Docker info timed out after %s`, timeout.String())}
		rawError = err
		jsonError = err
		humanError = err
		return err
	}
}
