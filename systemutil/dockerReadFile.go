package systemutil

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/docker/docker/client"
)

// DockerReadFile Read a file from a docker instance
// args: ["id", "path_to_file"]
func DockerReadFile(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/docker/readfile/"

	r, _ := regexp.Compile(`[^\w]`)

	for _, arg := range args {
		filename += r.ReplaceAllString(arg, "_")
	}

	var rawError, jsonError, humanError error = nil, nil, nil

	var datas []types.Data
	var paths []string

	completeChan := make(chan error, 1)

	go func() {
		cli, err := client.NewEnvClient()
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		readcloser, _, err := cli.CopyFromContainer(ctx, args[0], args[1])
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		// read everything
		response, err := ioutil.ReadAll(readcloser)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			completeChan <- err
			return
		}

		//close connection
		readcloser.Close()

		// Send the raw
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename+".tar"),
			Data:     response,
		})
		paths = append(paths, filepath.Join("/raw/", filename+".tar"))

		completeChan <- nil
	}()

	var err error

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf(`Reading a docker file at from host:%s and path:%s errored out with %s`, args[0], args[1], ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Task:       "dockerReadFile",
		Args:       args,
		Filenames:  paths,
		RawError:   rawError,
		JSONError:  jsonError,
		HumanError: humanError,
	}

	return datas, result, err
}
