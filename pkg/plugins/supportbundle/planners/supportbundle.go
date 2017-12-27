package planners

import (
	"github.com/replicatedcom/support-bundle/pkg/plugins/supportbundle/producers"
)

type SupportBundle struct {
	producers *producers.SupportBundle
}

func New(supportBundle *producers.SupportBundle) *SupportBundle {
	return &SupportBundle{
		producers: supportBundle,
	}
}
