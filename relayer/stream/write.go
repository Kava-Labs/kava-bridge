package stream

import (
	"encoding/binary"
	"io"

	protoio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
)

const MAX_MESSAGE_SIZE = 1024 * 1024

func WriteProtoMessage(w io.Writer, msg proto.Message) error {
	writer := protoio.NewUint32DelimitedWriter(w, binary.BigEndian)

	return writer.WriteMsg(msg)
}

func ReadProtoMessage(r io.Reader, msg proto.Message) error {
	writer := protoio.NewUint32DelimitedReader(r, binary.BigEndian, MAX_MESSAGE_SIZE)

	return writer.ReadMsg(msg)
}
