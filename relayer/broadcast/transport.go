package broadcast

import "github.com/gogo/protobuf/proto"

type Transporter interface {
	Send(pb proto.Message) error
	Receive(pb proto.Message) error
}
