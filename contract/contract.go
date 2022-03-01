package contract

import (
	_ "embed"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

var (
	//go:embed ethermint_json/ERC20MintableBurnable.json
	ERC20MintableBurnableJSON []byte // nolint: golint

	// ERC20MintableBurnableContract is the compiled erc20 contract
	ERC20MintableBurnableContract evmtypes.CompiledContract

	// ERC20MintableBurnableAddress is the erc20 module address
	ERC20MintableBurnableAddress common.Address
)

func init() {
	ERC20MintableBurnableAddress = types.ModuleEVMAddress

	err := json.Unmarshal(ERC20MintableBurnableJSON, &ERC20MintableBurnableContract)
	if err != nil {
		panic(err)
	}

	if len(ERC20MintableBurnableContract.Bin) == 0 {
		panic("load contract failed")
	}
}
