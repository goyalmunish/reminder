package logger

type Options struct {
	// PanicLevel 0
	// FatalLevel 1
	// ErrorLevel 2
	// WarnLevel 3
	// InfoLevel 4
	// DebugLevel 5
	// TraceLevel 6
	Level int8 `json:"level" yaml:"level" mapstructure:"level"`
	// Fields is a list of log fields that the logger
	// should become part of the log entry if available
	// in the context.
	LookupFields []string `json:"lookup_fields" yaml:"lookup_fields" mapstructure:"lookup_fields"`
}

func DefaultOptions() *Options {
	return &Options{
		// By default use debug level.
		Level: 5,
		// Value for the fields the logger should inject in the log
		// entry (if they are not nil)
		LookupFields: []string{"app", "run_id"},
	}
}
