package plans

import (
	"bufio"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	jww "github.com/spf13/jwalterweatherman"
)

var GlobalScrubbers []types.BytesScrubber
var GlobalScrubbersLock sync.Mutex

func SetCommonFieldsStreamsSource(task StreamsSource, spec types.Spec) (StreamsSource, error) {
	task.Spec = spec
	if task.JSONPath == "" && task.HumanPath == "" && task.RawPath == "" {
		task.RawPath = spec.Shared().OutputDir
	}
	scrubber, err := RawScrubber(spec.Shared().Scrub)
	if err != nil {
		return task, errors.Wrap(err, "create scrubber")
	}
	task.RawScrubber = scrubber
	if spec.Shared().TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.Shared().TimeoutSeconds) * time.Second
	}
	return task, nil
}

func SetCommonFieldsStreamSource(task StreamSource, spec types.Spec) (StreamSource, error) {
	task.Spec = spec
	if task.JSONPath == "" && task.HumanPath == "" && task.RawPath == "" {
		task.RawPath = spec.Shared().OutputDir
	}
	scrubber, err := RawScrubber(spec.Shared().Scrub)
	if err != nil {
		return task, errors.Wrap(err, "create scrubber")
	}
	task.RawScrubber = scrubber
	if spec.Shared().TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.Shared().TimeoutSeconds) * time.Second
	}
	return task, nil
}

func SetCommonFieldsStructuredSource(task StructuredSource, spec types.Spec) (StructuredSource, error) {
	task.Spec = spec
	if task.JSONPath == "" && task.HumanPath == "" && task.RawPath == "" {
		task.RawPath = spec.Shared().OutputDir
	}
	scrubber, err := RawScrubber(spec.Shared().Scrub)
	if err != nil {
		return task, errors.Wrap(err, "create scrubber")
	}
	task.RawScrubber = scrubber
	if spec.Shared().TimeoutSeconds != 0 {
		task.Timeout = time.Duration(spec.Shared().TimeoutSeconds) * time.Second
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

// closeLogErr
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

func filterStreams(readFrom io.Reader, writeTo io.Writer, scrubber func([]byte) []byte) error {
	// preserve newline at end
	var readCounter bytesCounter
	tee := io.TeeReader(readFrom, &readCounter)
	lineReader := bufio.NewReader(tee)
	var writeCounter bytesCounter
	writeTee := io.MultiWriter(writeTo, &writeCounter)

	i := 0
	var err error
	for err == nil {
		var line []byte
		line, err = readEntireLine(lineReader)
		if err != nil {
			break
		}
		line = scrubber(line)

		if i > 0 {
			line = append([]byte("\n"), line...)
		}

		_, werr := writeTee.Write(line)
		if werr != nil {
			return werr
		}
		i++
	}
	if err != nil && err != io.EOF {
		return err
	}
	for writeCounter.n < readCounter.n {
		_, err := writeTee.Write([]byte("\n"))
		return err
	}
	return nil
}

type bytesCounter struct {
	n int
}

func (w *bytesCounter) Write(p []byte) (int, error) {
	n := len(p)
	w.n += n
	return n, nil
}

func readEntireLine(reader *bufio.Reader) ([]byte, error) {
	var line []byte
	isPrefix := true
	var err error

	for isPrefix && err == nil {
		var newLine []byte
		newLine, isPrefix, err = reader.ReadLine()
		line = append(line, newLine...)
	}

	return line, err
}
