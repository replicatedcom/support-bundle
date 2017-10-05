package systemutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"
)

func ReadFile(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	readFile := args[0]

	// make a sanatized version of the filename we're searching for - replace forward slash, backslash colon and space with _
	r, _ := regexp.Compile(`[^\w]`)
	filename := "/system/readfile/" + r.ReplaceAllString(readFile, "_")

	var rawError, jsonError, humanError error = nil, nil, nil

	var datas []types.Data

	completeChan := make(chan error, 1)

	go func() {
		b, err := ioutil.ReadFile(readFile)
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

		human := fmt.Sprintf("Read file %q: %q", readFile, b)
		// Convert to human readable
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename),
			Data:     []byte(human),
		})

		type readFileStruct struct {
			File string `json:"file"`
		}
		u := readFileStruct{
			File: string(b),
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

		completeChan <- nil
	}()

	var err error

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf("Reading file at %s errored out due to %s", args[0], ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Name:        "dockerps",
		Description: "`docker ps` command outputs",
		Filename:    filename,
		RawError:    rawError,
		JSONError:   jsonError,
		HumanError:  humanError,
	}

	return datas, result, err
}
