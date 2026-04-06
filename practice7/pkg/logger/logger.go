package logger

type Interface interface {
    Info(args ...interface{})
    Error(args ...interface{})
}

type Logger struct{}

func New() *Logger { return &Logger{} }
func (l *Logger) Info(args ...interface{}) {}
func (l *Logger) Error(args ...interface{}) {}
