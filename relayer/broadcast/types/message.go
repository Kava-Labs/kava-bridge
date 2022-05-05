package types

// Message defines an interface for a message that can be sent to a peer.
type Message interface {
	Validate() error
}
