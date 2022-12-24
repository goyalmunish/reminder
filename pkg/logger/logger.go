package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

var (
	_log *logrus.Logger = logrus.New()
	// level   logrus.Level
	// _stdOut io.Writer
)

/*
Key type is used for keys of context.Context
*/
type Key string

func SetWithOptions(options *Options) {
	_log.SetLevel(logrus.Level(options.Level))
	_log.SetFormatter(&logrus.TextFormatter{})
}

func loggerWithContext(ctx context.Context) *logrus.Entry {
	logEntry := _log.WithFields(logrus.Fields{"app": "reminder"})
	if v := ctx.Value(Key("run_id")); v != nil {
		logEntry = logEntry.WithFields(logrus.Fields{"run_id": v})
	}
	return logEntry
}

func Trace(ctx context.Context, args ...interface{}) {
	le := loggerWithContext(ctx)
	le.Trace(args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	le := loggerWithContext(ctx)
	le.Debug(args...)
}

func Info(ctx context.Context, args ...interface{}) {
	le := loggerWithContext(ctx)
	le.Info(args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	le := loggerWithContext(ctx)
	le.Warn(args...)
}

func Error(ctx context.Context, args ...interface{}) {
	le := loggerWithContext(ctx)
	le.Error(args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	le := loggerWithContext(ctx)
	le.Fatal(args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	le := loggerWithContext(ctx)
	le.Panic(args...)
}
