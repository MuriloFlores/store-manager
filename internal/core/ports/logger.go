package ports

type Field map[string]interface{}

type Logger interface {
	ErrorLevel(message string, err error, tags ...Field)
	InfoLevel(message string, tags ...Field)
}
