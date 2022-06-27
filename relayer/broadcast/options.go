package broadcast

// BroadcasterOption defines an option for the Broadcaster.
type BroadcasterOption func(*P2PBroadcaster) error

// WithHandler sets the message handler for the BroadcasterOption.
func WithHandler(handler BroadcastHandler) BroadcasterOption {
	return func(b *P2PBroadcaster) error {
		b.handler = handler
		return nil
	}
}

// WithHook sets the message broadcast hook for the BroadcasterOption.
func WithHook(hook BroadcasterHook) BroadcasterOption {
	return func(b *P2PBroadcaster) error {
		b.broadcasterHook = hook
		return nil
	}
}
