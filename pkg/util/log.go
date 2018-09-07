package util

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	stdlog "log"
	"os"
)

func NewLogger(lv string) log.Logger {
	logger := log.NewLogfmtLogger(os.Stdout)
	var op level.Option
	switch lv {
	case "debug":
		op = level.AllowDebug()
	case "info":
		op = level.AllowInfo()
	case "warn":
		op = level.AllowWarn()
	case "error":
		op = level.AllowError()
	default:
		logger.Log("msg", fmt.Sprintf("invalid log level: %s. use `info` instead.", lv))
		op = level.AllowInfo()
	}
	logger = level.NewFilter(logger, op)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	return logger
}

func NewStdLogger(l log.Logger) *stdlog.Logger {
	return stdlog.New(log.NewStdlibAdapter(l), "", 0)
}
