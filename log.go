// Package tslog provides a flexible, high-performance logging library for Go applications.
// It supports multiple log levels, encoders, and output writers with a simple, clean API.
//
// The package offers both a default logger for quick usage and the ability to create
// custom logger instances with specific configurations. It's built on top of
// popular logging libraries like Zap for performance and reliability.
//
// Basic usage:
//
//	tslog.Info("Hello, World!")
//	tslog.Errorf("Error occurred: %v", err)
//
// Advanced usage:
//
//	logger := tslog.NewLogger(
//	    tslog.WithLevel(tslog.InfoLevel),
//	    tslog.WithEncoder(tslog.EncoderJSON),
//	)
//	logger.Info("Custom logger message")
package tslog

import (
	"fmt"
	"io"
	"strings"
)

// Log levels define the severity of log messages.
// The levels are ordered from least to most severe.
const (
	// NoneLevel represents no logging level (disabled logging)
	NoneLevel Level = iota
	// DebugLevel is used for detailed diagnostic information
	DebugLevel
	// InfoLevel is used for general information messages
	InfoLevel
	// WarnLevel is used for warning messages that don't halt execution
	WarnLevel
	// ErrorLevel is used for error messages that may affect functionality
	ErrorLevel
)

// Encoder types define the output format of log messages.
const (
	// EncoderJSON outputs logs in JSON format for structured logging
	EncoderJSON = "json"
	// EncoderConsole outputs logs in human-readable console format
	EncoderConsole = "console"
)

// T represents a map of key-value pairs for structured logging.
// It's used with the *t methods (Debugt, Infot, etc.) to add
// structured fields to log messages.
//
// Example:
//
//	tslog.Infot("User logged in", tslog.T{
//	    "user_id": 123,
//	    "ip": "192.168.1.1",
//	})
type T map[string]any

// Logger defines the interface for all logging operations.
// This interface provides a consistent API across different
// logging implementations and drivers.
type Logger interface {
	// Debug logs a message at Debug level
	Debug(args ...any)
	// Info logs a message at Info level
	Info(args ...any)
	// Warn logs a message at Warn level
	Warn(args ...any)
	// Error logs a message at Error level
	Error(args ...any)

	// Debugf logs a formatted message at Debug level
	Debugf(format string, args ...any)
	// Infof logs a formatted message at Info level
	Infof(format string, args ...any)
	// Warnf logs a formatted message at Warn level
	Warnf(format string, args ...any)
	// Errorf logs a formatted message at Error level
	Errorf(format string, args ...any)

	// Debugt logs a message with structured fields at Debug level
	Debugt(msg string, args T)
	// Infot logs a message with structured fields at Info level
	Infot(msg string, args T)
	// Warnt logs a message with structured fields at Warn level
	Warnt(msg string, args T)
	// Errort logs a message with structured fields at Error level
	Errort(msg string, args T)
}

// Level represents the logging level type.
// It's an int8 to minimize memory usage while providing
// enough range for all supported log levels.
type Level int8

// String returns the string representation of the log level.
// This method implements the fmt.Stringer interface.
func (l Level) String() string {
	switch l {
	case NoneLevel:
		return "none"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return fmt.Sprintf("Level(%d)", int(l))
	}
}

// Enabled returns true if the level is enabled for logging.
// NoneLevel is considered disabled.
func (l Level) Enabled() bool {
	return l != NoneLevel
}

// Options holds configuration settings for creating a new logger instance.
// All fields are private to ensure they can only be modified through
// the provided option functions, maintaining configuration integrity.
type Options struct {
	// lvl defines the minimum log level that will be processed
	lvl Level
	// w holds the list of writers where log messages will be output
	w []io.Writer
	// encoder specifies the format of log output (json or console)
	encoder string
	// caller determines whether to include caller information in logs
	caller bool
	// driver is the factory function used to create the actual logger implementation
	driver Driver
}

// Validate checks if the options are valid and returns an error if not.
func (o *Options) Validate() error {
	if o.driver == nil {
		return fmt.Errorf("driver cannot be nil")
	}
	if o.encoder != EncoderJSON && o.encoder != EncoderConsole {
		return fmt.Errorf("encoder must be either %q or %q", EncoderJSON, EncoderConsole)
	}
	if len(o.w) == 0 {
		return fmt.Errorf("at least one writer must be specified")
	}
	return nil
}

// FuncOption is a function type that modifies Options.
// This pattern allows for flexible and extensible configuration
// of logger instances without breaking existing API.
type FuncOption func(*Options)

// Driver is a factory function type that creates a Logger instance
// from the provided Options. This allows for different logging
// implementations to be used interchangeably.
type Driver func(*Options) Logger

// unmarshalLevelText maps string representations to Level values.
// This is used for parsing log levels from configuration files
// or command-line arguments.
var unmarshalLevelText = map[string]Level{
	"none":  NoneLevel,
	"debug": DebugLevel,
	"info":  InfoLevel,
	"warn":  WarnLevel,
	"error": ErrorLevel,
}

// defaultOptions returns a new Options instance with sensible defaults.
// This ensures that a logger can be created without any configuration
// and will work out of the box.
func defaultOptions() *Options {
	return &Options{
		lvl:     DebugLevel,
		encoder: EncoderJSON,
		caller:  false,
		driver:  NewZapDriver,
	}
}

// ParseLevel converts a string representation of a log level to a Level value.
// The comparison is case-insensitive. If the string doesn't match any
// known level, NoneLevel is returned.
//
// Supported level strings: "none", "debug", "info", "warn", "error"
func ParseLevel(text string) Level {
	text = strings.ToLower(strings.TrimSpace(text))
	if lvl, ok := unmarshalLevelText[text]; ok {
		return lvl
	}
	return NoneLevel
}

// WithLevel sets the minimum log level for the logger.
// Messages below this level will be ignored.
//
// Example:
//
//	logger := tslog.NewLogger(tslog.WithLevel(tslog.InfoLevel))
func WithLevel(lvl Level) FuncOption {
	return func(o *Options) {
		o.lvl = lvl
	}
}

// WithWriter sets one or more output writers for the logger.
// Multiple writers can be specified to output logs to different destinations
// simultaneously (e.g., stdout and a file).
//
// Example:
//
//	logger := tslog.NewLogger(tslog.WithWriter(os.Stdout, file))
func WithWriter(w ...io.Writer) FuncOption {
	return func(o *Options) {
		// Create a new slice to avoid modifying the original slice
		o.w = make([]io.Writer, len(w))
		copy(o.w, w)
	}
}

// WithCaller enables or disables the inclusion of caller information
// (file name and line number) in log messages.
//
// Example:
//
//	logger := tslog.NewLogger(tslog.WithCaller(true))
func WithCaller(caller bool) FuncOption {
	return func(o *Options) {
		o.caller = caller
	}
}

// WithEncoder sets the output format for log messages.
// Supported encoders are EncoderJSON and EncoderConsole.
//
// Example:
//
//	logger := tslog.NewLogger(tslog.WithEncoder(tslog.EncoderConsole))
func WithEncoder(encoder string) FuncOption {
	return func(o *Options) {
		o.encoder = encoder
	}
}

// WithDriver sets a custom driver for creating the logger implementation.
// This allows for different logging backends to be used.
//
// Example:
//
//	logger := tslog.NewLogger(tslog.WithDriver(MyCustomDriver))
func WithDriver(d Driver) FuncOption {
	return func(o *Options) {
		o.driver = d
	}
}

// NewLogger creates a new Logger instance with the specified options.
// If no options are provided, default options will be used.
// The function applies all options in order and then creates the logger
// using the configured driver.
//
// Example:
//
//	logger := tslog.NewLogger(
//	    tslog.WithLevel(tslog.InfoLevel),
//	    tslog.WithEncoder(tslog.EncoderJSON),
//	    tslog.WithCaller(true),
//	)
func NewLogger(funcOpts ...FuncOption) Logger {
	opts := defaultOptions()
	for _, f := range funcOpts {
		if f != nil {
			f(opts)
		}
	}

	// Validate options before creating logger
	if err := opts.Validate(); err != nil {
		// Fall back to a safe default if validation fails
		opts = defaultOptions()
	}

	return opts.driver(opts)
}
