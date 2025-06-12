package tslog

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDefaultLogger tests the default logger functionality
func TestDefaultLogger(t *testing.T) {
	// Get original default logger to restore later
	originalLogger := DefaultLogger()
	defer func() {
		UpdateDefaultLogger(originalLogger)
	}()

	t.Run("DefaultLoggerNotNil", func(t *testing.T) {
		logger := DefaultLogger()
		assert.NotNil(t, logger)
	})

	t.Run("DefaultLoggerThreadSafe", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make([]Logger, 10)

		// Access default logger from multiple goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = DefaultLogger()
			}(i)
		}

		wg.Wait()

		// All should return the same logger instance
		first := results[0]
		for i := 1; i < len(results); i++ {
			assert.Equal(t, first, results[i])
		}
	})
}

// TestUpdateDefaultLogger tests updating the default logger
func TestUpdateDefaultLogger(t *testing.T) {
	// Get original default logger to restore later
	originalLogger := DefaultLogger()
	defer func() {
		UpdateDefaultLogger(originalLogger)
	}()

	t.Run("UpdateWithValidLogger", func(t *testing.T) {
		var buf bytes.Buffer
		newLogger := NewLogger(WithWriter(&buf))

		UpdateDefaultLogger(newLogger)

		// Test that the default logger has been updated
		updatedLogger := DefaultLogger()
		assert.Equal(t, newLogger, updatedLogger)

		// Test that package-level functions use the new logger
		Info("test message")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("UpdateWithNilLogger", func(t *testing.T) {
		currentLogger := DefaultLogger()

		// Updating with nil should be ignored
		UpdateDefaultLogger(nil)

		// Default logger should remain unchanged
		assert.Equal(t, currentLogger, DefaultLogger())
	})

	t.Run("ConcurrentUpdate", func(t *testing.T) {
		var wg sync.WaitGroup
		loggers := make([]Logger, 5)

		// Create multiple loggers
		for i := 0; i < 5; i++ {
			var buf bytes.Buffer
			loggers[i] = NewLogger(WithWriter(&buf))
		}

		// Update default logger concurrently
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				UpdateDefaultLogger(loggers[idx])
			}(i)
		}

		wg.Wait()

		// Should not panic and default logger should be one of the loggers
		currentLogger := DefaultLogger()
		assert.NotNil(t, currentLogger)

		// Verify it's one of our loggers
		found := false
		for _, logger := range loggers {
			if currentLogger == logger {
				found = true
				break
			}
		}
		assert.True(t, found || currentLogger == originalLogger)
	})
}

// TestPackageLevelFunctions tests all package-level logging functions
func TestPackageLevelFunctions(t *testing.T) {
	// Get original default logger to restore later
	originalLogger := DefaultLogger()
	defer func() {
		UpdateDefaultLogger(originalLogger)
	}()

	var buf bytes.Buffer
	testLogger := NewLogger(WithWriter(&buf))
	UpdateDefaultLogger(testLogger)

	t.Run("DebugFunction", func(t *testing.T) {
		buf.Reset()
		Debug("debug message")
		assert.Contains(t, buf.String(), "debug message")
	})

	t.Run("InfoFunction", func(t *testing.T) {
		buf.Reset()
		Info("info message")
		assert.Contains(t, buf.String(), "info message")
	})

	t.Run("WarnFunction", func(t *testing.T) {
		buf.Reset()
		Warn("warn message")
		assert.Contains(t, buf.String(), "warn message")
	})

	t.Run("ErrorFunction", func(t *testing.T) {
		buf.Reset()
		Error("error message")
		assert.Contains(t, buf.String(), "error message")
	})

	t.Run("DebugfFunction", func(t *testing.T) {
		buf.Reset()
		Debugf("debug %s %d", "formatted", 123)
		output := buf.String()
		assert.Contains(t, output, "debug")
		assert.Contains(t, output, "formatted")
		assert.Contains(t, output, "123")
	})

	t.Run("InfofFunction", func(t *testing.T) {
		buf.Reset()
		Infof("info %s %d", "formatted", 456)
		output := buf.String()
		assert.Contains(t, output, "info")
		assert.Contains(t, output, "formatted")
		assert.Contains(t, output, "456")
	})

	t.Run("WarnfFunction", func(t *testing.T) {
		buf.Reset()
		Warnf("warn %s %d", "formatted", 789)
		output := buf.String()
		assert.Contains(t, output, "warn")
		assert.Contains(t, output, "formatted")
		assert.Contains(t, output, "789")
	})

	t.Run("ErrorfFunction", func(t *testing.T) {
		buf.Reset()
		Errorf("error %s %d", "formatted", 999)
		output := buf.String()
		assert.Contains(t, output, "error")
		assert.Contains(t, output, "formatted")
		assert.Contains(t, output, "999")
	})

	t.Run("DebugtFunction", func(t *testing.T) {
		buf.Reset()
		Debugt("debug structured", T{"key": "value", "num": 123})
		output := buf.String()
		assert.Contains(t, output, "debug structured")
		assert.Contains(t, output, "key")
		assert.Contains(t, output, "value")
	})

	t.Run("InfotFunction", func(t *testing.T) {
		buf.Reset()
		Infot("info structured", T{"key": "value", "num": 456})
		output := buf.String()
		assert.Contains(t, output, "info structured")
		assert.Contains(t, output, "key")
		assert.Contains(t, output, "value")
	})

	t.Run("WarntFunction", func(t *testing.T) {
		buf.Reset()
		Warnt("warn structured", T{"key": "value", "num": 789})
		output := buf.String()
		assert.Contains(t, output, "warn structured")
		assert.Contains(t, output, "key")
		assert.Contains(t, output, "value")
	})

	t.Run("ErrortFunction", func(t *testing.T) {
		buf.Reset()
		Errort("error structured", T{"key": "value", "num": 999})
		output := buf.String()
		assert.Contains(t, output, "error structured")
		assert.Contains(t, output, "key")
		assert.Contains(t, output, "value")
	})
}

// TestPackageLevelFunctionsConcurrency tests concurrent usage of package-level functions
func TestPackageLevelFunctionsConcurrency(t *testing.T) {
	// Get original default logger to restore later
	originalLogger := DefaultLogger()
	defer func() {
		UpdateDefaultLogger(originalLogger)
	}()

	var buf bytes.Buffer
	testLogger := NewLogger(WithWriter(&buf))
	UpdateDefaultLogger(testLogger)

	var wg sync.WaitGroup
	numGoroutines := 20
	numCallsPerGoroutine := 10

	// Test concurrent usage of all package-level functions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numCallsPerGoroutine; j++ {
				Debug("debug", id, j)
				Info("info", id, j)
				Warn("warn", id, j)
				Error("error", id, j)

				Debugf("debugf %d %d", id, j)
				Infof("infof %d %d", id, j)
				Warnf("warnf %d %d", id, j)
				Errorf("errorf %d %d", id, j)

				Debugt("debugt", T{"id": id, "j": j})
				Infot("infot", T{"id": id, "j": j})
				Warnt("warnt", T{"id": id, "j": j})
				Errort("errort", T{"id": id, "j": j})
			}
		}(i)
	}

	wg.Wait()

	// Verify we have output (exact content checking is difficult due to concurrency)
	output := buf.String()
	assert.True(t, len(output) > 0, "Expected output from concurrent logging")
	assert.Contains(t, output, "debug")
	assert.Contains(t, output, "info")
	assert.Contains(t, output, "warn")
	assert.Contains(t, output, "error")
}

// TestDefaultLoggerInitialization tests that the default logger is properly initialized
func TestDefaultLoggerInitialization(t *testing.T) {
	// This test verifies that the init() function works correctly
	logger := DefaultLogger()
	assert.NotNil(t, logger)

	// Test that we can use it immediately
	var buf bytes.Buffer
	testLogger := NewLogger(WithWriter(&buf))
	UpdateDefaultLogger(testLogger)

	Info("initialization test")
	assert.Contains(t, buf.String(), "initialization test")
}

// TestPackageLevelFunctionsWithEmptyFields tests package functions with empty structured fields
func TestPackageLevelFunctionsWithEmptyFields(t *testing.T) {
	// Get original default logger to restore later
	originalLogger := DefaultLogger()
	defer func() {
		UpdateDefaultLogger(originalLogger)
	}()

	var buf bytes.Buffer
	testLogger := NewLogger(WithWriter(&buf))
	UpdateDefaultLogger(testLogger)

	t.Run("EmptyFields", func(t *testing.T) {
		buf.Reset()
		Infot("message with empty fields", T{})
		assert.Contains(t, buf.String(), "message with empty fields")
	})

	t.Run("NilFields", func(t *testing.T) {
		buf.Reset()
		Infot("message with nil fields", nil)
		assert.Contains(t, buf.String(), "message with nil fields")
	})
}

// BenchmarkPackageLevelFunctions benchmarks the performance of package-level functions
func BenchmarkPackageLevelFunctions(b *testing.B) {
	// Use a no-op writer for benchmarking
	originalLogger := DefaultLogger()
	defer func() {
		UpdateDefaultLogger(originalLogger)
	}()

	testLogger := NewLogger(WithWriter(io.Discard))
	UpdateDefaultLogger(testLogger)

	b.Run("Info", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Info("benchmark message")
		}
	})

	b.Run("Infof", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Infof("benchmark message %d", i)
		}
	})

	b.Run("Infot", func(b *testing.B) {
		b.ReportAllocs()
		fields := T{"id": 123, "name": "benchmark"}
		for i := 0; i < b.N; i++ {
			Infot("benchmark message", fields)
		}
	})
}

// TestDefaultLoggerVariables tests the package-level variables
func TestDefaultLoggerVariables(t *testing.T) {
	assert.Equal(t, DebugLevel, defaultLogLevel)
}
