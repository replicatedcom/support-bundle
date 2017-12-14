package plans

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

// StreamsSource is a Task that gets its data as an io.Reader
type StreamsSource struct {
	// Producer provides the seed data for this task as an io.Reader array
	// Names of the sources are provided as a string array
	Producer func(context.Context) (map[string]io.Reader, error)
	// Parser, if defined, structures the raw data for json and human sinks
	Parser func(io.Reader) (interface{}, error)
	// StreamFormat describe stream format returned by Producer.  Only "" and "tar" are supported.
	StreamFormat string
	// RawScrubber, if defined, rewrites the raw data to to remove sensitive data
	RawScrubber func([]byte) []byte
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

func (task *StreamsSource) Exec(ctx context.Context, rootDir string) []*types.Result {
	if task.Producer == nil {
		err := errors.New("no data source defined for task")
		return task.resultsWithErr(err, "", "")
	}

	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""

	results := []*types.Result{}

	if !(raw || jsonify || human) {
		return results
	}

	if task.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, task.Timeout)
		defer cancel()
	}

	readers, err := task.Producer(ctx)
	if err != nil {
		return task.resultsWithErr(err, "", "")
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
				moreResults = task.execStream(ctx, rootDir, "", name, reader)

			case "tar":
				moreResults = task.execTarStream(ctx, rootDir, "", name, reader)

			default:
				if closer, ok := reader.(io.Closer); ok {
					closeLogErr(closer)
				}
				err := fmt.Errorf("unsupported stream format: %q", task.StreamFormat)
				moreResults = task.resultsWithErr(err, "", name)
			}

			resultsMut.Lock()
			results = append(results, moreResults...)
			resultsMut.Unlock()
		}(name, reader)
	}

	readerGroup.Wait()
	return results
}

func (task *StreamsSource) execTarStream(ctx context.Context, rootDir string, filePath string, fileExt string, reader io.Reader) []*types.Result {
	results := []*types.Result{}

	tarReader := tar.NewReader(reader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return task.resultsWithErr(err, filePath, fileExt)
		}

		if hdr.FileInfo().IsDir() {
			continue
		}

		moreResults := task.execStream(ctx, rootDir, hdr.Name, fileExt, tarReader)
		results = append(results, moreResults...)
	}

	return results
}

func (task *StreamsSource) execStream(ctx context.Context, rootDir string, filePath string, fileExt string, reader io.Reader) []*types.Result {
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

	rawResult := &types.Result{}
	jsonResult := &types.Result{}
	humanResult := &types.Result{}

	if raw {
		results = append(results, rawResult)
	}
	if jsonify {
		results = append(results, jsonResult)
	}
	if human {
		results = append(results, humanResult)
	}

	rawPath := mkPath(task.RawPath, filePath, fileExt)
	jsonPath := mkPath(task.JSONPath, filePath, fileExt)
	humanPath := mkPath(task.HumanPath, filePath, fileExt)

	var structured interface{}
	if parser {
		pr, pw := io.Pipe()
		tr := io.TeeReader(reader, pw)
		var err error
		structured, err = task.Parser(tr)
		reader = pr
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

func (task *StreamsSource) resultsWithErr(err error, filePath, fileExt string) []*types.Result {
	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""

	rawPath := mkPath(task.RawPath, filePath, fileExt)
	jsonPath := mkPath(task.JSONPath, filePath, fileExt)
	humanPath := mkPath(task.HumanPath, filePath, fileExt)

	results := []*types.Result{}

	if raw {
		results = append(results, &types.Result{Description: rawPath})
	}
	if jsonify {
		results = append(results, &types.Result{Description: jsonPath})
	}
	if human {
		results = append(results, &types.Result{Description: humanPath})
	}

	return resultsWithErr(err, results)
}

func mkPath(path, name, ext string) string {
	path = filepath.Join(path, name)
	if ext == "" {
		return path
	}
	return path + "." + ext
}
