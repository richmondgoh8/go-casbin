package utils

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type contextKey string

const (
	tracerCtxKey contextKey = "tracer"
)

type Tracer trace.Tracer

// WithTracer binds tracer to the context
func WithTracer(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, tracerCtxKey, tracer)
}

func ReturnTracer() trace.Tracer {
	return otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
}
