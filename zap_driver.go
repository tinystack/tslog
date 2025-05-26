package tslog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	zap *zap.SugaredLogger
}

var zapLevel = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

func NewZapDriver(opts *Options) Logger {
	if opts == nil {
		opts = defaultOptions()
	}

	lvl := zap.InfoLevel
	if l, ok := zapLevel[opts.lvl]; ok {
		lvl = l
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(lvl)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.FunctionKey = "fn"
	encoderConfig.MessageKey = "msg"

	var encoder zapcore.Encoder
	switch opts.encoder {
	case EncoderConsole:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var syncer []zapcore.WriteSyncer
	for _, w := range opts.w {
		syncer = append(syncer, zapcore.AddSync(w))
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(syncer...),
		atomicLevel,
	)
	z := zap.New(core, zap.WithCaller(opts.caller), zap.AddCallerSkip(2)).Sugar()
	return &zapLogger{
		zap: z,
	}
}

func (l *zapLogger) z() *zap.SugaredLogger {
	if l.zap == nil {
		panic("logger: Zap.zap not initialized")
	}
	return l.zap
}

func (l *zapLogger) Debug(args ...any) {
	l.z().Debug(args...)
}

func (l *zapLogger) Info(args ...any) {
	l.z().Info(args...)
}

func (l *zapLogger) Warn(args ...any) {
	l.z().Warn(args...)
}

func (l *zapLogger) Error(args ...any) {
	l.z().Error(args...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.z().Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...any) {
	l.z().Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.z().Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.z().Errorf(format, args...)
}

func (l *zapLogger) Debugt(msg string, args T) {
	l.z().Debugw(msg, l.keysAndValues(args)...)
}

func (l *zapLogger) Infot(msg string, args T) {
	l.z().Infow(msg, l.keysAndValues(args)...)
}

func (l *zapLogger) Warnt(msg string, args T) {
	l.z().Warnw(msg, l.keysAndValues(args)...)
}

func (l *zapLogger) Errort(msg string, args T) {
	l.z().Errorw(msg, l.keysAndValues(args)...)
}

func (l *zapLogger) keysAndValues(args T) []any {
	var keysAndValues []any
	for k, v := range args {
		keysAndValues = append(keysAndValues, k, v)
	}
	return keysAndValues
}
