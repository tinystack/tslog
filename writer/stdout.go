// Package writer provides various io.Writer implementations for logging output.
// This file contains stdout writer functionality.
package writer

import (
	"io"
	"os"
)

// NewStdoutWriter returns an io.Writer that writes to standard output.
// This is a simple wrapper around os.Stdout that provides a consistent
// interface for creating stdout writers in logging configurations.
//
// The returned writer is thread-safe as os.Stdout is thread-safe.
//
// Example:
//
//	writer := writer.NewStdoutWriter()
//	logger := tslog.NewLogger(tslog.WithWriter(writer))
func NewStdoutWriter() io.Writer {
	return os.Stdout
}
