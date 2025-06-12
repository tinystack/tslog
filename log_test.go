package tslog

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLevel tests the Level type and its methods
func TestLevel(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
		enabled  bool
	}{
		{NoneLevel, "none", false},
		{DebugLevel, "debug", true},
		{InfoLevel, "info", true},
		{WarnLevel, "warn", true},
		{ErrorLevel, "error", true},
		{Level(99), "Level(99)", true}, // Unknown level
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
			assert.Equal(t, tt.enabled, tt.level.Enabled())
		})
	}
}

// TestParseLevel tests level parsing from strings
func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DebugLevel},
		{"DEBUG", DebugLevel},
		{"  debug  ", DebugLevel},
		{"info", InfoLevel},
		{"INFO", InfoLevel},
		{"warn", WarnLevel},
		{"WARN", WarnLevel},
		{"error", ErrorLevel},
		{"ERROR", ErrorLevel},
		{"none", NoneLevel},
		{"NONE", NoneLevel},
		{"invalid", NoneLevel},
		{"", NoneLevel},
		{"   ", NoneLevel},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestOptions tests the Options struct and its validation
func TestOptions(t *testing.T) {
	t.Run("ValidOptions", func(t *testing.T) {
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{os.Stdout},
			encoder: EncoderJSON,
			caller:  true,
			driver:  NewZapDriver,
		}
		assert.NoError(t, opts.Validate())
	})

	t.Run("NilDriver", func(t *testing.T) {
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{os.Stdout},
			encoder: EncoderJSON,
			caller:  true,
			driver:  nil,
		}
		assert.Error(t, opts.Validate())
		assert.Contains(t, opts.Validate().Error(), "driver cannot be nil")
	})

	t.Run("InvalidEncoder", func(t *testing.T) {
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{os.Stdout},
			encoder: "invalid",
			caller:  true,
			driver:  NewZapDriver,
		}
		assert.Error(t, opts.Validate())
		assert.Contains(t, opts.Validate().Error(), "encoder must be either")
	})

	t.Run("NoWriters", func(t *testing.T) {
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{},
			encoder: EncoderJSON,
			caller:  true,
			driver:  NewZapDriver,
		}
		assert.Error(t, opts.Validate())
		assert.Contains(t, opts.Validate().Error(), "at least one writer must be specified")
	})
}

// TestWithLevel tests the WithLevel option function
func TestWithLevel(t *testing.T) {
	opt := &Options{
		lvl: DebugLevel,
	}
	WithLevel(ErrorLevel)(opt)
	assert.Equal(t, ErrorLevel, opt.lvl)
}

// TestWithWriter tests the WithWriter option function
func TestWithWriter(t *testing.T) {
	t.Run("SingleWriter", func(t *testing.T) {
		opt := &Options{
			w: []io.Writer{os.Stdout},
		}
		WithWriter(os.Stderr)(opt)
		assert.Equal(t, []io.Writer{os.Stderr}, opt.w)
	})

	t.Run("MultipleWriters", func(t *testing.T) {
		opt := &Options{
			w: []io.Writer{},
		}
		WithWriter(os.Stderr, os.Stdout)(opt)
		assert.Equal(t, []io.Writer{os.Stderr, os.Stdout}, opt.w)
	})

	t.Run("WriterSliceCopy", func(t *testing.T) {
		original := []io.Writer{os.Stdout}
		opt := &Options{}
		WithWriter(original...)(opt)

		// Modify original slice
		original[0] = os.Stderr

		// Options should still have os.Stdout (copy was made)
		assert.Equal(t, os.Stdout, opt.w[0])
	})
}

// TestWithCaller tests the WithCaller option function
func TestWithCaller(t *testing.T) {
	opt := &Options{
		caller: true,
	}
	WithCaller(false)(opt)
	assert.Equal(t, false, opt.caller)
}

// TestWithEncoder tests the WithEncoder option function
func TestWithEncoder(t *testing.T) {
	opt := &Options{
		encoder: EncoderConsole,
	}
	WithEncoder(EncoderJSON)(opt)
	assert.Equal(t, EncoderJSON, opt.encoder)
}

// TestWithDriver tests the WithDriver option function
func TestWithDriver(t *testing.T) {
	dummyDriver := func(*Options) Logger { return &NoneLogger{} }
	opt := &Options{
		driver: NewZapDriver,
	}
	WithDriver(dummyDriver)(opt)
	assert.NotNil(t, opt.driver)
}

// TestNewLogger tests logger creation with various options
func TestNewLogger(t *testing.T) {
	t.Run("DefaultOptions", func(t *testing.T) {
		logger := NewLogger()
		assert.NotNil(t, logger)
	})

	t.Run("WithOptions", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(
			WithLevel(InfoLevel),
			WithWriter(&buf),
			WithEncoder(EncoderJSON),
			WithCaller(true),
		)
		assert.NotNil(t, logger)

		// Test that logger works
		logger.Info("test message")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("WithNilOption", func(t *testing.T) {
		logger := NewLogger(nil, WithLevel(InfoLevel), nil)
		assert.NotNil(t, logger)
	})

	t.Run("WithInvalidOptions", func(t *testing.T) {
		// This should fall back to defaults when validation fails
		logger := NewLogger(WithDriver(nil))
		assert.NotNil(t, logger)
	})
}

// TestDefaultOptions tests the defaultOptions function
func TestDefaultOptions(t *testing.T) {
	opts := defaultOptions()
	assert.Equal(t, DebugLevel, opts.lvl)
	assert.Equal(t, EncoderJSON, opts.encoder)
	assert.Equal(t, false, opts.caller)
	assert.NotNil(t, opts.driver)
}

// TestT tests the T type (map[string]any)
func TestT(t *testing.T) {
	fields := T{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	assert.Equal(t, "value1", fields["key1"])
	assert.Equal(t, 42, fields["key2"])
	assert.Equal(t, true, fields["key3"])
}

// TestLoggerInterface tests that our implementations satisfy the Logger interface
func TestLoggerInterface(t *testing.T) {
	t.Run("ZapLogger", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(WithWriter(&buf), WithDriver(NewZapDriver))

		// Test all Logger interface methods
		logger.Debug("debug")
		logger.Info("info")
		logger.Warn("warn")
		logger.Error("error")

		logger.Debugf("debug %s", "formatted")
		logger.Infof("info %s", "formatted")
		logger.Warnf("warn %s", "formatted")
		logger.Errorf("error %s", "formatted")

		logger.Debugt("debug structured", T{"key": "value"})
		logger.Infot("info structured", T{"key": "value"})
		logger.Warnt("warn structured", T{"key": "value"})
		logger.Errort("error structured", T{"key": "value"})

		output := buf.String()
		assert.Contains(t, output, "debug")
		assert.Contains(t, output, "info")
		assert.Contains(t, output, "warn")
		assert.Contains(t, output, "error")
	})

	t.Run("NoneLogger", func(t *testing.T) {
		var logger Logger = &NoneLogger{}

		// All methods should be no-ops and not panic
		logger.Debug("debug")
		logger.Info("info")
		logger.Warn("warn")
		logger.Error("error")

		logger.Debugf("debug %s", "formatted")
		logger.Infof("info %s", "formatted")
		logger.Warnf("warn %s", "formatted")
		logger.Errorf("error %s", "formatted")

		logger.Debugt("debug structured", T{"key": "value"})
		logger.Infot("info structured", T{"key": "value"})
		logger.Warnt("warn structured", T{"key": "value"})
		logger.Errort("error structured", T{"key": "value"})
	})
}

// TestEncoderTypes tests encoder constants
func TestEncoderTypes(t *testing.T) {
	assert.Equal(t, "json", EncoderJSON)
	assert.Equal(t, "console", EncoderConsole)
}

// TestConcurrentUsage tests thread safety
func TestConcurrentUsage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(WithWriter(&buf))

	// Run multiple goroutines concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Infof("Message from goroutine %d", id)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Check that we have output (exact content doesn't matter due to concurrency)
	assert.True(t, buf.Len() > 0)
}

// BenchmarkLogger benchmarks logger performance
func BenchmarkLogger(b *testing.B) {
	logger := NewLogger(WithWriter(io.Discard))

	b.Run("Info", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Info("benchmark message")
		}
	})

	b.Run("Infof", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Infof("benchmark message %d", i)
		}
	})

	b.Run("Infot", func(b *testing.B) {
		b.ReportAllocs()
		fields := T{"id": 123, "name": "test"}
		for i := 0; i < b.N; i++ {
			logger.Infot("benchmark message", fields)
		}
	})
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("EmptyStructuredFields", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(WithWriter(&buf))
		logger.Infot("message with empty fields", T{})
		assert.Contains(t, buf.String(), "message with empty fields")
	})

	t.Run("NilStructuredFields", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(WithWriter(&buf))
		logger.Infot("message with nil fields", nil)
		// Should not panic and should contain the message
		assert.Contains(t, buf.String(), "message with nil fields")
	})

	t.Run("LargeStructuredFields", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(WithWriter(&buf))

		// Create a large T with many fields
		fields := make(T)
		for i := 0; i < 100; i++ {
			fields[strings.Repeat("key", i+1)] = i
		}

		logger.Infot("message with many fields", fields)
		assert.Contains(t, buf.String(), "message with many fields")
	})
}
