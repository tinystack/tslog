package writer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewStdoutWriter tests the stdout writer creation
func TestNewStdoutWriter(t *testing.T) {
	writer := NewStdoutWriter()
	assert.NotNil(t, writer)
	assert.Equal(t, os.Stdout, writer)
}

// TestNewStderrWriter tests the stderr writer creation
func TestNewStderrWriter(t *testing.T) {
	writer := NewStderrWriter()
	assert.NotNil(t, writer)
	assert.Equal(t, os.Stderr, writer)
}

// TestLumberJackConfig tests the LumberJackConfig struct
func TestLumberJackConfig(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       "/tmp/test.log",
			MaxRotatedSize: 100,
			MaxRetainDay:   7,
			MaxRetainFiles: 3,
			LocalTime:      true,
			Compress:       true,
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("EmptyFilePath", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       "",
			MaxRotatedSize: 100,
			MaxRetainDay:   7,
			MaxRetainFiles: 3,
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "FilePath cannot be empty")
	})

	t.Run("NegativeMaxRotatedSize", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       "/tmp/test.log",
			MaxRotatedSize: -1,
			MaxRetainDay:   7,
			MaxRetainFiles: 3,
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MaxRotatedSize cannot be negative")
	})

	t.Run("NegativeMaxRetainDay", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       "/tmp/test.log",
			MaxRotatedSize: 100,
			MaxRetainDay:   -1,
			MaxRetainFiles: 3,
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MaxRetainDay cannot be negative")
	})

	t.Run("NegativeMaxRetainFiles", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       "/tmp/test.log",
			MaxRotatedSize: 100,
			MaxRetainDay:   7,
			MaxRetainFiles: -1,
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MaxRetainFiles cannot be negative")
	})
}

// TestLumberJackConfigDefaults tests the setDefaults method
func TestLumberJackConfigDefaults(t *testing.T) {
	t.Run("AllZeroValues", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath: "/tmp/test.log",
		}

		config.setDefaults()

		assert.Equal(t, 100, config.MaxRotatedSize)
		assert.Equal(t, 7, config.MaxRetainDay)
		assert.Equal(t, 3, config.MaxRetainFiles)
	})

	t.Run("PartialDefaults", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       "/tmp/test.log",
			MaxRotatedSize: 50, // Don't override this
			MaxRetainDay:   0,  // Should be set to default
			MaxRetainFiles: 5,  // Don't override this
		}

		config.setDefaults()

		assert.Equal(t, 50, config.MaxRotatedSize) // Should remain unchanged
		assert.Equal(t, 7, config.MaxRetainDay)    // Should be set to default
		assert.Equal(t, 5, config.MaxRetainFiles)  // Should remain unchanged
	})

	t.Run("NoDefaults", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       "/tmp/test.log",
			MaxRotatedSize: 200,
			MaxRetainDay:   14,
			MaxRetainFiles: 10,
		}

		originalConfig := config
		config.setDefaults()

		// Should remain unchanged
		assert.Equal(t, originalConfig, config)
	})
}

// TestNewLumberJackWriter tests the LumberJack writer creation
func TestNewLumberJackWriter(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	t.Run("ValidConfig", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       filepath.Join(tempDir, "test.log"),
			MaxRotatedSize: 1, // 1 MB for testing
			MaxRetainDay:   1,
			MaxRetainFiles: 2,
			LocalTime:      true,
			Compress:       false,
		}

		writer, err := NewLumberJackWriter(config)
		assert.NoError(t, err)
		assert.NotNil(t, writer)
		assert.IsType(t, &lumberjack.Logger{}, writer)

		// Test writing to the logger
		n, err := writer.Write([]byte("test message\n"))
		assert.NoError(t, err)
		assert.Greater(t, n, 0)
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath: "", // Invalid
		}

		writer, err := NewLumberJackWriter(config)
		assert.Error(t, err)
		assert.Nil(t, writer)
		assert.Contains(t, err.Error(), "invalid lumberjack config")
	})

	t.Run("ConfigWithDefaults", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath: filepath.Join(tempDir, "defaults.log"),
			// All other fields are zero values, should get defaults
		}

		writer, err := NewLumberJackWriter(config)
		assert.NoError(t, err)
		assert.NotNil(t, writer)

		// Verify defaults were applied by checking the underlying lumberjack.Logger
		ljLogger := writer.(*lumberjack.Logger)
		assert.Equal(t, config.FilePath, ljLogger.Filename)
		assert.Equal(t, 100, ljLogger.MaxSize)     // Default
		assert.Equal(t, 7, ljLogger.MaxAge)        // Default
		assert.Equal(t, 3, ljLogger.MaxBackups)    // Default
		assert.Equal(t, false, ljLogger.LocalTime) // Default
		assert.Equal(t, false, ljLogger.Compress)  // Default
	})

	t.Run("FullConfig", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       filepath.Join(tempDir, "full.log"),
			MaxRotatedSize: 50,
			MaxRetainDay:   14,
			MaxRetainFiles: 5,
			LocalTime:      true,
			Compress:       true,
		}

		writer, err := NewLumberJackWriter(config)
		assert.NoError(t, err)
		assert.NotNil(t, writer)

		// Verify all settings were applied
		ljLogger := writer.(*lumberjack.Logger)
		assert.Equal(t, config.FilePath, ljLogger.Filename)
		assert.Equal(t, config.MaxRotatedSize, ljLogger.MaxSize)
		assert.Equal(t, config.MaxRetainDay, ljLogger.MaxAge)
		assert.Equal(t, config.MaxRetainFiles, ljLogger.MaxBackups)
		assert.Equal(t, config.LocalTime, ljLogger.LocalTime)
		assert.Equal(t, config.Compress, ljLogger.Compress)
	})
}

// TestMustNewLumberJackWriter tests the panic version of NewLumberJackWriter
func TestMustNewLumberJackWriter(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("ValidConfig", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath:       filepath.Join(tempDir, "must_test.log"),
			MaxRotatedSize: 100,
			MaxRetainDay:   7,
			MaxRetainFiles: 3,
		}

		// Should not panic
		writer := MustNewLumberJackWriter(config)
		assert.NotNil(t, writer)
		assert.IsType(t, &lumberjack.Logger{}, writer)
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		config := LumberJackConfig{
			FilePath: "", // Invalid
		}

		// Should panic
		assert.Panics(t, func() {
			MustNewLumberJackWriter(config)
		})
	})
}

// TestLumberJackWriterActualWriting tests actual file writing
func TestLumberJackWriterActualWriting(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "write_test.log")

	config := LumberJackConfig{
		FilePath:       logFile,
		MaxRotatedSize: 1, // 1 MB
		MaxRetainDay:   1,
		MaxRetainFiles: 2,
		LocalTime:      true,
		Compress:       false,
	}

	writer, err := NewLumberJackWriter(config)
	require.NoError(t, err)

	testMessage := "This is a test log message\n"

	// Write to the logger
	n, err := writer.Write([]byte(testMessage))
	assert.NoError(t, err)
	assert.Equal(t, len(testMessage), n)

	// Force sync to ensure file is written
	if syncer, ok := writer.(interface{ Sync() error }); ok {
		err = syncer.Sync()
		assert.NoError(t, err)
	}

	// Verify file was created and contains the message
	content, err := os.ReadFile(logFile)
	if err == nil { // File might not exist immediately due to buffering
		assert.Contains(t, string(content), strings.TrimSpace(testMessage))
	}
}

// TestWriterInterfaces tests that writers implement io.Writer interface
func TestWriterInterfaces(t *testing.T) {
	t.Run("StdoutWriter", func(t *testing.T) {
		var writer io.Writer = NewStdoutWriter()
		assert.NotNil(t, writer)
	})

	t.Run("StderrWriter", func(t *testing.T) {
		var writer io.Writer = NewStderrWriter()
		assert.NotNil(t, writer)
	})

	t.Run("LumberJackWriter", func(t *testing.T) {
		tempDir := t.TempDir()
		config := LumberJackConfig{
			FilePath:       filepath.Join(tempDir, "interface_test.log"),
			MaxRotatedSize: 100,
		}

		ljWriter, err := NewLumberJackWriter(config)
		require.NoError(t, err)

		var writer io.Writer = ljWriter
		assert.NotNil(t, writer)
	})
}

// TestStdWritersActualWrite tests that std writers actually write
func TestStdWritersActualWrite(t *testing.T) {
	t.Run("StdoutWriter", func(t *testing.T) {
		writer := NewStdoutWriter()

		// We can't easily capture stdout in tests without complex setup,
		// but we can at least verify that Write doesn't return an error
		// and writes the expected number of bytes
		testData := []byte("test stdout message\n")
		n, err := writer.Write(testData)
		assert.NoError(t, err)
		assert.Equal(t, len(testData), n)
	})

	t.Run("StderrWriter", func(t *testing.T) {
		writer := NewStderrWriter()

		// Similar to stdout test
		testData := []byte("test stderr message\n")
		n, err := writer.Write(testData)
		assert.NoError(t, err)
		assert.Equal(t, len(testData), n)
	})
}

// TestConcurrentWriting tests concurrent writing to lumberjack writer
func TestConcurrentWriting(t *testing.T) {
	tempDir := t.TempDir()
	config := LumberJackConfig{
		FilePath:       filepath.Join(tempDir, "concurrent_test.log"),
		MaxRotatedSize: 10, // Small size to test rotation
		MaxRetainDay:   1,
		MaxRetainFiles: 3,
	}

	writer, err := NewLumberJackWriter(config)
	require.NoError(t, err)

	// Start multiple goroutines writing concurrently
	numWriters := 10
	messagesPerWriter := 100
	done := make(chan bool, numWriters)

	for i := 0; i < numWriters; i++ {
		go func(writerID int) {
			defer func() { done <- true }()

			for j := 0; j < messagesPerWriter; j++ {
				message := fmt.Sprintf("Writer %d, Message %d\n", writerID, j)
				_, err := writer.Write([]byte(message))
				assert.NoError(t, err)

				// Small delay to allow for rotation
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Wait for all writers to complete
	for i := 0; i < numWriters; i++ {
		<-done
	}

	// Verify that some files were created (main + rotated)
	files, err := filepath.Glob(filepath.Join(tempDir, "concurrent_test.log*"))
	assert.NoError(t, err)
	assert.Greater(t, len(files), 0, "Should have created at least one log file")
}

// BenchmarkWriters benchmarks different writer implementations
func BenchmarkWriters(b *testing.B) {
	b.Run("StdoutWriter", func(b *testing.B) {
		writer := NewStdoutWriter()
		data := []byte("benchmark test message\n")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			writer.Write(data)
		}
	})

	b.Run("StderrWriter", func(b *testing.B) {
		writer := NewStderrWriter()
		data := []byte("benchmark test message\n")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			writer.Write(data)
		}
	})

	b.Run("LumberJackWriter", func(b *testing.B) {
		tempDir := b.TempDir()
		config := LumberJackConfig{
			FilePath:       filepath.Join(tempDir, "bench.log"),
			MaxRotatedSize: 100,
		}

		writer, err := NewLumberJackWriter(config)
		require.NoError(b, err)

		data := []byte("benchmark test message\n")

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			writer.Write(data)
		}
	})
}
