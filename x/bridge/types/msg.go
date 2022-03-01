package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgBridgeERC20FromEthereum{}
)

// NewMsgBridgeERC20FromEthereum returns a new MsgBridgeERC20FromEthereum
func NewMsgBridgeERC20FromEthereum(
	relayer string,
	ethereumERC20Address string,
	amount sdk.Int,
	receiver string,
	sequence sdk.Int,
) MsgBridgeERC20FromEthereum {
	return MsgBridgeERC20FromEthereum{
		Relayer:              relayer,
		EthereumERC20Address: ethereumERC20Address,
		Amount:               amount,
		Receiver:             receiver,
		Sequence:             sequence,
	}
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgBridgeERC20FromEthereum) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Relayer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgBridgeERC20FromEthereum) ValidateBasic() error {
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
