package variable

import (
	"reflect"
	"strings"
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
)

func TestDistiller_Distill(t *testing.T) {
	evalDistiller := distiller.Eval(`{{repl jq . ".key" | jqMessagesToString}}`)

	tests := []struct {
		name      string
		distiller Distiller
		input     string
		want      interface{}
		wantErr   bool
	}{
		{
			name: "match",
			distiller: Distiller{
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
			distiller: Distiller{
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
			distiller: Distiller{
				RegexpCapture: &distiller.RegexpCapture{
					Regexp: "NOMATCH",
				},
			},
			input: "BLAH BLAH",
			want:  nil,
		},
		{
			name: "no match scannable",
			distiller: Distiller{
				RegexpCapture: &distiller.RegexpCapture{
					Regexp: "NOMATCH",
				},
				Scannable: true,
			},
			input: "BLAH BLAH",
			want:  nil,
		},
		{
			name: "eval",
			distiller: Distiller{
				Eval:      &evalDistiller,
				Scannable: true,
			},
			input: `{"key": "val"}`,
			want:  "val",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := tt.distiller.Distill(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Distiller.Distill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Distiller.Distill() = %v, want %v", got, tt.want)
			}
		})
	}
}
