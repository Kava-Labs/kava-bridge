package broadcast

type Handler interface {
	HandleMessage(msg MessageWithPeerMetadata)
}

type HandlerFunc func(msg MessageWithPeerMetadata)

func (f HandlerFunc) HandleMessage(msg MessageWithPeerMetadata) {
	f(msg)
}
