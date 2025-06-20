package logger

import (
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

type log struct {
	Logger *zap.Logger
}

func NewLogger() ports.Logger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	fileWriter := zapcore.AddSync(lumberjackLogger)
	consoleWriter := zapcore.AddSync(os.Stdout)

	logLevel := zap.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, logLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, logLevel),
	)

	logger := zap.New(core, zap.AddCaller())

	return &log{
		Logger: logger,
	}
}

func (l *log) InfoLevel(message string, tags ...ports.Field) {
	zapFields := fieldsToZapFields(tags)

	l.Logger.Info(message, zapFields...)
	l.Logger.Sync()
}
func (l *log) ErrorLevel(message string, err error, tags ...ports.Field) {
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

func fieldsToZapFields(fields []ports.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		for k, v := range field {
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}

	return zapFields
}
