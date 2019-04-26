package reader

import (
	"archive/tar"
	"io"

	multierror "github.com/hashicorp/go-multierror"
)

type Scanner interface {
	Next() (*ScannerFile, error)
	Close() error
	Err() error
}

type BundleScanner struct {
	Prefix string

	tr      *tar.Reader
	err     error
	closers []io.Closer
}

type ScannerFile struct {
	Name   string
	Reader io.Reader
}

func (s *BundleScanner) Next() (*ScannerFile, error) {
	for {
		header, err := s.tr.Next()
		if err != nil {
			s.err = err
			return nil, err
		}
		if header == nil || header.Typeflag != tar.TypeReg {
			continue
		}

		name, ok := PathTrimPrefix(header.Name, s.Prefix)
		if !ok {
			continue
		}

		return &ScannerFile{Name: name, Reader: s.tr}, nil
	}
}

func (s *BundleScanner) Close() error {
	var multiErr *multierror.Error
	for _, closer := range s.closers {
		if err := closer.Close(); err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}
	return multiErr.ErrorOrNil()
}

func (s *BundleScanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}
