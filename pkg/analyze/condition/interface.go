package condition

type Interface interface {
	Eval(interface{}, map[string]interface{}, error) (bool, error)
}
