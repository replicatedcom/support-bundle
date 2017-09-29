package fileutil

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/replicatedcom/support-bundle/types"
	"github.com/replicatedcom/support-bundle/util"

	jww "github.com/spf13/jwalterweatherman"
)

func ReadFile(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	readFile := args[0]

	// make a sanatized version of the filename we're searching for - replace forward slash, backslash colon and space with _
	r, _ := regexp.Compile(`[\\\/:\s]`)
	filename := "/system/readfile/" + r.ReplaceAllString(readFile, "_")

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

	b, err := util.ReadFile(readFile)
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

	human := fmt.Sprintf("Read file %q: %q", readFile, b)
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
