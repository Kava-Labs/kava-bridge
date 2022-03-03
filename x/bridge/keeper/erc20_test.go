package keeper_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/server/config"
	etherminttests "github.com/tharsis/ethermint/tests"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type ERC20TestSuite struct {
	testutil.Suite
}

func TestERC20TestSuite(t *testing.T) {
	suite.Run(t, new(ERC20TestSuite))
}

func (suite *ERC20TestSuite) deployERC20() common.Address {
	// We can assume token is valid as it is from params and should be validated
	token := types.NewEnabledERC20Token(
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"Wrapped ETH",
		"WETH",
		18,
	)

	contractAddr, err := suite.App.BridgeKeeper.DeployMintableERC20Contract(suite.Ctx, token)
	suite.Require().NoError(err)
	suite.Require().Greater(len(contractAddr), 0)

	return contractAddr
}

func (suite *ERC20TestSuite) TestDeployERC20() {
	suite.deployERC20()
}

func (suite *ERC20TestSuite) queryContract(
	contractAbi abi.ABI,
	from common.Address,
	fromKey *ethsecp256k1.PrivKey,
	contract common.Address,
	method string,
	args ...interface{},
) ([]interface{}, error) {
	// Pack query args
	data, err := contractAbi.Pack(method, args...)
	suite.Require().NoError(err)

	// Send TX
	res := suite.sendTx(contract, from, fromKey, data)

	// Check for VM errors and unpack returned data
	switch res.VmError {
	case vm.ErrExecutionReverted.Error():
		response, err := abi.UnpackRevert(res.Ret)
		suite.Require().NoError(err)

		return nil, errors.New(response)
	case "": // No error, continue
	default:
		panic(fmt.Sprintf("unhandled vm error response: %v", res.VmError))
	}

	// Unpack response
	unpackedRes, err := contractAbi.Unpack(method, res.Ret)
	suite.Require().NoError(err)

	return unpackedRes, nil
}

func (suite *ERC20TestSuite) TestERC20Query() {
	contractAddr := suite.deployERC20()

	// Query ERC20.decimals()
	addr := common.BytesToAddress(suite.Key1.PubKey().Address())
	res, err := suite.queryContract(
		contract.ERC20MintableBurnableContract.ABI,
		addr,
		suite.Key1,
		contractAddr,
		"decimals",
	)
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)

	// Type should match abi json output
	decimals, ok := res[0].(uint8)
	suite.Require().True(ok, "decimals should respond with first uint8")
	suite.Require().Equal(uint8(18), decimals)
}

func (suite *ERC20TestSuite) TestERC20Mint_Unauthorized() {
	contractAddr := suite.deployERC20()

	// ERC20.mint() from key1 to key2
	addr := common.BytesToAddress(suite.Key1.PubKey().Address())
	receiver := common.BytesToAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(10)
	_, err := suite.queryContract(
		contract.ERC20MintableBurnableContract.ABI,
		addr,
		suite.Key1,
		contractAddr,
		"mint",
		receiver,
		&amount,
	)
	suite.Require().Error(err)
	suite.Require().Equal("Ownable: caller is not the owner", err.Error())
}

func (suite *ERC20TestSuite) TestERC20Mint() {
	contractAddr := suite.deployERC20()

	// We can't test mint by module account like the Unauthorized test as we
	// cannot sign as the module account. Instead, we test the keeper method for
	// minting.

	receiver := common.BytesToAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(1234)
	err := suite.App.BridgeKeeper.MintERC20(suite.Ctx, contractAddr, receiver, amount)
	suite.Require().NoError(err)

	// Query ERC20.balanceOf()
	addr := common.BytesToAddress(suite.Key1.PubKey().Address())
	res, err := suite.queryContract(
		contract.ERC20MintableBurnableContract.ABI,
		addr,
		suite.Key1,
		contractAddr,
		"balanceOf",
		receiver,
	)
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)

	balance, ok := res[0].(*big.Int)
	suite.Require().True(ok, "balanceOf should respond with *big.Int")
	suite.Require().Equal(big.NewInt(1234), balance)
}

func (suite *ERC20TestSuite) sendTx(
	contractAddr,
	from common.Address,
	signerKey *ethsecp256k1.PrivKey,
	transferData []byte,
) *evmtypes.MsgEthereumTxResponse {
	ctx := sdk.WrapSDKContext(suite.Ctx)
	chainID := suite.App.EvmKeeper.ChainID()

	args, err := json.Marshal(&evmtypes.TransactionArgs{
		To:   &contractAddr,
		From: &from,
		Data: (*hexutil.Bytes)(&transferData),
	})
	suite.Require().NoError(err)
	res, err := suite.QueryClientEvm.EstimateGas(ctx, &evmtypes.EthCallRequest{
		Args:   args,
		GasCap: uint64(config.DefaultGasCap),
	})
	suite.Require().NoError(err)

	nonce := suite.App.EvmKeeper.GetNonce(suite.Ctx, suite.Address)

	baseFee := suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx)
	suite.Require().NotNil(baseFee, "base fee is nil")

	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	suite.MintFeeCollector(sdk.NewCoins(
		sdk.NewCoin(
			"ukava",
			sdk.NewInt(baseFee.Int64()*int64(res.Gas)),
		)))

	ercTransferTx := evmtypes.NewTx(
		chainID,
		nonce,
		&contractAddr,
		nil,       // amount
		res.Gas*2, // gasLimit, TODO: runs out of gas with just res.Gas, ex: estimated was 21572 but used 24814
		nil,       // gasPrice
		suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx), // gasFeeCap
		big.NewInt(1), // gasTipCap
		transferData,
		&ethtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = hex.EncodeToString(signerKey.PubKey().Address())
	err = ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), etherminttests.NewSigner(signerKey))
	suite.Require().NoError(err)

	rsp, err := suite.App.EvmKeeper.EthereumTx(ctx, ercTransferTx)
	suite.Require().NoError(err)
	// Do not check vm error here since we want to check for errors later

	return rsp
}

func (suite *ERC20TestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.App.BankKeeper.MintCoins(suite.Ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.App.BankKeeper.SendCoinsFromModuleToModule(
		suite.Ctx,
		minttypes.ModuleName,
		authtypes.FeeCollectorName,
		coins,
	)
	suite.Require().NoError(err)
}
