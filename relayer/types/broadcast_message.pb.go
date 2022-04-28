// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: relayer/v1beta1/broadcast_message.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// BroadcastMessage is used between peers to wrap messages for each protocol
type BroadcastMessage struct {
	// Unique ID of this message.
	ID string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Selected recipients of the message, to partially restrict the broadcast to
	// a subset a peers.
	RecipientPeerIDs []string `protobuf:"bytes,2,rep,name=recipient_peer_ids,json=recipientPeerIds,proto3" json:"recipient_peer_ids,omitempty"`
	// Customtype workaround for not having to use a separate protocgen.sh script
	Payload github_com_gogo_protobuf_types.Any `protobuf:"bytes,3,opt,name=payload,proto3,customtype=github.com/gogo/protobuf/types.Any" json:"payload"`
	// Timestamp when the message was broadcasted.
	Created time.Time `protobuf:"bytes,4,opt,name=created,proto3,stdtime" json:"created"`
}

func (m *BroadcastMessage) Reset()         { *m = BroadcastMessage{} }
func (m *BroadcastMessage) String() string { return proto.CompactTextString(m) }
func (*BroadcastMessage) ProtoMessage()    {}
func (*BroadcastMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_8266a42107f80093, []int{0}
}
func (m *BroadcastMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BroadcastMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BroadcastMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BroadcastMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BroadcastMessage.Merge(m, src)
}
func (m *BroadcastMessage) XXX_Size() int {
	return m.Size()
}
func (m *BroadcastMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_BroadcastMessage.DiscardUnknown(m)
}

var xxx_messageInfo_BroadcastMessage proto.InternalMessageInfo

func (m *BroadcastMessage) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *BroadcastMessage) GetRecipientPeerIDs() []string {
	if m != nil {
		return m.RecipientPeerIDs
	}
	return nil
}

func (m *BroadcastMessage) GetCreated() time.Time {
	if m != nil {
		return m.Created
	}
	return time.Time{}
}

type HelloRequest struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8266a42107f80093, []int{1}
}
func (m *HelloRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HelloRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HelloRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloRequest.Merge(m, src)
}
func (m *HelloRequest) XXX_Size() int {
	return m.Size()
}
func (m *HelloRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HelloRequest proto.InternalMessageInfo

func (m *HelloRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*BroadcastMessage)(nil), "relayer.v1beta1.BroadcastMessage")
	proto.RegisterType((*HelloRequest)(nil), "relayer.v1beta1.HelloRequest")
}

func init() {
	proto.RegisterFile("relayer/v1beta1/broadcast_message.proto", fileDescriptor_8266a42107f80093)
}

var fileDescriptor_8266a42107f80093 = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x51, 0x4f, 0x6b, 0xa3, 0x40,
	0x14, 0x57, 0xb3, 0x24, 0x9b, 0xd9, 0x85, 0x0d, 0x12, 0x16, 0x37, 0x07, 0x0d, 0xb9, 0xac, 0x2c,
	0xac, 0x92, 0xdd, 0x7b, 0xa1, 0x92, 0x43, 0x72, 0x28, 0x14, 0xe9, 0x29, 0x97, 0x30, 0x3a, 0xaf,
	0x76, 0xa8, 0x66, 0xec, 0xcc, 0x24, 0xe0, 0x57, 0xe8, 0x29, 0x1f, 0x2b, 0xc7, 0x1c, 0x4b, 0x0f,
	0xb6, 0x98, 0x2f, 0x52, 0xa2, 0x0e, 0xa4, 0xed, 0xed, 0x3d, 0x7f, 0x3f, 0x7f, 0x7f, 0xe6, 0xa1,
	0xdf, 0x1c, 0x52, 0x5c, 0x00, 0xf7, 0xb7, 0xd3, 0x08, 0x24, 0x9e, 0xfa, 0x11, 0x67, 0x98, 0xc4,
	0x58, 0xc8, 0x55, 0x06, 0x42, 0xe0, 0x04, 0xbc, 0x9c, 0x33, 0xc9, 0xcc, 0x1f, 0x2d, 0xd1, 0x6b,
	0x89, 0xa3, 0x61, 0xc2, 0x12, 0x56, 0x63, 0xfe, 0x69, 0x6a, 0x68, 0xa3, 0x5f, 0x09, 0x63, 0x49,
	0x0a, 0x7e, 0xbd, 0x45, 0x9b, 0x5b, 0x1f, 0xaf, 0x8b, 0x16, 0x72, 0x3e, 0x42, 0x92, 0x66, 0x20,
	0x24, 0xce, 0xf2, 0x86, 0x30, 0x79, 0x34, 0xd0, 0x20, 0x50, 0xf6, 0x57, 0x8d, 0xbb, 0xf9, 0x13,
	0x19, 0x94, 0x58, 0xfa, 0x58, 0x77, 0xfb, 0x41, 0xb7, 0x2a, 0x1d, 0x63, 0x31, 0x0b, 0x0d, 0x4a,
	0xcc, 0x00, 0x99, 0x1c, 0x62, 0x9a, 0x53, 0x58, 0xcb, 0x55, 0x0e, 0xc0, 0x57, 0x94, 0x08, 0xcb,
	0x18, 0x77, 0xdc, 0x7e, 0x30, 0xac, 0x4a, 0x67, 0x10, 0x2a, 0xf4, 0x1a, 0x80, 0x2f, 0x66, 0x22,
	0x1c, 0xf0, 0x77, 0x5f, 0x88, 0x30, 0x97, 0xa8, 0x97, 0xe3, 0x22, 0x65, 0x98, 0x58, 0x9d, 0xb1,
	0xee, 0x7e, 0xfb, 0x37, 0xf4, 0x9a, 0x8c, 0x9e, 0xca, 0xe8, 0x5d, 0xae, 0x8b, 0xe0, 0xcf, 0xbe,
	0x74, 0xb4, 0xe7, 0xd2, 0x99, 0x24, 0x54, 0xde, 0x6d, 0x22, 0x2f, 0x66, 0x59, 0x5d, 0xf9, 0xac,
	0x49, 0x91, 0x83, 0x38, 0x71, 0x43, 0x25, 0x68, 0x5e, 0xa0, 0x5e, 0xcc, 0x01, 0x4b, 0x20, 0xd6,
	0x97, 0x5a, 0x7b, 0xf4, 0x49, 0xfb, 0x46, 0xf5, 0x0f, 0xbe, 0x9e, 0x1c, 0x76, 0x2f, 0x8e, 0x1e,
	0xaa, 0x9f, 0x26, 0x2e, 0xfa, 0x3e, 0x87, 0x34, 0x65, 0x21, 0x3c, 0x6c, 0x40, 0x48, 0xd3, 0x42,
	0xbd, 0xf6, 0x20, 0xcd, 0x63, 0x84, 0x6a, 0x0d, 0xe6, 0xfb, 0xca, 0xd6, 0x0f, 0x95, 0xad, 0xbf,
	0x56, 0xb6, 0xbe, 0x3b, 0xda, 0xda, 0xe1, 0x68, 0x6b, 0x4f, 0x47, 0x5b, 0x5b, 0x7a, 0x67, 0x81,
	0xef, 0xf1, 0x16, 0xff, 0x4d, 0x71, 0x24, 0x9a, 0x29, 0xe2, 0x94, 0x24, 0xe0, 0xab, 0xeb, 0xd7,
	0x05, 0xa2, 0x6e, 0x1d, 0xed, 0xff, 0x5b, 0x00, 0x00, 0x00, 0xff, 0xff, 0x69, 0x42, 0x4e, 0xa8,
	0x15, 0x02, 0x00, 0x00,
}

func (m *BroadcastMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BroadcastMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BroadcastMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.Created, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.Created):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintBroadcastMessage(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	{
		size := m.Payload.Size()
		i -= size
		if _, err := m.Payload.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBroadcastMessage(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.RecipientPeerIDs) > 0 {
		for iNdEx := len(m.RecipientPeerIDs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.RecipientPeerIDs[iNdEx])
			copy(dAtA[i:], m.RecipientPeerIDs[iNdEx])
			i = encodeVarintBroadcastMessage(dAtA, i, uint64(len(m.RecipientPeerIDs[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintBroadcastMessage(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *HelloRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HelloRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HelloRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintBroadcastMessage(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBroadcastMessage(dAtA []byte, offset int, v uint64) int {
	offset -= sovBroadcastMessage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BroadcastMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovBroadcastMessage(uint64(l))
	}
	if len(m.RecipientPeerIDs) > 0 {
		for _, s := range m.RecipientPeerIDs {
			l = len(s)
			n += 1 + l + sovBroadcastMessage(uint64(l))
		}
	}
	l = m.Payload.Size()
	n += 1 + l + sovBroadcastMessage(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.Created)
	n += 1 + l + sovBroadcastMessage(uint64(l))
	return n
}

func (m *HelloRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovBroadcastMessage(uint64(l))
	}
	return n
}

func sovBroadcastMessage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBroadcastMessage(x uint64) (n int) {
	return sovBroadcastMessage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BroadcastMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBroadcastMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BroadcastMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BroadcastMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBroadcastMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RecipientPeerIDs", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBroadcastMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RecipientPeerIDs = append(m.RecipientPeerIDs, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBroadcastMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Payload.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Created", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBroadcastMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.Created, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBroadcastMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *HelloRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBroadcastMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HelloRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HelloRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBroadcastMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBroadcastMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBroadcastMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBroadcastMessage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBroadcastMessage
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBroadcastMessage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBroadcastMessage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthBroadcastMessage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBroadcastMessage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBroadcastMessage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBroadcastMessage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBroadcastMessage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBroadcastMessage = fmt.Errorf("proto: unexpected end of group")
)
