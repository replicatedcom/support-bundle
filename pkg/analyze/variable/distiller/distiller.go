package distiller

import (
	"bufio"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

func Distill(distiller Interface, r io.Reader, scannable bool) (string, error) {
	if scannable {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			value, err := distiller.Distill(scanner.Text())
			if err != nil || value != "" {
				return value, err
			}
		}
		return "", nil
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.Wrap(err, "read all")
	}
	return distiller.Distill(string(b))
}
