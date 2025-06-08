package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer     trace.TracerProvider
	tracerName string
}

func NewTracer(serviceName string) *Tracer {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)

	return &Tracer{
		tracer:     tp,
		tracerName: serviceName,
	}
}

func (t *Tracer) StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tracer := t.tracer.Tracer(t.tracerName)
	return tracer.Start(ctx, spanName)
}

func (t *Tracer) EndSpan(span trace.Span) {
	span.End()
}

func (t *Tracer) Shutdown(ctx context.Context) error {
	if tp, ok := t.tracer.(*sdktrace.TracerProvider); ok {
		return tp.Shutdown(ctx)
	}
	return nil
}
