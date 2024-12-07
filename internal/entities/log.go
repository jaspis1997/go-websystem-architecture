package entities

import (
	"io"
	"time"

	"github.com/charmbracelet/log"
)

type LogLevel = log.Level

const (
	DebugLevel = log.DebugLevel
	InfoLevel  = log.InfoLevel
	WarnLevel  = log.WarnLevel
	ErrorLevel = log.ErrorLevel
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(msg interface{}, keyvals ...interface{})
	Info(msg interface{}, keyvals ...interface{})
	Warn(msg interface{}, keyvals ...interface{})
	Error(msg interface{}, keyvals ...interface{})
	Fatal(msg interface{}, keyvals ...interface{})

	SetPrefix(prefix string)
	SetLevel(level LogLevel)
	SetOutput(w io.Writer)
	SetTimeZone(loc *time.Location)
}
