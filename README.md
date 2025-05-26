# tslog
Go logger Package

## 安装

go get -u github.com/tinystack/tslog

## 示例

```go
import "github.com/tinystack/tslog"

// 使用默认实例
tslog.Debug("this is debug log message")
tslog.Info("this is info log message")
tslog.Warn("this is warn log message")
tslog.Error("this is error log message")

// 创建新实例
logger := tslog.NewLogger(
    tslog.WithLevel(tslog.ParseLevel("debug")),
    tslog.WithCaller(true),
    tslog.WithWriter(os.Stdout),
    tslog.WithEncoder(tslog.EncoderJSON),
)

logger.Debug("this is debug log message")
logger.Info("this is info log message")
logger.Warn("this is warn log message")
logger.Error("this is error log message")


// 更新全局默认实例
logger.UpdateDefaultLogger(logger)
```

## API

### Logger interface 

```golang
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

### 创建新Logger实例

```golang
NewLogger(funcOpts ...FuncOption) Logger
```