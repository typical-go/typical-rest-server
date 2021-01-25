package log

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/logruskit"
)

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
