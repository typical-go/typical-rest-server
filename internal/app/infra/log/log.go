package log

import (
	"log"

	"github.com/sirupsen/logrus"
)

// SetLogger set logger
func SetLogger(debug bool) *logrus.Logger {
	logger := logrus.StandardLogger()
	if debug {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{})
	} else {
		logger.SetLevel(logrus.WarnLevel)
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	log.SetOutput(logger.Writer())
	return logger
}

// Info ..
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warn ...
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warnf ...
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Error ...
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Fatal ...
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf ...
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
