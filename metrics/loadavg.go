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

type LoadAverage struct {
	minuteOne           float64
	minuteFive          float64
	minuteTen           float64
	processCountRunning int
	processCountTotal   int
}

func LoadAvg(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	filename := "/system/metrics/loadavg"

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "loadavg",
			Description: "System Load Average",
			Filename:    filename,
			RawError:    rawError,
			JSONError:   jsonError,
			HumanError:  humanError,
		}
		completeCh <- true
	}()

	timeoutChan := make(chan error, 1)

	go func() {
		b, err := util.ReadFile("/proc/loadavg")
		if err != nil {
			jww.ERROR.Print(err)
			rawError, jsonError, humanError = err, err, err
			timeoutChan <- err
			return
		}

		// Send the raw
		dataCh <- types.Data{
			Filename: filepath.Join("/raw/", filename),
			Data:     b,
		}

		loadAverage, err := parseLoadAvg(b)
		if err != nil {
			jsonError, humanError = err, err
			timeoutChan <- err
			return
		}

		human := fmt.Sprintf("%f %f %f", loadAverage.minuteOne, loadAverage.minuteFive, loadAverage.minuteTen)
		// Convert to human readable
		dataCh <- types.Data{
			Filename: filepath.Join("/human/", filename),
			Data:     []byte(human),
		}

		j, err := json.Marshal(loadAverage)
		if err != nil {
			jww.ERROR.Print(err)
			jsonError = err
			timeoutChan <- err
			return
		}

		dataCh <- types.Data{
			Filename: filepath.Join("/json/", filename),
			Data:     j,
		}
		timeoutChan <- nil
	}()

	select {
	case err := <-timeoutChan:
		//completed on time
		return err
	case <-time.After(timeout):
		//failed to complete on time
		err := types.TimeoutError{Message: fmt.Sprintf(`Fetching load averages timed out after %s`, timeout.String())}
		rawError = err
		jsonError = err
		humanError = err
		return err
	}
}

func parseLoadAvg(contents []byte) (*LoadAverage, error) {

	// # cat /proc/loadavg
	// 0.02 0.01 0.00 4/229 5

	parts := strings.Split(string(contents), " ")
	if len(parts) != 5 {
		err := fmt.Errorf("Expected 5 values in loadavg but found %d", len(parts))
		jww.ERROR.Print(err)
		return nil, err
	}

	oneMin, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}
	fiveMin, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}
	tenMin, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		jww.ERROR.Print(err)
		return nil, err
	}

	return &LoadAverage{
		minuteOne:  oneMin,
		minuteFive: fiveMin,
		minuteTen:  tenMin,
	}, nil
}
