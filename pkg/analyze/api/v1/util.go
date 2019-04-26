package v1

import (
	"reflect"
	"strings"
)

func getTagName(v reflect.Value, i int, key string) string {
	return strings.SplitN(v.Type().Field(i).Tag.Get(key), ",", 2)[0]
}
