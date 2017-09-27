package metrics

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"
)

func Hostname(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration) error {
	filename := "/system/metrics/hostname"

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "hostname",
			Description: "System Hostname",
			Filename:    filename,
			RawError:    rawError,
			JSONError:   jsonError,
			HumanError:  humanError,
		}
		completeCh <- true
	}()

	b, err := exec.Command("hostname").Output()
	if err != nil {
		log.Fatal(err)
	}

	// Send the raw
	dataCh <- types.Data{
		Filename: filepath.Join("/raw/", filename),
		Data:     b,
	}

	human := fmt.Sprintf("Hostname: %q", b)
	// Convert to human readable
	dataCh <- types.Data{
		Filename: filepath.Join("/human/", filename),
		Data:     []byte(human),
	}

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
		return err
	}

	dataCh <- types.Data{
		Filename: filepath.Join("/json/", filename),
		Data:     j,
	}

	return nil
}
