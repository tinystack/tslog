package tslog

type NoneLogger struct{}

func (*NoneLogger) Debug(args ...interface{}) {}
func (*NoneLogger) Info(args ...interface{})  {}
func (*NoneLogger) Warn(args ...interface{})  {}
func (*NoneLogger) Error(args ...interface{}) {}

func (*NoneLogger) Debugf(format string, args ...interface{}) {}
func (*NoneLogger) Infof(format string, args ...interface{})  {}
func (*NoneLogger) Warnf(format string, args ...interface{})  {}
func (*NoneLogger) Errorf(format string, args ...interface{}) {}

func (*NoneLogger) Debugt(msg string, args T) {}
func (*NoneLogger) Infot(msg string, args T)  {}
func (*NoneLogger) Warnt(msg string, args T)  {}
func (*NoneLogger) Errort(msg string, args T) {}
