package logger

import (
	"os"
	"strings"

	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envLogOutput = "LOG_OUTPUT"
	envLogLevel  = "LOG_LEVEL"
)

type zapLogger struct {
	sugar *zap.SugaredLogger
}

func New() (ports.Logger, func(), error) {
	logConfig := zap.Config{
		OutputPaths: []string{getOutputLogs()},
		Level:       zap.NewAtomicLevelAt(getLevelLogs()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "message",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	l, err := logConfig.Build()
	if err != nil {
		return nil, nil, err
	}

	return &zapLogger{sugar: l.Sugar()}, func() { _ = l.Sync() }, nil
}

func (l *zapLogger) Info(msg string, keysAndValues ...any) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, err error, keysAndValues ...any) {
	fields := append(keysAndValues, "error", err)
	l.sugar.Errorw(msg, fields...)
}

func (l *zapLogger) Debug(msg string, keysAndValues ...any) {
	l.sugar.Debugw(msg, keysAndValues...)
}

func getOutputLogs() string {
	output := strings.ToLower(strings.TrimSpace(os.Getenv(envLogOutput)))
	if output == "" {
		return "stdout"
	}
	return output
}

func getLevelLogs() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(envLogLevel))) {
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel // fallback seguro
	}
}
