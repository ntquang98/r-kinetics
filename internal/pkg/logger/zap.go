package logger

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = (*ZapLogger)(nil)

const (
	RequestIDKey     = "request.id"
	TraceIDKey       = "trace.id"
	SpanIDKey        = "span.id"
	TransactionIDKey = "transaction.id"
	UserKey          = "user_id"
	UserEmail        = "email"
	UserUID          = "user_uid"
)

type FilterLog func(key, value any) bool

// ZapLogger is a logger impl.
type ZapLogger struct {
	log    *zap.Logger
	Sync   func() error
	filter FilterLog
}

// NewZapLogger return a zap logger.
func NewZapLogger(encoder zapcore.EncoderConfig, level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
		), level)
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync, filter: defaultFilter()}
}

// Log Implementation of logger interface.
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	}
	return nil
}

func defaultFilter() FilterLog {
	return func(key, value any) bool {
		if key.(string) == RequestIDKey && value == "" {
			return true
		}
		if key.(string) == SpanIDKey && value == "" {
			return true
		}
		if key.(string) == TraceIDKey && value == "" {
			return true
		}
		if key.(string) == TransactionIDKey && value == "" {
			return true
		}
		if key.(string) == UserKey && value == "" {
			return true
		}
		return false
	}
}

func NewZapLoggerWrapper() *ZapLogger {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		FunctionKey:    "function",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logger := NewZapLogger(
		encoder,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.PanicLevel),
		),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)

	return logger
}
