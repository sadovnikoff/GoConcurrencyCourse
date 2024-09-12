package common

import (
	"log"
	"os"
)

type Logger struct {
	ELog *log.Logger
	DLog *log.Logger
	ILog *log.Logger
}

func NewLogger() *Logger {
	logger := &Logger{
		ILog: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		DLog: log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile),
		ELog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return logger
}
