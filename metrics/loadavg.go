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

type LoadAverage struct {
	minuteOne           float64
	minuteFive          float64
	minuteTen           float64
	processCountRunning int
	processCountTotal   int
}

func LoadAvg(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/system/metrics/loadavg"

	var err error

	var datas []types.Data
	var paths []string

	completeChan := make(chan error, 1)

	go func() {
		b, err := ioutil.ReadFile("/proc/loadavg")
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

		loadAverage, err := parseLoadAvg(b)
		if err != nil {
			completeChan <- err
			return
		}

		human := fmt.Sprintf("%f %f %f", loadAverage.minuteOne, loadAverage.minuteFive, loadAverage.minuteTen)
		// Convert to human readable
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename+".txt"),
			Data:     []byte(human),
		})
		paths = append(paths, filepath.Join("/human/", filename+".txt"))

		j, err := json.Marshal(loadAverage)
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
		err = types.TimeoutError{Message: fmt.Sprintf(`Fetching load averages failed due to: %s`, ctx.Err().Error())}
	}

	result := types.Result{
		Task:      "loadavg",
		Args:      args,
		Filenames: paths,
		Error:     err,
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
