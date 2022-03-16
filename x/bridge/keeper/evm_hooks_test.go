package keeper_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

type EVMHooksTestSuite struct {
	testutil.Suite
}

func TestEVMHooksTestSuite(t *testing.T) {
	suite.Run(t, new(EVMHooksTestSuite))
}

func (suite *EVMHooksTestSuite) deployERC20() types.InternalEVMAddress {
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

func (suite *EVMHooksTestSuite) TestERC20WithdrawUnpack() {
	erc20Abi := contract.ERC20MintableBurnableContract.ABI

	contractAddr := suite.deployERC20()

	key1Addr := common.BytesToAddress(suite.Key1.PubKey().Address())
	amount := big.NewInt(1234)
	err := suite.App.BridgeKeeper.MintERC20(suite.Ctx, contractAddr, key1Addr, amount)
	suite.Require().NoError(err)

	withdrawAmount := big.NewInt(10)
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())

	// method is lowercase but event is upper
	data, err := erc20Abi.Pack(
		"withdraw",
		withdrawToAddr,
		withdrawAmount,
	)
	suite.Require().NoError(err)

	// Send TX
	res := suite.SendTx(contractAddr, key1Addr, suite.Key1, data)
	suite.Require().False(res.Failed(), "evm tx should not fail")

	containsWithdrawEvent := false

	for _, log := range res.Logs {
		eventID := log.Topics[0]

		event, err := erc20Abi.EventByID(common.HexToHash(eventID))
		if err != nil {
			// invalid event for ERC20
			continue
		}

		if event.Name != types.ContractEventTypeWithdraw {
			continue
		}

		containsWithdrawEvent = true

		withdrawEvent, err := erc20Abi.Unpack(types.ContractEventTypeWithdraw, log.Data)
		suite.Require().NoError(err)

		suite.Require().Len(withdrawEvent, 1, "withdraw event data should only have 1 item for amount")

		loggedAmount, ok := withdrawEvent[0].(*big.Int)
		suite.Require().True(ok, "withdraw event data should be *big.Int")
		suite.Require().Equal(withdrawAmount, loggedAmount)

		// 3 topics:
		// - Keccak-256 hash of Withdraw(address,address,uint256)
		// - address indexed sender
		// - address indexed toAddr
		suite.Require().Len(log.Topics, 3, "withdraw event should have 3 topics")

		// log.Topics is padded to 32 bytes, addresses are 20 bytes.
		// common.HexToAddress handles this, crops from left.
		senderAddr := common.HexToAddress(log.Topics[1])
		suite.Require().Equal(key1Addr, senderAddr)

		logToAddr := common.HexToAddress(log.Topics[2])
		suite.Require().Equal(withdrawToAddr, logToAddr)
	}

	suite.Require().True(containsWithdrawEvent, "tx should contain Withdraw event")
}
