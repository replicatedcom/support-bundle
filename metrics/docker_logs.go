package metrics

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerLogs - reads the logs of a given docker container
// args[0] should be container id
func DockerLogs(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/docker/logs/" + args[0]

	var err error = nil

	var datas []types.Data
	var paths []string

	completeChan := make(chan error, 1)

	go func() {
		cli, err := client.NewEnvClient()
		if err != nil {
			jww.ERROR.Print(err)
			completeChan <- err
			return
		}

		options := dockertypes.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		}
		logsReader, err := cli.ContainerLogs(ctx, args[0], options)
		if err != nil {
			jww.ERROR.Print(err)
			completeChan <- err
			return
		}

		logs, err := ioutil.ReadAll(logsReader)
		if err != nil {
			jww.ERROR.Print(err)
			completeChan <- err
			return
		}

		// Send the raw
		datas = append(datas, types.Data{
			Filename: filepath.Join("/raw/", filename),
			Data:     logs,
		})
		paths = append(paths, filepath.Join("/raw/", filename))

		completeChan <- nil
	}()

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf(`Getting logs from host:%s failed due to: %s`, args[0], ctx.Err().Error())}
	}

	result := types.Result{
		Task:      "dockerLogs",
		Args:      args,
		Filenames: paths,
		Error:     err,
	}

	return datas, result, err
}
