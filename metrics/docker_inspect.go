package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/docker/docker/client"
)

// DockerInspect - inspects a given docker container
// args[0] should be container id
func DockerInspect(ctx context.Context, args []string) ([]types.Data, types.Result, error) {
	filename := "/docker/inspect/" + args[0]

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

		container, err := cli.ContainerInspect(ctx, args[0])
		if err != nil {
			jww.ERROR.Print(err)
			completeChan <- err
			return
		}

		containerJSON, err := json.Marshal(container)
		if err != nil {
			jww.ERROR.Print(err)
			completeChan <- err
			return
		}

		// Send the json
		datas = append(datas, types.Data{
			Filename: filepath.Join("/json/", filename+".json"),
			Data:     containerJSON,
		})
		paths = append(paths, filepath.Join("/json/", filename+".json"))

		containerIndentJSON, err := json.MarshalIndent(container, "", "  ")
		if err != nil {
			jww.ERROR.Print(err)
			completeChan <- err
			return
		}

		// Human readable version
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename+".json"),
			Data:     containerIndentJSON,
		})
		paths = append(paths, filepath.Join("/human/", filename+".json"))

		completeChan <- nil
	}()

	select {
	case err = <-completeChan:
		//completed on time
	case <-ctx.Done():
		//failed to complete on time
		err = types.TimeoutError{Message: fmt.Sprintf(`Inspecting host:%s failed due to: %s`, args[0], ctx.Err().Error())}
	}

	result := types.Result{
		Task:      "dockerInspect",
		Args:      args,
		Filenames: paths,
		Error:     err,
	}

	return datas, result, err
}
