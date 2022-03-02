package keeper_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	etherminttypes "github.com/tharsis/ethermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/kava-labs/kava-bridge/app"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type ERC20TestSuite struct {
	suite.Suite

	ctx            sdk.Context
	app            app.TestApp
	address        common.Address
	consAddress    sdk.ConsAddress
	relayerAddress sdk.AccAddress
}

func (suite *ERC20TestSuite) SetupTest() {
	suite.app = app.NewTestApp()

	// consensus key
	consPriv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.consAddress = sdk.ConsAddress(consPriv.PubKey().Address())

	// relayer key
	relayerPriv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.relayerAddress = sdk.AccAddress(relayerPriv.PubKey().Address())

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

	genesisState := app.NewDefaultGenesisState()
	genesisState[types.ModuleName] = suite.app.AppCodec().MustMarshalJSON(&bridgeGs)
	genesisState[evmtypes.ModuleName] = suite.app.AppCodec().MustMarshalJSON(evmGs)

	// Initialize the chain
	stateBytes, err := json.Marshal(genesisState)
	suite.Require().NoError(err)
	suite.app.InitChain(
		abci.RequestInitChain{
			Time:          time.Now().UTC(),
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
			ChainId:       "kavatest_8888-1",
			// Set consensus params, which is needed by x/feemarket
			ConsensusParams: &abci.ConsensusParams{
				Block: &abci.BlockParams{
					MaxBytes: 200000,
					MaxGas:   20000000,
				},
			},
		},
	)

	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		Height:          1,
		ChainID:         "kavatest_8888-1",
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

func (suite *ERC20TestSuite) TestDeployERC20() {
	token := types.NewEnabledERC20Token(
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"Wrapped ETH",
		"WETH",
		18,
	)

	addr, err := suite.app.BridgeKeeper.DeployMintableERC20Contract(suite.ctx, token)
	suite.Require().NoError(err)
	suite.Require().Greater(len(addr), 0)
}
