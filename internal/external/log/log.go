package log

import (
	"io"
	"os"
	"playground/internal/entities"
	"time"

	"github.com/charmbracelet/log"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	log.SetPrefix("[external] ")
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}

type wrapLogger struct {
	*log.Logger
}

func (l *wrapLogger) SetTimeZone(loc *time.Location) {
	if loc == nil {
		panic("SetTimeZone: loc == nil")
	}
	l.Logger.SetTimeFunction(func(t time.Time) time.Time {
		return t.In(loc)
	})
}

func NewDefaultLogger() entities.Logger {
	lumberjack := &lumberjack.Logger{
		Filename:   "./log/app.log",
		MaxAge:     30, //days
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		Compress:   true, // disabled by default
	}
	l := log.New(io.MultiWriter(os.Stdout, lumberjack))
	l.SetTimeFormat("2006-01-02 15:04:05")
	l.SetPrefix("[playground] ")
	l.SetReportCaller(true)
	return &wrapLogger{Logger: l}
}

var (
	Debug = log.Debug
	Info  = log.Info
	Warn  = log.Warn
	Error = log.Error
	Fatal = log.Fatal

	Debugf = log.Debugf
	Infof  = log.Infof
	Warnf  = log.Warnf
	Errorf = log.Errorf
	Fatalf = log.Fatalf
)
