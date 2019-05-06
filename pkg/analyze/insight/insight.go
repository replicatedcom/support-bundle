package insight

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api/common"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/replicatedcom/support-bundle/pkg/templates"
)

type Insight struct {
	meta.Meta `json:",inline" yaml:",inline" hcl:",inline"`

	Primary  string          `json:"primary" yaml:"primary" hcl:"primary"`
	Detail   string          `json:"detail" yaml:"detail" hcl:"detail"`
	Severity common.Severity `json:"severity,omitempty" yaml:"severity,omitempty" hcl:"severity,omitempty"`
}

func (m *Insight) Render(data interface{}) (*Insight, error) {
	var multiErr *multierror.Error
	var err error
	built := &Insight{
		Meta:     m.Meta,
		Severity: m.Severity,
	}
	built.Primary, err = templates.String(m.Primary, data)
	multiErr = multierror.Append(multiErr, errors.Wrap(err, "build primary"))
	built.Detail, err = templates.String(m.Detail, data)
	multiErr = multierror.Append(multiErr, errors.Wrap(err, "build detail"))
	return built, multiErr.ErrorOrNil()
}
