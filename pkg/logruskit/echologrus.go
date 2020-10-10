package logruskit

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type (
	// EchoLogrus is implementation of echo Logger
	EchoLogrus struct {
		*logrus.Logger
		prefix string
	}
)

var _ (echo.Logger) = (*EchoLogrus)(nil)

// EchoLogger logrus logger in echo log interface
func EchoLogger(logger *logrus.Logger) *EchoLogrus {
	return &EchoLogrus{
		Logger: logger,
	}
}

// SetHeader to set header (NOT SUPPORTED)
func (l *EchoLogrus) SetHeader(string) {}

// SetPrefix to set prefix
func (l *EchoLogrus) SetPrefix(prefix string) {
	l.prefix = prefix
}

// Prefix of echo logrus
func (l *EchoLogrus) Prefix() string {
	return l.prefix
}

// SetLevel set level to logger from given log.Lvl
func (l *EchoLogrus) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		l.Logger.SetLevel(logrus.DebugLevel)
	case log.WARN:
		l.Logger.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		l.Logger.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		l.Logger.SetLevel(logrus.InfoLevel)
	default:
		logrus.Warnf("Unknown level: %v", lvl)
		l.Logger.SetLevel(logrus.WarnLevel)
	}
}

// Level returns logger level
func (l *EchoLogrus) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	}
	return log.WARN
}

// Output logger output func
func (l *EchoLogrus) Output() io.Writer {
	return l.Out
}

// Printj print json log
func (l *EchoLogrus) Printj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Print()
}

// Debugj debug json log
func (l *EchoLogrus) Debugj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Debug()
}

// Infoj info json log
func (l *EchoLogrus) Infoj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Info()
}

// Warnj warning json log
func (l *EchoLogrus) Warnj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Warn()
}

// Errorj error json log
func (l *EchoLogrus) Errorj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Error()
}

// Fatalj fatal json log
func (l *EchoLogrus) Fatalj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Fatal()
}

// Panicj panic json log
func (l *EchoLogrus) Panicj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Panic()
}
