package testutil

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
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
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
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
	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

var (
	MinWETHWithdrawAmount  = sdk.NewInt(10_000_000_000_000_000)
	MinWKavaWithdrawAmount = sdk.NewInt(2_000_000)
	MinUSDCWithdrawAmount  = sdk.NewInt(10_000_000)
	MinWBTCWithdrawAmount  = sdk.NewInt(1_000_000_000)
)

type Suite struct {
	suite.Suite

	Ctx            sdk.Context
	App            app.TestApp
	Address        common.Address
	Key1           *ethsecp256k1.PrivKey
	Key1Addr       types.InternalEVMAddress
	Key2           *ethsecp256k1.PrivKey
	ConsAddress    sdk.ConsAddress
	RelayerAddress sdk.AccAddress
	RelayerKey     *ethsecp256k1.PrivKey

	QueryClientEvm    evmtypes.QueryClient
	QueryClientBridge types.QueryClient

	BankKeeper bankkeeper.Keeper
	Keeper     keeper.Keeper
}

func (suite *Suite) SetupTest() {
	suite.App = app.NewTestApp()
	cdc := suite.App.AppCodec()

	suite.BankKeeper = suite.App.BankKeeper
	suite.Keeper = suite.App.BridgeKeeper

	// consensus key
	consPriv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.ConsAddress = sdk.ConsAddress(consPriv.PubKey().Address())

	// relayer key
	relayerPriv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.RelayerAddress = sdk.AccAddress(relayerPriv.PubKey().Address())
	suite.RelayerKey = relayerPriv

	// test user keys that have no minting permissions
	suite.Key1, err = ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.Key1Addr = types.NewInternalEVMAddress(common.BytesToAddress(suite.Key1.PubKey().Address()))

	suite.Key2, err = ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	coins := sdk.NewCoins(sdk.NewInt64Coin("ukava", 1000_000_000_000_000_000))
	authGS := app.NewFundedGenStateWithSameCoins(cdc, coins, []sdk.AccAddress{
		sdk.AccAddress(suite.Key1.PubKey().Address()),
		sdk.AccAddress(suite.Key2.PubKey().Address()),
	})

	// Genesis states
	evmGs := evmtypes.NewGenesisState(
		evmtypes.NewParams("ukava", true, true, evmtypes.DefaultChainConfig()),
		nil,
	)

	bridgeGs := types.NewGenesisState(
		types.NewParams(
			true,
			types.EnabledERC20Tokens{
				types.NewEnabledERC20Token(
					MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
					"Wrapped Ether",
					"WETH",
					18,
					MinWETHWithdrawAmount,
				),
				types.NewEnabledERC20Token(
					MustNewExternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
					"Wrapped Kava",
					"WKAVA",
					6,
					MinWKavaWithdrawAmount,
				),
				types.NewEnabledERC20Token(
					MustNewExternalEVMAddressFromString("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
					"USD Coin",
					"USDC",
					6,
					MinUSDCWithdrawAmount,
				),
				types.NewEnabledERC20Token(
					MustNewExternalEVMAddressFromString("foobar-wbtc"),
					"Wrapped BTC",
					"WBTC",
					8,
					MinWBTCWithdrawAmount,
				),
			},
			suite.RelayerAddress,
			types.NewConversionPairs(
				types.NewConversionPair(
					// First contract bridge module deploys
					MustNewInternalEVMAddressFromString("0x404F9466d758eA33eA84CeBE9E444b06533b369e"),
					"erc20/usdc",
				),
				types.NewConversionPair(
					MustNewInternalEVMAddressFromString("foobar-wbtc"),
					"erc20/wbtc",
				),
			),
		),
		types.NewERC20BridgePairs(
			types.NewERC20BridgePair(
				MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
				MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
			),
			types.NewERC20BridgePair(
				MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
				MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000B"),
			),
		),
		types.DefaultNextWithdrawSequence,
	)

	feemarketGenesis := feemarkettypes.DefaultGenesisState()
	feemarketGenesis.Params.EnableHeight = 1
	feemarketGenesis.Params.NoBaseFee = false

	gs := app.GenesisState{
		types.ModuleName:          cdc.MustMarshalJSON(&bridgeGs),
		evmtypes.ModuleName:       cdc.MustMarshalJSON(evmGs),
		feemarkettypes.ModuleName: cdc.MustMarshalJSON(feemarketGenesis),
	}

	// Initialize the chain
	suite.App.InitializeFromGenesisStates(authGS, gs)

	// InitializeFromGenesisStates commits first block so we start at 2 here
	suite.Ctx = suite.App.NewContext(false, tmproto.Header{
		Height:          suite.App.LastBlockHeight() + 1,
		ChainID:         "kavatest_1-1",
		Time:            time.Now().UTC(),
		ProposerAddress: suite.ConsAddress.Bytes(),
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
		BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(suite.Address.Bytes()), nil, 0, 0),
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}

	suite.App.AccountKeeper.SetAccount(suite.Ctx, acc)

	valAddr := sdk.ValAddress(suite.Address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, consPriv.PubKey(), stakingtypes.Description{})
	suite.Require().NoError(err)

	err = suite.App.StakingKeeper.SetValidatorByConsAddr(suite.Ctx, validator)
	suite.Require().NoError(err)

	suite.App.StakingKeeper.SetValidator(suite.Ctx, validator)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelper, suite.App.EvmKeeper)
	suite.QueryClientEvm = evmtypes.NewQueryClient(queryHelper)

	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(suite.App.BridgeKeeper))
	suite.QueryClientBridge = types.NewQueryClient(queryHelper)

	// We need to commit so that the ethermint feemarket beginblock runs to set the minfee
	// feeMarketKeeper.GetBaseFee() will return nil otherwise
	suite.Commit()
}

func (suite *Suite) Commit() {
	_ = suite.App.Commit()
	header := suite.Ctx.BlockHeader()
	header.Height += 1
	suite.App.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.Ctx = suite.App.NewContext(false, header)
}

func (suite *Suite) DeployERC20() types.InternalEVMAddress {
	// We can assume token is valid as it is from params and should be validated
	token := types.NewEnabledERC20Token(
		MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		"Wrapped ETH",
		"WETH",
		18,
		MinWETHWithdrawAmount,
	)

	contractAddr, err := suite.App.BridgeKeeper.DeployMintableERC20Contract(suite.Ctx, token)
	suite.Require().NoError(err)
	suite.Require().Greater(len(contractAddr.Address), 0)

	return contractAddr
}

func (suite *Suite) TestDeployERC20() {
	suite.DeployERC20()
}

func (suite *Suite) GetERC20BalanceOf(
	contractAbi abi.ABI,
	contractAddr types.InternalEVMAddress,
	accountAddr types.InternalEVMAddress,
) *big.Int {
	// Query ERC20.balanceOf()
	addr := common.BytesToAddress(suite.Key1.PubKey().Address())
	res, err := suite.QueryContract(
		contract.ERC20MintableBurnableContract.ABI,
		addr,
		suite.Key1,
		contractAddr,
		"balanceOf",
		accountAddr.Address,
	)
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)

	balance, ok := res[0].(*big.Int)
	suite.Require().True(ok, "balanceOf should respond with *big.Int")
	return balance
}

func (suite *Suite) QueryContract(
	contractAbi abi.ABI,
	from common.Address,
	fromKey *ethsecp256k1.PrivKey,
	contract types.InternalEVMAddress,
	method string,
	args ...interface{},
) ([]interface{}, error) {
	// Pack query args
	data, err := contractAbi.Pack(method, args...)
	suite.Require().NoError(err)

	// Send TX
	res, err := suite.SendTx(contract, from, fromKey, data)
	suite.Require().NoError(err)

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
	suite.Require().NoErrorf(err, "failed to unpack method %v response", method)

	return unpackedRes, nil
}

// SendTx submits a transaction to the block.
func (suite *Suite) SendTx(
	contractAddr types.InternalEVMAddress,
	from common.Address,
	signerKey *ethsecp256k1.PrivKey,
	transferData []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	ctx := sdk.WrapSDKContext(suite.Ctx)
	chainID := suite.App.EvmKeeper.ChainID()

	args, err := json.Marshal(&evmtypes.TransactionArgs{
		To:   &contractAddr.Address,
		From: &from,
		Data: (*hexutil.Bytes)(&transferData),
	})
	if err != nil {
		return nil, err
	}
	gasRes, err := suite.QueryClientEvm.EstimateGas(ctx, &evmtypes.EthCallRequest{
		Args:   args,
		GasCap: config.DefaultGasCap,
	})
	if err != nil {
		return nil, err
	}

	nonce := suite.App.EvmKeeper.GetNonce(suite.Ctx, suite.Address)

	baseFee := suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx)
	suite.Require().NotNil(baseFee, "base fee is nil")

	// Mint the max gas to the FeeCollector to ensure balance in case of refund
	suite.MintFeeCollector(sdk.NewCoins(
		sdk.NewCoin(
			"ukava",
			sdk.NewInt(baseFee.Int64()*int64(gasRes.Gas*2)),
		)))

	ercTransferTx := evmtypes.NewTx(
		chainID,
		nonce,
		&contractAddr.Address,
		nil,          // amount
		gasRes.Gas*2, // gasLimit, TODO: runs out of gas with just res.Gas, ex: estimated was 21572 but used 24814
		nil,          // gasPrice
		suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx), // gasFeeCap
		big.NewInt(1), // gasTipCap
		transferData,
		&ethtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = hex.EncodeToString(signerKey.PubKey().Address())
	err = ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), etherminttests.NewSigner(signerKey))
	if err != nil {
		return nil, err
	}

	rsp, err := suite.App.EvmKeeper.EthereumTx(ctx, ercTransferTx)
	if err != nil {
		return nil, err
	}
	// Do not check vm error here since we want to check for errors later

	return rsp, nil
}

func (suite *Suite) MintFeeCollector(coins sdk.Coins) {
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

// GetEvents returns emitted events on the sdk context
func (suite *Suite) GetEvents() sdk.Events {
	return suite.Ctx.EventManager().Events()
}

// EventsContains asserts that the expected event is in the provided events
func (suite *Suite) EventsContains(events sdk.Events, expectedEvent sdk.Event) {
	foundMatch := false
	for _, event := range events {
		if event.Type == expectedEvent.Type {
			if reflect.DeepEqual(attrsToMap(expectedEvent.Attributes), attrsToMap(event.Attributes)) {
				foundMatch = true
			}
		}
	}

	suite.Truef(foundMatch, "event of type %s not found or did not match", expectedEvent.Type)
}

// EventsDoNotContain asserts that the event is **not** is in the provided events
func (suite *Suite) EventsDoNotContain(events sdk.Events, eventType string) {
	foundMatch := false
	for _, event := range events {
		if event.Type == eventType {
			foundMatch = true
		}
	}

	suite.Falsef(foundMatch, "event of type %s should not be found, but was found", eventType)
}

func attrsToMap(attrs []abci.EventAttribute) []sdk.Attribute {
	out := []sdk.Attribute{}

	for _, attr := range attrs {
		out = append(out, sdk.NewAttribute(string(attr.Key), string(attr.Value)))
	}

	return out
}

// MustNewExternalEVMAddressFromString returns a new ExternalEVMAddress from a
// hex string. This will panic if the input hex string is invalid.
func MustNewExternalEVMAddressFromString(addrStr string) types.ExternalEVMAddress {
	addr, err := types.NewExternalEVMAddressFromString(addrStr)
	if err != nil {
		panic(err)
	}

	return addr
}

// MustNewInternalEVMAddressFromString returns a new InternalEVMAddress from a
// hex string. This will panic if the input hex string is invalid.
func MustNewInternalEVMAddressFromString(addrStr string) types.InternalEVMAddress {
	addr, err := types.NewInternalEVMAddressFromString(addrStr)
	if err != nil {
		panic(err)
	}

	return addr
}
