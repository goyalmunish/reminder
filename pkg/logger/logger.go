package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

var (
	_log          *logrus.Logger = logrus.New()
	_options      *Options
	_globalFields logrus.Fields = logrus.Fields{}
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

// SetWithOptions setups the logger based on the passed options.
func SetWithOptions(options *Options) {
	_options = options
	_log.SetLevel(logrus.Level(options.Level))
	_log.SetFormatter(&logrus.TextFormatter{})
}

// SetGlobalFields setups the global fields.
func SetGlobalFields(fields map[string]interface{}) {
	for key, value := range fields {
		_globalFields[key] = value
	}
}

// entryWithGlobalFields sets global value to the logger
// The global values are values in the scope of whole run of the app.
// These are particular relevent in desktop apps.
func entryWithGlobalFields() *logrus.Entry {
	return _log.WithFields(_globalFields)
}

// addContext enhances the log entry with context and with additional fields.
func addContext(ctx context.Context, logEntry *logrus.Entry) *logrus.Entry {
	// Add context (for hooks)
	logEntry = logEntry.WithContext(ctx)
	// Add log fields (if they are available)
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

// Trace is logrus.Trace with Global settings.
func Trace(args ...interface{}) {
	le := entryWithGlobalFields()
	le.Trace(args...)
}

// Debug is logrus.Debug with Global settings.
func Debug(args ...interface{}) {
	le := entryWithGlobalFields()
	le.Debug(args...)
}

// Info is logrus.Info with Global settings.
func Info(args ...interface{}) {
	le := entryWithGlobalFields()
	le.Info(args...)
}

// Warn is logrus.Warn with Global settings.
func Warn(args ...interface{}) {
	le := entryWithGlobalFields()
	le.Warn(args...)
}

// Error is logrus.Error with Global settings.
func Error(args ...interface{}) {
	le := entryWithGlobalFields()
	le.Error(args...)
}

// Fatal is logrus.Fatal with Global settings.
func Fatal(args ...interface{}) {
	le := entryWithGlobalFields()
	le.Fatal(args...)
}

// Panic is logrus.Panic with Global settings.
func Panic(args ...interface{}) {
	le := entryWithGlobalFields()
	le.Panic(args...)
}

// TraceC is Trace with context.
func TraceC(ctx context.Context, args ...interface{}) {
	le := entryWithGlobalFields()
	le = addContext(ctx, le)
	le.Trace(args...)
}

// DebugC is Debug with context.
func DebugC(ctx context.Context, args ...interface{}) {
	le := entryWithGlobalFields()
	le = addContext(ctx, le)
	le.Debug(args...)
}

// InfoC is Info with context.
func InfoC(ctx context.Context, args ...interface{}) {
	le := entryWithGlobalFields()
	le = addContext(ctx, le)
	le.Info(args...)
}

// WarnC is Warn with context.
func WarnC(ctx context.Context, args ...interface{}) {
	le := entryWithGlobalFields()
	le = addContext(ctx, le)
	le.Warn(args...)
}

// ErrorC is Error with context.
func ErrorC(ctx context.Context, args ...interface{}) {
	le := entryWithGlobalFields()
	le = addContext(ctx, le)
	le.Error(args...)
}

// FatalC is Fatal with context.
func FatalC(ctx context.Context, args ...interface{}) {
	le := entryWithGlobalFields()
	le = addContext(ctx, le)
	le.Fatal(args...)
}

// PanicC is Panic with context.
func PanicC(ctx context.Context, args ...interface{}) {
	le := entryWithGlobalFields()
	le = addContext(ctx, le)
	le.Panic(args...)
}
