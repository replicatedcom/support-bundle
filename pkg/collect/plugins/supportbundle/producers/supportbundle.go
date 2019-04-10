package producers

import "github.com/replicatedcom/support-bundle/pkg/util"

type SupportBundle struct {
	LogOutput *util.Buffer
}

func New(logOutput *util.Buffer) *SupportBundle {
	return &SupportBundle{LogOutput: logOutput}
}
