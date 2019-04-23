package distiller

import "github.com/replicatedcom/support-bundle/pkg/templates"

var _ Interface = new(Eval)

type Eval string

func (d *Eval) Distill(input string) (interface{}, bool, error) {
	val, err := templates.String(string(*d), input)
	if err != nil {
		return val, false, err
	}
	return val, true, err
}
