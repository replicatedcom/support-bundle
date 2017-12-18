package pipeline

import "io"

type RawSource io.Reader

type HumanSource io.Reader

type JSONSource io.Reader

type StructuredSource interface{}
