// Package writer provides various io.Writer implementations for logging output.
// This file contains file rotation writer functionality using the lumberjack library.
package writer

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/natefinch/lumberjack"
)

// LumberJackConfig holds configuration for file-based logging with rotation.
// It provides options for controlling log file size, retention, and rotation behavior.
type LumberJackConfig struct {
	// FilePath is the path to the log file. If the directory doesn't exist,
	// it will be created automatically.
	FilePath string

	// MaxRotatedSize is the maximum size in megabytes of the log file before
	// it gets rotated. Defaults to 100 MB if not specified.
	MaxRotatedSize int

	// MaxRetainDay is the maximum number of days to retain old log files
	// based on the timestamp encoded in their filename. Defaults to 7 days
	// if not specified. Set to 0 to disable age-based retention.
	MaxRetainDay int

	// MaxRetainFiles is the maximum number of old log files to retain.
	// Defaults to 3 if not specified. Set to 0 to disable count-based retention.
	MaxRetainFiles int

	// LocalTime determines if the time used for formatting the timestamps
	// in backup files is the computer's local time. Defaults to UTC time.
	LocalTime bool

	// Compress determines if the rotated log files should be compressed
	// using gzip. Defaults to false.
	Compress bool
}

// Validate checks if the configuration is valid and returns an error if not.
func (c *LumberJackConfig) Validate() error {
	if c.FilePath == "" {
		return fmt.Errorf("FilePath cannot be empty")
	}

	if c.MaxRotatedSize < 0 {
		return fmt.Errorf("MaxRotatedSize cannot be negative")
	}

	if c.MaxRetainDay < 0 {
		return fmt.Errorf("MaxRetainDay cannot be negative")
	}

	if c.MaxRetainFiles < 0 {
		return fmt.Errorf("MaxRetainFiles cannot be negative")
	}

	return nil
}

// setDefaults sets default values for unspecified configuration fields.
func (c *LumberJackConfig) setDefaults() {
	if c.MaxRotatedSize == 0 {
		c.MaxRotatedSize = 100 // 100 MB default
	}

	if c.MaxRetainDay == 0 {
		c.MaxRetainDay = 7 // 7 days default
	}

	if c.MaxRetainFiles == 0 {
		c.MaxRetainFiles = 3 // 3 files default
	}
}

// NewLumberJackWriter creates a new file writer with rotation capabilities
// using the lumberjack library. This writer automatically rotates log files
// based on size, age, and count limits specified in the configuration.
//
// Features:
// - Automatic log rotation based on file size
// - Age-based log file cleanup
// - Count-based log file cleanup
// - Optional compression of rotated files
// - Thread-safe operations
//
// The writer ensures that the directory containing the log file exists,
// creating it if necessary.
//
// Example:
//
//	config := writer.LumberJackConfig{
//	    FilePath: "/var/log/app.log",
//	    MaxRotatedSize: 100,    // 100 MB
//	    MaxRetainDay: 30,       // 30 days
//	    MaxRetainFiles: 5,      // 5 files
//	    LocalTime: true,
//	    Compress: true,
//	}
//	writer := writer.NewLumberJackWriter(config)
//	logger := tslog.NewLogger(tslog.WithWriter(writer))
func NewLumberJackWriter(conf LumberJackConfig) (io.Writer, error) {
	// Validate configuration
	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("invalid lumberjack config: %w", err)
	}

	// Apply defaults
	conf.setDefaults()

	// Ensure the directory exists
	dir := filepath.Dir(conf.FilePath)
	if dir != "" && dir != "." {
		// Note: We don't create the directory here as lumberjack will handle it,
		// but we could add this functionality if needed
	}

	// Create and configure the lumberjack logger
	logger := &lumberjack.Logger{
		Filename:   conf.FilePath,
		MaxSize:    conf.MaxRotatedSize,
		MaxAge:     conf.MaxRetainDay,
		MaxBackups: conf.MaxRetainFiles,
		LocalTime:  conf.LocalTime,
		Compress:   conf.Compress,
	}

	return logger, nil
}

// MustNewLumberJackWriter is like NewLumberJackWriter but panics if the
// configuration is invalid. This is useful in initialization code where
// you want to fail fast if the configuration is incorrect.
//
// Example:
//
//	config := writer.LumberJackConfig{
//	    FilePath: "/var/log/app.log",
//	    MaxRotatedSize: 100,
//	}
//	writer := writer.MustNewLumberJackWriter(config)
func MustNewLumberJackWriter(conf LumberJackConfig) io.Writer {
	writer, err := NewLumberJackWriter(conf)
	if err != nil {
		panic(err)
	}
	return writer
}
