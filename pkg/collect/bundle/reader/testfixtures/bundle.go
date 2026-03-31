package testfixtures

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func WriteBundle(w io.Writer, bundlePath string) error {
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filename)
	dir := filepath.Join(basePath, bundlePath)

	gw := gzip.NewWriter(w)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = entry.Name()
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, f); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}
