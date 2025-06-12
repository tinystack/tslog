// Package writer provides various io.Writer implementations for logging output.
// This file contains stderr writer functionality.
package writer

import (
	"io"
	"os"
)

// NewStderrWriter returns an io.Writer that writes to standard error.
// This is a simple wrapper around os.Stderr that provides a consistent
// interface for creating stderr writers in logging configurations.
//
// The returned writer is thread-safe as os.Stderr is thread-safe.
// stderr is typically used for error messages and diagnostic output.
//
// Example:
//
//	writer := writer.NewStderrWriter()
//	logger := tslog.NewLogger(tslog.WithWriter(writer))
func NewStderrWriter() io.Writer {
	return os.Stderr
}
