package keeper_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/server/config"
	etherminttests "github.com/tharsis/ethermint/tests"
	etherminttypes "github.com/tharsis/ethermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	"github.com/kava-labs/kava-bridge/app"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type ERC20TestSuite struct {
	suite.Suite

	ctx            sdk.Context
	app            app.TestApp
	address        common.Address
	key1           *ethsecp256k1.PrivKey
	key2           *ethsecp256k1.PrivKey
	consAddress    sdk.ConsAddress
	relayerAddress sdk.AccAddress

	queryClientEvm evmtypes.QueryClient
}

func (suite *ERC20TestSuite) SetupTest() {
	suite.app = app.NewTestApp()
	cdc := suite.app.AppCodec()

	// consensus key
	consPriv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.consAddress = sdk.ConsAddress(consPriv.PubKey().Address())

	// relayer key
	relayerPriv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.relayerAddress = sdk.AccAddress(relayerPriv.PubKey().Address())

	// test user keys that have no minting permissions
	suite.key1, err = ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	suite.key2, err = ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	coins := sdk.NewCoins(sdk.NewInt64Coin("ukava", 1000_000_000_000_000_000))
	authGS := app.NewFundedGenStateWithSameCoins(cdc, coins, []sdk.AccAddress{
		sdk.AccAddress(suite.key1.PubKey().Address()),
		sdk.AccAddress(suite.key2.PubKey().Address()),
	})

	// Genesis states
	evmGs := evmtypes.NewGenesisState(
		evmtypes.NewParams("ukava", true, true, evmtypes.DefaultChainConfig()),
		nil,
	)

	bridgeGs := types.NewGenesisState(types.NewParams(
		types.EnabledERC20Tokens{
			types.NewEnabledERC20Token(
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				"Wrapped ETH",
				"WETH",
				18,
			),
		},
		suite.relayerAddress,
	))

	feemarketGenesis := feemarkettypes.DefaultGenesisState()
	feemarketGenesis.Params.EnableHeight = 1
	feemarketGenesis.Params.NoBaseFee = false

	gs := app.GenesisState{
		types.ModuleName:          cdc.MustMarshalJSON(&bridgeGs),
		evmtypes.ModuleName:       cdc.MustMarshalJSON(evmGs),
		feemarkettypes.ModuleName: cdc.MustMarshalJSON(feemarketGenesis),
	}

	// Initialize the chain
	suite.app.InitializeFromGenesisStates(authGS, gs)

	// InitializeFromGenesisStates commits first block so we start at 2 here
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		Height:          suite.app.LastBlockHeight() + 1,
		ChainID:         "kavatest_1-1",
		Time:            time.Now().UTC(),
		ProposerAddress: suite.consAddress.Bytes(),
		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	// We need to set the validator as calling the EVM looks up the validator address
	// https://github.com/tharsis/ethermint/blob/f21592ebfe74da7590eb42ed926dae970b2a9a3f/x/evm/keeper/state_transition.go#L487
	// evmkeeper.EVMConfig() will return error "failed to load evm config" if not set
	acc := &etherminttypes.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), nil, 0, 0),
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}

	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, consPriv.PubKey(), stakingtypes.Description{})
	suite.Require().NoError(err)

	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	suite.Require().NoError(err)

	suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
	suite.queryClientEvm = evmtypes.NewQueryClient(queryHelperEvm)

	// We need to commit so that the ethermint feemarket beginblock runs to set the minfee
	// feeMarketKeeper.GetBaseFee() will return nil otherwise
	suite.Commit()
}

func TestERC20TestSuite(t *testing.T) {
	suite.Run(t, new(ERC20TestSuite))
}

func (suite *ERC20TestSuite) Commit() {
	_ = suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height += 1
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.ctx = suite.app.NewContext(false, header)
}

func (suite *ERC20TestSuite) deployERC20() common.Address {
	// We can assume token is valid as it is from params and should be validated
	token := types.NewEnabledERC20Token(
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"Wrapped ETH",
		"WETH",
		18,
	)

	contractAddr, err := suite.app.BridgeKeeper.DeployMintableERC20Contract(suite.ctx, token)
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
	addr := common.BytesToAddress(suite.key1.PubKey().Address())
	res, err := suite.queryContract(
		contract.ERC20MintableBurnableContract.ABI,
		addr,
		suite.key1,
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
	addr := common.BytesToAddress(suite.key1.PubKey().Address())
	receiver := common.BytesToAddress(suite.key2.PubKey().Address())
	amount := big.NewInt(10)
	_, err := suite.queryContract(
		contract.ERC20MintableBurnableContract.ABI,
		addr,
		suite.key1,
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

	receiver := common.BytesToAddress(suite.key2.PubKey().Address())
	amount := big.NewInt(1234)
	err := suite.app.BridgeKeeper.MintERC20(suite.ctx, contractAddr, receiver, amount)
	suite.Require().NoError(err)

	// Query ERC20.balanceOf()
	addr := common.BytesToAddress(suite.key1.PubKey().Address())
	res, err := suite.queryContract(
		contract.ERC20MintableBurnableContract.ABI,
		addr,
		suite.key1,
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
	ctx := sdk.WrapSDKContext(suite.ctx)
	chainID := suite.app.EvmKeeper.ChainID()

	args, err := json.Marshal(&evmtypes.TransactionArgs{
		To:   &contractAddr,
		From: &from,
		Data: (*hexutil.Bytes)(&transferData),
	})
	suite.Require().NoError(err)
	res, err := suite.queryClientEvm.EstimateGas(ctx, &evmtypes.EthCallRequest{
		Args:   args,
		GasCap: uint64(config.DefaultGasCap),
	})
	suite.Require().NoError(err)

	nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, suite.address)

	baseFee := suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx)
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
		suite.app.FeeMarketKeeper.GetBaseFee(suite.ctx), // gasFeeCap
		big.NewInt(1), // gasTipCap
		transferData,
		&ethtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = hex.EncodeToString(signerKey.PubKey().Address())
	err = ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), etherminttests.NewSigner(signerKey))
	suite.Require().NoError(err)

	rsp, err := suite.app.EvmKeeper.EthereumTx(ctx, ercTransferTx)
	suite.Require().NoError(err)
	// Do not check vm error here since we want to check for errors later

	return rsp
}

func (suite *ERC20TestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(
		suite.ctx,
		minttypes.ModuleName,
		authtypes.FeeCollectorName,
		coins,
	)
	suite.Require().NoError(err)
}
