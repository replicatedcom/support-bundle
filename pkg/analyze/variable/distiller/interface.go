package distiller

type Interface interface {
	Distill(string) (string, error)
}
