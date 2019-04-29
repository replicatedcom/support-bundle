package render

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/analyze/insight"
	"github.com/replicatedcom/support-bundle/pkg/ui"
)

func TestHumanEncoder_Encode(t *testing.T) {
	tests := []struct {
		name    string
		results []api.Result
		want    string
		wantErr bool
	}{
		{
			name: "basic",
			results: []api.Result{
				{
					Insight: &insight.Insight{
						Primary:  "Ubuntu version is 14.04",
						Detail:   "Ubuntu version must be at least 16.04",
						Severity: common.SeverityWarn,
					},
					Severity: common.SeverityWarn,
				},
				{
					Insight: &insight.Insight{
						Primary:  "Ubuntu version is 18.04",
						Detail:   "Ubuntu version must be at least 16.04",
						Severity: common.SeverityInfo,
					},
					Severity: common.SeverityInfo,
				},
				{
					Insight: &insight.Insight{
						Primary:  "Ubuntu version is ",
						Detail:   "Ubuntu version must be at least 16.04",
						Severity: common.SeverityInfo,
					},
					Severity: common.SeverityError,
					Error:    "Panic!!!",
				},
			},
			want: `! Ubuntu version is 14.04
Ubuntu version must be at least 16.04

✓ Ubuntu version is 18.04
Ubuntu version must be at least 16.04

✗ Panic!!!
Ubuntu version must be at least 16.04

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			ui := ui.New(nil, &buf, &buf)
			e := &HumanEncoder{UI: ui}
			if err := e.Encode(tt.results); (err != nil) != tt.wantErr {
				t.Errorf("HumanEncoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := buf.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
