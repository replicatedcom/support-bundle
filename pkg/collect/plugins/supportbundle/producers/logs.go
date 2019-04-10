package producers

import (
	"bytes"
	"context"
	"io"
)

func (s *SupportBundle) Logs(ctx context.Context) (io.Reader, error) {
	str := s.LogOutput.String()
	return bytes.NewBufferString(str), nil
}
