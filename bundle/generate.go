package bundle

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/replicatedcom/support-bundle/metrics"
	"github.com/replicatedcom/support-bundle/systemutil"
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
			Description: "Get File",
			ExecFunc:    systemutil.ReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"C:/Go/VERSION"},
		},

		Task{
			Description: "Get Other File",
			ExecFunc:    systemutil.ReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"C:/Go/README.md"},
		},

		Task{
			Description: "Run Command",
			ExecFunc:    systemutil.RunCommand,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"docker", "ps"},
		},

		Task{
			Description: "Run Other Command",
			ExecFunc:    systemutil.RunCommand,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"docker", "info"},
		},

		Task{
			Description: "Docker run command in container",
			ExecFunc:    systemutil.DockerRunCommand,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057", "root", "ls", "-a"},
		},

		Task{
			Description: "Docker run command in container with large output",
			ExecFunc:    systemutil.DockerRunCommand,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057", "root", "ls", "-a", "-R"},
		},

		Task{
			Description: "Docker run command in container, timeout",
			ExecFunc:    systemutil.DockerRunCommand,
			Timeout:     time.Duration(time.Second * 1),
			Args:        []string{"7e47d28f0057", "root", "sleep", "1m"},
		},

		Task{
			Description: "Docker read file from container",
			ExecFunc:    systemutil.DockerReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057", "/usr/local/bin/docker-entrypoint.sh"},
		},

		Task{
			Description: "Docker read file from container",
			ExecFunc:    systemutil.DockerReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057", "/usr/local/bin/"},
		},

		Task{
			Description: "Docker ps",
			ExecFunc:    metrics.Dockerps,
			Timeout:     time.Duration(time.Second * 15),
		},

		Task{
			Description: "Docker info",
			ExecFunc:    metrics.DockerInfo,
			Timeout:     time.Duration(time.Second * 15),
		},

		Task{
			Description: "Docker logs",
			ExecFunc:    metrics.DockerLogs,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057"},
		},

		Task{
			Description: "Docker inspect",
			ExecFunc:    metrics.DockerInspect,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{"7e47d28f0057"},
		},

		Task{
			Description: "System hostname",
			ExecFunc:    metrics.Hostname,
			Timeout:     time.Duration(time.Second * 1),
		},

		Task{
			Description: "System loadavg",
			ExecFunc:    metrics.LoadAvg,
			Timeout:     time.Duration(time.Second * 1),
		},

		Task{
			Description: "System uptime",
			ExecFunc:    metrics.Uptime,
			Timeout:     time.Duration(time.Second * 1),
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
