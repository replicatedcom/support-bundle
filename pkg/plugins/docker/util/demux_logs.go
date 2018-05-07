package util

import (
	"context"
	"io"

	"github.com/docker/docker/pkg/stdcopy"
	jww "github.com/spf13/jwalterweatherman"
)

// Given a docker logs stream, return a map of the stdout and stderr components
func DemuxLogs(ctx context.Context, source io.Reader, prefix string) (map[string]io.Reader, error) {
	stdoutR, stdoutW := io.Pipe()
	stderrR, stderrW := io.Pipe()

	go func() {
		// TODO context interruptable
		_, err := stdcopy.StdCopy(stdoutW, stderrW, source)
		if err != nil {
			stdoutW.CloseWithError(err)
			stderrW.CloseWithError(err)
			return
		}
		if err := stdoutW.Close(); err != nil {
			jww.ERROR.Print(err)
		}
		if err := stderrW.Close(); err != nil {
			jww.ERROR.Print(err)
		}
	}()

	readers := make(map[string]io.Reader)
	readers[prefix+".stdout"] = stdoutR
	readers[prefix+".stderr"] = stderrR
	return readers, nil
}
