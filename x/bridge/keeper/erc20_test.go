package keeper_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/common"

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

func (suite *ERC20TestSuite) deployERC20() types.InternalEVMAddress {
	// We can assume token is valid as it is from params and should be validated
	token := types.NewEnabledERC20Token(
		testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		"Wrapped ETH",
		"WETH",
		18,
	)

	contractAddr, err := suite.App.BridgeKeeper.DeployMintableERC20Contract(suite.Ctx, token)
	suite.Require().NoError(err)
	suite.Require().Greater(len(contractAddr.Address), 0)

	return contractAddr
}

func (suite *ERC20TestSuite) TestDeployERC20() {
	suite.deployERC20()
}

func (suite *ERC20TestSuite) TestERC20Query() {
	contractAddr := suite.deployERC20()

	// Query ERC20.decimals()
	addr := common.BytesToAddress(suite.Key1.PubKey().Address())
	res, err := suite.QueryContract(
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
	_, err := suite.QueryContract(
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
	res, err := suite.QueryContract(
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

func (suite *ERC20TestSuite) TestERC20_NotEnabled() {
	// WETH but last char changed
	extAddr := testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc4")

	_, err := suite.App.BridgeKeeper.GetOrDeployInternalERC20(suite.Ctx, extAddr)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, types.ErrERC20NotEnabled)
}

func (suite *ERC20TestSuite) TestERC20SaveDeploy() {
	extAddr := testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")

	_, found := suite.App.BridgeKeeper.GetBridgePairFromExternal(suite.Ctx, extAddr)
	suite.Require().False(found, "internal ERC20 address should not be set before first bridge")

	firstInternalAddr, err := suite.App.BridgeKeeper.GetOrDeployInternalERC20(suite.Ctx, extAddr)
	suite.Require().NoError(err)

	// Fetch from store
	savedPair, found := suite.App.BridgeKeeper.GetBridgePairFromExternal(suite.Ctx, extAddr)
	suite.Require().True(found, "internal ERC20 address should be saved after first bridge")
	suite.Require().Equal(firstInternalAddr, savedPair.GetInternalAddress(), "deployed address should match saved internal ERC20 address")

	// Fetch addr again to make sure we get the same one and another ERC20 isn't deployed
	secondInternal, err := suite.App.BridgeKeeper.GetOrDeployInternalERC20(suite.Ctx, extAddr)
	suite.Require().NoError(err)

	suite.Require().Equal(firstInternalAddr, secondInternal, "second call should return the saved internal ERC20 address")
}
