package logger

import "log"

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func New() *Logger {
	return &Logger{
		Info:  log.Default(),
		Error: log.Default(),
	}
}
