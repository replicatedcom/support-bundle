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

func Dockerps(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	filename := "/docker/metrics/ps"

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "dockerps",
			Description: "`docker ps` command outputs",
			Filename:    filename,
			RawError:    rawError,
			JSONError:   jsonError,
			HumanError:  humanError,
		}
		completeCh <- true
	}()

	b, err := exec.Command("docker", "ps").Output()
	if err != nil {
		log.Fatal(err)
	}

	// Send the raw
	dataCh <- types.Data{
		Filename: filepath.Join("/raw/", filename),
		Data:     b,
	}

	human := fmt.Sprintf("Docker ps: %q", b)
	// Convert to human readable
	dataCh <- types.Data{
		Filename: filepath.Join("/human/", filename),
		Data:     []byte(human),
	}

	type dockerPS struct {
		Result string `json:"result"`
	}
	u := dockerPS{
		Result: string(b),
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
