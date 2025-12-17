package templates

import (
	"github.com/blang/semver"
	"github.com/pkg/errors"
)

func init() {
	RegisterFunc("semverCompare", SemverCompare)
}

func SemverCompare(a, b string) int {
	aV, err := semver.Parse(a)
	if err != nil {
		Panic("semverCompare", errors.Wrap(err, a))
	}
	bV, err := semver.Parse(b)
	if err != nil {
		Panic("semverCompare", errors.Wrap(err, b))
	}
	return aV.Compare(bV)
}
