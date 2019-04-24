package templates

import (
	"math"
	"reflect"
	"strconv"
)

func init() {
	RegisterFunc("addFloat", addFloat)
	RegisterFunc("subFloat", subFloat)
	RegisterFunc("divFloat", divFloat)
	RegisterFunc("mulFloat", mulFloat)
}

func addFloat(i ...interface{}) float64 {
	var a float64
	for _, b := range i {
		a += toFloat64(b)
	}
	return a
}

func subFloat(a, b interface{}) float64 {
	return toFloat64(a) - toFloat64(b)
}

func divFloat(a, b interface{}) float64 {
	return toFloat64(a) / toFloat64(b)
}

func mulFloat(a interface{}, v ...interface{}) float64 {
	val := toFloat64(a)
	for _, b := range v {
		val = val * toFloat64(b)
	}
	return val
}

// toInt64 converts integer types to 64-bit integers
func toInt64(v interface{}) int64 {
	if str, ok := v.(string); ok {
		iv, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return val.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		tv := val.Uint()
		if tv <= math.MaxInt64 {
			return int64(tv)
		}
		// TODO: What is the sensible thing to do here?
		return math.MaxInt64
	case reflect.Float32, reflect.Float64:
		return int64(val.Float())
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
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
