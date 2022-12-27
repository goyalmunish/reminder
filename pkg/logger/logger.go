package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

var (
	_log     *logrus.Logger = logrus.New()
	_options *Options
	// level   logrus.Level
	// _stdOut io.Writer
)

// to inspect the logger settings, you can expose this instance as following
// to the outside world, temporarily.
// var Instance = _log

/*
Key type is used for keys of context.Context.
*/
type Key string

// SetWithOptions setsup the logger based on the passed options.
func SetWithOptions(options *Options) {
	_options = options
	_log.SetLevel(logrus.Level(options.Level))
	_log.SetFormatter(&logrus.TextFormatter{})
}

// loggerWithContext enhances the log entry with context and with additional fields.
func loggerWithContext(ctx context.Context) *logrus.Entry {
	// add context (for hooks)
	logEntry := _log.WithContext(ctx)
	// add log fields (if they are available)
	// LookupFields can be nil while logger is being setup, make it blank
	// to mitigate log issues for such cases.
	if _options == nil || _options.LookupFields == nil {
		_options = DefaultOptions()
		_options.LookupFields = []string{}
	}
	for _, field := range _options.LookupFields {
		if v := ctx.Value(Key(field)); v != nil {
			logEntry = logEntry.WithFields(logrus.Fields{field: v})
		}
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
