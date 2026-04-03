package ports

type Logger interface {
	Info(msg string, keysAndValues ...any)
	Error(msg string, err error, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
}
