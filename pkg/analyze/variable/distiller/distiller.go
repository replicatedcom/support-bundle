package distiller

import (
	"bufio"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

func Distill(distiller Interface, r io.Reader, scannable bool) (interface{}, bool, error) {
	if scannable {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			value, ok, err := distiller.Distill(scanner.Text())
			if err != nil || ok {
				return value, ok, err
			}
		}
		return nil, false, nil
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", false, errors.Wrap(err, "read all")
	}
	return distiller.Distill(string(b))
}
