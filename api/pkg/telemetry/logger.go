package telemetry

import "context"

// Logger is an interface for creating customer logger implementations
type Logger interface {
	// Error logs an error
	Error(err error)

	// WithCodeNamespace creates a new structured logger instance with a service name
	WithCodeNamespace(string) Logger

	// WithAttribute creates a new structured logger instance with a string
	WithAttribute(key string, value any) Logger

	// Trace logs a new message with trace level.
	Trace(value string)

	// Info logs a new message with information level.
	Info(value string)

	// Warn logs a new message with warning level.
	Warn(err error)

	// Debug logs a new message with debug level.
	Debug(value string)

	// Fatal logs a new message with fatal level.
	Fatal(err error)

	// Printf makes the logger compatible with retryablehttp.Logger
	Printf(string, ...any)

	// Write makes the logger compatible with io.Writer
	Write([]byte) (int, error)

	// WithContext adds a context.Context to a logger
	WithContext(ctx context.Context) Logger
}
