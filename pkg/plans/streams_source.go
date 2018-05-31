package plans

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

// StreamsSource is a Task that gets its data as an io.Reader
type StreamsSource struct {
	Spec types.Spec

	// Producer provides the seed data for this task as an io.Reader array
	// Names of the sources are provided as a string array
	Producer func(context.Context) (map[string]io.Reader, error)
	// StreamFormat describe stream format returned by Producer. Only "" and "tar" are supported.
	StreamFormat string
	// RawScrubber, if defined, rewrites the raw data to to remove sensitive data
	RawScrubber func([]byte) []byte
	// Parser, if defined, structures the raw data for json and human sinks
	Parser func(io.Reader) (interface{}, error)
	// Template, if defined, renders structured data in a human-readable format
	Template string
	// If RawPath is defined it will get a copy of the data
	RawPath string
	// If JSONPath is defined it will get a copy of the data
	JSONPath string
	// If HumanPath is defined it will get a copy of the data
	HumanPath string
	// If Timeout is defined, it will be used rather than the context provided
	// to Exec.
	Timeout time.Duration
}

func (task *StreamsSource) GetSpec() types.Spec {
	return task.Spec
}

func (task *StreamsSource) Exec(ctx context.Context, rootDir string) []*types.Result {
	cancel := make(chan struct{})
	defer close(cancel)
	go func() {
		select {
		case _, _ = <-cancel:
		case <-ctx.Done():
			b, _ := json.Marshal(task.Spec)
			jww.WARN.Println("Task failed to complete before context was canceled:", string(b))
		}
	}()
	if task.Producer == nil {
		err := errors.New("no data source defined for task")
		return task.resultsWithErr(err, "")
	}

	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""

	results := []*types.Result{}

	if !(raw || jsonify || human) {
		return results
	}

	startTime := time.Now()
	specBytes, err := json.Marshal(task.Spec)
	if err != nil {
		jww.ERROR.Printf("Unable to parse spec due to %s", err.Error())
	} else {
		defer func() {
			jww.DEBUG.Printf("Task with spec %s completed in %s", string(specBytes), time.Since(startTime).String())
		}()
	}

	if task.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, task.Timeout)
		defer cancel()
	}

	readers, err := task.Producer(ctx)
	if err != nil {
		return task.resultsWithErr(err, "")
	}

	var resultsMut sync.Mutex
	var readerGroup sync.WaitGroup
	readerGroup.Add(len(readers))

	for name, reader := range readers {
		go func(name string, reader io.Reader) {
			defer readerGroup.Done()

			var moreResults []*types.Result
			switch task.StreamFormat {
			case "":
				moreResults = task.execStream(ctx, rootDir, name, reader)

			case StreamFormatTar:
				moreResults = task.execTarStream(ctx, rootDir, name, reader)

			default:
				if closer, ok := reader.(io.Closer); ok {
					closeLogErr(closer)
				}
				err := fmt.Errorf("unsupported stream format: %q", task.StreamFormat)
				moreResults = task.resultsWithErr(err, name)
			}

			resultsMut.Lock()
			results = append(results, moreResults...)
			resultsMut.Unlock()
		}(name, reader)
	}

	readerGroup.Wait()
	if !task.Spec.Shared().IncludeEmpty {
		task.cleanupResults(ctx, rootDir, results)
	}
	return results
}

func (task *StreamsSource) execTarStream(ctx context.Context, rootDir string, filePath string, reader io.Reader) []*types.Result {
	results := []*types.Result{}

	tarReader := tar.NewReader(reader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return task.resultsWithErr(err, filePath)
		}

		if hdr.FileInfo().IsDir() {
			continue
		}

		var fp string
		if filePath != "" {
			fp = path.Join(filePath, hdr.Name)
		} else {
			fp = hdr.Name
		}
		moreResults := task.execStream(ctx, rootDir, fp, tarReader)
		results = append(results, moreResults...)
	}

	return results
}

func (task *StreamsSource) execStream(ctx context.Context, rootDir string, filePath string, reader io.Reader) []*types.Result {
	if closer, ok := reader.(io.Closer); ok {
		defer closeLogErr(closer)
	}

	parser := task.Parser != nil
	templated := task.Template != ""

	raw := task.RawPath != ""

	jsonify := task.JSONPath != ""
	jsonParsed := jsonify && parser
	jsonRaw := jsonify && !jsonParsed

	human := task.HumanPath != ""
	humanTemplated := human && parser && templated
	humanYAML := human && parser && !humanTemplated
	humanRaw := human && !humanTemplated && !humanYAML

	results := []*types.Result{}

	if task.RawScrubber != nil {
		scrubbedReader, scrubbedWriter := io.Pipe()
		go filterStreams(reader, scrubbedWriter, task.RawScrubber)
		reader = scrubbedReader
	}

	rawResult := &types.Result{Spec: task.Spec}
	jsonResult := &types.Result{Spec: task.Spec}
	humanResult := &types.Result{Spec: task.Spec}

	if raw {
		results = append(results, rawResult)
	}
	if jsonify {
		results = append(results, jsonResult)
	}
	if human {
		results = append(results, humanResult)
	}

	rawPath := filepath.Join(task.RawPath, filePath)
	jsonPath := filepath.Join(task.JSONPath, filePath)
	humanPath := filepath.Join(task.HumanPath, filePath)

	var structured interface{}
	if parser {
		buf := bytes.NewBuffer(nil)
		tr := io.TeeReader(reader, buf)
		var err error
		structured, err = task.Parser(tr)
		reader = buf
		if err != nil {
			jsonResult.Error = err
			if humanTemplated || humanYAML {
				humanResult.Error = err
				// nothing left to do
				return results
			}
		} else {
			writeResultJSON(ctx, rootDir, task.JSONPath, jsonResult, structured)
		}
	}

	if raw {
		writeResult(ctx, rootDir, rawPath, rawResult, reader)
	} else if jsonRaw {
		writeResult(ctx, rootDir, jsonPath, jsonResult, reader)
	} else if humanRaw {
		writeResult(ctx, rootDir, humanPath, humanResult, reader)
	}

	if raw && jsonRaw {
		jsonResult = rawResult
		if rawResult.Path != "" {
			if err := os.Link(rawPath, jsonPath); err != nil {
				jsonResult.Error = err
			} else {
				jsonResult.Path = jsonPath
			}
		}
	}

	if humanTemplated {
		writeResultTemplate(ctx, rootDir, task.HumanPath, humanResult, task.Template, structured)
	}

	if humanYAML {
		writeResultYAML(ctx, rootDir, task.HumanPath, humanResult, structured)
	}

	if raw && humanRaw {
		humanResult = rawResult
		if rawResult.Path != "" {
			if err := os.Link(rawPath, humanPath); err != nil {
				humanResult.Error = err
			} else {
				humanResult.Path = humanPath
			}
		}
	} else if jsonRaw && humanRaw {
		humanResult = jsonResult
		if jsonResult.Path != "" {
			if err := os.Link(jsonPath, humanPath); err != nil {
				humanResult.Error = err
			} else {
				humanResult.Path = humanPath
			}
		}
	}

	return results
}

func (task *StreamsSource) resultsWithErr(err error, filePath string) []*types.Result {
	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""

	rawPath := filepath.Join(task.RawPath, filePath)
	jsonPath := filepath.Join(task.JSONPath, filePath)
	humanPath := filepath.Join(task.HumanPath, filePath)

	results := []*types.Result{}

	if raw {
		results = append(results, &types.Result{Spec: task.Spec, Path: rawPath})
	}
	if jsonify {
		results = append(results, &types.Result{Spec: task.Spec, Path: jsonPath})
	}
	if human {
		results = append(results, &types.Result{Spec: task.Spec, Path: humanPath})
	}

	return resultsWithErr(err, results)
}

func (task *StreamsSource) cleanupResults(ctx context.Context, rootDir string, results []*types.Result) {
	for _, result := range results {
		if result.Size == 0 {
			err := os.Remove(path.Join(rootDir, result.Path))
			if err != nil {
				jww.DEBUG.Printf("Unable to remove empty file %s within %s because of %s", result.Path, rootDir, err.Error())
			}
		}
	}
}
