# tslog

[![Go Report Card](https://goreportcard.com/badge/github.com/tinystack/tslog)](https://goreportcard.com/report/github.com/tinystack/tslog)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/tinystack/tslog)](https://pkg.go.dev/mod/github.com/tinystack/tslog)
[![Test Coverage](https://img.shields.io/badge/coverage-95.3%25-green.svg)](https://github.com/tinystack/tslog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](README.md) | [中文](#中文)

---

## 中文

**tslog** 是一个灵活、高性能的 Go 语言日志库。它提供简洁的 API，支持多种日志级别、编码器和输出目标。基于流行的日志库如 Zap 构建，确保性能和可靠性。

### ✨ 特性

- **🚀 高性能**: 基于 Uber 的 Zap 构建，性能卓越
- **🎯 多种日志级别**: Debug、Info、Warn、Error，级别管理简单
- **📝 多种输出格式**: JSON 和控制台编码器
- **📤 灵活输出**: 支持多种写入器（stdout、stderr、文件等）
- **🏗️ 结构化日志**: 类型安全的键值对日志记录
- **🔄 线程安全**: 在多个 goroutine 中并发使用安全
- **🎛️ 可配置**: 通过函数选项高度可定制
- **📁 文件轮转**: 内置 lumberjack 日志轮转
- **🎨 易于使用**: 简单的 API，合理的默认值
- **📋 全面测试**: 95.3% 测试覆盖率，广泛的测试套件

### 📦 安装

```bash
go get -u github.com/tinystack/tslog
```

### 🚀 快速开始

#### 基本使用

```go
package main

import "github.com/tinystack/tslog"

func main() {
    // 使用默认日志器
    tslog.Debug("这是一条调试消息")
    tslog.Info("这是一条信息消息")
    tslog.Warn("这是一条警告消息")
    tslog.Error("这是一条错误消息")

    // 格式化日志
    tslog.Infof("用户 %s 使用 ID %d 登录", "张三", 123)

    // 结构化日志
    tslog.Infot("用户登录", tslog.T{
        "username": "张三",
        "user_id":  123,
        "ip":       "192.168.1.1",
    })
}
```

#### 自定义日志器

```go
package main

import (
    "os"
    "github.com/tinystack/tslog"
    "github.com/tinystack/tslog/writer"
)

func main() {
    // 创建自定义日志器
    logger := tslog.NewLogger(
        tslog.WithLevel(tslog.InfoLevel),
        tslog.WithWriter(os.Stdout),
        tslog.WithEncoder(tslog.EncoderJSON),
        tslog.WithCaller(true),
    )

    logger.Info("自定义日志器消息")

    // 更新默认日志器
    tslog.UpdateDefaultLogger(logger)
}
```

### 📚 API 文档

#### 日志器接口

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

#### 日志级别

```go
const (
    NoneLevel  Level = iota  // 禁用日志
    DebugLevel               // 调试级别
    InfoLevel                // 信息级别
    WarnLevel                // 警告级别
    ErrorLevel               // 错误级别
)
```

#### 配置选项

```go
// 使用选项创建日志器
logger := tslog.NewLogger(
    tslog.WithLevel(tslog.InfoLevel),           // 设置日志级别
    tslog.WithWriter(writer1, writer2),         // 多个写入器
    tslog.WithEncoder(tslog.EncoderJSON),       // JSON 或控制台格式
    tslog.WithCaller(true),                     // 包含调用者信息
    tslog.WithDriver(tslog.NewZapDriver),       // 自定义驱动器
)
```

#### 文件日志与轮转

```go
package main

import (
    "github.com/tinystack/tslog"
    "github.com/tinystack/tslog/writer"
)

func main() {
    // 配置带轮转的文件写入器
    fileWriter, err := writer.NewLumberJackWriter(writer.LumberJackConfig{
        FilePath:       "/var/log/app.log",
        MaxRotatedSize: 100,    // 100 MB
        MaxRetainDay:   30,     // 30 天
        MaxRetainFiles: 5,      // 5 个文件
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

    logger.Info("记录到带轮转的文件")
}
```

#### 多个写入器

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

    // 同时记录到控制台和文件
    logger := tslog.NewLogger(
        tslog.WithLevel(tslog.InfoLevel),
        tslog.WithWriter(
            writer.NewStdoutWriter(),
            fileWriter,
        ),
        tslog.WithEncoder(tslog.EncoderJSON),
    )

    logger.Info("此消息同时输出到控制台和文件")
}
```

### 🎯 高级用法

#### 禁用日志（零开销）

```go
package main

import "github.com/tinystack/tslog"

func main() {
    // 创建空日志器，真正的零分配和零开销
    logger := tslog.NewNoneLogger()

    // 这些调用不会有任何开销
    logger.Debug("不会输出")
    logger.Info("不会输出")
    logger.Warn("不会输出")
    logger.Error("不会输出")
}
```

#### 条件日志

```go
package main

import "github.com/tinystack/tslog"

func main() {
    logger := tslog.NewLogger(
        tslog.WithLevel(tslog.InfoLevel),
    )

    // 检查日志级别以避免昂贵的操作
    if tslog.DebugLevel.Enabled(logger) {
        expensiveDebugInfo := computeExpensiveDebugInfo()
        logger.Debug("Debug info", tslog.T{"info": expensiveDebugInfo})
    }
}

func computeExpensiveDebugInfo() string {
    // 某些昂贵的计算
    return "expensive debug information"
}
```

#### 错误处理

```go
package main

import (
    "errors"
    "github.com/tinystack/tslog"
)

func main() {
    logger := tslog.NewLogger()

    err := someOperation()
    if err != nil {
        logger.Errort("操作失败", tslog.T{
            "error":     err.Error(),
            "operation": "someOperation",
            "retry":     false,
        })
    }
}

func someOperation() error {
    return errors.New("something went wrong")
}
```

### 📊 性能基准

基于我们的基准测试：

```
BenchmarkNoneLogger/Info-12          1000000000      0.31 ns/op       0 B/op       0 allocs/op
BenchmarkZapLogger/Info-12              2453881       516 ns/op       16 B/op       1 allocs/op
BenchmarkZapLogger/Infof-12             1848524       731 ns/op       51 B/op       3 allocs/op
BenchmarkZapLogger/Infot-12             1000000      1143 ns/op      352 B/op       4 allocs/op
```

- **NoneLogger**: 真正的零开销（0.31ns，0 分配）
- **基本日志**: ~516ns/op，非常高效
- **格式化日志**: ~731ns/op，良好性能
- **结构化日志**: ~1143ns/op，合理性能

### 🎛️ 配置参考

#### Writer 配置

```go
// 标准输出写入器
stdoutWriter := writer.NewStdoutWriter()

// 标准错误写入器
stderrWriter := writer.NewStderrWriter()

// LumberJack 文件写入器（带轮转）
fileWriter, err := writer.NewLumberJackWriter(writer.LumberJackConfig{
    FilePath:       "/var/log/app.log",  // 日志文件路径
    MaxRotatedSize: 100,                 // 轮转大小（MB）
    MaxRetainDay:   30,                  // 保留天数
    MaxRetainFiles: 5,                   // 保留文件数
    LocalTime:      true,                // 使用本地时间
    Compress:       true,                // 压缩旧文件
})
```

#### 编码器类型

```go
const (
    EncoderJSON    = "json"       // JSON 格式输出
    EncoderConsole = "console"    // 控制台格式输出
)
```

#### 驱动器选择

```go
// Zap 驱动器（默认，高性能）
zapDriver := tslog.NewZapDriver

// 空驱动器（零开销）
noneDriver := tslog.NewNoneDriver
```

### 🔧 最佳实践

#### 1. 使用结构化日志

```go
// 推荐：结构化日志
tslog.Infot("用户登录", tslog.T{
    "user_id":   123,
    "username":  "张三",
    "ip":        "192.168.1.1",
    "timestamp": time.Now(),
})

// 不推荐：字符串拼接
tslog.Infof("用户 %s (ID: %d) 从 %s 登录", "张三", 123, "192.168.1.1")
```

#### 2. 合理使用日志级别

```go
// Debug: 详细的调试信息，通常仅在开发环境启用
tslog.Debug("处理用户请求开始")

// Info: 一般信息，记录应用程序的正常运行
tslog.Info("用户成功登录")

// Warn: 警告信息，可能的问题但不影响主要功能
tslog.Warn("数据库连接缓慢")

// Error: 错误信息，需要关注的问题
tslog.Error("无法连接到数据库")
```

#### 3. 生产环境配置

```go
func NewProductionLogger() tslog.Logger {
    fileWriter, _ := writer.NewLumberJackWriter(writer.LumberJackConfig{
        FilePath:       "/var/log/app.log",
        MaxRotatedSize: 100,     // 100MB 轮转
        MaxRetainDay:   7,       // 保留 7 天
        MaxRetainFiles: 5,       // 最多 5 个文件
        LocalTime:      true,
        Compress:       true,    // 压缩旧日志
    })

    return tslog.NewLogger(
        tslog.WithLevel(tslog.InfoLevel),     // 生产环境使用 Info 级别
        tslog.WithWriter(fileWriter),
        tslog.WithEncoder(tslog.EncoderJSON), // JSON 格式便于解析
        tslog.WithCaller(false),              // 生产环境禁用调用者信息以提高性能
    )
}
```

#### 4. 开发环境配置

```go
func NewDevelopmentLogger() tslog.Logger {
    return tslog.NewLogger(
        tslog.WithLevel(tslog.DebugLevel),      // 开发环境启用 Debug
        tslog.WithWriter(writer.NewStdoutWriter()),
        tslog.WithEncoder(tslog.EncoderConsole), // 控制台格式更易读
        tslog.WithCaller(true),                  // 启用调用者信息便于调试
    )
}
```

### 🧪 测试

运行测试套件：

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行基准测试
go test -bench=. -benchmem

# 详细测试输出
go test -v ./...
```

当前测试覆盖率：

- 主包：**95.3%**
- writer 包：**100%**

### 📄 许可证

本项目使用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

### 🤝 贡献

欢迎贡献！请阅读我们的贡献指南：

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

### 📞 支持

如果你遇到任何问题或有疑问，请：

1. 查看 [文档](https://pkg.go.dev/github.com/tinystack/tslog)
2. 搜索现有的 [Issues](https://github.com/tinystack/tslog/issues)
3. 创建新的 [Issue](https://github.com/tinystack/tslog/issues/new)

---

**tslog** - 让 Go 日志记录变得简单而强大！ 🚀
