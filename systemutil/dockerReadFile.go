package systemutil

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerReadFile Read a file from a docker instance
// args: ["id", "path_to_file"]
func DockerReadFile(dataCh chan types.Data, completeCh chan bool, resultsCh chan types.Result, timeout time.Duration, args []string) error {
	filename := "/docker/readfile/"

	r, _ := regexp.Compile(`[^\w]`)

	for _, arg := range args {
		filename += r.ReplaceAllString(arg, "_")
	}

	var rawError, jsonError, humanError error = nil, nil, nil
	defer func() {
		resultsCh <- types.Result{
			Name:        "dockerReadFile",
			Description: "A file from a docker container",
			Filename:    filename,
			RawError:    rawError,
			JSONError:   jsonError,
			HumanError:  humanError,
		}
		completeCh <- true
	}()

	cli, err := client.NewEnvClient()
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	readcloser, fileinfo, err := cli.CopyFromContainer(context.Background(), args[0], args[1])
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	// read everything
	response, err := ioutil.ReadAll(readcloser)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		humanError = err
		return err
	}

	//close connection
	readcloser.Close()

	// Send the raw
	dataCh <- types.Data{
		Filename: filepath.Join("/raw/", filename),
		Data:     response,
	}

	// Human readable version
	dataCh <- types.Data{
		Filename: filepath.Join("/human/", filename),
		Data:     response,
	}

	type readFileStruct struct {
		File string                        `json:"file"`
		Info dockertypes.ContainerPathStat `json:"info"`
	}
	u := readFileStruct{
		File: string(response),
		Info: fileinfo,
	}
	j, err := json.Marshal(u)
	if err != nil {
		jww.ERROR.Print(err)
		rawError = err
		jsonError = err
		return err
	}

	// Send the json
	dataCh <- types.Data{
		Filename: filepath.Join("/json/", filename+".json"),
		Data:     j,
	}

	return nil
}
