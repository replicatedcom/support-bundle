package types

import (
	"context"
	"io"
)

type BytesProducer func(context.Context) ([]byte, error)
type StructuredProducer func(context.Context) (interface{}, error)
type StreamProducer func(context.Context) (io.Reader, error)

// probably stdout and stderr
type StreamsProducer func(context.Context) (map[string]io.Reader, error)

type Planner func(Spec) []Task

type Plugin map[string]Planner

type BytesScrubber func([]byte) []byte
