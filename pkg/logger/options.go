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
}

func DefaultOptions() *Options {
	return &Options{
		// by default use debug level
		Level: 5,
	}
}
