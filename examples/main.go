package main

import (
	"os"

	"github.com/tinystack/tslog"
)

func main() {

	// 使用预置Driver
	tslog.Debug("this is debug log message")
	tslog.Info("this is info log message")
	tslog.Warn("this is warn log message")
	tslog.Error("this is error log message")

	tslog.Debugf("this is debugf log message, %s: %s", "key1", "value1")
	tslog.Infof("this is infof log message, %s: %s", "key1", "value1")
	tslog.Warnf("this is warnf log message, %s: %s", "key1", "value1")
	tslog.Errorf("this is errorf log message, %s: %s", "key1", "value1")

	tslog.Debugt("this is debugt log message", tslog.T{"key1": "value1"})
	tslog.Infot("this is infot log message", tslog.T{"key1": "value1"})
	tslog.Warnt("this is warnt log message", tslog.T{"key1": "value1"})
	tslog.Errort("this is errort log message", tslog.T{"key1": "value1"})

	// 自定义Driver
	logger := tslog.NewLogger(
		tslog.WithLevel(tslog.ParseLevel("debug")),
		tslog.WithCaller(true),
		tslog.WithWriter(os.Stdout),
		tslog.WithEncoder(tslog.EncoderJSON),
		tslog.WithDriver(tslog.NewZapDriver),
	)

	logger.Debug("this is debug log message")
	logger.Info("this is info log message")
	logger.Warn("this is warn log message")
	logger.Error("this is error log message")

	logger.Debugf("this is debugf log message, %s: %s", "key1", "value1")
	logger.Infof("this is infof log message, %s: %s", "key1", "value1")
	logger.Warnf("this is warnf log message, %s: %s", "key1", "value1")
	logger.Errorf("this is errorf log message, %s: %s", "key1", "value1")

	logger.Debugt("this is debugt log message", tslog.T{"key1": "value1"})
	logger.Infot("this is infot log message", tslog.T{"key1": "value1"})
	logger.Warnt("this is warnt log message", tslog.T{"key1": "value1"})
	logger.Errort("this is errort log message", tslog.T{"key1": "value1"})
}
