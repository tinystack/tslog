// Package tslog provides Zap-based logger implementation.
// This file contains the Zap driver that implements the Logger interface
// using Uber's Zap logging library for high performance structured logging.
package tslog

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger is a wrapper around Zap's SugaredLogger that implements
// the tslog.Logger interface. It provides thread-safe logging operations
// with high performance and low allocation overhead.
type zapLogger struct {
	zap    *zap.SugaredLogger
	mutex  sync.RWMutex // Protects the zap field for safe concurrent access
	closed bool         // Indicates if the logger has been closed
}

// zapLevel maps tslog.Level to zapcore.Level for compatibility.
// This mapping ensures that log levels are correctly translated
// between the tslog interface and Zap's internal representation.
var zapLevel = map[Level]zapcore.Level{
	NoneLevel:  zapcore.FatalLevel + 1, // Effectively disables logging
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

// NewZapDriver creates a new Logger instance using Zap as the underlying
// logging implementation. This function configures Zap with the provided
// options and returns a Logger that implements the tslog.Logger interface.
//
// The driver supports:
// - Multiple log levels with efficient level checking
// - JSON and Console output encoders
// - Multiple output writers
// - Optional caller information
// - High-performance structured logging
//
// If opts is nil, default options will be used.
//
// Example:
//
//	logger := tslog.NewZapDriver(&tslog.Options{
//	    lvl: tslog.InfoLevel,
//	    encoder: tslog.EncoderJSON,
//	    w: []io.Writer{os.Stdout},
//	})
func NewZapDriver(opts *Options) Logger {
	if opts == nil {
		opts = defaultOptions()
	}

	// Validate options
	if err := opts.Validate(); err != nil {
		// Log validation error and use defaults
		fmt.Fprintf(os.Stderr, "tslog: invalid options (%v), using defaults\n", err)
		opts = defaultOptions()
	}

	// Convert tslog level to Zap level
	lvl := zapcore.InfoLevel
	if l, ok := zapLevel[opts.lvl]; ok {
		lvl = l
	}

	// Create atomic level for dynamic level changes
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(lvl)

	// Configure encoder with production-ready settings
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.FunctionKey = "func"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stacktrace"

	// Choose encoder based on configuration
	var encoder zapcore.Encoder
	switch opts.encoder {
	case EncoderConsole:
		// Console encoder for human-readable output
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case EncoderJSON:
		fallthrough
	default:
		// JSON encoder for structured logging (default)
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create write syncers from provided writers
	var syncers []zapcore.WriteSyncer
	if len(opts.w) == 0 {
		// Fallback to stdout if no writers provided
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	} else {
		for _, w := range opts.w {
			if w != nil {
				syncers = append(syncers, zapcore.AddSync(w))
			}
		}
	}

	// Create core with multi-writer support
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(syncers...),
		atomicLevel,
	)

	// Configure Zap options
	zapOpts := []zap.Option{
		zap.AddCallerSkip(2), // Skip tslog wrapper functions
	}

	if opts.caller {
		zapOpts = append(zapOpts, zap.AddCaller())
	}

	// Add stack traces for error level and above
	zapOpts = append(zapOpts, zap.AddStacktrace(zapcore.ErrorLevel))

	// Create the Zap logger
	z := zap.New(core, zapOpts...).Sugar()

	return &zapLogger{
		zap:    z,
		closed: false,
	}
}

// z returns the underlying Zap SugaredLogger in a thread-safe manner.
// It panics if the logger is not initialized or has been closed.
func (l *zapLogger) z() *zap.SugaredLogger {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if l.zap == nil {
		panic("tslog: zapLogger not initialized")
	}
	if l.closed {
		panic("tslog: zapLogger has been closed")
	}
	return l.zap
}

// Close flushes any buffered log entries and closes the logger.
// After calling Close, the logger should not be used.
func (l *zapLogger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.closed {
		return nil
	}

	var err error
	if l.zap != nil {
		err = l.zap.Sync()
		l.zap = nil
	}
	l.closed = true
	return err
}

// Debug logs a message at Debug level.
// Arguments are handled in the manner of fmt.Print.
func (l *zapLogger) Debug(args ...any) {
	l.z().Debug(args...)
}

// Info logs a message at Info level.
// Arguments are handled in the manner of fmt.Print.
func (l *zapLogger) Info(args ...any) {
	l.z().Info(args...)
}

// Warn logs a message at Warn level.
// Arguments are handled in the manner of fmt.Print.
func (l *zapLogger) Warn(args ...any) {
	l.z().Warn(args...)
}

// Error logs a message at Error level.
// Arguments are handled in the manner of fmt.Print.
func (l *zapLogger) Error(args ...any) {
	l.z().Error(args...)
}

// Debugf logs a formatted message at Debug level.
// Arguments are handled in the manner of fmt.Printf.
func (l *zapLogger) Debugf(format string, args ...any) {
	l.z().Debugf(format, args...)
}

// Infof logs a formatted message at Info level.
// Arguments are handled in the manner of fmt.Printf.
func (l *zapLogger) Infof(format string, args ...any) {
	l.z().Infof(format, args...)
}

// Warnf logs a formatted message at Warn level.
// Arguments are handled in the manner of fmt.Printf.
func (l *zapLogger) Warnf(format string, args ...any) {
	l.z().Warnf(format, args...)
}

// Errorf logs a formatted message at Error level.
// Arguments are handled in the manner of fmt.Printf.
func (l *zapLogger) Errorf(format string, args ...any) {
	l.z().Errorf(format, args...)
}

// Debugt logs a message with structured fields at Debug level.
// The structured fields are converted to key-value pairs for Zap.
func (l *zapLogger) Debugt(msg string, args T) {
	if len(args) == 0 {
		l.z().Debug(msg)
		return
	}
	l.z().Debugw(msg, l.keysAndValues(args)...)
}

// Infot logs a message with structured fields at Info level.
// The structured fields are converted to key-value pairs for Zap.
func (l *zapLogger) Infot(msg string, args T) {
	if len(args) == 0 {
		l.z().Info(msg)
		return
	}
	l.z().Infow(msg, l.keysAndValues(args)...)
}

// Warnt logs a message with structured fields at Warn level.
// The structured fields are converted to key-value pairs for Zap.
func (l *zapLogger) Warnt(msg string, args T) {
	if len(args) == 0 {
		l.z().Warn(msg)
		return
	}
	l.z().Warnw(msg, l.keysAndValues(args)...)
}

// Errort logs a message with structured fields at Error level.
// The structured fields are converted to key-value pairs for Zap.
func (l *zapLogger) Errort(msg string, args T) {
	if len(args) == 0 {
		l.z().Error(msg)
		return
	}
	l.z().Errorw(msg, l.keysAndValues(args)...)
}

// keysAndValues converts a T (map[string]any) to a slice of alternating
// keys and values that Zap's structured logging methods expect.
// This method is optimized for performance and minimal allocations.
func (l *zapLogger) keysAndValues(args T) []any {
	if len(args) == 0 {
		return nil
	}

	// Pre-allocate slice with exact capacity to avoid reallocations
	keysAndValues := make([]any, 0, len(args)*2)

	// Convert map to alternating key-value pairs
	for k, v := range args {
		keysAndValues = append(keysAndValues, k, v)
	}

	return keysAndValues
}
