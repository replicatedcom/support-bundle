package retraced

import (
	"github.com/replicatedcom/support-bundle/pkg/plugins/retraced/planners"
	"github.com/replicatedcom/support-bundle/pkg/plugins/retraced/producers"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type Retraced struct {
	planner *planners.Retraced
}

func New() (*Retraced, error) {
	producers := producers.New()
	return &Retraced{
		planner: planners.New(producers),
	}, nil
}

func (p *Retraced) Plan(spec types.Spec) types.Planner {
	switch {
	case spec.RetracedEvents != nil:
		return p.planner.Events
	default:
		return nil
	}
}
