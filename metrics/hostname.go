package metrics

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/types"
)

func Hostname(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/system/metrics/hostname"

	var rawError, jsonError, humanError error = nil, nil, nil

	var datas []types.Data
	var paths []string

	completeChan := make(chan error, 1)

	go func() {
		b, err := exec.Command("hostname").Output()
		if err != nil {
			log.Fatal(err)
		}

		// Send the raw
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename),
			Data:     b,
		})
		paths = append(paths, filepath.Join("/raw/", filename))

		completeChan <- nil
	}()

	var err error

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf(`Fetching hostname failed due to: %s`, ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Task:       "hostname",
		Args:       args,
		Filenames:  paths,
		RawError:   rawError,
		JSONError:  jsonError,
		HumanError: humanError,
	}

	return datas, result, err
}
