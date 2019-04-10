package producers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"runtime/pprof"
)

func (s *SupportBundle) Goroutines(ctx context.Context) (io.Reader, error) {
	profile := pprof.Lookup("goroutine")
	if profile == nil {
		return nil, errors.New("goroutine profile not found")
	}
	buf := bytes.NewBuffer(nil)
	err := profile.WriteTo(buf, 2)
	return buf, err
}
