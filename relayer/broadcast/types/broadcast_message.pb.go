// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: relayer/broadcast/v1beta1/broadcast_message.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	github_com_libp2p_go_libp2p_core_peer "github.com/libp2p/go-libp2p-core/peer"
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
	RecipientPeerIDs []github_com_libp2p_go_libp2p_core_peer.ID `protobuf:"bytes,2,rep,name=recipient_peer_ids,json=recipientPeerIds,proto3,customtype=github.com/libp2p/go-libp2p-core/peer.ID" json:"recipient_peer_ids,omitempty"`
	// Customtype workaround for not having to use a separate protocgen.sh script
	Payload github_com_gogo_protobuf_types.Any `protobuf:"bytes,3,opt,name=payload,proto3,customtype=github.com/gogo/protobuf/types.Any" json:"payload"`
	// Timestamp when the message was broadcasted.
	Created time.Time `protobuf:"bytes,4,opt,name=created,proto3,stdtime" json:"created"`
	// Seconds after created time until the message expires. This requires
	// roughly synced times between peers
	TTLSeconds uint64 `protobuf:"varint,5,opt,name=ttl_seconds,json=ttlSeconds,proto3" json:"ttl_seconds,omitempty"`
}

func (m *BroadcastMessage) Reset()         { *m = BroadcastMessage{} }
func (m *BroadcastMessage) String() string { return proto.CompactTextString(m) }
func (*BroadcastMessage) ProtoMessage()    {}
func (*BroadcastMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cba85a09bf262990, []int{0}
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

func (m *BroadcastMessage) GetCreated() time.Time {
	if m != nil {
		return m.Created
	}
	return time.Time{}
}

func (m *BroadcastMessage) GetTTLSeconds() uint64 {
	if m != nil {
		return m.TTLSeconds
	}
	return 0
}

type HelloRequest struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cba85a09bf262990, []int{1}
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
	proto.RegisterFile("relayer/broadcast/v1beta1/broadcast_message.proto", fileDescriptor_cba85a09bf262990)
}

var fileDescriptor_cba85a09bf262990 = []byte{
	// 432 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x52, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x8d, 0xdd, 0xd2, 0xd0, 0x2d, 0x82, 0xc8, 0xaa, 0x90, 0xc9, 0xc1, 0x8e, 0x72, 0xb2, 0x90,
	0xe2, 0x55, 0xca, 0x85, 0x13, 0x02, 0x2b, 0x07, 0x22, 0x81, 0x84, 0x96, 0x9c, 0x7a, 0x89, 0xd6,
	0xde, 0xc1, 0xac, 0x58, 0x7b, 0xcd, 0xee, 0xa6, 0x92, 0xff, 0x45, 0x7f, 0x56, 0x8f, 0x3d, 0xa2,
	0x1e, 0x0c, 0x72, 0x0e, 0xfc, 0x0d, 0xe4, 0x2f, 0x11, 0xc8, 0xed, 0xcd, 0xcc, 0x9b, 0xf7, 0xf4,
	0x66, 0x17, 0x2d, 0x15, 0x08, 0x5a, 0x82, 0xc2, 0xb1, 0x92, 0x94, 0x25, 0x54, 0x1b, 0x7c, 0xb3,
	0x8c, 0xc1, 0xd0, 0xe5, 0xdf, 0xce, 0x36, 0x03, 0xad, 0x69, 0x0a, 0x61, 0xa1, 0xa4, 0x91, 0xce,
	0xb3, 0x7e, 0x25, 0xec, 0x89, 0xd3, 0xcb, 0x54, 0xa6, 0xb2, 0x9d, 0xe1, 0x06, 0x75, 0xb4, 0xe9,
	0x8b, 0x54, 0xca, 0x54, 0x00, 0x6e, 0xab, 0x78, 0xf7, 0x05, 0xd3, 0xbc, 0xec, 0x47, 0xfe, 0xff,
	0x23, 0xc3, 0x33, 0xd0, 0x86, 0x66, 0x45, 0x47, 0x98, 0xff, 0xb6, 0xd1, 0x24, 0x1a, 0xec, 0x3f,
	0x76, 0xee, 0xce, 0x73, 0x64, 0x73, 0xe6, 0x5a, 0x33, 0x2b, 0x38, 0x8f, 0xce, 0xea, 0xca, 0xb7,
	0xd7, 0x2b, 0x62, 0x73, 0xe6, 0xe4, 0xc8, 0x51, 0x90, 0xf0, 0x82, 0x43, 0x6e, 0xb6, 0x05, 0x80,
	0xda, 0x72, 0xa6, 0x5d, 0x7b, 0x76, 0x12, 0x9c, 0x47, 0x6f, 0x1f, 0x2a, 0x3f, 0x48, 0xb9, 0xf9,
	0xba, 0x8b, 0xc3, 0x44, 0x66, 0x58, 0xf0, 0xb8, 0xb8, 0x2a, 0x70, 0x2a, 0x17, 0x1d, 0x5a, 0x24,
	0x52, 0x01, 0x6e, 0x96, 0xc2, 0xf5, 0xaa, 0xae, 0xfc, 0x09, 0x19, 0x94, 0x3e, 0x01, 0xa8, 0xf5,
	0x4a, 0x93, 0x89, 0xfa, 0xa7, 0xc3, 0xb4, 0x73, 0x8d, 0xc6, 0x05, 0x2d, 0x85, 0xa4, 0xcc, 0x3d,
	0x99, 0x59, 0xc1, 0xc5, 0xd5, 0x65, 0xd8, 0xe5, 0x09, 0x87, 0x3c, 0xe1, 0xbb, 0xbc, 0x8c, 0x5e,
	0xde, 0x55, 0xfe, 0xe8, 0xa1, 0xf2, 0xe7, 0x07, 0xf6, 0xcd, 0x79, 0x0e, 0x52, 0x97, 0x05, 0xe8,
	0x86, 0x4b, 0x06, 0x41, 0xe7, 0x0d, 0x1a, 0x27, 0x0a, 0xa8, 0x01, 0xe6, 0x9e, 0xb6, 0xda, 0xd3,
	0x23, 0xed, 0xcd, 0x70, 0xab, 0xe8, 0x71, 0xe3, 0x70, 0xfb, 0xd3, 0xb7, 0xc8, 0xb0, 0xe4, 0x60,
	0x74, 0x61, 0x8c, 0xd8, 0x6a, 0x48, 0x64, 0xce, 0xb4, 0xfb, 0x68, 0x66, 0x05, 0xa7, 0xd1, 0xd3,
	0xba, 0xf2, 0xd1, 0x66, 0xf3, 0xe1, 0x73, 0xd7, 0x25, 0xc8, 0x18, 0xd1, 0xe3, 0x79, 0x80, 0x9e,
	0xbc, 0x07, 0x21, 0x24, 0x81, 0xef, 0x3b, 0xd0, 0xc6, 0x71, 0xd1, 0xb8, 0x7f, 0xed, 0xee, 0xd2,
	0x64, 0x28, 0x23, 0x72, 0x57, 0x7b, 0xd6, 0x7d, 0xed, 0x59, 0xbf, 0x6a, 0xcf, 0xba, 0xdd, 0x7b,
	0xa3, 0xfb, 0xbd, 0x37, 0xfa, 0xb1, 0xf7, 0x46, 0xd7, 0xaf, 0x0f, 0x12, 0x7e, 0xa3, 0x37, 0x74,
	0x21, 0x68, 0xac, 0x3b, 0x14, 0x2b, 0xce, 0x52, 0xc0, 0xc7, 0x9f, 0xac, 0xcd, 0x1e, 0x9f, 0xb5,
	0xa9, 0x5e, 0xfd, 0x09, 0x00, 0x00, 0xff, 0xff, 0x29, 0xbe, 0xa5, 0x01, 0x86, 0x02, 0x00, 0x00,
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
	if m.TTLSeconds != 0 {
		i = encodeVarintBroadcastMessage(dAtA, i, uint64(m.TTLSeconds))
		i--
		dAtA[i] = 0x28
	}
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
			{
				size := m.RecipientPeerIDs[iNdEx].Size()
				i -= size
				if _, err := m.RecipientPeerIDs[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintBroadcastMessage(dAtA, i, uint64(size))
			}
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
		for _, e := range m.RecipientPeerIDs {
			l = e.Size()
			n += 1 + l + sovBroadcastMessage(uint64(l))
		}
	}
	l = m.Payload.Size()
	n += 1 + l + sovBroadcastMessage(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.Created)
	n += 1 + l + sovBroadcastMessage(uint64(l))
	if m.TTLSeconds != 0 {
		n += 1 + sovBroadcastMessage(uint64(m.TTLSeconds))
	}
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
			var v github_com_libp2p_go_libp2p_core_peer.ID
			m.RecipientPeerIDs = append(m.RecipientPeerIDs, v)
			if err := m.RecipientPeerIDs[len(m.RecipientPeerIDs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TTLSeconds", wireType)
			}
			m.TTLSeconds = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBroadcastMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TTLSeconds |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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