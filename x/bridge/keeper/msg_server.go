package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the bridge MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// BridgeEthereumToKava handles a bridge from Ethereum message.
func (s msgServer) BridgeEthereumToKava(
	goCtx context.Context,
	msg *types.MsgBridgeEthereumToKava,
) (*types.MsgBridgeEthereumToKavaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	relayer, err := sdk.AccAddressFromBech32(msg.Relayer)
	if err != nil {
		return nil, fmt.Errorf("invalid Relayer address: %w", err)
	}

	receiver, err := types.NewInternalEVMAddressFromString(msg.Receiver)
	if err != nil {
		return nil, fmt.Errorf("invalid Receiver address: %w", err)
	}

	externalAddress, err := types.NewExternalEVMAddressFromString(msg.EthereumERC20Address)
	if err != nil {
		return nil, fmt.Errorf("invalid EthereumERC20Address: %w", err)
	}

	if err := s.keeper.BridgeEthereumToKava(
		ctx,
		relayer,
		externalAddress,
		receiver,
		msg.Amount.BigInt(),
		msg.Sequence,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Relayer),
		),
	)

	return &types.MsgBridgeEthereumToKavaResponse{}, nil
}

// ConvertCoinToERC20 handles a MsgConvertCoinToERC20 message to convert
// sdk.Coin to Kava EVM tokens.
func (s msgServer) ConvertCoinToERC20(
	goCtx context.Context,
	msg *types.MsgConvertCoinToERC20,
) (*types.MsgConvertCoinToERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	initiator, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		return nil, fmt.Errorf("invalid Initiator address: %w", err)
	}

	receiver, err := types.NewInternalEVMAddressFromString(msg.Receiver)
	if err != nil {
		return nil, fmt.Errorf("invalid Receiver address: %w", err)
	}

	if err := s.keeper.ConvertCoinToERC20(
		ctx,
		initiator,
		receiver,
		*msg.Amount,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Initiator),
		),
	)

	return &types.MsgConvertCoinToERC20Response{}, nil
}
