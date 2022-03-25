package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgBridgeEthereumToKava{}
	_ sdk.Msg = &MsgConvertCoinToERC20{}
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
