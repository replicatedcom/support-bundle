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

	"github.com/replicatedcom/support-bundle/types"

	jww "github.com/spf13/jwalterweatherman"

	dockertypes "github.com/docker/docker/api/types"
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

		readcloser, fileinfo, err := cli.CopyFromContainer(ctx, args[0], args[1])
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

		// Human readable version
		datas = append(datas, types.Data{
			Filename: filepath.Join("/human/", filename+".tar"),
			Data:     response,
		})

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
				completeChan <- err
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
			completeChan <- err
			return
		}

		// Send the json
		datas = append(datas, types.Data{
			Filename: filepath.Join("/json/", filename+".json"),
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
		err = types.TimeoutError{Message: fmt.Sprintf(`Reading a docker file at from host:%s and path:%s errored out with %s`, args[0], args[1], ctx.Err().Error())}
		rawError = err
		jsonError = err
		humanError = err
	}

	result := types.Result{
		Name:        "dockerReadFile",
		Description: "A file from a docker container",
		Filename:    filename,
		RawError:    rawError,
		JSONError:   jsonError,
		HumanError:  humanError,
	}

	return datas, result, err
}
