package variable

import (
	"reflect"
	"testing"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func TestFileMatch_MatchResults(t *testing.T) {
	tests := []struct {
		name        string
		pathRegexps []string
		index       []collecttypes.Result
		want        []collecttypes.Result
		wantErr     bool
	}{
		{
			name:        "match",
			pathRegexps: []string{`/docker_info\.json`},
			index: []collecttypes.Result{
				{
					Path: "default/docker/docker_info.json",
					Size: 1,
					Spec: collecttypes.Spec{DockerInfo: &collecttypes.DockerInfoOptions{
						SpecShared: collecttypes.SpecShared{
							OutputDir: "default/docker",
						},
					},
					},
				},
			},
			want: []collecttypes.Result{
				{
					Path: "default/docker/docker_info.json",
					Size: 1,
					Spec: collecttypes.Spec{DockerInfo: &collecttypes.DockerInfoOptions{
						SpecShared: collecttypes.SpecShared{
							OutputDir: "default/docker",
						},
					},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &FileMatch{
				PathRegexps: tt.pathRegexps,
			}
			got, err := v.MatchResults(tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileMatch.MatchResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileMatch.MatchResults() = %v, want %v", got, tt.want)
			}
		})
	}
}
