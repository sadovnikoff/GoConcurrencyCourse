package common

import (
	"errors"
	"io"
	"log"
	"os"
)

const (
	infoLevel  = "info"
	debugLevel = "debug"
	errorLevel = "error"
)

type Logger struct {
	eLog *log.Logger
	dLog *log.Logger
	iLog *log.Logger
	fd   *os.File
}

func NewLogger(level, output string) (*Logger, error) {
	logger := &Logger{}

	iLog := log.New(io.Discard, "", 0)
	dLog := log.New(io.Discard, "", 0)
	eLog := log.New(io.Discard, "", 0)

	out := os.Stdout
	if output != "" {
		f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		out = f
		logger.fd = f
	}

	switch level {
	case "":

	case infoLevel:
		iLog = log.New(out, "INFO\t", log.Ldate|log.Ltime)
	case debugLevel:
		dLog = log.New(out, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	case errorLevel:
		eLog = log.New(out, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	default:
		return nil, errors.New("invalid logging level")
	}

	logger.iLog = iLog
	logger.dLog = dLog
	logger.eLog = eLog

	return logger, nil
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.eLog.Printf(format, v...)
	if l.fd != nil {
		// TODO need to implement log writing with buffer
		if err := l.fd.Sync(); err != nil {
			log.Printf("Error syncing file: %v", err)
			l.eLog.SetOutput(os.Stderr)
			l.eLog.Printf(format, v...)
			l.eLog.SetOutput(l.fd)
		}
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.iLog.Printf(format, v...)
	if l.fd != nil {
		if err := l.fd.Sync(); err != nil {
			log.Printf("Error syncing file: %v", err)
			l.iLog.SetOutput(os.Stdout)
			l.iLog.Printf(format, v...)
			l.iLog.SetOutput(l.fd)
		}
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.dLog.Printf(format, v...)
	if l.fd != nil {
		if err := l.fd.Sync(); err != nil {
			log.Printf("Error syncing file: %v", err)
			l.dLog.SetOutput(os.Stdout)
			l.dLog.Printf(format, v...)
			l.dLog.SetOutput(l.fd)
		}
	}
}

func (l *Logger) Close() error {
	if l.fd != nil {
		if err := l.fd.Close(); err != nil {
			return err
		}
	}
	return nil
}
