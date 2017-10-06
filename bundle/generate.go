package bundle

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
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
	Task  string   `json:"task"`
	Args  []string `json:"arguments"`
	Error string   `json:"error"`
}

// Generate is called to start a new support bundle generation
func Generate(tasks []Task, timeout time.Duration) (string, error) {
	var wg sync.WaitGroup

	var resultMutex = &sync.Mutex{}
	var allResultInfo []resultInfo
	var allErrorInfo []errorInfo

	wg.Add(len(tasks))

	collectDir, err := ioutil.TempDir("", "support-bundle")
	if err != nil {
		err = errors.Wrap(err, "Creating a temporary directory to store results failed")
		jww.ERROR.Print(err)
		return "", err
	}
	defer os.RemoveAll(collectDir)

	ctx := context.Background()
	defaultCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, task := range tasks {
		go func(task Task) {
			var datas []types.Data
			var result types.Result
			var err error

			if task.Timeout == 0 {
				// use the default context for this task
				datas, result, err = task.ExecFunc(defaultCtx, task.Args)
			} else {
				// use a unique context+timeout for this task
				uniqueCtx, cancel := context.WithTimeout(ctx, task.Timeout)
				defer cancel()
				datas, result, err = task.ExecFunc(uniqueCtx, task.Args)
			}

			if err != nil {
				jww.ERROR.Printf(err.Error() + "\n")
			}

			for _, data := range datas {
				dataFile := filepath.Join(collectDir, data.Filename)
				if err := os.MkdirAll(filepath.Dir(dataFile), 0700); err != nil {
					jww.ERROR.Print(err)
				}
				if err := ioutil.WriteFile(dataFile, data.Data, 0666); err != nil {
					jww.ERROR.Print(err)
				}
			}

			resultMutex.Lock()
			allResultInfo = append(allResultInfo, resultInfo{
				Paths: result.Filenames,
				Task:  result.Task,
				Args:  result.Args,
			})
			if result.Error != nil {
				allErrorInfo = append(allErrorInfo, errorInfo{
					Task:  result.Task,
					Args:  result.Args,
					Error: result.Error.Error(),
				})
			}
			resultMutex.Unlock()

			wg.Done()
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
		err = errors.Wrap(err, "Creating a temporary file to compress results failed")
		jww.ERROR.Print(err)
		return "", err
	}

	comp := compressor.NewTgz()
	comp.SetTarConfig(compressor.Tar{TruncateLongFiles: true})
	if err := comp.Compress(collectDir, archiveFile.Name()); err != nil {
		err = errors.Wrap(err, "Compressing results directory failed")
		jww.ERROR.Print(err)
		return "", err
	}

	jww.TRACE.Printf("Created support bundle at %q\n", archiveFile.Name())

	return archiveFile.Name(), nil
}
