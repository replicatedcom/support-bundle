package distiller

type Interface interface {
	Distill(string) (interface{}, bool, error)
}
