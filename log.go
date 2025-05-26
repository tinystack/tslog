package tslog

import (
	"io"
	"strings"
)

const (
	NoneLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

const (
	EncoderJSON    = "json"
	EncoderConsole = "console"
)

type T map[string]any

type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Debugt(msg string, args T)
	Infot(msg string, args T)
	Warnt(msg string, args T)
	Errort(msg string, args T)
}

type Level int8

type Options struct {
	lvl     Level
	w       []io.Writer
	encoder string
	caller  bool
	driver  Driver
}

type FuncOption func(*Options)

type Driver func(*Options) Logger

var unmarshalLevelText = map[string]Level{
	"debug": DebugLevel,
	"info":  InfoLevel,
	"warn":  WarnLevel,
	"error": ErrorLevel,
}

func defaultOptions() *Options {
	return &Options{
		lvl:     DebugLevel,
		encoder: EncoderJSON,
		caller:  false,
		driver:  NewZapDriver,
	}
}

func ParseLevel(text string) Level {
	text = strings.ToLower(text)
	lvl, _ := unmarshalLevelText[text]
	return lvl
}

func WithLevel(lvl Level) FuncOption {
	return func(o *Options) {
		o.lvl = lvl
	}
}

func WithWriter(w ...io.Writer) FuncOption {
	return func(o *Options) {
		o.w = w
	}
}

func WithCaller(caller bool) FuncOption {
	return func(o *Options) {
		o.caller = caller
	}
}

func WithEncoder(encoder string) FuncOption {
	return func(o *Options) {
		o.encoder = encoder
	}
}

func WithDriver(d Driver) FuncOption {
	return func(o *Options) {
		o.driver = d
	}
}

func NewLogger(funcOpts ...FuncOption) Logger {
	opts := defaultOptions()
	for _, f := range funcOpts {
		f(opts)
	}
	return opts.driver(opts)
}
