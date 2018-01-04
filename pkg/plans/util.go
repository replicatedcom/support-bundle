package plans

import (
	"bufio"
	"io"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

func SetCommonFieldsStreamsSource(task StreamsSource, spec types.Spec) (StreamsSource, error) {
	task.Spec = spec
	if task.RawPath == "" {
		task.RawPath = spec.OutputDir
	}
	scrubber, err := RawScrubber(spec.Scrub)
	if err != nil {
		return task, errors.Wrap(err, "create scrubber")
	}
	task.RawScrubber = scrubber
	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}
	return task, nil
}

func SetCommonFieldsStreamSource(task StreamSource, spec types.Spec) (StreamSource, error) {
	task.Spec = spec
	if task.RawPath == "" {
		task.RawPath = spec.OutputDir
	}
	scrubber, err := RawScrubber(spec.Scrub)
	if err != nil {
		return task, errors.Wrap(err, "create scrubber")
	}
	task.RawScrubber = scrubber
	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}
	return task, nil
}

func SetCommonFieldsStructuredSource(task StructuredSource, spec types.Spec) (StructuredSource, error) {
	task.Spec = spec
	if task.RawPath == "" {
		task.RawPath = spec.OutputDir
	}
	scrubber, err := RawScrubber(spec.Scrub)
	if err != nil {
		return task, errors.Wrap(err, "create scrubber")
	}
	task.RawScrubber = scrubber
	if spec.TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.TimeoutSeconds) * time.Second
	}
	return task, nil
}

// add an error to every result, returning the results argument
// skips results that already have an error
func resultsWithErr(err error, results []*types.Result) []*types.Result {
	for _, r := range results {
		if r.Error == nil {
			r.Error = err
		}
	}

	return results
}

// closseLogErr
func closeLogErr(c io.Closer) {
	if err := c.Close(); err != nil {
		jww.ERROR.Printf("Failed to close closer: %v", err)
	}
}

// RawScrubber creates a scrubber function from a scrubSpec
func RawScrubber(scrubSpec *types.Scrub) (types.BytesScrubber, error) {
	if scrubSpec == nil || scrubSpec.Regex == "" {
		return nil, nil
	}

	regex, err := regexp.Compile(scrubSpec.Regex)
	if err != nil {
		return nil, errors.Wrapf(err, "parse regex %s", scrubSpec.Regex)
	}

	return func(in []byte) []byte {
		return regex.ReplaceAll(in, []byte(scrubSpec.Replace))
	}, nil
}

func filterStreams(readFrom io.Reader, writeTo *io.PipeWriter, scrubber func([]byte) []byte) {
	lineScanner := bufio.NewScanner(readFrom)
	for lineScanner.Scan() {
		line := lineScanner.Bytes()
		line = scrubber(line)

		_, err := writeTo.Write(line)
		if err != nil {
			writeTo.CloseWithError(err)
			return
		}
	}
	writeTo.Close()
}
