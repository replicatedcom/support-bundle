package producers

import (
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

type Docker struct {
	client *docker.Client
}

func New(client *docker.Client) *Docker {
	client.NegotiateAPIVersion(context.Background())

	if client.ClientVersion() == "1.24" {
		// there is a possibility that negotiation failed as this is the default value for that case
		// so we send a ping and check ourselves
		ping, _ := client.Ping(context.Background())
		if ping.APIVersion == "" {
			// negotiation failed, so we get to fake it
			log.Printf("Docker API version negotiation failed. Attempting fallback...")
			_, err := client.ServerVersion(context.Background())

			if err != nil {
				// this matches server version within error strings like this:
				// `Error response from daemon: client is newer than server (client API version: 1.24, server API version: 1.19)`
				r := regexp.MustCompile(`server API version:\s*(\d\.\d+)\s*\)`)
				matches := r.FindStringSubmatch(err.Error())

				if len(matches) < 2 {
					log.Printf("Docker API version negotiation fallback failed")
				} else {
					log.Printf("Fallback API version detection: %+v", matches[1])
					var fakePing types.Ping
					fakePing.APIVersion = matches[1]
					client.NegotiateAPIVersionPing(fakePing)
				}
			}
		}
	}

	return &Docker{client}
}

// get the ID of a container given the name
func (d *Docker) getContainerID(ctx context.Context, containerName string) (string, error) {
	matchingID := ""

	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return matchingID, err
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == containerName || name == "/"+containerName {
				matchingID = container.ID
				return matchingID, nil
			}
		}
	}

	err = errors.New("unable to find container with name " + containerName)
	return matchingID, err
}
