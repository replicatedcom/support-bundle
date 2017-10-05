package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"
)

func Hostname(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/system/metrics/hostname"

	var rawError, jsonError, humanError error = nil, nil, nil

	var datas []types.Data

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

		human := fmt.Sprintf("Hostname: %q", b)
		// Convert to human readable
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename),
			Data:     []byte(human),
		})

		type hostname struct {
			Hostname string `json:"hostname"`
		}
		u := hostname{
			Hostname: string(b),
		}
		j, err := json.Marshal(u)
		if err != nil {
			jww.ERROR.Print(err)
			jsonError = err
			completeChan <- err
			return
		}

		datas = append(datas, types.Data{
			Filename: filepath.Join("/json/", filename),
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
		err = types.TimeoutError{Message: fmt.Sprintf(`Fetching hostname failed due to: %s`, ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Name:        "hostname",
		Description: "System Hostname",
		Filename:    filename,
		RawError:    rawError,
		JSONError:   jsonError,
		HumanError:  humanError,
	}

	return datas, result, err
}
