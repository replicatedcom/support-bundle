package producers

import (
	"context"
	"regexp"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	jww "github.com/spf13/jwalterweatherman"
)

type Docker struct {
	client docker.CommonAPIClient
}

// this matches server version within error strings like this:
// `Error response from daemon: client is newer than server (client API version: 1.24, server API version: 1.19)`
var dockerErrorVersionRegexp = regexp.MustCompile(`server API version:\s*(\d\.\d+)\s*\)`)

func New(client *docker.Client) *Docker {
	client.NegotiateAPIVersion(context.Background())

	if client.ClientVersion() != "1.24" {
		return &Docker{client}
	}

	// there is a possibility that negotiation failed as this is the default value for that case
	// so we send a ping and check ourselves
	ping, _ := client.Ping(context.Background())
	if ping.APIVersion != "" {
		return &Docker{client}
	}

	// negotiation failed, so we get to fake it
	jww.INFO.Printf("Docker API version negotiation failed. Attempting fallback...")
	_, err := client.ServerVersion(context.Background())

	if err == nil {
		// ironically, this is actually a bit of a failure
		return &Docker{client}
	}

	matches := dockerErrorVersionRegexp.FindStringSubmatch(err.Error())

	if len(matches) < 2 {
		jww.INFO.Printf("Docker API version negotiation fallback failed")
	} else {
		jww.INFO.Printf("Fallback API version detection: %+v", matches[1])
		var fakePing types.Ping
		fakePing.APIVersion = matches[1]
		client.NegotiateAPIVersionPing(fakePing)
	}

	return &Docker{client}
}
