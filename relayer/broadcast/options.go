package broadcast

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// BroadcasterOption defines an option for the Broadcaster.
type BroadcasterOption func(*Broadcaster) error

// WithHandler sets the message handler for the BroadcasterOption.
func WithHandler(handler BroadcastHandler) BroadcasterOption {
	return func(b *Broadcaster) error {
		b.handler = handler
		return nil
	}
}

// WithHook sets the message broadcast hook for the BroadcasterOption.
func WithHook(hook BroadcasterHook) BroadcasterOption {
	return func(b *Broadcaster) error {
		b.broadcasterHook = hook
		return nil
	}
}

// WithTracer enables tracing for the given Jaeger provider.
func WithTracer(jaegerUrl string) BroadcasterOption {
	return func(b *Broadcaster) error {
		tp, err := tracerProvider("http://localhost:14268/api/traces", b.host.ID())
		if err != nil {
			return err
		}

		log.Infow("Enabled tracing", "url", jaegerUrl)

		tc := propagation.TraceContext{}
		// Register the TraceContext propagator globally.
		otel.SetTextMapPropagator(tc)
		otel.SetTracerProvider(tp)

		return nil
	}
}
