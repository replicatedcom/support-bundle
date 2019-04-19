package distiller

var _ Interface = new(Identity)

type Identity struct {
}

func (d *Identity) Distill(input string) (string, error) {
	return input, nil
}
