package types

import (
	"context"
	"io"
	"io/ioutil"
)

type BytesProducer func(context.Context) ([]byte, error)
type StructuredProducer func(context.Context) (interface{}, error)
type StreamProducer func(context.Context) (io.Reader, error)

// probably stdout and stderr
type StreamsProducer func(context.Context) (map[string]io.Reader, error)

type Planner func(Spec) []Task

type Plugin map[string]Planner

type BytesScrubber func([]byte) []byte

func BytesProducerFromStreamProducer(ss StreamProducer) BytesProducer {
	return func(ctx context.Context) ([]byte, error) {
		r, err := ss(ctx)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(r)
	}
}
