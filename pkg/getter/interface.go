package getter

type Interface interface {
	Get(src string) (dst string, err error)
	DstDir() string
}
