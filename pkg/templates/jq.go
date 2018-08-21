package templates

import (
	"encoding/json"
	"errors"
	"fmt"

	jqpipe "github.com/andrewchambers/go-jqpipe"
)

func init() {
	RegisterFunc("jq", JSONQuery)

	RegisterFunc("jqMessagesToStringSlice", JqMessagesToStringSlice)
	RegisterFunc("jqMessagesToString", JqMessagesToString)
}

func JSONQuery(obj, exp string) []interface{} {
	seq, err := jqpipe.Eval(obj, exp)
	if err != nil {
		Panic("jq", err)
	}
	var out []interface{}
	for _, s := range seq {
		o, err := unmarshalArbitraryJSON(s)
		if err != nil {
			Panic("jq", err)
		}
		out = append(out, o)
	}
	return out
}

func JqMessagesToString(msgs []interface{}) string {
	if len(msgs) == 0 {
		Panic("jqMessagesToString", errors.New("jq messages empty"))
	}
	switch e := msgs[0].(type) {
	case string:
		return e
	default:
		Panic("jqMessagesToString", fmt.Errorf("type %T unsupported, must be of type string", msgs[0]))
		return "" // panic above
	}
}

func JqMessagesToStringSlice(msgs []interface{}) []string {
	var ss []string
	for _, msg := range msgs {
		switch e := msg.(type) {
		case string:
			ss = append(ss, e)
		default:
			Panic("jqMessagesToStringSlice", fmt.Errorf("type %T unsupported, must be of type string", msg))
		}
	}
	return ss
}

// "number", "bool", "string", "array", "object"
func unmarshalArbitraryJSON(obj []byte) (interface{}, error) {
	a := []interface{}{}
	obj = append([]byte("["), obj...)
	obj = append(obj, []byte("]")...)
	if err := json.Unmarshal(obj, &a); err != nil {
		return nil, err
	}
	if len(a) == 0 {
		return nil, errors.New("unmarshal json unexpected empty array")
	}
	return a[0], nil
}
