//go:build integration
// +build integration

package testing

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/kava-labs/kava-bridge/app"
	"github.com/kava-labs/kava-bridge/relayer"

	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	bridgeBinName  = "kava-bridged"
	relayerBinName = "kava-relayer"

	testUserMnemonic    = "news tornado sponsor drastic dolphin awful plastic select true lizard width idle ability pigeon runway lift oppose isolate maple aspect safe jungle author hole"
	testRelayerMnemonic = "never reject sniff east arctic funny twin feed upper series stay shoot vivid adapt defense economy pledge fetch invite approve ceiling admit gloom exit"
)

func buildBin(binName string) (func(), error) {
	build := exec.Command("go", "build", "-o", binName, fmt.Sprintf("../cmd/%s", binName))
	if err := build.Run(); err != nil {
		return nil, fmt.Errorf("Failed to build %s: %s", binName, err)
	}

	return func() { os.Remove(binName) }, nil
}

func TestMain(m *testing.M) {
	// build the kava test chain w/ birdge module
	cleanupBridge, err := buildBin(bridgeBinName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// build the kava relayer
	cleanupRelayer, err := buildBin(relayerBinName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// run test suite
	r := m.Run()

	// clean up binaries
	cleanupBridge()
	cleanupRelayer()

	// exit
	os.Exit(r)
}

func TestEthToKavaRelaying(t *testing.T) {
	//
	// Initialize Ethereum Client
	//
	client, err := ethclient.Dial("http://localhost:8555")
	require.NoError(t, err)
	chainID, err := client.ChainID(context.Background())
	require.NoError(t, err)

	//
	// Initialize Kava Connection
	//
	app.SetSDKConfig()
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	//
	// Initialize EVM Client
	//
	evmclient, err := ethclient.Dial("http://localhost:8545")
	require.NoError(t, err)
	//evmchainID, err := evmclient.ChainID(context.Background())
	//require.NoError(t, err)

	//
	// Start relayer
	//
	cmd := exec.Command(relayerBinName,
		"start",
		"--eth-rpc", "http://localhost:8555",
		"--eth-bridge-address", "0xb588617416D0B0A3C29618bf8Fb6aC0cAd4Ede7f",
		"--kava-grpc", "http://localhost:9090",
		// NOTE: this will be replaced with TSS
		"--relayer-mnemonic", testRelayerMnemonic,
	)
	err = cmd.Start()
	require.NoError(t, err)
	defer func() {
		err := cmd.Process.Kill()
		require.NoError(t, err)
	}()

	//
	// Get User Private Key From Mnemonic
	//
	wallet, err := hdwallet.NewFromMnemonic(testUserMnemonic)
	require.NoError(t, err)
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	user, err := wallet.Derive(path, false)
	require.NoError(t, err)
	privateKey, err := wallet.PrivateKey(user)
	require.NoError(t, err)

	//
	// Create instance to communicate with weth contract
	//
	wethTokenAddress := common.HexToAddress("0x6098c27D41ec6dc280c2200A737D443b0AaA2E8F")
	tokenInstance, err := relayer.NewERC20(wethTokenAddress, client)
	require.NoError(t, err)

	//
	// Create instance to communicate with bridge contract
	//
	bridgeAddress := common.HexToAddress("0xb588617416D0B0A3C29618bf8Fb6aC0cAd4Ede7f")
	bridgeInstance, err := relayer.NewBridge(bridgeAddress, client)
	require.NoError(t, err)

	//
	// Initiate a transfer
	//
	nonce, err := client.PendingNonceAt(context.Background(), user.Address)
	require.NoError(t, err)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)

	transferAmount := big.NewInt(1e15)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice

	// approve token transfer
	tx, err := tokenInstance.Approve(auth, bridgeAddress, transferAmount)
	require.NoError(t, err)

	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)
	auth.Nonce = big.NewInt(int64(nonce + 1))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice

	// lock tokens in bridge
	tx, err = bridgeInstance.Lock(auth, wethTokenAddress, user.Address, big.NewInt(1e15))
	require.NoError(t, err)

	// wait for transaction receipt
	lockReceipt := (*ethtypes.Receipt)(nil)
	for lockReceipt == nil {
		lockReceipt, err = client.TransactionReceipt(context.Background(), tx.Hash())

		if err != nil && !errors.Is(err, ethereum.NotFound) {
			t.Log(err)
			t.FailNow()
		}

		time.Sleep(500 * time.Millisecond)
	}

	lockLog, err := bridgeInstance.ParseLock(*lockReceipt.Logs[1])
	require.NoError(t, err)

	//
	// Search Kava Chain For Transaction With Lock Sequence
	//
	txClient := txtypes.NewServiceClient(conn)

	txRequest := txtypes.GetTxsEventRequest{
		Events: []string{fmt.Sprintf("bridge_ethereum_to_kava.sequence='%s'", lockLog.LockSequence)},
	}

	timeout := time.Now().Add(60 * time.Second)
	var txsResponse *txtypes.GetTxsEventResponse

	for {
		if time.Now().After(timeout) {
			t.Fatal("timed out waiting for successful kava transaction")
		}

		txsResponse, err = txClient.GetTxsEvent(context.Background(), &txRequest)
		require.NoError(t, err)

		if len(txsResponse.TxResponses) > 0 && txsResponse.TxResponses[0].Code == 0 {
			break
		}
	}

	//
	// Create instance to communicate with weth contract on kava evm
	//
	// TODO: we should query this address from module, for now we rely on weth deployment being the
	// first evm tx signed by the bridge module
	evmWethTokenAddress := common.HexToAddress("0x404f9466d758ea33ea84cebe9e444b06533b369e")
	evmTokenInstance, err := relayer.NewERC20(evmWethTokenAddress, evmclient)
	require.NoError(t, err)

	//
	// Assert balance on kava evm matches value locked
	//
	balance, err := evmTokenInstance.BalanceOf(nil, user.Address)
	require.Equal(t, int64(1e15), balance.Int64())
}
