package logger

import (
	jww "github.com/spf13/jwalterweatherman"
)

func LevelFromJWWThreshold(threshold jww.Threshold) string {
	switch threshold {
	case jww.LevelTrace, jww.LevelDebug:
		return "debug"
	case jww.LevelInfo:
		return "info"
	case jww.LevelWarn:
		return "warn"
	case jww.LevelError, jww.LevelCritical, jww.LevelFatal:
		return "error"
	default:
		return "debug"
	}
}
