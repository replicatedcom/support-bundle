package v1

import (
	"reflect"

	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

type Variable struct {
	meta.Meta `json:",inline" yaml:",inline" hcl:",inline"`

	// built-in
	CollectRef *variable.CollectRef `json:"collectRef,omitempty" yaml:"collectRef,omitempty" hcl:"collectRef,omitempty"`
	Eval       *variable.Eval       `json:"eval,omitempty" yaml:"eval,omitempty" hcl:"eval,omitempty"`
	FileMatch  *variable.FileMatch  `json:"fileMatch,omitempty" yaml:"fileMatch,omitempty" hcl:"fileMatch,omitempty"`
	Os         *variable.Os         `json:"os,omitempty" yaml:"os,omitempty" hcl:"os,omitempty"`
	OsUptime   *variable.OsUptime   `json:"os.uptime,omitempty" yaml:"os.uptime,omitempty" hcl:"os.uptime,omitempty"`
}

func (v *Variable) GetName() string {
	return v.Name
}

func (v *Variable) GetVariable() (variable.Interface, string, bool) {
	val := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < val.NumField(); i++ {
		if v, ok := val.Field(i).Interface().(variable.Interface); ok && !reflect.ValueOf(v).IsNil() {
			return v, getTagName(val, i, "yaml"), true
		}
	}
	return nil, "", false
}
