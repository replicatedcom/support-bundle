package condition

import "fmt"

var _ Interface = new(Empty)

type Empty struct {
}

func (c *Empty) Eval(ref interface{}, data map[string]interface{}) (bool, error) {
	var str string
	if ref != nil {
		str = fmt.Sprintf("%v", ref)
	}
	return str == "", nil
}
