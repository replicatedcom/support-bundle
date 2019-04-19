package variable

import (
	"io"
	"reflect"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	"github.com/replicatedcom/support-bundle/pkg/util"
)

var identity = &distiller.Identity{}

type Distiller struct {
	Scannable bool `json:"scannable,omitempty" yaml:"scannable,omitempty" hcl:"scannable,omitempty"`

	RegexpCapture *distiller.RegexpCapture `json:"regexpCapture,omitempty" yaml:"regexpCapture,omitempty" hcl:"regexpCapture,omitempty"`
	Identity      *distiller.Identity      `json:"identity,omitempty" yaml:"identity,omitempty" hcl:"identity,omitempty"`
}

func (v *Distiller) Distill(r io.Reader) (string, error) {
	val := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < val.NumField(); i++ {
		if d, ok := val.Field(i).Interface().(distiller.Interface); ok && !reflect.ValueOf(d).IsNil() {
			str, err := distiller.Distill(d, r, v.Scannable)
			return str, errors.Wrapf(err, "distill %q", util.StructTagName(val, i, "yaml"))
		}
	}
	// return the whole string if there is no distiller
	str, err := distiller.Distill(identity, r, v.Scannable)
	return str, errors.Wrap(err, "distill identity")
}
