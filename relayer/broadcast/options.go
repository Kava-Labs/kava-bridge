package broadcast

type BroadcasterOption func(*Broadcaster) error

func WithHandler(handler BroadcastHandler) BroadcasterOption {
	return func(b *Broadcaster) error {
		b.handler = handler
		return nil
	}
}
