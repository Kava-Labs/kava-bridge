// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bridge/v1beta1/event.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// EventBridgeEthereumToKava is emitted on Msg/BridgeEthereumToKava
type EventBridgeEthereumToKava struct {
	Relayer              string `protobuf:"bytes,1,opt,name=relayer,proto3" json:"relayer,omitempty"`
	EthereumErc20Address string `protobuf:"bytes,2,opt,name=ethereum_erc20_address,json=ethereumErc20Address,proto3" json:"ethereum_erc20_address,omitempty"`
	Receiver             string `protobuf:"bytes,3,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Amount               string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Sequence             string `protobuf:"bytes,5,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

func (m *EventBridgeEthereumToKava) Reset()         { *m = EventBridgeEthereumToKava{} }
func (m *EventBridgeEthereumToKava) String() string { return proto.CompactTextString(m) }
func (*EventBridgeEthereumToKava) ProtoMessage()    {}
func (*EventBridgeEthereumToKava) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4ba029c8787a2be, []int{0}
}
func (m *EventBridgeEthereumToKava) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventBridgeEthereumToKava) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventBridgeEthereumToKava.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventBridgeEthereumToKava) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventBridgeEthereumToKava.Merge(m, src)
}
func (m *EventBridgeEthereumToKava) XXX_Size() int {
	return m.Size()
}
func (m *EventBridgeEthereumToKava) XXX_DiscardUnknown() {
	xxx_messageInfo_EventBridgeEthereumToKava.DiscardUnknown(m)
}

var xxx_messageInfo_EventBridgeEthereumToKava proto.InternalMessageInfo

func (m *EventBridgeEthereumToKava) GetRelayer() string {
	if m != nil {
		return m.Relayer
	}
	return ""
}

func (m *EventBridgeEthereumToKava) GetEthereumErc20Address() string {
	if m != nil {
		return m.EthereumErc20Address
	}
	return ""
}

func (m *EventBridgeEthereumToKava) GetReceiver() string {
	if m != nil {
		return m.Receiver
	}
	return ""
}

func (m *EventBridgeEthereumToKava) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func (m *EventBridgeEthereumToKava) GetSequence() string {
	if m != nil {
		return m.Sequence
	}
	return ""
}

func init() {
	proto.RegisterType((*EventBridgeEthereumToKava)(nil), "bridge.v1beta1.EventBridgeEthereumToKava")
}

func init() { proto.RegisterFile("bridge/v1beta1/event.proto", fileDescriptor_e4ba029c8787a2be) }

var fileDescriptor_e4ba029c8787a2be = []byte{
	// 252 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x90, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x85, 0x6b, 0x7e, 0x0a, 0x78, 0x60, 0xb0, 0x50, 0x65, 0x3a, 0x58, 0x88, 0x89, 0x85, 0x98,
	0x02, 0x2f, 0x40, 0xa5, 0x0c, 0x88, 0x0d, 0x31, 0xb1, 0x54, 0x76, 0x72, 0xd5, 0x46, 0x34, 0x71,
	0xb9, 0x76, 0x2c, 0xfa, 0x16, 0xbc, 0x0f, 0x2f, 0xc0, 0xd8, 0x91, 0x11, 0x25, 0x2f, 0x82, 0xe2,
	0x38, 0x6c, 0xf7, 0xbb, 0xdf, 0x39, 0xcb, 0xa1, 0x53, 0x8d, 0x45, 0xbe, 0x04, 0xe9, 0x67, 0x1a,
	0x9c, 0x9a, 0x49, 0xf0, 0x50, 0xb9, 0x64, 0x83, 0xc6, 0x19, 0x76, 0xda, 0xbb, 0x24, 0xba, 0xcb,
	0x2f, 0x42, 0xcf, 0xd3, 0xce, 0xcf, 0xc3, 0x3f, 0x75, 0x2b, 0x40, 0xa8, 0xcb, 0x17, 0xf3, 0xa4,
	0xbc, 0x62, 0x9c, 0x1e, 0x21, 0xac, 0xd5, 0x16, 0x90, 0x93, 0x0b, 0x72, 0x75, 0xf2, 0x3c, 0x20,
	0xbb, 0xa7, 0x13, 0x88, 0xd9, 0x05, 0x60, 0x76, 0x7b, 0xb3, 0x50, 0x79, 0x8e, 0x60, 0x2d, 0xdf,
	0x0b, 0xc1, 0xb3, 0xc1, 0xa6, 0x9d, 0x7c, 0xe8, 0x1d, 0x9b, 0xd2, 0x63, 0x84, 0x0c, 0x0a, 0x0f,
	0xc8, 0xf7, 0x43, 0xee, 0x9f, 0xd9, 0x84, 0x8e, 0x55, 0x69, 0xea, 0xca, 0xf1, 0x83, 0x60, 0x22,
	0x75, 0x1d, 0x0b, 0xef, 0x35, 0x54, 0x19, 0xf0, 0xc3, 0xbe, 0x33, 0xf0, 0xfc, 0xf1, 0xbb, 0x11,
	0x64, 0xd7, 0x08, 0xf2, 0xdb, 0x08, 0xf2, 0xd9, 0x8a, 0xd1, 0xae, 0x15, 0xa3, 0x9f, 0x56, 0x8c,
	0x5e, 0xe5, 0xb2, 0x70, 0xab, 0x5a, 0x27, 0x99, 0x29, 0xe5, 0x9b, 0xf2, 0xea, 0x7a, 0xad, 0xb4,
	0xed, 0xaf, 0xb8, 0xce, 0x87, 0x8c, 0x87, 0xdb, 0x6e, 0xc0, 0xea, 0x71, 0xd8, 0xe7, 0xee, 0x2f,
	0x00, 0x00, 0xff, 0xff, 0xc2, 0xa6, 0xc9, 0x82, 0x3d, 0x01, 0x00, 0x00,
}

func (m *EventBridgeEthereumToKava) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventBridgeEthereumToKava) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventBridgeEthereumToKava) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Sequence) > 0 {
		i -= len(m.Sequence)
		copy(dAtA[i:], m.Sequence)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Sequence)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Amount) > 0 {
		i -= len(m.Amount)
		copy(dAtA[i:], m.Amount)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Amount)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Receiver) > 0 {
		i -= len(m.Receiver)
		copy(dAtA[i:], m.Receiver)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Receiver)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.EthereumErc20Address) > 0 {
		i -= len(m.EthereumErc20Address)
		copy(dAtA[i:], m.EthereumErc20Address)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.EthereumErc20Address)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Relayer) > 0 {
		i -= len(m.Relayer)
		copy(dAtA[i:], m.Relayer)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Relayer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventBridgeEthereumToKava) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Relayer)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.EthereumErc20Address)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Receiver)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Amount)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Sequence)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventBridgeEthereumToKava) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: EventBridgeEthereumToKava: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventBridgeEthereumToKava: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relayer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relayer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthereumErc20Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EthereumErc20Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receiver", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receiver = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sequence = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
