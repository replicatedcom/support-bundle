package variable

import (
	"reflect"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type matcher func(index []collecttypes.Result) ([]collecttypes.Result, error)

func matchAny(m ...matcher) matcher {
	return func(index []collecttypes.Result) ([]collecttypes.Result, error) {
		var results []collecttypes.Result
		for _, matcher := range m {
			matches, err := matcher(index)
			if err != nil {
				return results, err
			}
			results = append(results, matches...)
		}
		return results, nil
	}
}

func matchCollector(field interface{}) matcher {
	return func(index []collecttypes.Result) (results []collecttypes.Result, err error) {
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
