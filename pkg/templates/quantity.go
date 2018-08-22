package templates

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
)

func init() {
	RegisterFunc("parseQuantity", ParseQuantity)
	RegisterFunc("formatQuantity", FormatQuantity)
	RegisterFunc("sumQuantities", SumQuantities)
}

func SumQuantities(ii []string) int {
	var total int
	for _, i := range ii {
		total += ParseQuantity(i)
	}
	return total
}

func ParseQuantity(str string) int {
	q, err := resource.ParseQuantity(str)
	if err != nil {
		Panic("parseQuantity", errors.Wrapf(err, "parse %v", str))
	}
	return int(q.Value())
}

func FormatQuantity(value int, args ...string) string {
	var format resource.Format
	if len(args) > 0 {
		format = resource.Format(args[0])
	} else {
		format = resource.BinarySI
	}
	q := resource.NewQuantity(int64(value), format)
	return q.String()
}
