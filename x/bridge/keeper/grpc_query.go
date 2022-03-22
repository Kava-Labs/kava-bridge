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

// ERC20BridgePairs queries ERC20 bridge pair addresses.
func (s queryServer) ERC20BridgePairs(
	stdCtx context.Context,
	req *types.QueryERC20BridgePairsRequest,
) (*types.QueryERC20BridgePairsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	var bridgePairs types.ERC20BridgePairs
	s.keeper.IterateBridgePairs(ctx, func(pair types.ERC20BridgePair) bool {
		bridgePairs = append(bridgePairs, pair)
		return false
	})

	return &types.QueryERC20BridgePairsResponse{
		ERC20BridgePairs: bridgePairs,
	}, nil
}
