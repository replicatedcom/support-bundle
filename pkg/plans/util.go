package plans

import (
	"bufio"
	"io"
	"regexp"

	"github.com/pkg/errors"

	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

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
		jww.ERROR.Print(err)
	}
}

// RawScrubber creates a scrubber function from a scrubSpec
func RawScrubber(scrubSpec types.Scrub) (types.BytesScrubber, error) {
	if scrubSpec.Regex == "" {
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
		}
	}
	writeTo.Close()
}
