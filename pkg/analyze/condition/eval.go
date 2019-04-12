package condition

import "github.com/pkg/errors"

func Eval(condition Interface, ref string, data map[string]interface{}) (bool, error) {
	value, err := resolveRef(ref, data)
	if err != nil {
		return false, errors.Wrapf(err, "resolve ref %s", ref)
	}
	b, err := condition.Eval(value, data)
	return b, errors.Wrap(err, "eval")
}

func resolveRef(ref string, data map[string]interface{}) (interface{}, error) {
	if ref == "" {
		return nil, nil
	}
	value, ok := data[ref]
	if !ok {
		return nil, errors.New("not found")
	}
	return value, nil
}
