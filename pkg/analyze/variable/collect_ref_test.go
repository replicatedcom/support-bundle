package variable

import (
	"reflect"
	"strings"
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
)

func TestCollectRef_ExtractValue(t *testing.T) {
	tests := []struct {
		name       string
		collectRef CollectRef
		input      string
		want       interface{}
		wantErr    bool
	}{
		{
			name: "match",
			collectRef: CollectRef{
				RegexpCapture: &distiller.RegexpCapture{
					Regexp: `(?m)^VERSION_ID="([^"]+)"`,
					Index:  1,
				},
			},
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
			want: "16.04",
		},
		{
			name: "match scannable",
			collectRef: CollectRef{
				RegexpCapture: &distiller.RegexpCapture{
					Regexp: `(?m)^VERSION_ID="([^"]+)"`,
					Index:  0,
				},
				Scannable: true,
			},
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
			want: "VERSION_ID=\"16.04\"",
		},
		{
			name: "no match",
			collectRef: CollectRef{
				RegexpCapture: &distiller.RegexpCapture{
					Regexp: `NOMATCH`,
				},
			},
			input: `BLAH BLAH`,
			want:  nil,
		},
		{
			name: "no match scannable",
			collectRef: CollectRef{
				RegexpCapture: &distiller.RegexpCapture{
					Regexp: `NOMATCH`,
				},
				Scannable: true,
			},
			input: `BLAH BLAH`,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := tt.collectRef.ExtractValue(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectRef.ExtractValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectRef.ExtractValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
