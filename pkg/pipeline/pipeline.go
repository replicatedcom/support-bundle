package pipeline

import (
	"context"
)

type Pipe func(ctx context.Context, i interface{}) (interface{}, error)

type Banzai struct {
	Pipeline []Pipe
}

func (b Banzai) Pipe(ctx context.Context, s interface{}) (interface{}, error) {
	var err error
	for _, pipe := range b.Pipeline {
		s, err = pipe(ctx, s)
		if err != nil {
			break
		}
	}
	return s, err
}
