package variable

import (
	"reflect"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type matcher func(index []collecttypes.Result) []collecttypes.Result

func matchAny(m ...matcher) matcher {
	return func(index []collecttypes.Result) (results []collecttypes.Result) {
		for _, matcher := range m {
			results = append(results, matcher(index)...)
		}
		return
	}
}

func matchCollector(field interface{}) matcher {
	return func(index []collecttypes.Result) (results []collecttypes.Result) {
		for _, result := range index {
			if result.Size == 0 {
				continue
			}
			val := reflect.ValueOf(result.Spec)
			for i := 0; i < val.NumField(); i++ {
				if v := val.Field(i); v.Type() == reflect.ValueOf(field).Type() {
					if !reflect.ValueOf(v.Interface()).IsNil() {
						results = append(results, result)
					}
				}
			}
		}
		return
	}
}
