package systemutil

import (
	"context"
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

	var err error = nil

	var datas []types.Data
	var paths []string

	completeChan := make(chan error, 1)

	go func() {
		b, err := ioutil.ReadFile(readFile)
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

		completeChan <- nil
	}()

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf("Reading file at %s errored out due to %s", args[0], ctx.Err().Error())}
	}

	result := types.Result{
		Task:      "readFile",
		Args:      args,
		Filenames: paths,
		Error:     err,
	}

	return datas, result, err
}
