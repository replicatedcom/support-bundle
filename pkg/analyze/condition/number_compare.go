package condition

import (
	"errors"
	"reflect"
	"strconv"
)

var _ Interface = new(StringCompare)

type NumberCompare struct {
	Compare `json:",inline" yaml:",inline" hcl:",inline"`
}

func (c *NumberCompare) Eval(ref interface{}, data map[string]interface{}) (bool, error) {
	v := toFloat64(ref)
	switch {
	case c.Eq != nil:
		return v == toFloat64(c.Eq), nil
	case c.Lt != nil:
		return v < toFloat64(c.Lt), nil
	case c.Lte != nil:
		return v <= toFloat64(c.Lte), nil
	case c.Gt != nil:
		return v > toFloat64(c.Gt), nil
	case c.Gte != nil:
		return v >= toFloat64(c.Gte), nil
	case c.In != nil:
		for _, i := range c.In {
			if v == toFloat64(i) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.New("comparison function expected")
	}
}

// toFloat64 converts 64-bit floats
func toFloat64(v interface{}) float64 {
	if str, ok := v.(string); ok {
		iv, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return float64(val.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return float64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		return float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
	}
}
