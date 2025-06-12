package tslog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewNoneLogger tests the NoneLogger constructor
func TestNewNoneLogger(t *testing.T) {
	logger := NewNoneLogger()
	assert.NotNil(t, logger)
	assert.IsType(t, &NoneLogger{}, logger)
}

// TestNewNoneDriver tests the NoneDriver function
func TestNewNoneDriver(t *testing.T) {
	logger := NewNoneDriver(nil)
	assert.NotNil(t, logger)
	assert.IsType(t, &NoneLogger{}, logger)

	// Test with actual options (should be ignored)
	opts := &Options{
		lvl:     InfoLevel,
		encoder: EncoderJSON,
		caller:  true,
	}
	logger2 := NewNoneDriver(opts)
	assert.NotNil(t, logger2)
	assert.IsType(t, &NoneLogger{}, logger2)
}

// TestNoneLoggerInterface tests that NoneLogger implements Logger interface
func TestNoneLoggerInterface(t *testing.T) {
	var logger Logger = &NoneLogger{}
	assert.NotNil(t, logger)

	// All methods should be callable without panic
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	logger.Debugf("debug %s", "formatted")
	logger.Infof("info %s", "formatted")
	logger.Warnf("warn %s", "formatted")
	logger.Errorf("error %s", "formatted")

	logger.Debugt("debug", T{"key": "value"})
	logger.Infot("info", T{"key": "value"})
	logger.Warnt("warn", T{"key": "value"})
	logger.Errort("error", T{"key": "value"})
}

// TestNoneLoggerMethods tests all NoneLogger methods individually
func TestNoneLoggerMethods(t *testing.T) {
	logger := &NoneLogger{}

	t.Run("BasicMethods", func(t *testing.T) {
		// These should not panic
		logger.Debug("test")
		logger.Info("test")
		logger.Warn("test")
		logger.Error("test")
	})

	t.Run("FormattedMethods", func(t *testing.T) {
		// These should not panic
		logger.Debugf("test %s %d", "string", 123)
		logger.Infof("test %s %d", "string", 123)
		logger.Warnf("test %s %d", "string", 123)
		logger.Errorf("test %s %d", "string", 123)
	})

	t.Run("StructuredMethods", func(t *testing.T) {
		fields := T{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		}

		// These should not panic
		logger.Debugt("test", fields)
		logger.Infot("test", fields)
		logger.Warnt("test", fields)
		logger.Errort("test", fields)
	})

	t.Run("StructuredMethodsNilFields", func(t *testing.T) {
		// These should not panic with nil fields
		logger.Debugt("test", nil)
		logger.Infot("test", nil)
		logger.Warnt("test", nil)
		logger.Errort("test", nil)
	})

	t.Run("StructuredMethodsEmptyFields", func(t *testing.T) {
		// These should not panic with empty fields
		logger.Debugt("test", T{})
		logger.Infot("test", T{})
		logger.Warnt("test", T{})
		logger.Errort("test", T{})
	})
}

// TestNoneLoggerVariableArguments tests methods with various argument types
func TestNoneLoggerVariableArguments(t *testing.T) {
	logger := &NoneLogger{}

	// Test with no arguments
	logger.Debug()
	logger.Info()
	logger.Warn()
	logger.Error()

	// Test with single argument
	logger.Debug("single")
	logger.Info("single")
	logger.Warn("single")
	logger.Error("single")

	// Test with multiple arguments
	logger.Debug("multiple", "args", 123, true, nil)
	logger.Info("multiple", "args", 123, true, nil)
	logger.Warn("multiple", "args", 123, true, nil)
	logger.Error("multiple", "args", 123, true, nil)

	// Test with different types
	type customType struct {
		field string
	}
	custom := customType{field: "test"}

	logger.Debug("custom type", custom)
	logger.Info("custom type", custom)
	logger.Warn("custom type", custom)
	logger.Error("custom type", custom)
}

// TestNoneLoggerFormattedArguments tests formatted methods with various arguments
func TestNoneLoggerFormattedArguments(t *testing.T) {
	logger := &NoneLogger{}

	// Test with no format arguments
	logger.Debugf("no args")
	logger.Infof("no args")
	logger.Warnf("no args")
	logger.Errorf("no args")

	// Test with mismatched format and arguments (should not panic)
	logger.Debugf("format %s %d")
	logger.Infof("format %s %d", "only one")
	logger.Warnf("format %s", "one", "two", "three")
	logger.Errorf("format", "extra", "args")

	// Test with nil arguments
	logger.Debugf("nil %v", nil)
	logger.Infof("nil %v", nil)
	logger.Warnf("nil %v", nil)
	logger.Errorf("nil %v", nil)
}

// TestNoneLoggerConcurrency tests thread safety of NoneLogger
func TestNoneLoggerConcurrency(t *testing.T) {
	logger := &NoneLogger{}

	// Run multiple goroutines concurrently
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func(id int) {
			// Call all methods multiple times
			for j := 0; j < 10; j++ {
				logger.Debug("debug", id, j)
				logger.Info("info", id, j)
				logger.Warn("warn", id, j)
				logger.Error("error", id, j)

				logger.Debugf("debugf %d %d", id, j)
				logger.Infof("infof %d %d", id, j)
				logger.Warnf("warnf %d %d", id, j)
				logger.Errorf("errorf %d %d", id, j)

				logger.Debugt("debugt", T{"id": id, "j": j})
				logger.Infot("infot", T{"id": id, "j": j})
				logger.Warnt("warnt", T{"id": id, "j": j})
				logger.Errort("errort", T{"id": id, "j": j})
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// If we reach here without panic, the test passes
}

// TestNoneLoggerAsDriver tests using NoneLogger as a driver
func TestNoneLoggerAsDriver(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(
		WithDriver(NewNoneDriver),
		WithWriter(&buf), // Add a writer to pass validation
	)
	assert.NotNil(t, logger)
	assert.IsType(t, &NoneLogger{}, logger)

	// Test that it works as expected (no output should be produced)
	logger.Info("test message")
	logger.Infof("test %s", "formatted")
	logger.Infot("test", T{"key": "value"})

	// Buffer should be empty since NoneLogger discards everything
	assert.Empty(t, buf.String())
}

// BenchmarkNoneLogger benchmarks NoneLogger performance
func BenchmarkNoneLogger(b *testing.B) {
	logger := &NoneLogger{}

	b.Run("Debug", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Debug("benchmark message")
		}
	})

	b.Run("Info", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Info("benchmark message")
		}
	})

	b.Run("Warn", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Warn("benchmark message")
		}
	})

	b.Run("Error", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Error("benchmark message")
		}
	})

	b.Run("Debugf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Debugf("benchmark message %d", i)
		}
	})

	b.Run("Infof", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Infof("benchmark message %d", i)
		}
	})

	b.Run("Warnf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Warnf("benchmark message %d", i)
		}
	})

	b.Run("Errorf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Errorf("benchmark message %d", i)
		}
	})

	b.Run("Debugt", func(b *testing.B) {
		b.ReportAllocs()
		fields := T{"id": 123, "name": "benchmark"}
		for i := 0; i < b.N; i++ {
			logger.Debugt("benchmark message", fields)
		}
	})

	b.Run("Infot", func(b *testing.B) {
		b.ReportAllocs()
		fields := T{"id": 123, "name": "benchmark"}
		for i := 0; i < b.N; i++ {
			logger.Infot("benchmark message", fields)
		}
	})

	b.Run("Warnt", func(b *testing.B) {
		b.ReportAllocs()
		fields := T{"id": 123, "name": "benchmark"}
		for i := 0; i < b.N; i++ {
			logger.Warnt("benchmark message", fields)
		}
	})

	b.Run("Errort", func(b *testing.B) {
		b.ReportAllocs()
		fields := T{"id": 123, "name": "benchmark"}
		for i := 0; i < b.N; i++ {
			logger.Errort("benchmark message", fields)
		}
	})
}

// TestNoneLoggerMemoryUsage tests that NoneLogger doesn't allocate memory
func TestNoneLoggerMemoryUsage(t *testing.T) {
	logger := &NoneLogger{}

	// Create a function that does logging
	logFunc := func() {
		logger.Info("test message")
		logger.Infof("test %s %d", "formatted", 123)
		logger.Infot("test", T{"key": "value", "num": 42})
	}

	// Measure allocations
	allocs := testing.AllocsPerRun(1000, logFunc)

	// NoneLogger should not allocate any memory
	assert.Equal(t, float64(0), allocs, "NoneLogger should not allocate memory")
}
