package variable

import (
	"testing"

	"github.com/stretchr/testify/assert"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func Test_matcherCollector(t *testing.T) {
	tests := []struct {
		name    string
		field   interface{}
		results []collecttypes.Result
		want    []collecttypes.Result
	}{
		{
			name:    "match",
			field:   &collecttypes.CoreUptimeOptions{},
			results: []collecttypes.Result{{Size: 1, Spec: collecttypes.Spec{CoreUptime: &collecttypes.CoreUptimeOptions{}}}},
			want:    []collecttypes.Result{{Size: 1, Spec: collecttypes.Spec{CoreUptime: &collecttypes.CoreUptimeOptions{}}}},
		},
		{
			name:    "no match",
			field:   &collecttypes.CoreUptimeOptions{},
			results: []collecttypes.Result{{Size: 1, Spec: collecttypes.Spec{CoreHostname: &collecttypes.CoreHostnameOptions{}}}},
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchCollector(tt.field)(tt.results)
			assert.Equal(t, tt.want, got)
		})
	}
}
