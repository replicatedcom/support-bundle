package distiller

import (
	"bufio"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

func Distill(distiller Interface, r io.Reader, scannable bool) (interface{}, error) {
	if scannable {
		scanner := bufio.NewScanner(r)
		defer func() {
			// read the rest
			for scanner.Scan() {
			}
		}()
		for scanner.Scan() {
			value, ok, err := distiller.Distill(scanner.Text())
			if err != nil || ok {
				return value, err
			}
		}
		return nil, errors.Wrap(scanner.Err(), "scan error")
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read all")
	}
	value, _, err := distiller.Distill(string(b))
	return value, err
}
