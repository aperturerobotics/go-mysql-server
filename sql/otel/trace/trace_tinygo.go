//go:build tinygo || js

package trace

import (
	"context"

	"github.com/dolthub/go-mysql-server/sql/otel/attribute"
)

type SpanStartOption interface{}

type EventOption interface{}

type SpanEndOption interface{}

type SpanStartEventOption struct{}

type Span interface {
	IsRecording() bool
	AddEvent(string, ...EventOption)
	End(...SpanEndOption)
	RecordError(error, ...EventOption)
	SetAttributes(...attribute.KeyValue)
}

type noopSpan struct{}

func (noopSpan) IsRecording() bool {
	return false
}

func (noopSpan) AddEvent(_ string, _ ...EventOption) {}

func (noopSpan) End(_ ...SpanEndOption) {}

func (noopSpan) RecordError(_ error, _ ...EventOption) {}

func (noopSpan) SetAttributes(_ ...attribute.KeyValue) {}

type Tracer interface {
	Start(context.Context, string, ...SpanStartOption) (context.Context, Span)
}

type noopTracer struct{}

func (noopTracer) Start(ctx context.Context, _ string, _ ...SpanStartOption) (context.Context, Span) {
	return ctx, noopSpan{}
}

type TracerProvider struct{}

func (TracerProvider) Tracer(_ string, _ ...TracerOption) Tracer {
	return noopTracer{}
}

type TracerOption struct{}

func NewNoopTracerProvider() TracerProvider {
	return TracerProvider{}
}

func WithAttributes(_ ...attribute.KeyValue) SpanStartEventOption {
	return SpanStartEventOption{}
}
