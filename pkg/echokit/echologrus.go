package echokit

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

// Modification from https://github.com/plutov/echo-logrus/blob/master/middleware.go

type (
	// EchoLogrus : implement Logger
	EchoLogrus struct {
		*logrus.Logger
	}
)

var _ (echo.Logger) = (*EchoLogrus)(nil)

// WrapLogrus logrus logger in echo log interface
func WrapLogrus(logger *logrus.Logger) *EchoLogrus {
	return &EchoLogrus{logger}
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

	return log.INFO
}

// SetHeader is a stub to satisfy interface
// It's controlled by Logger
func (l *EchoLogrus) SetHeader(_ string) {}

// SetPrefix It's controlled by Logger
func (l *EchoLogrus) SetPrefix(s string) {}

// Prefix It's controlled by Logger
func (l *EchoLogrus) Prefix() string {
	return ""
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
		l.Panic("Invalid level")
	}
}

// Output logger output func
func (l *EchoLogrus) Output() io.Writer {
	return l.Out
}

// SetOutput change output, default os.Stdout
func (l *EchoLogrus) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
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

// Print string log
func (l *EchoLogrus) Print(i ...interface{}) {
	l.Logger.Print(i[0].(string))
}

// Debug string log
func (l *EchoLogrus) Debug(i ...interface{}) {
	l.Logger.Debug(i[0].(string))
}

// Info string log
func (l *EchoLogrus) Info(i ...interface{}) {
	l.Logger.Info(i[0].(string))
}

// Warn string log
func (l *EchoLogrus) Warn(i ...interface{}) {
	l.Logger.Warn(i[0].(string))
}

// Error string log
func (l *EchoLogrus) Error(i ...interface{}) {
	l.Logger.Error(i[0].(string))
}

// Fatal string log
func (l *EchoLogrus) Fatal(i ...interface{}) {
	l.Logger.Fatal(i[0].(string))
}

// Panic string log
func (l *EchoLogrus) Panic(i ...interface{}) {
	l.Logger.Panic(i[0].(string))
}
