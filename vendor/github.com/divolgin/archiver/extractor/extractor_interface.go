package extractor

import (
	"io"
)

type Extractor interface {
	Extract(src, dest string) error
	ExtractFromReader(inputReader io.Reader, dest string) error
}
