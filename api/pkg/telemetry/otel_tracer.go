package telemetry

import (
	"context"
	"math"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type otelTracer struct {
	logger Logger
}

// NewOtelLogger creates a new Tracer
func NewOtelLogger(logger Logger) Tracer {
	return &otelTracer{
		logger: logger,
	}
}

// Redact replaces the middle of a string with ****
func (tracer *otelTracer) Redact(secret string) string {
	visible := int(math.Ceil(float64(len(secret)) / 5))

	result := strings.Builder{}

	if len(secret) >= visible {
		result.WriteString(secret[:visible])
	}

	result.WriteString("*****")

	if len(secret) >= visible {
		result.WriteString(secret[len(secret)-visible:])
	}

	return result.String()
}

func (tracer *otelTracer) StartFromFiberCtxWithLogger(c *fiber.Ctx, logger Logger, name ...string) (context.Context, trace.Span, Logger) {
	ctx, span := tracer.StartFromFiberCtx(c, getName(name...))
	return ctx, span, tracer.CtxLogger(logger, span)
}

func (tracer *otelTracer) StartFromFiberCtx(c *fiber.Ctx, name ...string) (context.Context, trace.Span) {
	return tracer.Start(c.UserContext(), getName(name...))
}

func (tracer *otelTracer) CtxLogger(logger Logger, span trace.Span) Logger {
	return logger.WithSpan(span.SpanContext())
}

func (tracer *otelTracer) StartWithLogger(c context.Context, logger Logger, name ...string) (context.Context, trace.Span, Logger) {
	ctx, span := tracer.Start(c, getName(name...))
	return ctx, span, tracer.CtxLogger(logger, span)
}

func (tracer *otelTracer) Start(c context.Context, name ...string) (context.Context, trace.Span) {
	parentSpan := trace.SpanFromContext(c)
	ctx, span := parentSpan.TracerProvider().Tracer("").Start(c, getName(name...))

	span.SetAttributes(attribute.Key("traceID").String(parentSpan.SpanContext().TraceID().String()))
	span.SetAttributes(attribute.Key("spanID").String(span.SpanContext().SpanID().String()))
	span.SetAttributes(attribute.Key("traceFlags").String(parentSpan.SpanContext().TraceFlags().String()))

	return ctx, span
}

// Span returns the trace.Span from context.Context
func (tracer *otelTracer) Span(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func (tracer *otelTracer) WrapErrorSpan(span trace.Span, err error) error {
	if err == nil {
		return nil
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, strings.Split(err.Error(), "\n")[0])

	return err
}

func getName(name ...string) string {
	if len(name) > 0 {
		return name[0]
	}
	return functionName()
}

func functionName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return strings.ReplaceAll(frame.Function, "github.com/NdoleStudio/httpmock/", "")
}
