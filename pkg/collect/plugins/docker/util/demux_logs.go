package util

import (
	"context"
	"io"

	"github.com/docker/docker/pkg/stdcopy"
	jww "github.com/spf13/jwalterweatherman"
)

// DemuxLogs will return a map of the stdout and stderr components given a docker logs stream.
func DemuxLogs(ctx context.Context, source io.Reader, prefix string) (map[string]io.Reader, error) {
	stdoutR, stdoutW := io.Pipe()
	stderrR, stderrW := io.Pipe()

	go func() {
		// TODO context interruptable
		_, err := stdcopy.StdCopy(stdoutW, stderrW, source)
		if closer, ok := source.(io.Closer); ok {
			closer.Close()
		}
		if err := stdoutW.CloseWithError(err); err != nil {
			jww.ERROR.Print(err)
		}
		if err := stderrW.CloseWithError(err); err != nil {
			jww.ERROR.Print(err)
		}
	}()

	readers := make(map[string]io.Reader)
	readers[prefix+".stdout"] = stdoutR
	readers[prefix+".stderr"] = stderrR
	return readers, nil
}
