package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Tracer wraps OpenTelemetry tracer
type Tracer struct {
	tracer trace.Tracer
}

// NewTracer creates a new tracer instance
func NewTracer(serviceName string) *Tracer {
	tracer := otel.Tracer(serviceName)
	return &Tracer{
		tracer: tracer,
	}
}

// StartSpan starts a new span
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// TraceFunction wraps a function with tracing
func (t *Tracer) TraceFunction(ctx context.Context, name string, fn func(context.Context) error) error {
	ctx, span := t.StartSpan(ctx, name)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// TraceFunction traces the execution of a function that returns an error.
func TraceFunction(tracer *Tracer, ctx context.Context, name string, fn func(context.Context) error) error {
	ctx, span := tracer.StartSpan(ctx, name)
	defer span.End()
	return fn(ctx)
}

// TraceFunctionWithResult wraps a function with tracing and returns a result
func TraceFunctionWithResult[T any](tracer *Tracer, ctx context.Context, name string, fn func(context.Context) (T, error)) (T, error) {
	ctx, span := tracer.StartSpan(ctx, name)
	defer span.End()

	result, err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}

// AddSpanEvent adds an event to the current span
func (t *Tracer) AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetSpanAttributes sets attributes on the current span
func (t *Tracer) SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}
