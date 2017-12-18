package pipeline

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"
)

type Humanifier struct {
	Template *template.Template
}

func (p *Humanifier) Pipe(ctx context.Context, i interface{}) (interface{}, error) {
	if s, ok := i.(StructuredSource); ok {
		return p.pipe(ctx, s)
	}
	return nil, errors.New("Humanifier requires StructuredSource")

}

func (p *Humanifier) pipe(ctx context.Context, s StructuredSource) (HumanSource, error) {
	pr, pw := io.Pipe()
	go func() {
		pw.CloseWithError(p.Template.Execute(pw, s))
	}()
	return HumanSource(pr), nil
}
