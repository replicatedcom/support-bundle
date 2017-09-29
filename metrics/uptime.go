package metrics

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/replicatedcom/support-bundle/types"
	"github.com/replicatedcom/support-bundle/util"

	jww "github.com/spf13/jwalterweatherman"
)

func Uptime(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	filename := "/system/metrics/uptime"

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "uptime",
			Description: "System Uptime",
			Filename:    filename,
			RawError:    rawError,
			JSONError:   jsonError,
			HumanError:  humanError,
		}
		completeCh <- true
	}()

	b, err := util.ReadFile("/proc/uptime")
	if err != nil {
		jww.ERROR.Print(err)
		rawError, jsonError, humanError = err, err, err
		return err
	}

	// Send the raw
	dataCh <- types.Data{
		Filename: filepath.Join("/raw/", filename),
		Data:     b,
	}

	uptimeSeconds, err := parseUptime(b)
	if err != nil {
		jsonError, humanError = err, err
		return err
	}

	human := fmt.Sprintf("Total Time (seconds): %f\nIdle Time (seconds): %f", uptimeSeconds[0], uptimeSeconds[1])
	// Convert to human readable
	dataCh <- types.Data{
		Filename: filepath.Join("/human/", filename),
		Data:     []byte(human),
	}

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
		return err
	}

	dataCh <- types.Data{
		Filename: filepath.Join("/json/", filename),
		Data:     j,
	}

	return nil
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
