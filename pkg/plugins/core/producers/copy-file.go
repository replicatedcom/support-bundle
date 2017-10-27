package producers

import "os"

func CopyFile(dest, src string) error {
	return os.Link(src, dest)
}
