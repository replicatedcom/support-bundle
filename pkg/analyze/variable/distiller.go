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

	Eval             *distiller.Eval             `json:"eval,omitempty" yaml:"eval,omitempty" hcl:"eval,omitempty"`
	Identity         *distiller.Identity         `json:"identity,omitempty" yaml:"identity,omitempty" hcl:"identity,omitempty"`
	RegexpCapture    *distiller.RegexpCapture    `json:"regexpCapture,omitempty" yaml:"regexpCapture,omitempty" hcl:"regexpCapture,omitempty"`
	RegexpCaptureAll *distiller.RegexpCaptureAll `json:"regexpCaptureAll,omitempty" yaml:"regexpCaptureAll,omitempty" hcl:"regexpCaptureAll,omitempty"`
}

func (v *Distiller) Distill(r io.Reader) (interface{}, error) {
	val := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < val.NumField(); i++ {
		if d, ok := val.Field(i).Interface().(distiller.Interface); ok && !reflect.ValueOf(d).IsNil() {
			str, err := distiller.Distill(d, r, v.Scannable)
			return str, errors.Wrapf(err, "distill %q", util.StructTagName(val, i, "yaml"))
		}
	}
	// return the whole string if there is no distiller
	distilled, err := distiller.Distill(identity, r, v.Scannable)
	return distilled, errors.Wrap(err, "distill identity")
}
