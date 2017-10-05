package bundle

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/replicatedcom/support-bundle/types"

	"github.com/divolgin/archiver/compressor"
	jww "github.com/spf13/jwalterweatherman"
)

type resultInfo struct {
	Paths []string `json:"paths"`
	Task  string   `json:"task"`
	Args  []string `json:"arguments"`
}

type errorInfo struct {
	Task   string   `json:"task"`
	Args   []string `json:"arguments"`
	Errors []string `json:"errors"`
}

// Generate is called to start a new support bundle generation
func Generate(tasks []Task) (io.Reader, error) {
	var wg sync.WaitGroup

	resultsCh := make(chan types.Result)
	dataCh := make(chan types.Data)
	completeCh := make(chan bool)

	var resultMutex = &sync.Mutex{}
	var allResultInfo []resultInfo
	var allErrorInfo []errorInfo

	wg.Add(len(tasks))

	collectDir, err := ioutil.TempDir("", "support-bundle")
	if err != nil {
		jww.ERROR.Fatal(err)
		return nil, err
	}
	defer os.RemoveAll(collectDir)

	go func() {
		for {
			select {
			case <-completeCh:
				wg.Done()
				break

			case result := <-resultsCh:
				resultMutex.Lock()
				allResultInfo = append(allResultInfo, resultInfo{
					Paths: result.Filenames,
					Task:  result.Task,
					Args:  result.Args,
				})
				if result.RawError != nil || result.HumanError != nil || result.JSONError != nil {
					newErrorInfo := errorInfo{
						Task:   result.Task,
						Args:   result.Args,
						Errors: []string{},
					}
					if result.RawError != nil {
						newErrorInfo.Errors = append(newErrorInfo.Errors, "Raw: "+result.RawError.Error())
					}
					if result.HumanError != nil {
						newErrorInfo.Errors = append(newErrorInfo.Errors, "Human: "+result.HumanError.Error())
					}
					if result.JSONError != nil {
						newErrorInfo.Errors = append(newErrorInfo.Errors, "JSON: "+result.JSONError.Error())
					}

					allErrorInfo = append(allErrorInfo, newErrorInfo)
				}
				resultMutex.Unlock()
				break

			case data := <-dataCh:
				dataFile := filepath.Join(collectDir, data.Filename)
				if err := os.MkdirAll(filepath.Dir(dataFile), 0700); err != nil {
					jww.ERROR.Print(err)
					continue
				}
				if err := ioutil.WriteFile(dataFile, data.Data, 0666); err != nil {
					jww.ERROR.Print(err)
					continue
				}
				break
			}
		}
	}()

	for _, task := range tasks {
		go func(task Task) {
			ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
			defer cancel()
			datas, results, err := task.ExecFunc(ctx, task.Args)

			if err != nil {
				jww.ERROR.Printf(err.Error() + "\n")

				//replace errors with marshallable versions (most error types will return {} when converted to json)
				// results.RawError = types.MarshallableError{Message: results.RawError.Error()}
				// results.HumanError = types.MarshallableError{Message: results.HumanError.Error()}
				// results.JSONError = types.MarshallableError{Message: results.JSONError.Error()}
			}

			for _, data := range datas {
				dataCh <- data
			}

			resultsCh <- results
			completeCh <- true
		}(task)
	}

	wg.Wait()

	//write index and error json files
	indexJSON, err := json.MarshalIndent(allResultInfo, "", "  ")
	if err != nil {
		jww.ERROR.Print(err)
	}
	ioutil.WriteFile(filepath.Join(collectDir, "index.json"), indexJSON, 0666)

	errorJSON, err := json.MarshalIndent(allErrorInfo, "", "  ")
	if err != nil {
		jww.ERROR.Print(err)
	}
	ioutil.WriteFile(filepath.Join(collectDir, "error.json"), errorJSON, 0666)

	// Build the output tar file
	archiveFile, err := ioutil.TempFile("", "support-bundle")
	if err != nil {
		jww.ERROR.Fatal(err)
		return nil, err
	}

	comp := compressor.NewTgz()
	comp.SetTarConfig(compressor.Tar{TruncateLongFiles: true})
	if err := comp.Compress(collectDir, archiveFile.Name()); err != nil {
		jww.ERROR.Fatal(err)
		return nil, err
	}

	jww.TRACE.Printf("Created support bundle at %q\n", archiveFile.Name())

	return os.Open(archiveFile.Name())
}
