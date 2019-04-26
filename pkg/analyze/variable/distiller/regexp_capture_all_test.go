package distiller

import (
	"reflect"
	"testing"
)

func TestRegexpCaptureAll_Distill(t *testing.T) {
	tests := []struct {
		name             string
		regexpCaptureAll RegexpCaptureAll
		input            string
		want             interface{}
		wantOk           bool
		wantErr          bool
	}{
		{
			name: "num proc",
			regexpCaptureAll: RegexpCaptureAll{
				Regexp: `processor\s+:\s+(\d+)`,
				Index:  1,
			},
			input: `processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 62

processor	: 1
vendor_id	: GenuineIntel
cpu family	: 6
model		: 62
`,
			want:   []string{"0", "1"},
			wantOk: true,
		},
		{
			name: "no match",
			regexpCaptureAll: RegexpCaptureAll{
				Regexp: `NOMATCH`,
				Index:  1,
			},
			input:  "BLAH BLAH",
			want:   nil,
			wantOk: false,
		},
		{
			name: "error",
			regexpCaptureAll: RegexpCaptureAll{
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
			got, ok, err := tt.regexpCaptureAll.Distill(tt.input)
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
