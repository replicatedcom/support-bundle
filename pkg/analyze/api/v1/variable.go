package v1

import (
	"reflect"

	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"
	"github.com/replicatedcom/support-bundle/pkg/meta"
)

type Variable struct {
	meta.Meta `json:",inline" yaml:",inline" hcl:",inline"`

	// built-in
	CollectRef           *variable.CollectRef           `json:"collectRef,omitempty" yaml:"collectRef,omitempty" hcl:"collectRef,omitempty"`
	CoreReadFilePath     *variable.CoreReadFilePath     `json:"coreReadFilePath,omitempty" yaml:"coreReadFilePath,omitempty" hcl:"coreReadFilePath,omitempty"`
	CPUCores             *variable.CPUCores             `json:"cpuCores,omitempty" yaml:"cpuCores,omitempty" hcl:"cpuCores,omitempty"`
	Eval                 *variable.Eval                 `json:"eval,omitempty" yaml:"eval,omitempty" hcl:"eval,omitempty"`
	FileMatch            *variable.FileMatch            `json:"fileMatch,omitempty" yaml:"fileMatch,omitempty" hcl:"fileMatch,omitempty"`
	Loadavg              *variable.Loadavg              `json:"loadavg,omitempty" yaml:"loadavg,omitempty" hcl:"loadavg,omitempty"`
	Loadavg1             *variable.Loadavg1             `json:"loadavg1,omitempty" yaml:"loadavg1,omitempty" hcl:"loadavg1,omitempty"`
	Loadavg5             *variable.Loadavg5             `json:"loadavg5,omitempty" yaml:"loadavg5,omitempty" hcl:"loadavg5,omitempty"`
	Loadavg15            *variable.Loadavg15            `json:"loadavg15,omitempty" yaml:"loadavg15,omitempty" hcl:"loadavg15,omitempty"`
	MemoryUsageAvailable *variable.MemoryUsageAvailable `json:"memoryUsageAvailable,omitempty" yaml:"memoryUsageAvailable,omitempty" hcl:"memoryUsageAvailable,omitempty"`
	MemoryUsageTotal     *variable.MemoryUsageTotal     `json:"memoryUsageTotal,omitempty" yaml:"memoryUsageTotal,omitempty" hcl:"memoryUsageTotal,omitempty"`
	Os                   *variable.Os                   `json:"os,omitempty" yaml:"os,omitempty" hcl:"os,omitempty"`
	OsUptime             *variable.OsUptime             `json:"osUptime,omitempty" yaml:"osUptime,omitempty" hcl:"osUptime,omitempty"`
	WhichCloud           *variable.WhichCloud           `json:"whichCloud,omitempty" yaml:"whichCloud,omitempty" hcl:"whichCloud,omitempty"`
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
