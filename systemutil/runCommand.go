package systemutil

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"
)

func RunCommand(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	command := args[0]
	arg := args[1]

	// make a sanatized version of the filename we're searching for - replace forward slash, backslash colon and space with _
	r, _ := regexp.Compile(`[^\w]`)
	filename := "/system/runcommand/" + r.ReplaceAllString(command, "_") + "_" + r.ReplaceAllString(arg, "_")

	var err error

	var datas []types.Data
	var paths []string

	completeChan := make(chan error, 1)

	go func() {
		b, err := exec.Command(command, arg).Output()
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
		err = types.TimeoutError{Message: fmt.Sprintf(`Command "%s" errored out with %s`, command+"_"+arg, ctx.Err().Error())}
	}

	result := types.Result{
		Task:      "runCommand",
		Args:      args,
		Filenames: paths,
		Error:     err,
	}

	return datas, result, err
}
