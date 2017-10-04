package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/replicatedcom/support-bundle/types"
	"github.com/replicatedcom/support-bundle/util"

	jww "github.com/spf13/jwalterweatherman"
)

func Uptime(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/system/metrics/uptime"

	var rawError, jsonError, humanError error = nil, nil, nil

	var datas []types.Data

	completeChan := make(chan error, 1)

	go func() {
		b, err := util.ReadFile("/proc/uptime")
		if err != nil {
			jww.ERROR.Print(err)
			rawError, jsonError, humanError = err, err, err
			completeChan <- err
			return
		}

		// Send the raw
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename),
			Data:     b,
		})

		uptimeSeconds, err := parseUptime(b)
		if err != nil {
			jsonError, humanError = err, err
			completeChan <- err
			return
		}

		human := fmt.Sprintf("Total Time (seconds): %f\nIdle Time (seconds): %f", uptimeSeconds[0], uptimeSeconds[1])
		// Convert to human readable
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename),
			Data:     []byte(human),
		})

		type uptime struct {
			TotalSeconds float64 `json:"total_seconds"`
			IdleSeconds  float64 `json:"idle_seconds"`
		}
		u := uptime{
			TotalSeconds: uptimeSeconds[0],
			IdleSeconds:  uptimeSeconds[1],
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
	}()

	var err error

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf(`Fetching uptime failed due to: %s`, ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Name:        "uptime",
		Description: "System Uptime",
		Filename:    filename,
		RawError:    rawError,
		JSONError:   jsonError,
		HumanError:  humanError,
	}

	return datas, result, err
}

func parseUptime(contents []byte) ([]float64, error) {

	// # cat /proc/uptime
	// 33524.72 66785.42

	parts := strings.Split(string(contents), " ")
	if len(parts) != 2 {
		err := fmt.Errorf("Expected 2 values in uptime but found %d", len(parts))
		jww.ERROR.Print(err)
		return nil, err
	}

	totalSeconds, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}
	idleSeconds, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}

	return []float64{
		totalSeconds,
		idleSeconds,
	}, nil
}
