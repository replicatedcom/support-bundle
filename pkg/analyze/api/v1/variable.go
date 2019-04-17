package v1

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

type Variable struct {
	meta.Meta `json:",inline" yaml:",inline" hcl:",inline"`

	// built-in
	CollectRef *variable.CollectRef `json:"collectRef,omitempty" yaml:"collectRef,omitempty" hcl:"collectRef,omitempty"`
	FileMatch  *variable.FileMatch  `json:"fileMatch,omitempty" yaml:"fileMatch,omitempty" hcl:"fileMatch,omitempty"`
	Os         *variable.Os         `json:"os,omitempty" yaml:"os,omitempty" hcl:"os,omitempty"`
	OsUptime   *variable.OsUptime   `json:"os.uptime,omitempty" yaml:"os.uptime,omitempty" hcl:"os.uptime,omitempty"`
}

func (v *Variable) Register(bundleReader bundlereader.BundleReader, data interface{}) (map[string]interface{}, error) {
	var value interface{}
	val := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < val.NumField(); i++ {
		if v, ok := val.Field(i).Interface().(variable.Interface); ok && !reflect.ValueOf(v).IsNil() {
			var err error
			value, err = variable.Extract(v, bundleReader, data)
			if err != nil {
				return nil, errors.Wrapf(err, "variable %q", getTagName(val, i, "yaml"))
			} else if value != nil {
				break
			}
		}
	}
	return map[string]interface{}{v.Name: value}, nil
}
