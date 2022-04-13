package stream

import (
	"encoding/binary"
	"io"

	protoio "github.com/gogo/protobuf/io"
)

const MAX_MESSAGE_SIZE = 1024 * 1024

// NewProtoMessageWriter returns a new writer for writing proto messages.
func NewProtoMessageWriter(w io.Writer) protoio.WriteCloser {
	return protoio.NewUint32DelimitedWriter(w, binary.BigEndian)
}

// NewProtoMessageReader returns a new reader for reading proto messages.
func NewProtoMessageReader(r io.Reader) protoio.ReadCloser {
	return protoio.
		NewUint32DelimitedReader(r, binary.BigEndian, MAX_MESSAGE_SIZE)
}
