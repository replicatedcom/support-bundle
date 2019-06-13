package producers

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/util"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"

	dockertypes "github.com/docker/docker/api/types"
)

const (
	LogsReaderIdleTimeout = 20 * time.Second
)

func (d *Docker) TaskLogs(taskID string, opts *dockertypes.ContainerLogsOptions) types.StreamsProducer {
	return func(ctx context.Context) (map[string]io.Reader, error) {
		if opts == nil {
			opts = &dockertypes.ContainerLogsOptions{}
		}
		if !opts.ShowStdout && !opts.ShowStderr {
			opts.ShowStdout = true
			opts.ShowStderr = true
		}
		opts.Timestamps = true
		opts.Follow = false

		reader, err := d.client.TaskLogs(ctx, taskID, *opts)
		if err != nil {
			return nil, err
		}
		readerWithTimeout := logsReaderWithTimeout(reader, LogsReaderIdleTimeout)
		return util.DemuxLogs(ctx, readerWithTimeout, taskID)
	}
}

func logsReaderWithTimeout(r io.ReadCloser, timeout time.Duration) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer r.Close()

		size := 32 * 1024
		buf := make([]byte, size)

		for {
			var n int
			var err error
			done := make(chan struct{})
			go func() {
				n, err = r.Read(buf)
				close(done)
			}()

			select {
			case <-done:
				if err != nil {
					pw.CloseWithError(err) // this will handle the EOF
					return
				}
				if _, err := pw.Write(buf[:n]); err != nil {
					pw.CloseWithError(err)
					return
				}

			case <-time.After(timeout):
				// Fix issue https://github.com/moby/moby/issues/38640
				// Service and task logs hang for large logs
				pw.CloseWithError(fmt.Errorf("reader timeout after %s", timeout))
				return
			}
		}
	}()
	return pr
}
