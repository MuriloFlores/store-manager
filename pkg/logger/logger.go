package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type LogInterface interface {
	ErrorLevel(message string, err error, tags ...Field)
	InfoLevel(message string, tags ...Field)
}

type Field map[string]interface{}

type log struct {
	Logger *zap.Logger
}

func NewLogger(LogOutput, LogLevel string) LogInterface {
	logConfig := zap.Config{
		OutputPaths: []string{getOutputLogs(LogOutput)},
		Level:       zap.NewAtomicLevelAt(getLevelLogs(LogLevel)),
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

	Logger, _ := logConfig.Build()

	return &log{
		Logger: Logger,
	}
}

func (l *log) InfoLevel(message string, tags ...Field) {
	zapFields := fieldsToZapFields(tags)

	l.Logger.Info(message, zapFields...)
	l.Logger.Sync()
}
func (l *log) ErrorLevel(message string, err error, tags ...Field) {
	zapFields := fieldsToZapFields(tags)
	zapFields = append(zapFields, zap.NamedError("error", err))

	l.Logger.Error(message, zapFields...)
	l.Logger.Sync()
}

func getOutputLogs(OutputLog string) string {
	output := strings.ToLower(strings.TrimSpace(OutputLog))
	if output == "" {
		return "stdout"
	}

	return output
}

func getLevelLogs(LevelLog string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(LevelLog)) {
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.ErrorLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func fieldsToZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		for k, v := range field {
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}

	return zapFields
}
