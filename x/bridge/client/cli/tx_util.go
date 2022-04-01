package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// PrintTx outputs a signing.Tx in JSON format, ie. when the GenerateOnly flag
// is enabled.
func PrintTx(clientCtx client.Context, signingTx signing.Tx) error {
	json, err := clientCtx.TxConfig.TxJSONEncoder()(signingTx)
	if err != nil {
		return err
	}

	return clientCtx.PrintString(fmt.Sprintf("%s\n", json))
}

// CheckConfirm outputs the transaction to be signed and requests confirmation
// if the SkipConfirm flag is not enabled.
func CheckConfirm(clientCtx client.Context, signingTx signing.Tx) error {
	if clientCtx.SkipConfirm {
		return nil
	}

	out, err := clientCtx.TxConfig.TxJSONEncoder()(signingTx)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

	buf := bufio.NewReader(os.Stdin)
	ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

	if err != nil || !ok {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
		return err
	}

	return nil
}
