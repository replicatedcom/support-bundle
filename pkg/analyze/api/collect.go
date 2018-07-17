package api

import collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"

// Collect is the top level collect object
type Collect struct {
	V1 []collecttypes.Spec `json:"v1,omitempty" yaml:"v1,omitempty" hcl:"v1,omitempty"`
}
