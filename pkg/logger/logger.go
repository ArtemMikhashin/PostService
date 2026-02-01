package logger

import "log"

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

// TODO: кидать дополнительно код ошибки
func New() *Logger {
	return &Logger{
		Info:  log.Default(),
		Error: log.Default(),
	}
}
