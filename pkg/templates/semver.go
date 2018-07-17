package templates

import (
	"github.com/blang/semver"
	"github.com/pkg/errors"
)

func init() {
	RegisterFunc("SemverCompare", SemverCompare)
}

func SemverCompare(a, b string) int {
	aV, err := semver.Parse(a)
	if err != nil {
		Panic("SemverCompare", errors.Wrap(err, a))
	}
	bV, err := semver.Parse(b)
	if err != nil {
		Panic("SemverCompare", errors.Wrapf(err, b))
	}
	return aV.Compare(bV)
}
