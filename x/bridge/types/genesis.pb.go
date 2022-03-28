// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bridge/v1beta1/genesis.proto

package types

import (
	bytes "bytes"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
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

// GenesisState defines the bridge module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// erc20_bridge_pairs defines all of the bridged erc20 tokens.
	ERC20BridgePairs ERC20BridgePairs `protobuf:"bytes,2,rep,name=erc20_bridge_pairs,json=erc20BridgePairs,proto3,castrepeated=ERC20BridgePairs" json:"erc20_bridge_pairs"`
	// next_withdraw_sequence defines the unique incrementing sequence per withdraw tx.
	NextWithdrawSequence github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=next_withdraw_sequence,json=nextWithdrawSequence,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"next_withdraw_sequence"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a343f3772a97af9, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetERC20BridgePairs() ERC20BridgePairs {
	if m != nil {
		return m.ERC20BridgePairs
	}
	return nil
}

// Params defines the bridge module params
type Params struct {
	// List of ERC20Tokens that are allowed to be bridged to Kava
	EnabledERC20Tokens EnabledERC20Tokens `protobuf:"bytes,1,rep,name=enabled_erc20_tokens,json=enabledErc20Tokens,proto3,castrepeated=EnabledERC20Tokens" json:"enabled_erc20_tokens"`
	// Permissioned relayer address that is allowed to submit bridge messages
	Relayer github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,2,opt,name=relayer,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"relayer,omitempty"`
	// enabled_conversion_pairs defines the list of conversion pairs allowed to be
	// converted between Kava ERC20 and sdk.Coin
	EnabledConversionPairs ConversionPairs `protobuf:"bytes,3,rep,name=enabled_conversion_pairs,json=enabledConversionPairs,proto3,castrepeated=ConversionPairs" json:"enabled_conversion_pairs"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a343f3772a97af9, []int{1}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetEnabledERC20Tokens() EnabledERC20Tokens {
	if m != nil {
		return m.EnabledERC20Tokens
	}
	return nil
}

func (m *Params) GetRelayer() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Relayer
	}
	return nil
}

func (m *Params) GetEnabledConversionPairs() ConversionPairs {
	if m != nil {
		return m.EnabledConversionPairs
	}
	return nil
}

// EnabledERC20Token defines an external ERC20 that is allowed to be bridged to Kava
type EnabledERC20Token struct {
	// Address of the contract on Ethereum
	Address HexBytes `protobuf:"bytes,1,opt,name=address,proto3,casttype=HexBytes" json:"address,omitempty"`
	// Name of the token.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Symbol of the ERC20 token, usually a shorter version of the name.
	Symbol string `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// Number of decimals the ERC20 uses to get its user representation. The max
	// value is an unsigned 8 bit integer, but is an uint32 as the smallest
	// protobuf integer type.
	Decimals uint32 `protobuf:"varint,4,opt,name=decimals,proto3" json:"decimals,omitempty"`
}

func (m *EnabledERC20Token) Reset()         { *m = EnabledERC20Token{} }
func (m *EnabledERC20Token) String() string { return proto.CompactTextString(m) }
func (*EnabledERC20Token) ProtoMessage()    {}
func (*EnabledERC20Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a343f3772a97af9, []int{2}
}
func (m *EnabledERC20Token) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EnabledERC20Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EnabledERC20Token.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EnabledERC20Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnabledERC20Token.Merge(m, src)
}
func (m *EnabledERC20Token) XXX_Size() int {
	return m.Size()
}
func (m *EnabledERC20Token) XXX_DiscardUnknown() {
	xxx_messageInfo_EnabledERC20Token.DiscardUnknown(m)
}

var xxx_messageInfo_EnabledERC20Token proto.InternalMessageInfo

func (m *EnabledERC20Token) GetAddress() HexBytes {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *EnabledERC20Token) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EnabledERC20Token) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *EnabledERC20Token) GetDecimals() uint32 {
	if m != nil {
		return m.Decimals
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "bridge.v1beta1.GenesisState")
	proto.RegisterType((*Params)(nil), "bridge.v1beta1.Params")
	proto.RegisterType((*EnabledERC20Token)(nil), "bridge.v1beta1.EnabledERC20Token")
}

func init() { proto.RegisterFile("bridge/v1beta1/genesis.proto", fileDescriptor_6a343f3772a97af9) }

var fileDescriptor_6a343f3772a97af9 = []byte{
	// 580 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x31, 0x6f, 0xd3, 0x40,
	0x18, 0xcd, 0xa5, 0x55, 0xda, 0x5e, 0x03, 0x94, 0x23, 0x0a, 0x26, 0x42, 0x76, 0x88, 0x50, 0x95,
	0x81, 0xd8, 0x69, 0x60, 0x64, 0xa9, 0xab, 0x8a, 0x76, 0x2b, 0x2e, 0x12, 0x12, 0x8b, 0x75, 0xb6,
	0x4f, 0xa9, 0x95, 0xd8, 0x97, 0xde, 0x5d, 0xd3, 0x44, 0x8c, 0x88, 0x9d, 0x9f, 0x81, 0x3a, 0x97,
	0xff, 0xd0, 0xb1, 0xea, 0x84, 0x18, 0xd2, 0x92, 0x0c, 0xfc, 0x87, 0x4e, 0xc8, 0x77, 0xd7, 0xaa,
	0x75, 0x18, 0x98, 0xfc, 0xdd, 0x7b, 0x9f, 0xde, 0xf7, 0xbd, 0xe7, 0x3b, 0xf8, 0x3c, 0x60, 0x71,
	0xd4, 0x25, 0xce, 0x70, 0x23, 0x20, 0x02, 0x6f, 0x38, 0x5d, 0x92, 0x12, 0x1e, 0x73, 0x7b, 0xc0,
	0xa8, 0xa0, 0xe8, 0xa1, 0x62, 0x6d, 0xcd, 0xd6, 0x2a, 0x5d, 0xda, 0xa5, 0x92, 0x72, 0xb2, 0x4a,
	0x75, 0xd5, 0x9e, 0x85, 0x94, 0x27, 0x94, 0xfb, 0x8a, 0x50, 0x07, 0x4d, 0xd5, 0x72, 0xf2, 0x84,
	0x85, 0x9d, 0xb6, 0xe6, 0x5e, 0xe6, 0xb8, 0x90, 0xa6, 0x43, 0xc2, 0x78, 0x4c, 0x53, 0x7f, 0x80,
	0x63, 0xa6, 0xba, 0x1a, 0x3f, 0x8a, 0xb0, 0xfc, 0x4e, 0x2d, 0xb5, 0x2f, 0xb0, 0x20, 0xe8, 0x0d,
	0x2c, 0x0d, 0x30, 0xc3, 0x09, 0x37, 0x40, 0x1d, 0x34, 0x57, 0x3b, 0x55, 0xfb, 0xfe, 0x92, 0xf6,
	0x9e, 0x64, 0xdd, 0xc5, 0xb3, 0x89, 0x55, 0xf0, 0x74, 0x2f, 0xfa, 0x0c, 0x91, 0x9c, 0xed, 0xab,
	0x66, 0x39, 0x81, 0x1b, 0xc5, 0xfa, 0x42, 0x73, 0xb5, 0x63, 0xe5, 0x15, 0xb6, 0xbd, 0xad, 0x4e,
	0xdb, 0x95, 0xd8, 0x1e, 0x8e, 0x99, 0xfb, 0x2a, 0x93, 0x9a, 0x4e, 0xac, 0xb5, 0x1c, 0xc1, 0x4f,
	0x2e, 0xe7, 0x31, 0x6f, 0x4d, 0x0e, 0xba, 0x83, 0x20, 0x06, 0xab, 0x29, 0x19, 0x09, 0xff, 0x38,
	0x16, 0x07, 0x11, 0xc3, 0xc7, 0x3e, 0x27, 0x87, 0x47, 0x24, 0x0d, 0x89, 0xb1, 0x50, 0x07, 0xcd,
	0x15, 0xf7, 0x6d, 0xa6, 0xff, 0x6b, 0x62, 0xad, 0x77, 0x63, 0x71, 0x70, 0x14, 0xd8, 0x21, 0x4d,
	0x74, 0x8c, 0xfa, 0xd3, 0xe2, 0x51, 0xcf, 0x11, 0xe3, 0x01, 0xe1, 0xf6, 0x6e, 0x2a, 0x2e, 0x4e,
	0x5b, 0x50, 0xa7, 0xbc, 0x9b, 0x0a, 0xaf, 0x92, 0x69, 0x7f, 0xd4, 0xd2, 0xfb, 0x5a, 0xb9, 0xf1,
	0xa7, 0x08, 0x4b, 0x2a, 0x09, 0xf4, 0x15, 0xc0, 0x0a, 0x49, 0x71, 0xd0, 0x27, 0x91, 0xaf, 0x42,
	0x10, 0xb4, 0x47, 0xd2, 0x2c, 0xc0, 0xcc, 0xfe, 0x8b, 0x39, 0xfb, 0xaa, 0x57, 0x1a, 0xfb, 0x90,
	0x75, 0xba, 0x6d, 0x1d, 0x00, 0x9a, 0xa3, 0xb2, 0x08, 0xfe, 0x81, 0x7a, 0x48, 0x0f, 0xdc, 0xce,
	0xe6, 0x29, 0x0c, 0x05, 0x70, 0x89, 0x91, 0x3e, 0x1e, 0x13, 0x66, 0x14, 0xeb, 0xa0, 0x59, 0x76,
	0x77, 0xae, 0x27, 0x56, 0xeb, 0x3f, 0x3c, 0x6f, 0x86, 0xe1, 0x66, 0x14, 0x31, 0xc2, 0xf9, 0xc5,
	0x69, 0xeb, 0x89, 0xb6, 0xae, 0x11, 0x77, 0x2c, 0x08, 0xf7, 0x6e, 0x84, 0xd1, 0x21, 0x34, 0x6e,
	0xac, 0xe6, 0xee, 0x13, 0x37, 0x16, 0xa4, 0x5d, 0x33, 0x6f, 0x77, 0xeb, 0xb6, 0x4f, 0xfe, 0xec,
	0xa7, 0x99, 0xd7, 0x93, 0x4b, 0xeb, 0xd1, 0x7d, 0x9c, 0x7b, 0x55, 0x2d, 0x9c, 0xc3, 0x1b, 0x5f,
	0x00, 0x7c, 0x3c, 0x97, 0x00, 0x5a, 0x87, 0x4b, 0x58, 0x6d, 0x28, 0xef, 0x69, 0xd9, 0x2d, 0x5f,
	0x4f, 0xac, 0xe5, 0x1d, 0x32, 0xd2, 0x0b, 0x6b, 0x12, 0x21, 0xb8, 0x98, 0xe2, 0x84, 0xc8, 0x44,
	0x56, 0x3c, 0x59, 0xa3, 0x2a, 0x2c, 0xf1, 0x71, 0x12, 0xd0, 0xbe, 0xba, 0x1f, 0x9e, 0x3e, 0xa1,
	0x1a, 0x5c, 0x8e, 0x48, 0x18, 0x27, 0xb8, 0xcf, 0x8d, 0xc5, 0x3a, 0x68, 0x3e, 0xf0, 0x6e, 0xcf,
	0xee, 0xfb, 0xab, 0xdf, 0x26, 0xf8, 0x3e, 0x35, 0xc1, 0xd9, 0xd4, 0x04, 0xe7, 0x53, 0x13, 0x5c,
	0x4d, 0x4d, 0xf0, 0x6d, 0x66, 0x16, 0xce, 0x67, 0x66, 0xe1, 0xe7, 0xcc, 0x2c, 0x7c, 0x72, 0xee,
	0x24, 0xdd, 0xc3, 0x43, 0xdc, 0xea, 0xe3, 0x80, 0xab, 0x4a, 0xbf, 0xc4, 0x91, 0xa3, 0x0b, 0x19,
	0x7b, 0x50, 0x92, 0x2f, 0xf0, 0xf5, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x16, 0x25, 0xe9, 0x87,
	0x24, 0x04, 0x00, 0x00,
}

func (this *GenesisState) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*GenesisState)
	if !ok {
		that2, ok := that.(GenesisState)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *GenesisState")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *GenesisState but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *GenesisState but is not nil && this == nil")
	}
	if !this.Params.Equal(&that1.Params) {
		return fmt.Errorf("Params this(%v) Not Equal that(%v)", this.Params, that1.Params)
	}
	if len(this.ERC20BridgePairs) != len(that1.ERC20BridgePairs) {
		return fmt.Errorf("ERC20BridgePairs this(%v) Not Equal that(%v)", len(this.ERC20BridgePairs), len(that1.ERC20BridgePairs))
	}
	for i := range this.ERC20BridgePairs {
		if !this.ERC20BridgePairs[i].Equal(&that1.ERC20BridgePairs[i]) {
			return fmt.Errorf("ERC20BridgePairs this[%v](%v) Not Equal that[%v](%v)", i, this.ERC20BridgePairs[i], i, that1.ERC20BridgePairs[i])
		}
	}
	if !this.NextWithdrawSequence.Equal(that1.NextWithdrawSequence) {
		return fmt.Errorf("NextWithdrawSequence this(%v) Not Equal that(%v)", this.NextWithdrawSequence, that1.NextWithdrawSequence)
	}
	return nil
}
func (this *GenesisState) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GenesisState)
	if !ok {
		that2, ok := that.(GenesisState)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Params.Equal(&that1.Params) {
		return false
	}
	if len(this.ERC20BridgePairs) != len(that1.ERC20BridgePairs) {
		return false
	}
	for i := range this.ERC20BridgePairs {
		if !this.ERC20BridgePairs[i].Equal(&that1.ERC20BridgePairs[i]) {
			return false
		}
	}
	if !this.NextWithdrawSequence.Equal(that1.NextWithdrawSequence) {
		return false
	}
	return true
}
func (this *Params) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *Params")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *Params but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *Params but is not nil && this == nil")
	}
	if len(this.EnabledERC20Tokens) != len(that1.EnabledERC20Tokens) {
		return fmt.Errorf("EnabledERC20Tokens this(%v) Not Equal that(%v)", len(this.EnabledERC20Tokens), len(that1.EnabledERC20Tokens))
	}
	for i := range this.EnabledERC20Tokens {
		if !this.EnabledERC20Tokens[i].Equal(&that1.EnabledERC20Tokens[i]) {
			return fmt.Errorf("EnabledERC20Tokens this[%v](%v) Not Equal that[%v](%v)", i, this.EnabledERC20Tokens[i], i, that1.EnabledERC20Tokens[i])
		}
	}
	if !bytes.Equal(this.Relayer, that1.Relayer) {
		return fmt.Errorf("Relayer this(%v) Not Equal that(%v)", this.Relayer, that1.Relayer)
	}
	if len(this.EnabledConversionPairs) != len(that1.EnabledConversionPairs) {
		return fmt.Errorf("EnabledConversionPairs this(%v) Not Equal that(%v)", len(this.EnabledConversionPairs), len(that1.EnabledConversionPairs))
	}
	for i := range this.EnabledConversionPairs {
		if !this.EnabledConversionPairs[i].Equal(&that1.EnabledConversionPairs[i]) {
			return fmt.Errorf("EnabledConversionPairs this[%v](%v) Not Equal that[%v](%v)", i, this.EnabledConversionPairs[i], i, that1.EnabledConversionPairs[i])
		}
	}
	return nil
}
func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.EnabledERC20Tokens) != len(that1.EnabledERC20Tokens) {
		return false
	}
	for i := range this.EnabledERC20Tokens {
		if !this.EnabledERC20Tokens[i].Equal(&that1.EnabledERC20Tokens[i]) {
			return false
		}
	}
	if !bytes.Equal(this.Relayer, that1.Relayer) {
		return false
	}
	if len(this.EnabledConversionPairs) != len(that1.EnabledConversionPairs) {
		return false
	}
	for i := range this.EnabledConversionPairs {
		if !this.EnabledConversionPairs[i].Equal(&that1.EnabledConversionPairs[i]) {
			return false
		}
	}
	return true
}
func (this *EnabledERC20Token) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*EnabledERC20Token)
	if !ok {
		that2, ok := that.(EnabledERC20Token)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *EnabledERC20Token")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *EnabledERC20Token but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *EnabledERC20Token but is not nil && this == nil")
	}
	if !bytes.Equal(this.Address, that1.Address) {
		return fmt.Errorf("Address this(%v) Not Equal that(%v)", this.Address, that1.Address)
	}
	if this.Name != that1.Name {
		return fmt.Errorf("Name this(%v) Not Equal that(%v)", this.Name, that1.Name)
	}
	if this.Symbol != that1.Symbol {
		return fmt.Errorf("Symbol this(%v) Not Equal that(%v)", this.Symbol, that1.Symbol)
	}
	if this.Decimals != that1.Decimals {
		return fmt.Errorf("Decimals this(%v) Not Equal that(%v)", this.Decimals, that1.Decimals)
	}
	return nil
}
func (this *EnabledERC20Token) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*EnabledERC20Token)
	if !ok {
		that2, ok := that.(EnabledERC20Token)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.Address, that1.Address) {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Symbol != that1.Symbol {
		return false
	}
	if this.Decimals != that1.Decimals {
		return false
	}
	return true
}
func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.NextWithdrawSequence.Size()
		i -= size
		if _, err := m.NextWithdrawSequence.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.ERC20BridgePairs) > 0 {
		for iNdEx := len(m.ERC20BridgePairs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ERC20BridgePairs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EnabledConversionPairs) > 0 {
		for iNdEx := len(m.EnabledConversionPairs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.EnabledConversionPairs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Relayer) > 0 {
		i -= len(m.Relayer)
		copy(dAtA[i:], m.Relayer)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Relayer)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.EnabledERC20Tokens) > 0 {
		for iNdEx := len(m.EnabledERC20Tokens) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.EnabledERC20Tokens[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *EnabledERC20Token) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EnabledERC20Token) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EnabledERC20Token) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Decimals != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Decimals))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.ERC20BridgePairs) > 0 {
		for _, e := range m.ERC20BridgePairs {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.NextWithdrawSequence.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.EnabledERC20Tokens) > 0 {
		for _, e := range m.EnabledERC20Tokens {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = len(m.Relayer)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.EnabledConversionPairs) > 0 {
		for _, e := range m.EnabledConversionPairs {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *EnabledERC20Token) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Decimals != 0 {
		n += 1 + sovGenesis(uint64(m.Decimals))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ERC20BridgePairs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ERC20BridgePairs = append(m.ERC20BridgePairs, ERC20BridgePair{})
			if err := m.ERC20BridgePairs[len(m.ERC20BridgePairs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextWithdrawSequence", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.NextWithdrawSequence.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnabledERC20Tokens", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EnabledERC20Tokens = append(m.EnabledERC20Tokens, EnabledERC20Token{})
			if err := m.EnabledERC20Tokens[len(m.EnabledERC20Tokens)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relayer", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relayer = append(m.Relayer[:0], dAtA[iNdEx:postIndex]...)
			if m.Relayer == nil {
				m.Relayer = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnabledConversionPairs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EnabledConversionPairs = append(m.EnabledConversionPairs, ConversionPair{})
			if err := m.EnabledConversionPairs[len(m.EnabledConversionPairs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *EnabledERC20Token) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: EnabledERC20Token: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EnabledERC20Token: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decimals", wireType)
			}
			m.Decimals = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Decimals |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
