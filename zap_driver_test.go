package tslog

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewZapDriver tests the creation of Zap-based loggers
func TestNewZapDriver(t *testing.T) {
	t.Run("DefaultOptions", func(t *testing.T) {
		logger := NewZapDriver(nil)
		assert.NotNil(t, logger)

		// Test that the logger works
		var buf bytes.Buffer
		logger2 := NewZapDriver(&Options{
			lvl:     InfoLevel,
			w:       []io.Writer{&buf},
			encoder: EncoderJSON,
			caller:  false,
			driver:  NewZapDriver,
		})

		logger2.Info("test message")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("JSONEncoder", func(t *testing.T) {
		var buf bytes.Buffer
		opts := &Options{
			lvl:     DebugLevel,
			w:       []io.Writer{&buf},
			encoder: EncoderJSON,
			caller:  false,
			driver:  NewZapDriver,
		}

		logger := NewZapDriver(opts)
		logger.Info("json test")

		output := buf.String()
		assert.Contains(t, output, "json test")
		// JSON output should contain structured format
		assert.Contains(t, output, `"msg"`)
	})

	t.Run("ConsoleEncoder", func(t *testing.T) {
		var buf bytes.Buffer
		opts := &Options{
			lvl:     DebugLevel,
			w:       []io.Writer{&buf},
			encoder: EncoderConsole,
			caller:  false,
			driver:  NewZapDriver,
		}

		logger := NewZapDriver(opts)
		logger.Info("console test")

		output := buf.String()
		assert.Contains(t, output, "console test")
	})

	t.Run("WithCaller", func(t *testing.T) {
		var buf bytes.Buffer
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{&buf},
			encoder: EncoderJSON,
			caller:  true,
			driver:  NewZapDriver,
		}

		logger := NewZapDriver(opts)
		logger.Info("caller test")

		output := buf.String()
		assert.Contains(t, output, "caller test")
		// With caller enabled, should include file info
		assert.Contains(t, output, "caller")
	})

	t.Run("MultipleWriters", func(t *testing.T) {
		var buf1, buf2 bytes.Buffer
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{&buf1, &buf2},
			encoder: EncoderJSON,
			caller:  false,
			driver:  NewZapDriver,
		}

		logger := NewZapDriver(opts)
		logger.Info("multi writer test")

		// Both buffers should contain the message
		assert.Contains(t, buf1.String(), "multi writer test")
		assert.Contains(t, buf2.String(), "multi writer test")
	})

	t.Run("InvalidOptions", func(t *testing.T) {
		// Test with invalid options should fall back to defaults
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{}, // Invalid: no writers
			encoder: "invalid",     // Invalid encoder
			caller:  false,
			driver:  NewZapDriver,
		}

		logger := NewZapDriver(opts)
		assert.NotNil(t, logger)

		// Should still work (fallback to stdout)
		logger.Info("test message")
	})

	t.Run("NoWriters", func(t *testing.T) {
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{},
			encoder: EncoderJSON,
			caller:  false,
			driver:  NewZapDriver,
		}

		logger := NewZapDriver(opts)
		assert.NotNil(t, logger)

		// Should fallback to stdout and not panic
		logger.Info("no writers test")
	})

	t.Run("NilWriters", func(t *testing.T) {
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{nil, os.Stdout, nil},
			encoder: EncoderJSON,
			caller:  false,
			driver:  NewZapDriver,
		}

		logger := NewZapDriver(opts)
		assert.NotNil(t, logger)

		// Should skip nil writers and work with valid ones
		logger.Info("nil writers test")
	})
}

// TestZapLoggerLevels tests logging at different levels
func TestZapLoggerLevels(t *testing.T) {
	levels := []struct {
		level Level
		name  string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
	}

	for _, lvl := range levels {
		t.Run(lvl.name, func(t *testing.T) {
			var buf bytes.Buffer
			opts := &Options{
				lvl:     lvl.level,
				w:       []io.Writer{&buf},
				encoder: EncoderJSON,
				caller:  false,
				driver:  NewZapDriver,
			}

			logger := NewZapDriver(opts)

			// Test that the level works
			switch lvl.level {
			case DebugLevel:
				logger.Debug("debug message")
			case InfoLevel:
				logger.Info("info message")
			case WarnLevel:
				logger.Warn("warn message")
			case ErrorLevel:
				logger.Error("error message")
			}

			output := buf.String()
			assert.Contains(t, output, lvl.name+" message")
		})
	}
}

// TestZapLoggerMethods tests all logger interface methods
func TestZapLoggerMethods(t *testing.T) {
	var buf bytes.Buffer
	opts := &Options{
		lvl:     DebugLevel,
		w:       []io.Writer{&buf},
		encoder: EncoderJSON,
		caller:  false,
		driver:  NewZapDriver,
	}

	logger := NewZapDriver(opts)

	t.Run("BasicLogging", func(t *testing.T) {
		buf.Reset()
		logger.Debug("debug msg")
		logger.Info("info msg")
		logger.Warn("warn msg")
		logger.Error("error msg")

		output := buf.String()
		assert.Contains(t, output, "debug msg")
		assert.Contains(t, output, "info msg")
		assert.Contains(t, output, "warn msg")
		assert.Contains(t, output, "error msg")
	})

	t.Run("FormattedLogging", func(t *testing.T) {
		buf.Reset()
		logger.Debugf("debug %s %d", "formatted", 123)
		logger.Infof("info %s %d", "formatted", 456)
		logger.Warnf("warn %s %d", "formatted", 789)
		logger.Errorf("error %s %d", "formatted", 999)

		output := buf.String()
		assert.Contains(t, output, "debug formatted 123")
		assert.Contains(t, output, "info formatted 456")
		assert.Contains(t, output, "warn formatted 789")
		assert.Contains(t, output, "error formatted 999")
	})

	t.Run("StructuredLogging", func(t *testing.T) {
		buf.Reset()
		fields := T{"key1": "value1", "key2": 42, "key3": true}

		logger.Debugt("debug structured", fields)
		logger.Infot("info structured", fields)
		logger.Warnt("warn structured", fields)
		logger.Errort("error structured", fields)

		output := buf.String()
		assert.Contains(t, output, "debug structured")
		assert.Contains(t, output, "info structured")
		assert.Contains(t, output, "warn structured")
		assert.Contains(t, output, "error structured")
		assert.Contains(t, output, "key1")
		assert.Contains(t, output, "value1")
		assert.Contains(t, output, "key2")
		assert.Contains(t, output, "42")
	})

	t.Run("StructuredLoggingEmptyFields", func(t *testing.T) {
		buf.Reset()
		logger.Infot("empty fields", T{})
		assert.Contains(t, buf.String(), "empty fields")
	})

	t.Run("StructuredLoggingNilFields", func(t *testing.T) {
		buf.Reset()
		logger.Infot("nil fields", nil)
		assert.Contains(t, buf.String(), "nil fields")
	})
}

// TestZapLoggerKeysAndValues tests the keysAndValues optimization
func TestZapLoggerKeysAndValues(t *testing.T) {
	var buf bytes.Buffer
	opts := &Options{
		lvl:     InfoLevel,
		w:       []io.Writer{&buf},
		encoder: EncoderJSON,
		caller:  false,
		driver:  NewZapDriver,
	}

	zapLogger := NewZapDriver(opts).(*zapLogger)

	t.Run("EmptyFields", func(t *testing.T) {
		result := zapLogger.keysAndValues(T{})
		assert.Nil(t, result)
	})

	t.Run("NilFields", func(t *testing.T) {
		result := zapLogger.keysAndValues(nil)
		assert.Nil(t, result)
	})

	t.Run("ValidFields", func(t *testing.T) {
		fields := T{"key1": "value1", "key2": 42}
		result := zapLogger.keysAndValues(fields)

		assert.Equal(t, 4, len(result)) // 2 fields * 2 (key + value)

		// Check that all keys and values are present
		keys := make(map[string]bool)
		values := make(map[interface{}]bool)

		for i := 0; i < len(result); i += 2 {
			keys[result[i].(string)] = true
			values[result[i+1]] = true
		}

		assert.True(t, keys["key1"])
		assert.True(t, keys["key2"])
		assert.True(t, values["value1"])
		assert.True(t, values[42])
	})

	t.Run("LargeFields", func(t *testing.T) {
		fields := make(T)
		for i := 0; i < 100; i++ {
			fields[strings.Repeat("k", i+1)] = i
		}

		result := zapLogger.keysAndValues(fields)
		assert.Equal(t, 200, len(result)) // 100 fields * 2
	})
}

// TestZapLoggerConcurrency tests thread safety
func TestZapLoggerConcurrency(t *testing.T) {
	var buf bytes.Buffer
	opts := &Options{
		lvl:     InfoLevel,
		w:       []io.Writer{&buf},
		encoder: EncoderJSON,
		caller:  false,
		driver:  NewZapDriver,
	}

	logger := NewZapDriver(opts)

	var wg sync.WaitGroup
	numGoroutines := 50
	numMessages := 10

	// Test concurrent logging
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numMessages; j++ {
				logger.Infof("goroutine %d message %d", id, j)
				logger.Infot("structured", T{"goroutine": id, "message": j})
			}
		}(i)
	}

	wg.Wait()

	// Verify we have output
	output := buf.String()
	assert.True(t, len(output) > 0)
	assert.Contains(t, output, "goroutine")
	assert.Contains(t, output, "message")
}

// TestZapLoggerClose tests the Close method
func TestZapLoggerClose(t *testing.T) {
	var buf bytes.Buffer
	opts := &Options{
		lvl:     InfoLevel,
		w:       []io.Writer{&buf},
		encoder: EncoderJSON,
		caller:  false,
		driver:  NewZapDriver,
	}

	zapLogger := NewZapDriver(opts).(*zapLogger)

	// Should work before closing
	zapLogger.Info("before close")
	assert.Contains(t, buf.String(), "before close")

	// Close the logger
	err := zapLogger.Close()
	assert.NoError(t, err)

	// Second close should be safe
	err = zapLogger.Close()
	assert.NoError(t, err)

	// Using after close should panic
	assert.Panics(t, func() {
		zapLogger.Info("after close")
	})
}

// TestZapLoggerPanics tests panic conditions
func TestZapLoggerPanics(t *testing.T) {
	t.Run("UninitializedLogger", func(t *testing.T) {
		zapLogger := &zapLogger{}
		assert.Panics(t, func() {
			zapLogger.Info("test")
		})
	})

	t.Run("ClosedLogger", func(t *testing.T) {
		var buf bytes.Buffer
		opts := &Options{
			lvl:     InfoLevel,
			w:       []io.Writer{&buf},
			encoder: EncoderJSON,
			caller:  false,
			driver:  NewZapDriver,
		}

		zapLogger := NewZapDriver(opts).(*zapLogger)
		zapLogger.Close()

		assert.Panics(t, func() {
			zapLogger.Info("test")
		})
	})
}

// TestZapLoggerLevelMapping tests level mapping
func TestZapLoggerLevelMapping(t *testing.T) {
	tests := []struct {
		tslogLevel Level
		shouldLog  bool
		levelName  string
	}{
		{NoneLevel, false, "none"},
		{DebugLevel, true, "debug"},
		{InfoLevel, true, "info"},
		{WarnLevel, true, "warn"},
		{ErrorLevel, true, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.levelName, func(t *testing.T) {
			var buf bytes.Buffer
			opts := &Options{
				lvl:     tt.tslogLevel,
				w:       []io.Writer{&buf},
				encoder: EncoderJSON,
				caller:  false,
				driver:  NewZapDriver,
			}

			logger := NewZapDriver(opts)
			logger.Info("test message")

			if tt.shouldLog && tt.tslogLevel <= InfoLevel {
				assert.Contains(t, buf.String(), "test message")
			}
		})
	}
}

// BenchmarkZapLogger benchmarks Zap logger performance
func BenchmarkZapLogger(b *testing.B) {
	opts := &Options{
		lvl:     InfoLevel,
		w:       []io.Writer{io.Discard},
		encoder: EncoderJSON,
		caller:  false,
		driver:  NewZapDriver,
	}

	logger := NewZapDriver(opts)

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
		fields := T{"id": 123, "name": "benchmark"}
		for i := 0; i < b.N; i++ {
			logger.Infot("benchmark message", fields)
		}
	})

	b.Run("InfotLarge", func(b *testing.B) {
		b.ReportAllocs()
		fields := make(T)
		for i := 0; i < 10; i++ {
			fields[strings.Repeat("key", i+1)] = i
		}
		for i := 0; i < b.N; i++ {
			logger.Infot("benchmark message", fields)
		}
	})
}

// TestZapLevelMapping tests the zapLevel mapping
func TestZapLevelMapping(t *testing.T) {
	// Test that all tslog levels have corresponding zap levels
	for tslogLvl := NoneLevel; tslogLvl <= ErrorLevel; tslogLvl++ {
		zapLvl, exists := zapLevel[tslogLvl]
		assert.True(t, exists, "Level %v should have zap mapping", tslogLvl)
		assert.NotNil(t, zapLvl)
	}
}
