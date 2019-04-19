package util

import (
	"io"

	multierror "github.com/hashicorp/go-multierror"
)

func MultiWriterAppend(writers ...io.Writer) (multiWriter io.Writer) {
	for _, w := range writers {
		if w == nil {
			continue
		}
		if multiWriter == nil {
			multiWriter = w
		} else {
			multiWriter = io.MultiWriter(multiWriter, w)
		}
	}
	return
}

func MultiCloserAppend(closers ...io.Closer) io.Closer {
	var allClosers []io.Closer
	for _, c := range closers {
		if c == nil {
			continue
		}
		if mc, ok := c.(*multiCloser); ok {
			allClosers = append(allClosers, mc.closers...)
		} else {
			allClosers = append(allClosers, c)
		}
	}
	return &multiCloser{closers: allClosers}
}

type multiCloser struct {
	closers []io.Closer
}

func (c *multiCloser) Close() error {
	var multiErr *multierror.Error
	for _, close := range c.closers {
		if err := close.Close(); err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}
	return multiErr.ErrorOrNil()
}
