package log

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/logruskit"
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
func Info(ctx context.Context, args ...interface{}) {
	logrus.WithFields(logruskit.GetFields(ctx)).Info(args...)
}

// Infof ...
func Infof(ctx context.Context, format string, args ...interface{}) {
	logrus.WithFields(logruskit.GetFields(ctx)).Infof(format, args...)
}

// Warn ...
func Warn(ctx context.Context, args ...interface{}) {
	logrus.WithFields(logruskit.GetFields(ctx)).Warn(args...)
}

// Warnf ...
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithFields(logruskit.GetFields(ctx)).Warnf(format, args...)
}
