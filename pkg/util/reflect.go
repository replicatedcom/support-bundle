package util

import (
	"reflect"
	"strings"
)

func StructTagName(v reflect.Value, i int, key string) string {
	return strings.SplitN(v.Type().Field(i).Tag.Get(key), ",", 2)[0]
}
