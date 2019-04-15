package variable

import (
	"testing"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func Test_matcherCollector(t *testing.T) {
	tests := []struct {
		name   string
		field  interface{}
		result collecttypes.Result
		want   bool
	}{
		{
			name:   "match",
			field:  &collecttypes.CoreUptimeOptions{},
			result: collecttypes.Result{Size: 1, Spec: collecttypes.Spec{CoreUptime: &collecttypes.CoreUptimeOptions{}}},
			want:   true,
		},
		{
			name:   "no match",
			field:  &collecttypes.CoreUptimeOptions{},
			result: collecttypes.Result{Size: 1, Spec: collecttypes.Spec{CoreHostname: &collecttypes.CoreHostnameOptions{}}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matcherCollector(tt.field)(tt.result); got != tt.want {
				t.Errorf("matcherCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}
