package templates

import (
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/docker"
)

func init() {
	RegisterFunc("dockerServerVersionCompare", DockerServerVersionCompare)
}

func DockerServerVersionCompare(a, b string) int {
	aV, err := docker.ParseServerVersion(a)
	if err != nil {
		Panic("dockerServerVersionCompare", errors.Wrap(err, a))
	}
	bV, err := docker.ParseServerVersion(b)
	if err != nil {
		Panic("dockerServerVersionCompare", errors.Wrapf(err, b))
	}
	return aV.Compare(bV)
}
