package tracing

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	logger *log.Logger
	tracer trace.Tracer
}

// NewTracer creates a tracer
func NewTracer(logger *log.Logger, tracer trace.Tracer) *Tracer {
	return &Tracer{
		logger: logger,
		tracer: tracer,
	}
}

// Start starts a new trace and logs the trace & span ID.
func (t *Tracer) Start(ctx context.Context, spanName, logInfo string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(
		ctx,
		spanName,
		opts...)

	t.LogTrace(span).Infof("%s", logInfo)
	return ctx, span
}

// LogTrace logs the trace and span ID of the provided span
func (t *Tracer) LogTrace(span trace.Span) *log.Entry {
	spanCtx := span.SpanContext()

	return t.logger.WithFields(log.Fields{
		"traceId": spanCtx.TraceID().String(),
		"spanId":  spanCtx.SpanID().String(),
	})
}

// FromContext Gets Tracing from context.
func FromContext(ctx context.Context) (*Tracer, error) {
	if ctx == nil {
		panic("nil context")
	}

	tracer, ok := ctx.Value("tracer").(*Tracer)
	if !ok {
		return nil, errors.New("unable to find tracing in the context")
	}

	return tracer, nil
}

// TracerToContext adds a tracer to a certain context
func TracerToContext(ctx context.Context, tracer *Tracer) context.Context {
	return context.WithValue(ctx, "tracer", tracer)
}
