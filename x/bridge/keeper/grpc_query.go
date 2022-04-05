package keeper

import (
	"bytes"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

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

// ERC20BridgePair queries for a ERC20 bridge pair's addresses.
func (s queryServer) ERC20BridgePair(
	stdCtx context.Context,
	req *types.QueryERC20BridgePairRequest,
) (*types.QueryERC20BridgePairResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	if !common.IsHexAddress(req.Address) {
		return nil, status.Error(codes.InvalidArgument, "not a valid hex address")
	}
	addrBytes := common.HexToAddress(req.Address)

	var bridgePair types.ERC20BridgePair
	found := false
	s.keeper.IterateBridgePairs(ctx, func(pair types.ERC20BridgePair) bool {
		// Match either internal or external
		if bytes.Equal(pair.ExternalERC20Address, addrBytes.Bytes()) || bytes.Equal(pair.InternalERC20Address, addrBytes.Bytes()) {
			bridgePair = pair
			found = true
			return true
		}

		return false
	})

	if !found {
		return nil, status.Error(codes.NotFound, "could not find bridge pair with provided address")
	}

	return &types.QueryERC20BridgePairResponse{
		ERC20BridgePair: bridgePair,
	}, nil
}
