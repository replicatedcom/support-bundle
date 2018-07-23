package docker

import (
	"context"
	"regexp"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// this matches server version within error strings like this:
// `Error response from daemon: client is newer than server (client API version: 1.24, server API version: 1.19)`
var dockerErrorVersionRegexp = regexp.MustCompile(`server API version:\s*(\d\.\d+)\s*\)`)

func NewEnvClient(ctx context.Context, logger log.Logger) (docker.CommonAPIClient, error) {
	client, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}

	client.NegotiateAPIVersion(ctx)

	if client.ClientVersion() != "1.24" {
		return client, nil
	}

	// there is a possibility that negotiation failed as this is the default value for that case
	// so we send a ping and check ourselves
	ping, _ := client.Ping(ctx)
	if ping.APIVersion != "" {
		return client, nil
	}

	// negotiation failed, so we get to fake it
	if _, err := client.ServerVersion(ctx); err == nil {
		// ironically, this is actually a bit of a failure
		return client, nil
	}

	matches := dockerErrorVersionRegexp.FindStringSubmatch(err.Error())

	if len(matches) < 2 {
		info := level.Info(log.With(logger, "method", "docker.NewEnvClient"))
		info.Log(
			"version-fallback", client.ClientVersion(),
			"err", "version negotiation failed",
		)
	} else {
		var fakePing types.Ping
		fakePing.APIVersion = matches[1]
		client.NegotiateAPIVersionPing(fakePing)

		debug := level.Info(log.With(logger, "method", "docker.NewEnvClient"))
		debug.Log(
			"version-fallback", client.ClientVersion(),
		)
	}
	return client, nil
}
