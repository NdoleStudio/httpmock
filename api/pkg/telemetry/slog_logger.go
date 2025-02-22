package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type slogLogger struct {
	ctx            context.Context
	slog           *slog.Logger
	name           string
	skipFrameCount int
	attributes     []any
}

// NewSlogLogger creates a new instance of the slogLogger
func NewSlogLogger(ctx context.Context, skipFrameCount int, handler slog.Handler, attributes []any) Logger {
	logger := &slogLogger{
		skipFrameCount: skipFrameCount,
		attributes:     attributes,
		ctx:            ctx,
		slog:           slog.New(handler),
	}

	return logger
}

// WithCodeNamespace creates a new structured logger instance with a service name
func (logger *slogLogger) WithCodeNamespace(codeNamespace string) Logger {
	logger.slog.Handler()
	return NewSlogLogger(
		logger.ctx,
		logger.skipFrameCount,
		logger.slog.Handler(),
		logger.addAttribute(string(semconv.CodeNamespaceKey), codeNamespace),
	)
}

func (logger *slogLogger) Write(bytes []byte) (int, error) {
	logger.slog.InfoContext(logger.ctx, strings.TrimSpace(string(bytes)), logger.attributes...)
	return len(bytes), nil
}

func (logger *slogLogger) Printf(s string, i ...any) {
	logger.slog.InfoContext(logger.ctx, fmt.Sprintf(s, i...), logger.attributes...)
}

// WithAttribute creates a new slogLogger instance with a key value pair
func (logger *slogLogger) WithAttribute(key string, value any) Logger {
	return NewSlogLogger(logger.ctx, logger.skipFrameCount, logger.slog.Handler(), logger.addAttribute(key, value))
}

// Info logs a new message with information level.
func (logger *slogLogger) Info(value string) {
	logger.slog.InfoContext(logger.ctx, value, logger.attributesWithCaller()...)
}

// Trace logs a new message with trace level.
func (logger *slogLogger) Trace(value string) {
	logger.slog.DebugContext(logger.ctx, value, logger.attributesWithCaller()...)
}

// Warn logs a new message with warning level.
func (logger *slogLogger) Warn(err error) {
	logger.slog.WarnContext(logger.ctx, err.Error(), logger.attributesWithCaller()...)
}

// Debug logs a new message with debug level.
func (logger *slogLogger) Debug(value string) {
	logger.slog.DebugContext(logger.ctx, value, logger.attributesWithCaller()...)
}

// Fatal logs a new message with fatal level.
func (logger *slogLogger) Fatal(err error) {
	logger.slog.ErrorContext(logger.ctx, err.Error(), logger.attributesWithCaller()...)
	os.Exit(1)
}

// Error logs an error
func (logger *slogLogger) Error(err error) {
	logger.slog.WarnContext(logger.ctx, err.Error(), logger.attributesWithCaller()...)
}

// WithContext adds a context.Context to a logger
func (logger *slogLogger) WithContext(ctx context.Context) Logger {
	return NewSlogLogger(ctx, logger.skipFrameCount, logger.slog.Handler(), logger.cloneAttributes())
}

func (logger *slogLogger) cloneAttributes() []any {
	var attributes []any
	for _, attribute := range logger.attributes {
		attributes = append(attributes, attribute)
	}
	return attributes
}

func (logger *slogLogger) attributesWithCaller() []any {
	_, filename, line, _ := runtime.Caller(logger.skipFrameCount)
	return logger.addAttribute(string(semconv.CodeFilepathKey), fmt.Sprintf("%s:%d", filename, line))
}

func (logger *slogLogger) addAttribute(key string, value any) []any {
	return append(logger.cloneAttributes(), key, value)
}
