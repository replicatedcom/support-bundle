package templates

import (
	"strconv"

	"github.com/pkg/errors"
)

func init() {
	RegisterFunc("ParseInt", ParseInt)
}

func ParseInt(str string, args ...string) int {
	base := 10
	bitSize := 64
	i, err := strconv.ParseInt(str, base, bitSize) // TODO: args
	if err != nil {
		Panic("ParseInt", errors.Wrapf(err, "parse %v (base=%d, bitSize=%d)", str, base, bitSize))
	}
	return int(i)
}
