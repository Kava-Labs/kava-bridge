package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/ethereum/go-ethereum/common"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg            = &MsgBridgeEthereumToKava{}
	_ sdk.Msg            = &MsgConvertCoinToERC20{}
	_ sdk.Msg            = &MsgConvertERC20ToCoin{}
	_ legacytx.LegacyMsg = &MsgConvertERC20ToCoin{}
)

// NewMsgBridgeEthereumToKava returns a new MsgBridgeEthereumToKava
func NewMsgBridgeEthereumToKava(
	relayer string,
	ethereumERC20Address string,
	amount sdk.Int,
	receiver string,
	sequence sdk.Int,
) MsgBridgeEthereumToKava {
	return MsgBridgeEthereumToKava{
		Relayer:              relayer,
		EthereumERC20Address: ethereumERC20Address,
		Amount:               amount,
		Receiver:             receiver,
		Sequence:             sequence,
	}
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgBridgeEthereumToKava) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Relayer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgBridgeEthereumToKava) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Relayer)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !common.IsHexAddress(msg.EthereumERC20Address) {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidAddress,
			"ethereum ERC20 address is not a valid hex address",
		)
	}

	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "amount must be positive non-zero %s", msg.Amount)
	}

	if !common.IsHexAddress(msg.Receiver) {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidAddress,
			"receiver address is not a valid hex address",
		)
	}

	if msg.Sequence.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidSequence, "sequence is negative %s", msg.Sequence)
	}

	return nil
}

// NewMsgConvertCoinToERC20 returns a new MsgConvertCoinToERC20
func NewMsgConvertCoinToERC20(
	initiator string,
	receiver string,
	amount sdk.Coin,
) MsgConvertCoinToERC20 {
	return MsgConvertCoinToERC20{
		Initiator: initiator,
		Receiver:  receiver,
		Amount:    &amount,
	}
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgConvertCoinToERC20) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgConvertCoinToERC20) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !common.IsHexAddress(msg.Receiver) {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidAddress,
			"Receiver is not a valid hex address",
		)
	}

	if msg.Amount.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "amount cannot be zero")
	}

	// Checks for negative
	return msg.Amount.Validate()
}

// NewMsgConvertERC20ToCoin returns a new MsgConvertERC20ToCoin
func NewMsgConvertERC20ToCoin(
	initiator InternalEVMAddress,
	receiver sdk.AccAddress,
	contractAddr InternalEVMAddress,
	amount sdk.Int,
) MsgConvertERC20ToCoin {
	return MsgConvertERC20ToCoin{
		Initiator:        initiator.String(),
		Receiver:         receiver.String(),
		KavaERC20Address: contractAddr.String(),
		Amount:           amount,
	}
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgConvertERC20ToCoin) GetSigners() []sdk.AccAddress {
	addr := common.HexToAddress(msg.Initiator)
	sender := sdk.AccAddress(addr.Bytes())
	return []sdk.AccAddress{sender}
}

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgConvertERC20ToCoin) ValidateBasic() error {
	if !common.IsHexAddress(msg.Initiator) {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidAddress,
			"initiator is not a valid hex address",
		)
	}

	if !common.IsHexAddress(msg.KavaERC20Address) {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidAddress,
			"erc20 contract address is not a valid hex address",
		)
	}

	_, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "receiver is not a valid bech32 address")
	}

	if msg.Amount.LTE(sdk.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "amount cannot be zero or less")
	}

	return nil
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (msg MsgConvertERC20ToCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// Route implements the LegacyMsg.Route method.
func (msg MsgConvertERC20ToCoin) Route() string {
	return RouterKey
}

// Type implements the LegacyMsg.Type method.
func (msg MsgConvertERC20ToCoin) Type() string {
	return "bridge_convert_erc20_to_coin"
}
