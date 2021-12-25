package observability

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	TraceIdLabel  = "traceId"
	SpanIdLabel   = "spanId"
	DatabaseLabel = "database"
)

func TraceLabel(traceId string) attribute.KeyValue {
	return attribute.String(TraceIdLabel, traceId)
}

func SpanLabel(spanId string) attribute.KeyValue {
	return attribute.String(SpanIdLabel, spanId)
}

func DatabaseTypeLabel(dbType string) attribute.KeyValue {
	return attribute.String(DatabaseLabel, dbType)
}
