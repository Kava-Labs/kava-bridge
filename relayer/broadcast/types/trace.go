package types

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func NewTraceContext() *TraceContext {
	return &TraceContext{
		Carrier: make(propagation.MapCarrier),
	}
}

func (t *TraceContext) GetCarrier() propagation.MapCarrier {
	return t.Carrier
}

func (tctx *TraceContext) Inject(ctx context.Context) {
	carrier := propagation.MapCarrier(tctx.Carrier)
	otel.GetTextMapPropagator().Inject(ctx, &carrier)
}
