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
	seq, err := jqpipe.Eval(exp, obj)
	if err != nil {
		Panic("jq", err)
	}
	var out []interface{}
	for _, s := range seq {
		var o interface{}
		err := json.Unmarshal(s, &o)
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
	case int:
		return fmt.Sprintf("%d", e)
	case float32, float64:
		return fmt.Sprintf("%f", e)
	case bool:
		return fmt.Sprintf("%v", e)
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
