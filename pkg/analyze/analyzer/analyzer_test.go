package analyzer

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	"github.com/replicatedcom/support-bundle/pkg/analyze/condition"
	"github.com/replicatedcom/support-bundle/pkg/analyze/message"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	collectbundle "github.com/replicatedcom/support-bundle/pkg/test-mocks/collect/bundle"
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
				Name: "os",
				Os:   &variable.Os{},
			},
			{
				Name: "osVersion",
				CollectRef: &variable.CollectRef{
					Ref: meta.Ref{
						Selector: meta.Selector{"analyze": "/etc/os-release"},
					},
					RegexpCapture: &distiller.RegexpCapture{
						Regexp: `(?m)^VERSION_ID="([^"]+)"`,
						Index:  1,
					},
				},
			},
		},
		Precondition: &v1.Condition{
			StringCompare: &condition.StringCompare{
				Compare: condition.Compare{Eq: "ubuntu"},
			},
			Ref: "os",
		},
		Condition: &v1.Condition{
			EvalCondition: &osVersionGte1604Eval,
			Ref:           "osVersion",
		},
		Messages: v1.Messages{
			ConditionTrue: &message.Message{
				Primary:  "Ubuntu version is {{repl .osVersion}}",
				Detail:   "Ubuntu version must be at least 16.04",
				Severity: common.SeverityInfo,
			},
			ConditionFalse: &message.Message{
				Primary:  "Ubuntu version is {{repl .osVersion}}",
				Detail:   "Ubuntu version must be at least 16.04",
				Severity: common.SeverityWarn,
			},
			PreconditionFalse: &message.Message{
				Primary:  "OS is not Ubuntu",
				Detail:   "Ubuntu version must be at least 16.04",
				Severity: common.SeverityDebug,
			},
			ConditionError: &message.Message{
				Primary:  "Failed to detect Ubuntu version",
				Detail:   "Ubuntu version must be at least 16.04",
				Severity: common.SeverityError,
			},
			PreconditionError: &message.Message{
				Primary:  "Failed to detect OS",
				Detail:   "Ubuntu version must be at least 16.04",
				Severity: common.SeverityError,
			},
		},
	}

	tests := []struct {
		name            string
		analyzerSpec    v1.Analyzer
		registerExpects func(*collectbundle.MockBundleReader)
		want            api.Result
		wantErr         bool
	}{
		{
			name:         "condition true",
			analyzerSpec: analyzerSpec,
			registerExpects: func(bundleReader *collectbundle.MockBundleReader) {
				bundleReader.
					EXPECT().
					GetIndex().
					Return(collectResultEtcOsRelease)

				bundleReader.
					EXPECT().
					Open("/default/etc/os-release").
					Return(ioutil.NopCloser(strings.NewReader(ubuntu1804EtcOsReleaseContents)), nil)

				bundleReader.
					EXPECT().
					ResultsFromRef(meta.Ref{
						Selector: meta.Selector{"analyze": "/etc/os-release"},
					}).
					Return(collectResultEtcOsRelease)

				bundleReader.
					EXPECT().
					Open("/default/etc/os-release").
					Return(ioutil.NopCloser(strings.NewReader(ubuntu1804EtcOsReleaseContents)), nil)
			},
			want: api.Result{
				Message: &message.Message{
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
			registerExpects: func(bundleReader *collectbundle.MockBundleReader) {
				bundleReader.
					EXPECT().
					GetIndex().
					Return(collectResultEtcOsRelease)

				bundleReader.
					EXPECT().
					Open("/default/etc/os-release").
					Return(ioutil.NopCloser(strings.NewReader(ubuntu1404EtcOsReleaseContents)), nil)

				bundleReader.
					EXPECT().
					ResultsFromRef(meta.Ref{
						Selector: meta.Selector{"analyze": "/etc/os-release"},
					}).
					Return(collectResultEtcOsRelease)

				bundleReader.
					EXPECT().
					Open("/default/etc/os-release").
					Return(ioutil.NopCloser(strings.NewReader(ubuntu1404EtcOsReleaseContents)), nil)
			},
			want: api.Result{
				Message: &message.Message{
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
			registerExpects: func(bundleReader *collectbundle.MockBundleReader) {
				bundleReader.
					EXPECT().
					GetIndex().
					Return(collectResultEtcOsRelease)

				bundleReader.
					EXPECT().
					Open("/default/etc/os-release").
					Return(ioutil.NopCloser(strings.NewReader(centos7EtcOsReleaseContents)), nil)

				bundleReader.
					EXPECT().
					ResultsFromRef(meta.Ref{
						Selector: meta.Selector{"analyze": "/etc/os-release"},
					}).
					Return(collectResultEtcOsRelease)

				bundleReader.
					EXPECT().
					Open("/default/etc/os-release").
					Return(ioutil.NopCloser(strings.NewReader(centos7EtcOsReleaseContents)), nil)
			},
			want: api.Result{
				Message: &message.Message{
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
			bundleReader := collectbundle.NewMockBundleReader(mc)
			defer mc.Finish()

			if tt.registerExpects != nil {
				tt.registerExpects(bundleReader)
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
			tt.want.AnalyzerSpec = api.Analyze{V1: []v1.Analyzer{tt.analyzerSpec}}
			assert.Equal(t, got, tt.want)
		})
	}
}
