package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bridgekeeper "github.com/kava-labs/kava-bridge/x/bridge/keeper"
	bridgetypes "github.com/kava-labs/kava-bridge/x/bridge/types"
)

// BridgeRelayerDecorator will validate bridge messages are from permissioned
// signer.
type BridgeRelayerDecorator struct {
	bk bridgekeeper.Keeper
}

func NewBridgeRelayerDecorator(bk bridgekeeper.Keeper) BridgeRelayerDecorator {
	return BridgeRelayerDecorator{
		bk: bk,
	}
}

func (brd BridgeRelayerDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	relayer := brd.bk.GetRelayer(ctx)

	for _, msg := range tx.GetMsgs() {
		if _, ok := msg.(*bridgetypes.MsgBridgeERC20FromEthereum); ok {
			signerAddrs := msg.GetSigners()

			if len(signerAddrs) != 1 {
				return ctx, sdkerrors.Wrapf(
					sdkerrors.ErrUnauthorized,
					"invalid number of signer; expected: 1, got %d",
					len(signerAddrs),
				)
			}

			if !relayer.Equals(signerAddrs[0]) {
				return ctx, sdkerrors.Wrapf(
					sdkerrors.ErrUnauthorized,
					"signer not authorized for bridge message",
				)
			}
		}
	}

	return next(ctx, tx, simulate)
}
