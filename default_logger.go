package tslog

import (
	"github.com/tinystack/tslog/writer"
)

var defaultLogger Logger

var (
	defaultLogLevel = DebugLevel
)

func init() {

	funcOpts := []FuncOption{
		WithLevel(defaultLogLevel),
		WithWriter(writer.NewStdoutWriter()),
		WithEncoder(EncoderConsole),
		WithCaller(false),
	}

	opts := defaultOptions()
	for _, f := range funcOpts {
		f(opts)
	}

	defaultLogger = opts.driver(opts)
}

func DefaultLogger() Logger {
	return defaultLogger
}

func UpdateDefaultLogger(l Logger) {
	defaultLogger = l
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Debugt(msg string, args T) {
	defaultLogger.Debugt(msg, args)
}

func Infot(msg string, args T) {
	defaultLogger.Infot(msg, args)
}

func Warnt(msg string, args T) {
	defaultLogger.Warnt(msg, args)
}

func Errort(msg string, args T) {
	defaultLogger.Errort(msg, args)
}
