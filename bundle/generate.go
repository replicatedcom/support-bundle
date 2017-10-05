package bundle

import (
	"context"
	"encoding/json"
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
func Generate(tasks []Task) error {
	var wg sync.WaitGroup

	resultsCh := make(chan types.Result)
	dataCh := make(chan types.Data)
	completeCh := make(chan bool)

	wg.Add(len(tasks))

	collectDir, err := ioutil.TempDir("", "support-bundle")
	if err != nil {
		jww.ERROR.Fatal(err)
		return err
	}
	defer os.RemoveAll(collectDir)

	go func() {
		for {
			select {
			case <-completeCh:
				wg.Done()
				break

			case result := <-resultsCh:
				dataFile := filepath.Join(collectDir, "/results/", result.Filename)
				if err := os.MkdirAll(filepath.Dir(dataFile), 0700); err != nil {
					jww.ERROR.Print(err)
					continue
				}
				b, err := json.Marshal(result)
				if err != nil {
					jww.ERROR.Print(err)
					continue
				}
				if err := ioutil.WriteFile(dataFile, b, 0666); err != nil {
					jww.ERROR.Print(err)
					continue
				}
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
				results.RawError = types.MarshallableError{Message: results.RawError.Error()}
				results.HumanError = types.MarshallableError{Message: results.HumanError.Error()}
				results.JSONError = types.MarshallableError{Message: results.JSONError.Error()}
			}

			for _, data := range datas {
				dataCh <- data
			}

			resultsCh <- results
			completeCh <- true
		}(task)
	}

	wg.Wait()

	// Build the output tar file
	archiveFilename := "./support-bundle.tar.gz"
	comp := compressor.NewTgz()
	comp.SetTarConfig(compressor.Tar{TruncateLongFiles: true})
	if err := comp.Compress(collectDir, archiveFilename); err != nil {
		jww.ERROR.Fatal(err)
		return err
	}

	jww.TRACE.Printf("Created support bundle at %q\n", archiveFilename)
	return nil
}
