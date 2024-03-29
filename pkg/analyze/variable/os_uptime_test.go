package variable

import (
	"reflect"
	"strings"
	"testing"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func TestOsUptime_Distill(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "basic",
			input: "54332.07 207669.53",
			want:  "54332.07",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &OsUptime{}
			got, err := v.DistillReader(strings.NewReader(tt.input), collecttypes.Result{
				Path: "/default/proc/uptime",
				Spec: collecttypes.Spec{
					CoreReadFile: &collecttypes.CoreReadFileOptions{
						Filepath: "/proc/uptime",
					},
				},
				Size: 1,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("OsUptime.Distill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OsUptime.Distill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOsUptime_ExtractValue(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:  "basic",
			input: "54332.07",
			want:  54332.07,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &OsUptime{}
			got, err := v.ExtractValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("OsUptime.ExtractValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OsUptime.ExtractValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
