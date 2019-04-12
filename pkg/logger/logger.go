package logger

import (
	"fmt"
	golog "log"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-stack/stack"
	"github.com/spf13/viper"
)

var (
	fullPathCaller = pathCaller(3)
	globalLogger   log.Logger
	logMtx         sync.Mutex
)

// FromViper builds a logger from env using viper
func FromViper(v *viper.Viper) log.Logger {
	return New(GetLevel(v))
}

func GetLevel(v *viper.Viper) string {
	if v.GetBool("verbose") {
		return "debug"
	}
	return v.GetString("log-level")
}

// New instantiates the global logger instance if not already created
func New(logLevel string) log.Logger {

	// one at a time plz
	logMtx.Lock()
	defer logMtx.Unlock()

	if globalLogger != nil {
		return globalLogger
	}

	globalLogger = log.NewLogfmtLogger(os.Stderr) // stderr since this is a cli
	globalLogger = log.With(globalLogger, "ts", log.DefaultTimestampUTC)
	globalLogger = withLevel(globalLogger, logLevel)
	globalLogger = log.With(globalLogger, "caller", fullPathCaller)
	golog.SetOutput(log.NewStdlibAdapter(level.Debug(globalLogger)))
	return globalLogger
}

func withLevel(logger log.Logger, lvl string) log.Logger {
	switch lvl {
	case "debug":
		return level.NewFilter(logger, level.AllowDebug())
	case "info":
		return level.NewFilter(logger, level.AllowInfo())
	case "warn":
		return level.NewFilter(logger, level.AllowWarn())
	case "error":
		return level.NewFilter(logger, level.AllowError())
	case "off":
		return level.NewFilter(logger, level.AllowNone())
	default:
		logger.Log("msg", "Unknown log level, using debug", "received", lvl)
		return level.NewFilter(logger, level.AllowDebug())
	}
}

func pathCaller(depth int) log.Valuer {
	return func() interface{} {
		return fmt.Sprintf("%+s", stack.Caller(depth))
	}
}
