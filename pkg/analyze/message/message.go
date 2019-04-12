package message

import "github.com/replicatedcom/support-bundle/pkg/analyze/api/common"

type Message struct {
	Primary  string          `json:"primary" yaml:"primary" hcl:"primary"`
	Detail   string          `json:"detail" yaml:"detail" hcl:"detail"`
	Severity common.Severity `json:"severity,omitempty" yaml:"severity,omitempty" hcl:"severity,omitempty"`
}
