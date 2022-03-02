package keeper_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	tmtime "github.com/tendermint/tendermint/types/time"
	"github.com/tendermint/tendermint/version"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	etherminttypes "github.com/tharsis/ethermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/kava-labs/kava-bridge/app"
	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type ERC20TestSuite struct {
	suite.Suite

	App          app.TestApp
	bridgeKeeper keeper.Keeper
	Ctx          sdk.Context
}

func (suite *ERC20TestSuite) SetupTest() {
	suite.App = app.NewTestApp()

	// consensus key
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	consAddress := sdk.ConsAddress(priv.PubKey().Address())

	suite.bridgeKeeper = suite.App.GetBridgeKeeper()
	// First is validator, second is bridge relayer
	keys, addrs := app.GeneratePrivKeyAddressPairs(3)

	suite.Ctx = suite.App.NewContext(true, tmproto.Header{
		Height:          1,
		ChainID:         "kavatest_1-1",
		Time:            tmtime.Now(),
		ProposerAddress: consAddress.Bytes(),
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
		BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(addrs[0].Bytes()), nil, 0, 0),
		CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
	}

	suite.App.AccountKeeper.SetAccount(suite.Ctx, acc)

	valAddr := sdk.ValAddress(addrs[0].Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, keys[0].PubKey(), stakingtypes.Description{})
	suite.Require().NoError(err)
	err = suite.App.StakingKeeper.SetValidatorByConsAddr(suite.Ctx, validator)
	suite.Require().NoError(err)
	suite.App.StakingKeeper.SetValidator(suite.Ctx, validator)

	bridgeGs := types.NewGenesisState(types.NewParams(
		types.EnabledERC20Tokens{
			types.NewEnabledERC20Token(
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				"Wrapped ETH",
				"WETH",
				18,
			),
		},
		addrs[1],
	))

	evmGs := evmtypes.NewGenesisState(
		evmtypes.NewParams("ukava", true, true, evmtypes.DefaultChainConfig()),
		nil,
	)

	evmGsBytes := app.GenesisState{evmtypes.ModuleName: suite.App.AppCodec().MustMarshalJSON(evmGs)}

	moduleGs := suite.App.AppCodec().MustMarshalJSON(&bridgeGs)
	gs := app.GenesisState{types.ModuleName: moduleGs}

	genesisStates := []app.GenesisState{evmGsBytes, gs}

	genesisState := app.NewDefaultGenesisState()
	for _, state := range genesisStates {
		for k, v := range state {
			genesisState[k] = v
		}
	}

	// Initialize the chain
	stateBytes, err := json.Marshal(genesisState)
	if err != nil {
		panic(err)
	}
	suite.App.InitChain(
		abci.RequestInitChain{
			Time:          time.Now().UTC(),
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
			ChainId:       "kavatest_1-1",
			// Set consensus params, which is needed by x/feemarket
			ConsensusParams: &abci.ConsensusParams{
				Block: &abci.BlockParams{
					MaxBytes: 200000,
					MaxGas:   20000000,
				},
			},
		},
	)

	suite.Commit()

	fmt.Printf("consAddr: %v\n", consAddress)
	fmt.Printf("consAddr from block: %v\n", sdk.ConsAddress(suite.Ctx.BlockHeader().ProposerAddress))
	validator, found := suite.App.StakingKeeper.GetValidatorByConsAddr(suite.Ctx, consAddress)
	fmt.Printf("GetValidatorByConsAddr (found: %v): %v\n", found, validator)
}

func TestERC20TestSuite(t *testing.T) {
	suite.Run(t, new(ERC20TestSuite))
}

func (suite *ERC20TestSuite) Commit() {
	_ = suite.App.Commit()
	header := suite.Ctx.BlockHeader()
	header.Height += 1
	suite.App.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.Ctx = suite.App.NewContext(false, header)
}

func (suite *ERC20TestSuite) TestDeployERC20() {
	token := types.NewEnabledERC20Token(
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"Wrapped ETH",
		"WETH",
		18,
	)

	addr, err := suite.bridgeKeeper.DeployMintableERC20Contract(suite.Ctx, token)
	suite.Require().NoError(err)
	suite.Require().Greater(0, len(addr))
}
