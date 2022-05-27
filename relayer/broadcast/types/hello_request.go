package types

var _ PeerMessage = (*HelloRequest)(nil)

func (msg *HelloRequest) ValidateBasic() error {
	return nil
}
