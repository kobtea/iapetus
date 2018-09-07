package util

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	stdlog "log"
	"os"
)

var logger = newLogger()

func newLogger() *log.Logger {
	l := log.NewLogfmtLogger(os.Stdout)
	return &l
}

func SetLogFilterLevel(lv string) {
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
		(*logger).Log("msg", fmt.Sprintf("invalid log level: %s. use `info` instead.", lv))
		op = level.AllowInfo()
	}
	l := level.NewFilter(*logger, op)
	logger = &l // FIXME
}

func GetLogger() *log.Logger {
	l := log.With(*logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	return &l
}

func GetStdErrorLogger() *stdlog.Logger {
	l := log.With(*logger, "ts", log.DefaultTimestampUTC)
	return stdlog.New(log.NewStdlibAdapter(level.Error(l)), "", stdlog.Lshortfile)
}
