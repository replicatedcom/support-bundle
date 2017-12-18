package pipeline

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

type Parser struct {
	Parse func(io.Reader) (interface{}, error)
}

func (p *Parser) Pipe(ctx context.Context, i interface{}) (interface{}, error) {
	if s, ok := i.(RawSource); ok {
		out, err := p.pipe(ctx, s)
		return out, err
	}
	if s, ok := i.(HumanSource); ok {
		out, err := p.pipe(ctx, s)
		return out, err
	}
	if s, ok := i.(JSONSource); ok {
		out, err := p.pipe(ctx, s)
		return out, err
	}
	return nil, errors.New("Parser requires RawSource")
}

func (p *Parser) pipe(ctx context.Context, s io.Reader) (StructuredSource, error) {
	out, err := p.Parse(s)
	if err != nil {
		return StructuredSource(nil), errors.Wrap(err, "parse")
	}
	return StructuredSource(out), nil
}
