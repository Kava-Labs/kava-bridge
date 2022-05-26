package types

import "github.com/libp2p/go-libp2p-core/peer"

var _ PeerMessage = (*HelloRequest)(nil)

func (msg *HelloRequest) ValidateBasic() error {
	return nil
}

func (msg *HelloRequest) GetSenderPeerID() peer.ID {
	return msg.PeerID
}
