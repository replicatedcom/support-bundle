package ginkgo

import (
	"archive/tar"
	"compress/gzip"
	. "github.com/onsi/gomega"
	jww "github.com/spf13/jwalterweatherman"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var tmpdir string
var cwd string
var err error

func EnterNewTempDir() {
	cwd, err = os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	tmpdir, err = ioutil.TempDir("", "support-bundle")
	Expect(err).NotTo(HaveOccurred())
	err = os.Chdir(tmpdir)
	Expect(err).NotTo(HaveOccurred())
}

func CleanupDir() {
	err = os.Chdir(cwd)
	Expect(err).NotTo(HaveOccurred())
	err = os.RemoveAll(tmpdir)
	Expect(err).NotTo(HaveOccurred())
}

func WriteFile(path string, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0666)
	Expect(err).NotTo(HaveOccurred())
}

func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())
	return data
}

func ReadFileFromBundle(archivePath, targetFile string) string {
	file, err := os.Open(archivePath)
	defer CloseLogErr(file)
	Expect(err).NotTo(HaveOccurred())

	gzr, err := gzip.NewReader(file)
	defer CloseLogErr(gzr)
	Expect(err).NotTo(HaveOccurred())

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		Expect(err).NotTo(HaveOccurred())
		if header == nil {
			continue
		}

		filePath := strings.TrimLeft(header.Name, "0123456789")
		jww.DEBUG.Printf("reading tar entry %s looking for %s", filePath, targetFile)

		if filePath == targetFile && header.Typeflag == tar.TypeReg {
			contents, err := ioutil.ReadAll(tr)
			Expect(err).NotTo(HaveOccurred())
			return string(contents)
		}
	}
}

func CloseLogErr(c io.Closer) {
	if err := c.Close(); err != nil {
		jww.ERROR.Print(err)
	}
}