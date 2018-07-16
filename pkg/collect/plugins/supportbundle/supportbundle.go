package supportbundle

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/supportbundle/planners"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/supportbundle/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type SupportBundle struct {
	planner *planners.SupportBundle
}

func New() (*SupportBundle, error) {
	producers := producers.New()
	return &SupportBundle{
		planner: planners.New(producers),
	}, nil
}

func (p *SupportBundle) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.SupportBundleVersion != nil:
		return p.planner.Version
	case spec.CustomerMeta != nil:
		return p.planner.CustomerMeta
	default:
		return nil
	}
}
