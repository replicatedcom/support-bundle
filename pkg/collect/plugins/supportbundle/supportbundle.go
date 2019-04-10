package supportbundle

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/supportbundle/planners"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/supportbundle/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/util"
)

type SupportBundle struct {
	planner *planners.SupportBundle
}

func New(logOutput *util.Buffer) *SupportBundle {
	return &SupportBundle{
		planner: planners.New(
			producers.New(logOutput),
		),
	}
}

func (p *SupportBundle) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.SupportBundleVersion != nil:
		return p.planner.Version
	case spec.SupportBundleGoroutines != nil:
		return p.planner.Goroutines
	case spec.SupportBundleLogs != nil:
		return p.planner.Logs
	case spec.CustomerMeta != nil:
		return p.planner.CustomerMeta
	case spec.GlobalRedaction != nil:
		return p.planner.GlobalRedaction
	default:
		return nil
	}
}
