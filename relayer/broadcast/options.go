package broadcast

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
func WithHook(hook broadcasterHook) BroadcasterOption {
	return func(b *Broadcaster) error {
		b.broadcasterHook = hook
		return nil
	}
}
