package templates

import (
	"time"
)

func init() {
	RegisterFunc("seconds", Seconds)
	RegisterFunc("duration", Duration)
}

func Seconds(v interface{}) time.Duration {
	return Duration(v) * time.Second
}

func Duration(v interface{}) time.Duration {
	return time.Duration(toInt64(v))
}
