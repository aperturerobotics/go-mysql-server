//go:build !tinygo

package trace

import (
	"context"

	"github.com/dolthub/go-mysql-server/sql/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Tracer = oteltrace.Tracer
type Span = oteltrace.Span
type SpanStartOption = oteltrace.SpanStartOption
type EventOption = oteltrace.EventOption
type SpanEndOption = oteltrace.SpanEndOption
type SpanStartEventOption = oteltrace.SpanStartEventOption
type TracerProvider = oteltrace.TracerProvider

func NewNoopTracerProvider() TracerProvider {
	return oteltrace.NewNoopTracerProvider()
}

func WithAttributes(attrs ...attribute.KeyValue) SpanStartEventOption {
	return oteltrace.WithAttributes(attrs...)
}

func Start(tracer Tracer, ctx context.Context, name string, opts ...SpanStartOption) (context.Context, Span) {
	return tracer.Start(ctx, name, opts...)
}
