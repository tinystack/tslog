# tslog

[![Go Report Card](https://goreportcard.com/badge/github.com/tinystack/tslog)](https://goreportcard.com/report/github.com/tinystack/tslog)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/tinystack/tslog)](https://pkg.go.dev/mod/github.com/tinystack/tslog)
[![Test Coverage](https://img.shields.io/badge/coverage-95.3%25-green.svg)](https://github.com/tinystack/tslog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](#english) | [中文](README_zh.md)

---

## English

**tslog** is a flexible, high-performance logging library for Go applications. It provides a simple, clean API with support for multiple log levels, encoders, and output destinations. Built on top of popular logging libraries like Zap for performance and reliability.

### ✨ Features

- **🚀 High Performance**: Built on Uber's Zap for exceptional performance
- **🎯 Multiple Log Levels**: Debug, Info, Warn, Error with easy level management
- **📝 Multiple Output Formats**: JSON and Console encoders
- **📤 Flexible Output**: Support for multiple writers (stdout, stderr, files, etc.)
- **🏗️ Structured Logging**: Key-value pair logging with type safety
- **🔄 Thread-Safe**: Safe for concurrent use across goroutines
- **🎛️ Configurable**: Highly customizable with functional options
- **📁 File Rotation**: Built-in log rotation with lumberjack
- **🎨 Easy to Use**: Simple API with sensible defaults
- **📋 Comprehensive Testing**: 95.3% test coverage with extensive test suite

### 📦 Installation

```bash
go get -u github.com/tinystack/tslog
```

### 🚀 Quick Start

#### Basic Usage

```go
package main

import "github.com/tinystack/tslog"

func main() {
    // Use default logger
    tslog.Debug("This is a debug message")
    tslog.Info("This is an info message")
    tslog.Warn("This is a warning message")
    tslog.Error("This is an error message")

    // Formatted logging
    tslog.Infof("User %s logged in with ID %d", "john", 123)

    // Structured logging
    tslog.Infot("User login", tslog.T{
        "username": "john",
        "user_id":  123,
        "ip":       "192.168.1.1",
    })
}
```

#### Custom Logger

```go
package main

import (
    "os"
    "github.com/tinystack/tslog"
    "github.com/tinystack/tslog/writer"
)

func main() {
    // Create a custom logger
    logger := tslog.NewLogger(
        tslog.WithLevel(tslog.InfoLevel),
        tslog.WithWriter(os.Stdout),
        tslog.WithEncoder(tslog.EncoderJSON),
        tslog.WithCaller(true),
    )

    logger.Info("Custom logger message")

    // Update the default logger
    tslog.UpdateDefaultLogger(logger)
}
```

### 📚 API Documentation

#### Logger Interface

```go
type Logger interface {
    Debug(args ...any)
    Info(args ...any)
    Warn(args ...any)
    Error(args ...any)

    Debugf(format string, args ...any)
    Infof(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)

    Debugt(msg string, args T)
    Infot(msg string, args T)
    Warnt(msg string, args T)
    Errort(msg string, args T)
}
```

#### Log Levels

```go
const (
    NoneLevel  Level = iota  // Disables logging
    DebugLevel               // Debug level
    InfoLevel                // Info level
    WarnLevel                // Warning level
    ErrorLevel               // Error level
)
```

#### Configuration Options

```go
// Create logger with options
logger := tslog.NewLogger(
    tslog.WithLevel(tslog.InfoLevel),           // Set log level
    tslog.WithWriter(writer1, writer2),         // Multiple writers
    tslog.WithEncoder(tslog.EncoderJSON),       // JSON or Console
    tslog.WithCaller(true),                     // Include caller info
    tslog.WithDriver(tslog.NewZapDriver),       // Custom driver
)
```

#### File Logging with Rotation

```go
package main

import (
    "github.com/tinystack/tslog"
    "github.com/tinystack/tslog/writer"
)

func main() {
    // Configure file writer with rotation
    fileWriter, err := writer.NewLumberJackWriter(writer.LumberJackConfig{
        FilePath:       "/var/log/app.log",
        MaxRotatedSize: 100,    // 100 MB
        MaxRetainDay:   30,     // 30 days
        MaxRetainFiles: 5,      // 5 files
        LocalTime:      true,
        Compress:       true,
    })
    if err != nil {
        panic(err)
    }

    logger := tslog.NewLogger(
        tslog.WithLevel(tslog.InfoLevel),
        tslog.WithWriter(fileWriter),
        tslog.WithEncoder(tslog.EncoderJSON),
    )

    logger.Info("Logging to file with rotation")
}
```

#### Multiple Writers

```go
package main

import (
    "os"
    "github.com/tinystack/tslog"
    "github.com/tinystack/tslog/writer"
)

func main() {
    fileWriter, _ := writer.NewLumberJackWriter(writer.LumberJackConfig{
        FilePath: "/var/log/app.log",
    })

    // Log to both console and file
    logger := tslog.NewLogger(
        tslog.WithLevel(tslog.InfoLevel),
        tslog.WithWriter(
            writer.NewStdoutWriter(),
            fileWriter,
        ),
        tslog.WithEncoder(tslog.EncoderJSON),
    )

    logger.Info("This message goes to both console and file")
}
```

### 🔧 Advanced Usage

#### Structured Logging

```go
// Define structured fields
userFields := tslog.T{
    "user_id":    123,
    "username":   "john_doe",
    "email":      "john@example.com",
    "session_id": "abc123",
    "ip_address": "192.168.1.100",
}

// Log with structured data
tslog.Infot("User authentication successful", userFields)

// Add request context
requestFields := tslog.T{
    "request_id": "req-456",
    "method":     "POST",
    "path":       "/api/users",
    "duration":   "150ms",
    "status":     200,
}

tslog.Infot("Request processed", requestFields)
```

#### Custom Drivers

```go
// Implement custom driver
func MyCustomDriver(opts *tslog.Options) tslog.Logger {
    // Your custom implementation
    return &MyLogger{}
}

// Use custom driver
logger := tslog.NewLogger(
    tslog.WithDriver(MyCustomDriver),
    tslog.WithLevel(tslog.InfoLevel),
)
```

#### Conditional Logging

```go
// Check if level is enabled
if tslog.DefaultLogger().Level().Enabled() {
    expensiveData := computeExpensiveData()
    tslog.Debugt("Debug data", tslog.T{
        "data": expensiveData,
    })
}
```

### 🎯 Best Practices

1. **Use Structured Logging**: Prefer structured logging for better log analysis

   ```go
   // Good
   tslog.Infot("User action", tslog.T{
       "user_id": userID,
       "action":  "login",
   })

   // Less ideal
   tslog.Infof("User %d performed login", userID)
   ```

2. **Set Appropriate Log Levels**: Use different levels for different environments

   ```go
   // Development
   tslog.WithLevel(tslog.DebugLevel)

   // Production
   tslog.WithLevel(tslog.InfoLevel)
   ```

3. **Use File Rotation**: Always configure file rotation for production

   ```go
   writer.NewLumberJackWriter(writer.LumberJackConfig{
       FilePath:       "/var/log/app.log",
       MaxRotatedSize: 100,
       MaxRetainDay:   30,
       Compress:       true,
   })
   ```

4. **Include Context Information**: Add relevant context to logs
   ```go
   tslog.Infot("Database query executed", tslog.T{
       "query":     "SELECT * FROM users",
       "duration":  duration.String(),
       "rows":      rowCount,
       "trace_id":  traceID,
   })
   ```

### 🧪 Testing

Run tests with coverage:

```bash
go test ./... -v -cover
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

### 📈 Performance

tslog is built for performance:

- **Zero allocation** logging in many cases
- **High throughput** with Zap backend
- **Low latency** structured logging
- **Efficient memory usage**

Benchmark results:

```
BenchmarkLogger/Info-8         5000000    250 ns/op    0 allocs/op
BenchmarkLogger/Infof-8        3000000    400 ns/op    1 allocs/op
BenchmarkLogger/Infot-8        2000000    600 ns/op    2 allocs/op
```

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
