package variable

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func TestOs_Distill(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		filepath string
		want     string
		wantErr  bool
	}{
		{
			name: "ubuntu",
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
			filepath: "/etc/os-release",
			want:     "ubuntu",
		},
		{
			name:     "rhel",
			input:    `Red Hat Enterprise Linux Server release 6.10 (Santiago)`,
			filepath: "/etc/system-release",
			want:     "rhel",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Os{}
			got, err := v.DistillReader(strings.NewReader(tt.input), collecttypes.Result{
				Path: filepath.Join("/default", tt.filepath),
				Spec: collecttypes.Spec{
					CoreReadFile: &collecttypes.CoreReadFileOptions{
						Filepath: tt.filepath,
					},
				},
				Size: 1,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Os.DistillReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Os.DistillReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
