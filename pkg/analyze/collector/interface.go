package collector

import (
	"context"
	"time"
)

type Interface interface {
	CollectBundle(
		ctx context.Context,
		customerID string,
		specs []string,
		specFiles []string,
		dest string,
		opts Options,
	) error
}

type Options struct {
	EnableCore       bool
	EnableDocker     bool
	EnableJournald   bool
	EnableKubernetes bool
	EnableRetraced   bool
	Timeout          time.Duration
	CustomerEndpoint string
}
