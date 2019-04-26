package distiller

var _ Interface = new(Identity)

type Identity struct {
}

func (d *Identity) Distill(input string) (interface{}, bool, error) {
	return input, true, nil
}
