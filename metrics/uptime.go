package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"
)

func Uptime(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/system/metrics/uptime"

	var err error

	var datas []types.Data
	var paths []string

	completeChan := make(chan error, 1)

	go func() {
		b, err := ioutil.ReadFile("/proc/uptime")
		if err != nil {
			jww.ERROR.Print(err)
			completeChan <- err
			return
		}

		// Send the raw
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename),
			Data:     b,
		})
		paths = append(paths, filepath.Join("/raw/", filename))

		uptimeSeconds, err := parseUptime(b)
		if err != nil {
			completeChan <- err
			return
		}

		human := fmt.Sprintf("Total Time (seconds): %f\nIdle Time (seconds): %f", uptimeSeconds[0], uptimeSeconds[1])
		// Convert to human readable
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename+".txt"),
			Data:     []byte(human),
		})
		paths = append(paths, filepath.Join("/human/", filename+".txt"))

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
			completeChan <- err
			return
		}

		datas = append(datas, types.Data{
			Filename: filepath.Join("/json/", filename+".json"),
			Data:     j,
		})
		paths = append(paths, filepath.Join("/json/", filename+".json"))

		completeChan <- nil
	}()

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf(`Fetching uptime failed due to: %s`, ctx.Err().Error())}
	}

	result := types.Result{
		Task:      "uptime",
		Args:      args,
		Filenames: paths,
		Error:     err,
	}

	return datas, result, err
}

func parseUptime(contents []byte) ([]float64, error) {

	// # cat /proc/uptime
	// 33524.72 66785.42

	parts := strings.Split(strings.TrimSpace(string(contents)), " ")
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
