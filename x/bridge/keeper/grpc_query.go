package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type queryServer struct {
	keeper Keeper
}

// NewQueryServerImpl creates a new server for handling gRPC queries.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{keeper: k}
}

var _ types.QueryServer = queryServer{}

// Params queries module params
func (s queryServer) Params(stdCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(stdCtx)
	params := s.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (s queryServer) BridgedERC20Pairs(
	stdCtx context.Context,
	req *types.QueryBridgedERC20PairsRequest,
) (*types.QueryBridgedERC20PairsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	var bridgePairs types.ERC20BridgePairs
	s.keeper.IterateERC20BridgePairs(ctx, func(pair types.ERC20BridgePair) bool {
		bridgePairs = append(bridgePairs, pair)
		return false
	})

	return &types.QueryBridgedERC20PairsResponse{
		ERC20BridgePairs: bridgePairs,
	}, nil
}
