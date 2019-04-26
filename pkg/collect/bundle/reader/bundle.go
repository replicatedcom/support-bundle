package reader

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/spf13/afero"
)

type BundleReader interface {
	GetIndex() []collecttypes.Result
	GetErrorIndex() []collecttypes.Result
	NewScanner() (Scanner, error)
}

type Piper struct {
	Name string
	W    *io.PipeWriter
}

type Bundle struct {
	Fs     afero.Fs
	Path   string
	Prefix string
	index  []collecttypes.Result
	errors []collecttypes.Result
}

func NewBundle(fs afero.Fs, path, prefix string) (*Bundle, error) {
	b := &Bundle{
		Fs:     fs,
		Path:   path,
		Prefix: prefix,
	}
	var err error
	b.index, b.errors, err = b.scanIndexes()
	return b, err
}

func (b *Bundle) GetIndex() []collecttypes.Result {
	return b.index
}

func (b *Bundle) GetErrorIndex() []collecttypes.Result {
	return b.errors
}

func (b *Bundle) NewScanner() (Scanner, error) {
	if _, err := b.Fs.Stat(b.Path); os.IsNotExist(err) {
		return nil, err
	}

	s := &BundleScanner{Prefix: b.Prefix}

	file, err := b.Fs.Open(b.Path)
	if err != nil {
		return nil, err
	}
	s.closers = append(s.closers, file)

	gzr, err := gzip.NewReader(file)
	if err != nil {
		s.Close()
		return nil, err
	}
	s.closers = append(s.closers, gzr)

	s.tr = tar.NewReader(gzr)

	return s, nil
}

func (b *Bundle) scanIndexes() (indexR []collecttypes.Result, errorR []collecttypes.Result, returnErr error) {
	scanner, err := b.NewScanner()
	if err != nil {
		returnErr = errors.Wrap(err, "new scanner")
		return
	}
	defer scanner.Close()

	var foundIndex, foundError bool
	for {
		f, err := scanner.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			returnErr = errors.Wrap(err, "scan next")
			return
		}

		switch {
		case PathsMatch(f.Name, "index.json"):
			foundIndex = true
			indexR, err = b.getResultsFromIndex(f.Reader)
			if err != nil {
				returnErr = errors.Wrap(err, "get results from index.json")
				return
			}

		case PathsMatch(f.Name, "error.json"):
			foundError = true
			errorR, err = b.getResultsFromIndex(f.Reader)
			if err != nil {
				returnErr = errors.Wrap(err, "get results from error.json")
				return
			}
		}

		if foundIndex && foundError {
			break
		}
	}

	if !foundIndex {
		returnErr = &os.PathError{Op: "open", Path: "index.json", Err: os.ErrNotExist}
	} else if !foundError {
		returnErr = &os.PathError{Op: "open", Path: "error.json", Err: os.ErrNotExist}
	}

	return
}

func (b *Bundle) getResultsFromIndex(reader io.Reader) ([]collecttypes.Result, error) {
	var results []collecttypes.Result
	return results, json.NewDecoder(reader).Decode(&results)
}
