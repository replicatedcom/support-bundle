package analyzer

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	"github.com/replicatedcom/support-bundle/pkg/analyze/condition"
	"github.com/replicatedcom/support-bundle/pkg/analyze/insight"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	mockbundlereader "github.com/replicatedcom/support-bundle/pkg/test-mocks/collect/bundle/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzer_analyze(t *testing.T) {
	osVersionGte1604Eval := condition.EvalCondition(`{{repl lt .osVersion "16.04" | not}}`)

	collectResultEtcOsRelease := []collecttypes.Result{
		{
			Path: "/default/etc/os-release",
			Spec: collecttypes.Spec{
				CoreReadFile: &collecttypes.CoreReadFileOptions{
					Filepath: "/etc/os-release",
				},
				SpecShared: collecttypes.SpecShared{
					Meta: &meta.Meta{
						Labels: map[string]string{
							"analyze": "etc.os-release",
						},
					},
				},
			},
			Size: 1,
		},
	}

	ubuntu1404EtcOsReleaseContents := `NAME="Ubuntu"
VERSION="14.04.6 LTS, Trusty Tahr"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 14.04.6 LTS"
VERSION_ID="14.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"`

	ubuntu1804EtcOsReleaseContents := `NAME="Ubuntu"
VERSION="18.04.2 LTS (Bionic Beaver)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 18.04.2 LTS"
VERSION_ID="18.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=bionic
UBUNTU_CODENAME=bionic`

	centos7EtcOsReleaseContents := `NAME="CentOS Linux"
VERSION="7 (Core)"
ID="centos"
ID_LIKE="rhel fedora"
VERSION_ID="7"
PRETTY_NAME="CentOS Linux 7 (Core)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:centos:centos:7"
HOME_URL="https://www.centos.org/"
BUG_REPORT_URL="https://bugs.centos.org/"

CENTOS_MANTISBT_PROJECT="CentOS-7"
CENTOS_MANTISBT_PROJECT_VERSION="7"
REDHAT_SUPPORT_PRODUCT="centos"
REDHAT_SUPPORT_PRODUCT_VERSION="7"
`

	analyzerSpec := v1.Analyzer{
		RegisterVariables: []v1.Variable{
			{
				Meta: meta.Meta{
					Name: "os",
				},
				Os: &variable.Os{},
			},
			{
				Meta: meta.Meta{
					Name: "osVersion",
				},
				CollectRef: &variable.CollectRef{
					Ref: meta.Ref{
						Selector: meta.Selector{"analyze": "etc.os-release"},
					},
					Distiller: variable.Distiller{
						RegexpCapture: &distiller.RegexpCapture{
							Regexp: `(?m)^VERSION_ID="([^"]+)"`,
							Index:  1,
						},
					},
				},
			},
		},
		EvaluateConditions: []v1.EvaluateCondition{
			{
				Condition: v1.Condition{
					StringCompare: &condition.StringCompare{
						Compare: condition.Compare{Eq: "ubuntu"},
					},
					VariableRef: "os",
				},
				InsightOnError: &insight.Insight{
					Primary:  "Failed to detect OS",
					Detail:   "Ubuntu version must be at least 16.04",
					Severity: common.SeverityError,
				},
				InsightOnFalse: &insight.Insight{
					Primary:  "OS is not Ubuntu",
					Detail:   "Ubuntu version must be at least 16.04",
					Severity: common.SeverityDebug,
				},
			},
			{
				Condition: v1.Condition{
					EvalCondition: &osVersionGte1604Eval,
					VariableRef:   "osVersion",
				},
				InsightOnError: &insight.Insight{
					Primary:  "Failed to detect Ubuntu version",
					Detail:   "Ubuntu version must be at least 16.04",
					Severity: common.SeverityError,
				},
				InsightOnFalse: &insight.Insight{
					Primary:  "Ubuntu version is {{repl .osVersion}}",
					Detail:   "Ubuntu version must be at least 16.04",
					Severity: common.SeverityWarn,
				},
			},
		},
		Insight: &insight.Insight{
			Primary:  "Ubuntu version is {{repl .osVersion}}",
			Detail:   "Ubuntu version must be at least 16.04",
			Severity: common.SeverityInfo,
		},
	}

	tests := []struct {
		name            string
		analyzerSpec    v1.Analyzer
		registerExpects func(*mockbundlereader.MockBundleReader, *gomock.Controller)
		want            *api.Result
		wantErr         bool
	}{
		{
			name:         "condition true",
			analyzerSpec: analyzerSpec,
			registerExpects: func(bundleReader *mockbundlereader.MockBundleReader, mc *gomock.Controller) {
				bundleReader.
					EXPECT().
					GetIndex().
					Return(collectResultEtcOsRelease)

				scanner := mockbundlereader.NewMockScanner(mc)

				scanner.
					EXPECT().
					Next().
					Return(&bundlereader.ScannerFile{
						Name:   "/blah/blah",
						Reader: strings.NewReader("blah blah"),
					}, nil)

				scanner.
					EXPECT().
					Next().
					Return(&bundlereader.ScannerFile{
						Name:   "/default/etc/os-release",
						Reader: strings.NewReader(ubuntu1804EtcOsReleaseContents),
					}, nil)

				scanner.
					EXPECT().
					Next().
					Return(&bundlereader.ScannerFile{
						Name:   "/blah/blah/blah",
						Reader: strings.NewReader("blah blah blah"),
					}, nil)

				scanner.
					EXPECT().
					Next().
					Return(nil, io.EOF)

				bundleReader.
					EXPECT().
					NewScanner().
					Return(scanner, nil)

				scanner.
					EXPECT().
					Close().
					Return(nil)
			},
			want: &api.Result{
				Insight: &insight.Insight{
					Primary:  "Ubuntu version is 18.04",
					Detail:   "Ubuntu version must be at least 16.04",
					Severity: common.SeverityInfo,
				},
				Severity: common.SeverityInfo,
				Variables: map[string]interface{}{
					"os":        "ubuntu",
					"osVersion": "18.04",
				},
				Error: "",
			},
		},
		{
			name:         "condition false",
			analyzerSpec: analyzerSpec,
			registerExpects: func(bundleReader *mockbundlereader.MockBundleReader, mc *gomock.Controller) {
				bundleReader.
					EXPECT().
					GetIndex().
					Return(collectResultEtcOsRelease)

				scanner := mockbundlereader.NewMockScanner(mc)

				scanner.
					EXPECT().
					Next().
					Return(&bundlereader.ScannerFile{
						Name:   "/default/etc/os-release",
						Reader: strings.NewReader(ubuntu1404EtcOsReleaseContents),
					}, nil)

				scanner.
					EXPECT().
					Next().
					Return(nil, io.EOF)

				bundleReader.
					EXPECT().
					NewScanner().
					Return(scanner, nil)

				scanner.
					EXPECT().
					Close().
					Return(nil)
			},
			want: &api.Result{
				Insight: &insight.Insight{
					Primary:  "Ubuntu version is 14.04",
					Detail:   "Ubuntu version must be at least 16.04",
					Severity: common.SeverityWarn,
				},
				Severity: common.SeverityWarn,
				Variables: map[string]interface{}{
					"os":        "ubuntu",
					"osVersion": "14.04",
				},
				Error: "",
			},
		},
		{
			name:         "precondition false",
			analyzerSpec: analyzerSpec,
			registerExpects: func(bundleReader *mockbundlereader.MockBundleReader, mc *gomock.Controller) {
				bundleReader.
					EXPECT().
					GetIndex().
					Return(collectResultEtcOsRelease)

				scanner := mockbundlereader.NewMockScanner(mc)

				scanner.
					EXPECT().
					Next().
					Return(&bundlereader.ScannerFile{
						Name:   "/default/etc/os-release",
						Reader: strings.NewReader(centos7EtcOsReleaseContents),
					}, nil)

				scanner.
					EXPECT().
					Next().
					Return(nil, io.EOF)

				bundleReader.
					EXPECT().
					NewScanner().
					Return(scanner, nil)

				scanner.
					EXPECT().
					Close().
					Return(nil)
			},
			want: &api.Result{
				Insight: &insight.Insight{
					Primary:  "OS is not Ubuntu",
					Detail:   "Ubuntu version must be at least 16.04",
					Severity: common.SeverityDebug,
				},
				Severity: common.SeverityDebug,
				Variables: map[string]interface{}{
					"os":        "centos",
					"osVersion": "7",
				},
				Error: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			bundleReader := mockbundlereader.NewMockBundleReader(mc)
			defer mc.Finish()

			if tt.registerExpects != nil {
				tt.registerExpects(bundleReader, mc)
			}

			a := &Analyzer{
				Logger: log.NewNopLogger(),
			}
			got, err := a.analyze(context.Background(), bundleReader, tt.analyzerSpec)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			b, _ := json.Marshal(api.Analyze{V1: []v1.Analyzer{tt.analyzerSpec}})
			tt.want.AnalyzerSpec = string(b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_resultFromAnalysis(t *testing.T) {
	type args struct {
		insight      *insight.Insight
		analysisErr  error
		analyzerSpec v1.Analyzer
		data         map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantResult *api.Result
		wantErr    bool
	}{
		{
			name: "insight with label override",
			args: args{
				insight: &insight.Insight{
					Meta: meta.Meta{
						Name: "the warning insight",
						Labels: map[string]string{
							"iconKey": "oh_no",
						},
					},
					Primary:  "{{repl .someVar}} primary",
					Detail:   "{{repl .someVar}} detail",
					Severity: common.SeverityWarn,
				},
				analysisErr: nil,
				analyzerSpec: v1.Analyzer{
					Meta: meta.Meta{
						Name: "the analyzer",
						Labels: map[string]string{
							"desiredPosition": "1",
							"iconKey":         "oh_yes",
						},
					},
				},
				data: map[string]interface{}{
					"someVar": "THE VALUE",
				},
			},
			wantResult: &api.Result{
				Meta: meta.Meta{
					Name: "the analyzer",
					Labels: map[string]string{
						"desiredPosition": "1",
						"iconKey":         "oh_no",
					},
				},
				Insight: &insight.Insight{
					Meta: meta.Meta{
						Name: "the warning insight",
						Labels: map[string]string{
							"iconKey": "oh_no",
						},
					},
					Primary:  "THE VALUE primary",
					Detail:   "THE VALUE detail",
					Severity: common.SeverityWarn,
				},
				Severity: common.SeverityWarn,
				Variables: map[string]interface{}{
					"someVar": "THE VALUE",
				},
				Error: "",
			},
		},
		{
			name: "error",
			args: args{
				insight:     nil,
				analysisErr: errors.New("THIS IS THE ERROR"),
				analyzerSpec: v1.Analyzer{
					Meta: meta.Meta{
						Name: "the analyzer",
						Labels: map[string]string{
							"desiredPosition": "1",
							"iconKey":         "oh_yes",
						},
					},
				},
				data: map[string]interface{}{
					"someVar": "THE VALUE",
				},
			},
			wantResult: &api.Result{
				Meta: meta.Meta{
					Name: "the analyzer",
					Labels: map[string]string{
						"desiredPosition": "1",
						"iconKey":         "oh_yes",
					},
				},
				Insight:  nil,
				Severity: common.SeverityError,
				Variables: map[string]interface{}{
					"someVar": "THE VALUE",
				},
				Error: "THIS IS THE ERROR",
			},
			wantErr: true,
		},
		{
			name: "insight and error",
			args: args{
				insight: &insight.Insight{
					Meta: meta.Meta{
						Name: "the warning insight",
						Labels: map[string]string{
							"iconKey": "oh_no",
						},
					},
					Primary:  "{{repl .someVar}} primary",
					Detail:   "{{repl .someVar}} detail",
					Severity: common.SeverityWarn,
				},
				analysisErr: errors.New("THIS IS THE ERROR"),
				analyzerSpec: v1.Analyzer{
					Meta: meta.Meta{
						Name: "the analyzer",
						Labels: map[string]string{
							"desiredPosition": "1",
							"iconKey":         "oh_yes",
						},
					},
				},
				data: map[string]interface{}{
					"someVar": "THE VALUE",
				},
			},
			wantResult: &api.Result{
				Meta: meta.Meta{
					Name: "the analyzer",
					Labels: map[string]string{
						"desiredPosition": "1",
						"iconKey":         "oh_no",
					},
				},
				Insight: &insight.Insight{
					Meta: meta.Meta{
						Name: "the warning insight",
						Labels: map[string]string{
							"iconKey": "oh_no",
						},
					},
					Primary:  "THE VALUE primary",
					Detail:   "THE VALUE detail",
					Severity: common.SeverityWarn,
				},
				Severity: common.SeverityWarn,
				Variables: map[string]interface{}{
					"someVar": "THE VALUE",
				},
				Error: "THIS IS THE ERROR",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := resultFromAnalysis(tt.args.insight, tt.args.analysisErr, tt.args.analyzerSpec, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("resultFromAnalysis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			b, _ := json.Marshal(api.Analyze{V1: []v1.Analyzer{tt.args.analyzerSpec}})
			tt.wantResult.AnalyzerSpec = string(b)
			assert.Equal(t, tt.wantResult, gotResult)
		})
	}
}
