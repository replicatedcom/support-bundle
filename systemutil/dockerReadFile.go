package systemutil

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	timeoutChan := make(chan error, 1)

	go func() {
		cli, err := client.NewEnvClient()
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			timeoutChan <- err
			return
		}

		readcloser, fileinfo, err := cli.CopyFromContainer(context.Background(), args[0], args[1])
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			timeoutChan <- err
			return
		}

		// read everything
		response, err := ioutil.ReadAll(readcloser)
		if err != nil {
			jww.ERROR.Print(err)
			rawError = err
			jsonError = err
			humanError = err
			timeoutChan <- err
			return
		}

		//close connection
		readcloser.Close()

		// Send the raw
		dataCh <- types.Data{
			Filename: filepath.Join("/raw/", filename+".tar"),
			Data:     response,
		}

		// Human readable version
		dataCh <- types.Data{
			Filename: filepath.Join("/human/", filename+".tar"),
			Data:     response,
		}

		//make a new buffer of the read file
		newReader := bytes.NewReader(response)

		type readFileStruct struct {
			FileContent string     `json:"filecontent"`
			Header      tar.Header `json:"fileheader"`
		}
		readFiles := []readFileStruct{}

		//remove the tar header & store all files
		tr := tar.NewReader(newReader)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				//end of tar archive
				break
			} else if err != nil {
				jww.ERROR.Print(err)
				jsonError = err
				timeoutChan <- err
				return
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(tr)

			readFiles = append(readFiles, readFileStruct{
				FileContent: buf.String(),
				Header:      *hdr,
			})
		}

		type readFilesStruct struct {
			Files []readFileStruct              `json:"files"`
			Info  dockertypes.ContainerPathStat `json:"info"`
		}
		u := readFilesStruct{
			Files: readFiles,
			Info:  fileinfo,
		}
		j, err := json.Marshal(u)
		if err != nil {
			jww.ERROR.Print(err)
			jsonError = err
			timeoutChan <- err
			return
		}

		// Send the json
		dataCh <- types.Data{
			Filename: filepath.Join("/json/", filename+".json"),
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
		err := types.TimeoutError{Message: fmt.Sprintf(`Reading a docker file at from host:%s and path:%s timed out after %s`, args[0], args[1], timeout.String())}
		rawError = err
		jsonError = err
		humanError = err
		return err
	}
}
