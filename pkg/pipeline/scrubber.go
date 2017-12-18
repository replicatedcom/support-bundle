package pipeline

import (
	"bufio"
	"context"
	"io"

	"github.com/pkg/errors"
)

type LineScrubber struct {
	Scrubber func([]byte) []byte
}

func (p *LineScrubber) Pipe(ctx context.Context, i interface{}) (interface{}, error) {
	if s, ok := i.(RawSource); ok {
		out, err := p.pipe(ctx, s)
		return RawSource(out), err
	}
	if s, ok := i.(HumanSource); ok {
		out, err := p.pipe(ctx, s)
		return HumanSource(out), err
	}
	if s, ok := i.(JSONSource); ok {
		out, err := p.pipe(ctx, s)
		return JSONSource(out), err
	}
	return nil, errors.New("LineScrubber requires RawSource, HumanSource or JSONSource")
}

func (p *LineScrubber) pipe(ctx context.Context, s io.Reader) (io.Reader, error) {
	pr, pw := io.Pipe()
	go func() {
		lineScanner := bufio.NewScanner(s)
		for lineScanner.Scan() {
			// NOTE: this appends a final newline even if there isn't one
			_, err := pw.Write(append(p.Scrubber(lineScanner.Bytes()), []byte("\n")...))
			if err != nil {
				pw.CloseWithError(err)
				return
			}
		}
		pw.CloseWithError(lineScanner.Err())
	}()
	return pr, nil
}
