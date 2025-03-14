package telemetry

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is used for tracing
type Tracer interface {
	// StartFromFiberCtx creates a spanContext and a context.Context containing the newly-created spanContext.
	StartFromFiberCtx(c *fiber.Ctx, name ...string) (context.Context, trace.Span)

	// StartFromFiberCtxWithLogger creates a spanContext and a context.Context containing the newly-created spanContext with a logger
	StartFromFiberCtxWithLogger(c *fiber.Ctx, logger Logger, name ...string) (context.Context, trace.Span, Logger)

	// Start creates a spanContext and a context.Context containing the newly-created spanContext.
	Start(c context.Context, name ...string) (context.Context, trace.Span)

	// StartWithLogger creates a spanContext and a context.Context containing the newly-created spanContext with a logger
	StartWithLogger(c context.Context, logger Logger, name ...string) (context.Context, trace.Span, Logger)

	// WrapErrorSpan sets a spanContext as error
	WrapErrorSpan(span trace.Span, err error) error

	// Span returns the trace.Span from context.Context
	Span(ctx context.Context) trace.Span

	// Redact replaces the middle of a string with ****
	Redact(string) string
}
