package analyzer

import (
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	collectbundle "github.com/replicatedcom/support-bundle/pkg/test-mocks/collect/bundle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzer_analyze(t *testing.T) {
	type args struct {
		analyzerSpec v1.Analyzer
	}
	tests := []struct {
		name         string
		args         args
		returnReader io.ReadCloser
		want         api.Result
		wantErr      bool
	}{
		{
			name: "basic",
			args: args{
				analyzerSpec: v1.Analyzer{
					KubernetesTotalMemory: &v1.KubernetesTotalMemoryRequirement{
						Min: "10Gi",
					},
					AnalyzerShared: v1.AnalyzerShared{
						CollectRefs: []meta.Ref{
							{
								Selector: meta.Selector{
									"analyze": "kubernetes-total-memory",
								},
							},
						},
					},
				},
			},
			returnReader: ioutil.NopCloser(strings.NewReader(`{"items": [
				{"status": {"capacity": {"memory": "3000000Ki"},"allocatable": {"memory": "2708864Ki"}}},
				{"status": {"capacity": {"memory": "3000000Ki"},"allocatable": {"memory": "2708864Ki"}}}
			]}`)),
			want: api.Result{
				Requirement: "Kubernetes total memory must be at least 10GiB",
				Message:     "Total memory for your Kubernetes cluster 6000000KiB is less than the minimum required memory of 10GiB",
				Severity:    common.SeverityError,
				Vars: []map[string]interface{}{
					{"Empty": "false"},
					{"Memory": "6144000000", "Min": "10737418240"},
				},
				Error: "",
			},
		},
		{
			name: "basic",
			args: args{
				analyzerSpec: v1.AnalyzerSpec{
					KubernetesTotalMemory: &v1.KubernetesTotalMemoryRequirement{
						Min: "1Gi",
					},
					AnalyzerShared: v1.AnalyzerShared{
						CollectRefs: []meta.Ref{
							{
								Selector: meta.Selector{
									"analyze": "kubernetes-total-memory",
								},
							},
						},
					},
				},
			},
			returnReader: ioutil.NopCloser(strings.NewReader(`{"items": [
				{"status": {"capacity": {"memory": "3000000Ki"},"allocatable": {"memory": "2708864Ki"}}},
				{"status": {"capacity": {"memory": "3000000Ki"},"allocatable": {"memory": "2708864Ki"}}}
			]}`)),
			want: api.Result{
				Requirement: "Kubernetes total memory must be at least 1GiB",
				Vars: []map[string]interface{}{
					{"Empty": "false"},
					{"Memory": "6144000000", "Min": "1073741824"},
				},
				Error: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			bundleReader := collectbundle.NewMockBundleReader(mc)
			defer mc.Finish()

			// TODO: support for multiple refs
			ref := tt.args.analyzerSpec.CollectRefs[0]
			if ref.Ref == "" {
				ref.Ref = "_Ref"
			}

			bundleReader.
				EXPECT().
				ReaderFromRef(ref).
				Return(tt.returnReader, nil)

			a := &Analyzer{
				Logger: log.NewNopLogger(),
			}
			got, err := a.analyze(context.Background(), bundleReader, tt.args.analyzerSpec)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			tt.want.AnalyzerSpec = api.Analyze{V1: []v1.AnalyzerSpec{tt.args.analyzerSpec}}
			assert.Equal(t, got, tt.want)
		})
	}
}
