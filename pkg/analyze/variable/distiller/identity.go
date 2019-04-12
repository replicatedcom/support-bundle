package distiller

var _ Interface = new(Identity)

type Identity struct {
}

func (d *Identity) Distill(input string) (interface{}, error) {
	return input, nil
}
