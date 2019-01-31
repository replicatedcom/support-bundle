package collector

import (
	"context"
	"time"
)

type Interface interface {
	CollectBundle(ctx context.Context, input CollectorInput) error
}

type Options struct {
	EnableCore       bool
	EnableDocker     bool
	EnableJournald   bool
	EnableKubernetes bool
	EnableRetraced   bool
	Timeout          time.Duration
	Endpoint         string
}
