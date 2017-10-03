package systemutil

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"
)

func RunCommand(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	command := args[0]
	arg := args[1]

	// make a sanatized version of the filename we're searching for - replace forward slash, backslash colon and space with _
	r, _ := regexp.Compile(`[^\w]`)
	filename := "/system/runcommand/" + r.ReplaceAllString(command, "_") + "_" + r.ReplaceAllString(arg, "_")

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

	timeoutChan := make(chan error, 1)

	go func() {
		b, err := exec.Command(command, arg).Output()
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

		human := fmt.Sprintf("Run command %q: %q", command, b)
		// Convert to human readable
		dataCh <- types.Data{
			Filename: filepath.Join("/human/", filename),
			Data:     []byte(human),
		}

		type runCommandStruct struct {
			Output string `json:"output"`
		}
		u := runCommandStruct{
			Output: string(b),
		}
		j, err := json.Marshal(u)
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
		err := types.TimeoutError{Message: fmt.Sprintf(`Command "%s" timed out after %s`, command+"_"+arg, timeout.String())}
		rawError = err
		jsonError = err
		humanError = err
		return err
	}
}
