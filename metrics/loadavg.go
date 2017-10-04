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

type LoadAverage struct {
	minuteOne           float64
	minuteFive          float64
	minuteTen           float64
	processCountRunning int
	processCountTotal   int
}

func LoadAvg(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/system/metrics/loadavg"

	var rawError, jsonError, humanError error = nil, nil, nil

	var datas []types.Data

	completeChan := make(chan error, 1)

	go func() {
		b, err := util.ReadFile("/proc/loadavg")
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

		loadAverage, err := parseLoadAvg(b)
		if err != nil {
			jsonError, humanError = err, err
			completeChan <- err
			return
		}

		human := fmt.Sprintf("%f %f %f", loadAverage.minuteOne, loadAverage.minuteFive, loadAverage.minuteTen)
		// Convert to human readable
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename),
			Data:     []byte(human),
		})

		j, err := json.Marshal(loadAverage)
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
		err = types.TimeoutError{Message: fmt.Sprintf(`Fetching load averages failed due to: %s`, ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Name:        "loadavg",
		Description: "System Load Average",
		Filename:    filename,
		RawError:    rawError,
		JSONError:   jsonError,
		HumanError:  humanError,
	}

	return datas, result, err
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
