package variable

import (
	"reflect"
	"strings"
	"testing"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func TestOs_ExtractValue(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{
			name: "match",
			input: `NAME="Ubuntu"
VERSION="16.04.5 LTS (Xenial Xerus)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 16.04.5 LTS"
VERSION_ID="16.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"
VERSION_CODENAME=xenial
UBUNTU_CODENAME=xenial`,
			want: "ubuntu",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Os{}
			got, err := v.ExtractValue(strings.NewReader(tt.input), collecttypes.Result{
				Path: "/default/etc/os-release",
				Spec: collecttypes.Spec{
					CoreReadFile: &collecttypes.CoreReadFileOptions{
						Filepath: "/etc/os-release",
					},
				},
				Size: 1,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Os.ExtractValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Os.ExtractValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
