package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"

	"github.com/tharsis/ethermint/server/config"
	etherminttypes "github.com/tharsis/ethermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	bridgeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmds := []*cobra.Command{
		getCmdMsgBridgeEthereumToKava(),
		getCmdMsgConvertCoinToERC20(),
		getCmdBridgeKavaToEthereum(),
	}

	for _, cmd := range cmds {
		flags.AddTxFlagsToCmd(cmd)
	}

	bridgeTxCmd.AddCommand(cmds...)

	return bridgeTxCmd
}

func getCmdMsgBridgeEthereumToKava() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge-eth-to-kava [token] [receiver] [amount] [sequence]",
		Short: "mints ERC20 tokens locked on Ethereum to Kava EVM co-chain",
		Example: fmt.Sprintf(
			`%s tx %s bridge-eth-to-kava 0xc778417e063141139fce010982780140aa0cd5ab 0x6B1088f788b412Ad1280F95240d56B886A64bc05 1000000000000000 1 --from <key>`,
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			token := args[0]
			receiver := args[1]

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("amount '%s' is an invalid int", args[2])
			}

			sequence, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("amount '%s' is an invalid int", args[2])
			}

			signer := clientCtx.GetFromAddress()
			msg := types.NewMsgBridgeEthereumToKava(signer.String(), token, amount, receiver, sequence)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}

func getCmdMsgConvertCoinToERC20() *cobra.Command {
	return &cobra.Command{
		Use:   "convert-coin-to-erc20 [receiver] [coin]",
		Short: "converts sdk.Coin to erc20 tokens on Kava eth co-chain",
		Example: fmt.Sprintf(
			`%s tx %s convert-coin-to-erc20 0x6B1088f788b412Ad1280F95240d56B886A64bc05 100000000weth --from <key>`,
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			receiver := args[0]
			if !common.IsHexAddress(receiver) {
				return fmt.Errorf("receiver '%s' is an invalid hex address", args[0])
			}

			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			msg := types.NewMsgConvertCoinToERC20(signer.String(), receiver, coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}

func getCmdBridgeKavaToEthereum() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge-kava-to-eth [Kava ERC20 address] [Ethereum receiver address] [amount]",
		Short: "burns ERC20 tokens on Kava EVM co-chain and unlocks on Ethereum",
		Example: fmt.Sprintf(
			`%s tx %s bridge-kava-to-eth 0x8223259205A3E31C54469fCbfc9F7Cf83D515ff6 0x21E360e198Cde35740e88572B59f2CAdE421E6b1 1000000000000000 --from <key>`,
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("contract address '%s' is an invalid hex address", args[0])
			}
			contractAddr := common.HexToAddress(args[0])

			if !common.IsHexAddress(args[1]) {
				return fmt.Errorf("receiver '%s' is an invalid hex address", args[1])
			}
			receiver := common.HexToAddress(args[1])

			amount, ok := new(big.Int).SetString(args[2], 10)
			if !ok {
				return fmt.Errorf("amount '%s' is invalid", args[2])
			}

			data, err := createContractCallData(
				contract.ERC20MintableBurnableContract.ABI,
				"withdraw",
				receiver,
				amount,
			)
			if err != nil {
				return err
			}

			ethTx, err := CreateEthCallContractTx(
				clientCtx,
				&contractAddr,
				data,
			)
			if err != nil {
				return err
			}

			txBytes, err := clientCtx.TxConfig.TxEncoder()(ethTx)
			if err != nil {
				return err
			}

			res, err := clientCtx.BroadcastTx(txBytes)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func createContractCallData(abi abi.ABI, method string, args ...interface{}) ([]byte, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction data: %w", err)
	}

	return data, nil
}

func CreateEthCallContractTx(
	ctx client.Context,
	contractAddr *common.Address,
	data []byte,
) (signing.Tx, error) {
	evmQueryClient := evmtypes.NewQueryClient(ctx)
	feemarketQueryClient := feemarkettypes.NewQueryClient(ctx)

	chainID, err := etherminttypes.ParseChainID(ctx.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chain ID: %w", err)
	}

	// Estimate Gas
	from := common.BytesToAddress(ctx.FromAddress.Bytes())
	transactionArgs := evmtypes.TransactionArgs{
		From: &from,
		To:   contractAddr,
		Data: (*hexutil.Bytes)(&data),
	}

	args, err := json.Marshal(transactionArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction args for gas estimate: %w", err)
	}

	res, err := evmQueryClient.EstimateGas(context.Background(), &evmtypes.EthCallRequest{
		Args:   args,
		GasCap: config.DefaultGasCap,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas from EVM: %w", err)
	}

	// Fetch base fee
	basefeeRes, err := feemarketQueryClient.BaseFee(
		context.Background(),
		&feemarkettypes.QueryBaseFeeRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch basefee from feemarket: %w", err)
	}

	// Fetch account nonce, ignore error to use use 0 nonce if first tx
	_, accSeq, _ := ctx.AccountRetriever.GetAccountNumberSequence(ctx, ctx.FromAddress)

	// Create MsgEthereumTx
	ethTx := evmtypes.NewTx(
		chainID,
		accSeq,                      // nonce
		contractAddr,                // to
		nil,                         // amount
		res.Gas,                     // gasLimit
		nil,                         // gasPrice
		basefeeRes.BaseFee.BigInt(), // gasFeeCap
		big.NewInt(1),               // gasTipCap
		data,                        // input
		&ethtypes.AccessList{},
	)

	// Must set from address before signing
	ethTx.From = from.String()

	// Sign Ethereum TX (not the cosmos Msg)
	signer := ethtypes.LatestSignerForChainID(chainID)

	// key := ctx.Keyring.KeyByAddress(ctx.FromAddress)

	if err := ethTx.Sign(signer, ctx.Keyring); err != nil {
		return nil, err
	}

	return ethTx.BuildTx(ctx.TxConfig.NewTxBuilder(), "ukava")
}
