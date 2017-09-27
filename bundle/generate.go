package bundle

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/replicatedcom/support-bundle/types"

	"github.com/divolgin/archiver/compressor"
	jww "github.com/spf13/jwalterweatherman"
)

// Generate is called to start a new support bundle generation
func Generate() error {
	var wg sync.WaitGroup

	resultsCh := make(chan types.Result)
	dataCh := make(chan types.Data)
	completeCh := make(chan bool)

	var tasks = []Task{
		Task{
			Description: "System Log Files",
			ExecFunc:    systemLogFiles,
			Timeout:     time.Duration(time.Second * 15),
		},

		Task{
			Description: "System Metrics",
			ExecFunc:    systemMetrics,
			Timeout:     time.Duration(time.Second * 15),
		},
	}
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
		_ = task.ExecFunc(dataCh, completeCh, resultsCh, task.Timeout)
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
