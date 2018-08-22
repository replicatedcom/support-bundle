package templates

import (
	"github.com/Masterminds/sprig"
)

func init() {
	for key, fn := range sprig.FuncMap() {
		RegisterFunc(key, fn)
	}
}
