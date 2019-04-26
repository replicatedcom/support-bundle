package testfixtures

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver"
)

func WriteBundle(w io.Writer, bundlePath string) error {
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filename)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(filepath.Join(basePath, bundlePath)); err != nil {
		return err
	}
	defer os.Chdir(cwd)

	var filePaths []string
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, info := range files {
		filePaths = append(filePaths, info.Name())
	}

	return archiver.TarGz.Write(w, filePaths)
}
