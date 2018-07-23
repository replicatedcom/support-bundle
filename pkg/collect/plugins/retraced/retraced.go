package retraced

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/retraced/planners"
	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/retraced/producers"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type Retraced struct {
	planner *planners.Retraced
}

func New() *Retraced {
	return &Retraced{
		planner: planners.New(
			producers.New(),
		),
	}
}

func (p *Retraced) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.RetracedEvents != nil:
		return p.planner.Events
	default:
		return nil
	}
}
