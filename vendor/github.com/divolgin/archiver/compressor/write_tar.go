package compressor

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Tar struct {
	ExcludeList       []string
	TruncateLongFiles bool // Set this to true to prevent "write too long" errors
}

func WriteTar(srcPath string, dest io.Writer) error {
	return Tar{}.writeTar(srcPath, dest)
}

func (t Tar) writeTar(srcPath string, dest io.Writer) error {
	absPath, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(dest)
	defer tw.Close()

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var relative string
		if os.IsPathSeparator(srcPath[len(srcPath)-1]) {
			relative, err = filepath.Rel(absPath, path)
		} else {
			relative, err = filepath.Rel(filepath.Dir(absPath), path)
		}

		relative = filepath.ToSlash(relative)

		if err != nil {
			return err
		}

		if t.fileIsExcluded(path) {
			return nil
		}

		return t.addTarFile(path, relative, tw)
	})

	return err
}

func (t Tar) addTarFile(path, name string, tw *tar.Writer) error {
	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}

	link := ""
	if fi.Mode()&os.ModeSymlink != 0 {
		if link, err = os.Readlink(path); err != nil {
			return err
		}
	}

	hdr, err := tar.FileInfoHeader(fi, link)
	if err != nil {
		return err
	}

	if fi.IsDir() && !os.IsPathSeparator(name[len(name)-1]) {
		name = name + "/"
	}

	if hdr.Typeflag == tar.TypeReg && name == "." {
		// archiving a single file
		hdr.Name = filepath.ToSlash(filepath.Base(path))
	} else {
		hdr.Name = filepath.ToSlash(name)
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	if hdr.Typeflag == tar.TypeReg {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		if t.TruncateLongFiles {
			_, err = io.CopyN(tw, file, hdr.Size)
		} else {
			_, err = io.Copy(tw, file)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (t Tar) fileIsExcluded(path string) bool {
	for _, f := range t.ExcludeList {
		if strings.HasPrefix(path, f) {
			return true
		}
	}
	return false
}
