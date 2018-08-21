package templates

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/docker"
)

func init() {
	RegisterFunc("DockerServerVersionCompare", DockerServerVersionCompare)
}

func DockerServerVersionCompare(a, b string) int {
	aV, err := docker.ParseServerVersion(a)
	if err != nil {
		Panic("DockerServerVersionCompare", errors.Wrap(err, a))
	}
	bV, err := docker.ParseServerVersion(b)
	if err != nil {
		Panic("DockerServerVersionCompare", errors.Wrapf(err, b))
	}
	return aV.Compare(bV)
}
