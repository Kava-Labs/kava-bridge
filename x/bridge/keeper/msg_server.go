package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the auction MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) BridgeERC20FromEthereum(
	goCtx context.Context,
	msg *types.MsgBridgeERC20FromEthereum,
) (*types.MsgBridgeERC20FromEthereumResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	receiver := common.Address{}
	err := receiver.UnmarshalText([]byte(msg.Receiver))
	if err != nil {
		return nil, fmt.Errorf("invalid receiver: %w", err)
	}

	externalAddress := types.ExternalEVMAddress{}
	err = externalAddress.UnmarshalText([]byte(msg.EthereumERC20Address))
	if err != nil {
		return nil, fmt.Errorf("invalid EthereumERC20Address: %w", err)
	}

	// TODO: Antehandler to check if this is made by a permissioned signer
	internalAddress, found := s.keeper.GetBridgedInternalEVMAddress(ctx, externalAddress)
	if !found {
		enabledToken, err := s.keeper.GetEnabledERC20Token(ctx, externalAddress.String())
		if err != nil {
			return nil, err
		}

		internalAddress, err = s.keeper.DeployMintableERC20Contract(ctx, enabledToken)
		if err != nil {
			return nil, err
		}
	}

	err = s.keeper.MintERC20(ctx, internalAddress, receiver, msg.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyRelayer, msg.Relayer),
			sdk.NewAttribute(types.AttributeKeyEthereumERC20Address, msg.EthereumERC20Address),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySequence, msg.Sequence.String()),
		),
	)
	return &types.MsgBridgeERC20FromEthereumResponse{}, nil
}
