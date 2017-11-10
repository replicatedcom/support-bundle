package plans

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

// StreamSource is a Task that gets its data as an io.Reader
type StreamSource struct {
	// Producer provides the seed data for this task as an io.Reader
	Producer func(context.Context) (io.Reader, error)
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

func (task *StreamSource) Exec(ctx context.Context, rootDir string) []*types.Result {
	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""

	results := []*types.Result{}

	if !(raw || jsonify || human) {
		return results
	}

	if raw {
		results = append(results, &types.Result{Description: task.RawPath})
	}
	if jsonify {
		results = append(results, &types.Result{Description: task.JSONPath})
	}
	if human {
		results = append(results, &types.Result{Description: task.HumanPath})
	}

	if task.Producer == nil {
		err := errors.New("no data source defined for task")
		return resultsWithErr(err, results)
	}

	if task.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, task.Timeout)
		defer cancel()
	}

	data, err := task.Producer(ctx)
	if err != nil {
		return resultsWithErr(err, results)
	}
	if closer, ok := data.(io.Closer); ok {
		defer closeLogErr(closer)
	}

	if task.StreamFormat == "" {
		return task.execStream(ctx, rootDir, "", data)
	}

	if task.StreamFormat != "tar" {
		err := fmt.Errorf("unsupported stream format: %q", task.StreamFormat)
		return resultsWithErr(err, results)
	}

	tarReader := tar.NewReader(data)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return resultsWithErr(err, results)
		}

		if hdr.FileInfo().IsDir() {
			continue
		}

		moreResults := task.execStream(ctx, rootDir, hdr.Name, tarReader)
		results = append(results, moreResults...)
	}

	return results
}

func (task *StreamSource) execStream(ctx context.Context, rootDir string, filePath string, reader io.Reader) []*types.Result {
	raw := task.RawPath != ""
	jsonify := task.JSONPath != ""
	human := task.HumanPath != ""

	results := []*types.Result{}

	if task.RawScrubber != nil {
		scrubbedReader, scrubbedWriter := io.Pipe()
		go filterStreams(reader, scrubbedWriter, task.RawScrubber)
		reader = scrubbedReader
	}

	rawResult := types.Result{}
	jsonResult := types.Result{}
	humanResult := types.Result{}

	rawPath := filepath.Join(task.RawPath, filePath)
	jsonPath := filepath.Join(task.JSONPath, filePath)
	humanPath := filepath.Join(task.HumanPath, filePath)

	// first write to one file
	if raw {
		writeResult(ctx, rootDir, rawPath, &rawResult, reader)
	} else if jsonify {
		writeResult(ctx, rootDir, jsonPath, &jsonResult, reader)
	} else if human {
		writeResult(ctx, rootDir, humanPath, &humanResult, reader)
	}

	// then link to any other requested paths
	if raw && jsonify {
		jsonResult = rawResult
		if rawResult.Path != "" {
			os.Link(rawPath, jsonPath)
			jsonResult.Path = jsonPath
		}
	}
	if raw && human {
		humanResult = rawResult
		if rawResult.Path != "" {
			os.Link(rawPath, humanPath)
			humanResult.Path = humanPath
		}
	}
	if jsonify && human {
		humanResult = jsonResult
		if jsonResult.Path != "" {
			os.Link(jsonPath, humanPath)
			humanResult.Path = humanPath
		}
	}

	if raw {
		results = append(results, &rawResult)
	}
	if jsonify {
		results = append(results, &jsonResult)
	}
	if human {
		results = append(results, &humanResult)
	}

	return results
}
