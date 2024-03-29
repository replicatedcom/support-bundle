package distiller

import (
	"reflect"
	"testing"
)

func TestRegexpCapture_Distill(t *testing.T) {
	tests := []struct {
		name          string
		regexpCapture RegexpCapture
		input         string
		want          interface{}
		wantOk        bool
		wantErr       bool
	}{
		{
			name: "os release version",
			regexpCapture: RegexpCapture{
				Regexp: `(?m)^VERSION_ID="([^"]+)"`,
				Index:  1,
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
			want:   "16.04",
			wantOk: true,
		},
		{
			name: "no match",
			regexpCapture: RegexpCapture{
				Regexp: `NOMATCH`,
				Index:  1,
			},
			input:  "BLAH BLAH",
			want:   nil,
			wantOk: false,
		},
		{
			name: "error",
			regexpCapture: RegexpCapture{
				Regexp: `(`,
				Index:  1,
			},
			input:   "BLAH BLAH",
			wantOk:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := tt.regexpCapture.Distill(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegexpCapture.Distill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegexpCapture.Distill() = %v, want %v", got, tt.want)
			}
			if ok != tt.wantOk {
				t.Errorf("RegexpCaptureAll.Distill() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}
