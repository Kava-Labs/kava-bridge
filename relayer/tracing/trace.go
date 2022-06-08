package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

const (
	service = "kava-relayer"
)

func TracerProvider(url string, production bool) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	options := []tracesdk.TracerProviderOption{
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		)),
	}

	if production {
		// Always be sure to batch in production.
		options = append(options, tracesdk.WithBatcher(exp))
	} else {
		options = append(options, tracesdk.WithSyncer(exp))
	}

	tp := tracesdk.NewTracerProvider(options...)

	return tp, nil
}

func RegisterProvider(tp *tracesdk.TracerProvider) {
	tc := propagation.TraceContext{}
	// Register the TraceContext propagator globally.
	otel.SetTextMapPropagator(tc)
	otel.SetTracerProvider(tp)
}
