package variable

import (
	"reflect"

	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type matcher func(result collecttypes.Result) bool

func matchAny(result collecttypes.Result, matchers ...matcher) bool {
	for _, matcher := range matchers {
		if matcher(result) {
			return true
		}
	}
	return false
}

func matcherCollector(field interface{}) matcher {
	return func(result collecttypes.Result) bool {
		if result.Size == 0 {
			return false
		}
		val := reflect.ValueOf(result.Spec)
		for i := 0; i < val.NumField(); i++ {
			if v := val.Field(i); v.Type() == reflect.ValueOf(field).Type() {
				return !reflect.ValueOf(v.Interface()).IsNil()
			}
		}
		return false
	}
}

func matcherCoreReadFileFilepath(filepath string) matcher {
	return func(result collecttypes.Result) bool {
		return result.Spec.CoreReadFile != nil &&
			result.Spec.CoreReadFile.Filepath == filepath &&
			result.Size > 0
	}
}
