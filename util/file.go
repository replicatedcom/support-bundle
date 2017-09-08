package util

import (
	"io/ioutil"
)

func ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
