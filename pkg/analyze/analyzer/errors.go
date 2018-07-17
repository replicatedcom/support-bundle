package analyzer

import "github.com/replicatedcom/support-bundle/pkg/meta"

type RefNotFoundError struct {
	Ref meta.Ref
}
