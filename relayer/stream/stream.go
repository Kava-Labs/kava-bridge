package stream

import (
	"encoding/binary"
	"io"

	protoio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
)

const MAX_MESSAGE_SIZE = 1024 * 1024

// WriteProtoMessage writes a proto message to a stream.
func WriteProtoMessage(w io.Writer, msg proto.Message) error {
	return protoio.NewUint32DelimitedWriter(w, binary.BigEndian).WriteMsg(msg)
}

// ReadProtoMessage reads a proto message from a stream with a max size of 1MB.
func ReadProtoMessage(r io.Reader, msg proto.Message) error {
	return protoio.
		NewUint32DelimitedReader(r, binary.BigEndian, MAX_MESSAGE_SIZE).
		ReadMsg(msg)
}
