// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package relayer

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"relayer_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"toAddr\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lockSequence\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"toAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockSequence\",\"type\":\"uint256\"}],\"name\":\"Unlock\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"toAddr\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockSequence\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610e77380380610e77833981810160405281019061003291906100ee565b60006001600081905550806001819055505080600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505061011b565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100bb82610090565b9050919050565b6100cb816100b0565b81146100d657600080fd5b50565b6000815190506100e8816100c2565b92915050565b6000602082840312156101045761010361008b565b5b6000610112848285016100d9565b91505092915050565b610d4d8061012a6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634202e907146100465780638406c07914610062578063a80de0e814610080575b600080fd5b610060600480360381019061005b919061077f565b61009c565b005b61006a61021a565b60405161007791906107f5565b60405180910390f35b61009a60048036038101906100959190610846565b610244565b005b600260005414156100e2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100d9906108f6565b60405180910390fd5b6002600081905550600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461017a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161017190610962565b60405180910390fd5b6101a583838673ffffffffffffffffffffffffffffffffffffffff166103439092919063ffffffff16565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fc1640cf787ea538af4a68163da63c6da8b4577194278080f5b040a75df0038998484604051610204929190610991565b60405180910390a3600160008190555050505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6002600054141561028a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610281906108f6565b60405180910390fd5b600260008190555061029a6103c9565b6102c73330838673ffffffffffffffffffffffffffffffffffffffff166103d6909392919063ffffffff16565b813373ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f04751865a034480d02848e192ee5d6a2b5fe77e3440602cc36a28ff559978c6e8461032061045f565b60405161032e929190610991565b60405180910390a46001600081905550505050565b6103c48363a9059cbb60e01b84846040516024016103629291906109ba565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050610469565b505050565b6001805401600181905550565b610459846323b872dd60e01b8585856040516024016103f7939291906109e3565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050610469565b50505050565b6000600154905090565b60006104cb826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166105309092919063ffffffff16565b905060008151111561052b57808060200190518101906104eb9190610a52565b61052a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161052190610af1565b60405180910390fd5b5b505050565b606061053f8484600085610548565b90509392505050565b60608247101561058d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161058490610b83565b60405180910390fd5b6105968561065c565b6105d5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105cc90610bef565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516105fe9190610c89565b60006040518083038185875af1925050503d806000811461063b576040519150601f19603f3d011682016040523d82523d6000602084013e610640565b606091505b509150915061065082828661067f565b92505050949350505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b6060831561068f578290506106df565b6000835111156106a25782518084602001fd5b816040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106d69190610cf5565b60405180910390fd5b9392505050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610716826106eb565b9050919050565b6107268161070b565b811461073157600080fd5b50565b6000813590506107438161071d565b92915050565b6000819050919050565b61075c81610749565b811461076757600080fd5b50565b60008135905061077981610753565b92915050565b60008060008060808587031215610799576107986106e6565b5b60006107a787828801610734565b94505060206107b887828801610734565b93505060406107c98782880161076a565b92505060606107da8782880161076a565b91505092959194509250565b6107ef8161070b565b82525050565b600060208201905061080a60008301846107e6565b92915050565b6000819050919050565b61082381610810565b811461082e57600080fd5b50565b6000813590506108408161081a565b92915050565b60008060006060848603121561085f5761085e6106e6565b5b600061086d86828701610734565b935050602061087e86828701610831565b925050604061088f8682870161076a565b9150509250925092565b600082825260208201905092915050565b7f5265656e7472616e637947756172643a207265656e7472616e742063616c6c00600082015250565b60006108e0601f83610899565b91506108eb826108aa565b602082019050919050565b6000602082019050818103600083015261090f816108d3565b9050919050565b7f4272696467653a20756e74727573746564206164647265737300000000000000600082015250565b600061094c601983610899565b915061095782610916565b602082019050919050565b6000602082019050818103600083015261097b8161093f565b9050919050565b61098b81610749565b82525050565b60006040820190506109a66000830185610982565b6109b36020830184610982565b9392505050565b60006040820190506109cf60008301856107e6565b6109dc6020830184610982565b9392505050565b60006060820190506109f860008301866107e6565b610a0560208301856107e6565b610a126040830184610982565b949350505050565b60008115159050919050565b610a2f81610a1a565b8114610a3a57600080fd5b50565b600081519050610a4c81610a26565b92915050565b600060208284031215610a6857610a676106e6565b5b6000610a7684828501610a3d565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b6000610adb602a83610899565b9150610ae682610a7f565b604082019050919050565b60006020820190508181036000830152610b0a81610ace565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b6000610b6d602683610899565b9150610b7882610b11565b604082019050919050565b60006020820190508181036000830152610b9c81610b60565b9050919050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b6000610bd9601d83610899565b9150610be482610ba3565b602082019050919050565b60006020820190508181036000830152610c0881610bcc565b9050919050565b600081519050919050565b600081905092915050565b60005b83811015610c43578082015181840152602081019050610c28565b83811115610c52576000848401525b50505050565b6000610c6382610c0f565b610c6d8185610c1a565b9350610c7d818560208601610c25565b80840191505092915050565b6000610c958284610c58565b915081905092915050565b600081519050919050565b6000601f19601f8301169050919050565b6000610cc782610ca0565b610cd18185610899565b9350610ce1818560208601610c25565b610cea81610cab565b840191505092915050565b60006020820190508181036000830152610d0f8184610cbc565b90509291505056fea264697066735822122024f08beb2e26ad582a14c1d7cafc7f6e934288a39b26c3101797285d56b5fcb364736f6c63430008090033",
}

// BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeMetaData.ABI instead.
var BridgeABI = BridgeMetaData.ABI

// BridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeMetaData.Bin instead.
var BridgeBin = BridgeMetaData.Bin

// DeployBridge deploys a new Ethereum contract, binding an instance of Bridge to it.
func DeployBridge(auth *bind.TransactOpts, backend bind.ContractBackend, relayer_ common.Address) (common.Address, *types.Transaction, *Bridge, error) {
	parsed, err := BridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeBin), backend, relayer_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// Bridge is an auto generated Go binding around an Ethereum contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_Bridge *BridgeCaller) Relayer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "relayer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_Bridge *BridgeSession) Relayer() (common.Address, error) {
	return _Bridge.Contract.Relayer(&_Bridge.CallOpts)
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_Bridge *BridgeCallerSession) Relayer() (common.Address, error) {
	return _Bridge.Contract.Relayer(&_Bridge.CallOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xa80de0e8.
//
// Solidity: function lock(address token, bytes32 toAddr, uint256 amount) returns()
func (_Bridge *BridgeTransactor) Lock(opts *bind.TransactOpts, token common.Address, toAddr [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "lock", token, toAddr, amount)
}

// Lock is a paid mutator transaction binding the contract method 0xa80de0e8.
//
// Solidity: function lock(address token, bytes32 toAddr, uint256 amount) returns()
func (_Bridge *BridgeSession) Lock(token common.Address, toAddr [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Lock(&_Bridge.TransactOpts, token, toAddr, amount)
}

// Lock is a paid mutator transaction binding the contract method 0xa80de0e8.
//
// Solidity: function lock(address token, bytes32 toAddr, uint256 amount) returns()
func (_Bridge *BridgeTransactorSession) Lock(token common.Address, toAddr [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Lock(&_Bridge.TransactOpts, token, toAddr, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x4202e907.
//
// Solidity: function unlock(address token, address toAddr, uint256 amount, uint256 unlockSequence) returns()
func (_Bridge *BridgeTransactor) Unlock(opts *bind.TransactOpts, token common.Address, toAddr common.Address, amount *big.Int, unlockSequence *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "unlock", token, toAddr, amount, unlockSequence)
}

// Unlock is a paid mutator transaction binding the contract method 0x4202e907.
//
// Solidity: function unlock(address token, address toAddr, uint256 amount, uint256 unlockSequence) returns()
func (_Bridge *BridgeSession) Unlock(token common.Address, toAddr common.Address, amount *big.Int, unlockSequence *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Unlock(&_Bridge.TransactOpts, token, toAddr, amount, unlockSequence)
}

// Unlock is a paid mutator transaction binding the contract method 0x4202e907.
//
// Solidity: function unlock(address token, address toAddr, uint256 amount, uint256 unlockSequence) returns()
func (_Bridge *BridgeTransactorSession) Unlock(token common.Address, toAddr common.Address, amount *big.Int, unlockSequence *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Unlock(&_Bridge.TransactOpts, token, toAddr, amount, unlockSequence)
}

// BridgeLockIterator is returned from FilterLock and is used to iterate over the raw logs and unpacked data for Lock events raised by the Bridge contract.
type BridgeLockIterator struct {
	Event *BridgeLock // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeLock)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeLock)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeLock represents a Lock event raised by the Bridge contract.
type BridgeLock struct {
	Token        common.Address
	Sender       common.Address
	ToAddr       [32]byte
	Amount       *big.Int
	LockSequence *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLock is a free log retrieval operation binding the contract event 0x04751865a034480d02848e192ee5d6a2b5fe77e3440602cc36a28ff559978c6e.
//
// Solidity: event Lock(address indexed token, address indexed sender, bytes32 indexed toAddr, uint256 amount, uint256 lockSequence)
func (_Bridge *BridgeFilterer) FilterLock(opts *bind.FilterOpts, token []common.Address, sender []common.Address, toAddr [][32]byte) (*BridgeLockIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var toAddrRule []interface{}
	for _, toAddrItem := range toAddr {
		toAddrRule = append(toAddrRule, toAddrItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Lock", tokenRule, senderRule, toAddrRule)
	if err != nil {
		return nil, err
	}
	return &BridgeLockIterator{contract: _Bridge.contract, event: "Lock", logs: logs, sub: sub}, nil
}

// WatchLock is a free log subscription operation binding the contract event 0x04751865a034480d02848e192ee5d6a2b5fe77e3440602cc36a28ff559978c6e.
//
// Solidity: event Lock(address indexed token, address indexed sender, bytes32 indexed toAddr, uint256 amount, uint256 lockSequence)
func (_Bridge *BridgeFilterer) WatchLock(opts *bind.WatchOpts, sink chan<- *BridgeLock, token []common.Address, sender []common.Address, toAddr [][32]byte) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var toAddrRule []interface{}
	for _, toAddrItem := range toAddr {
		toAddrRule = append(toAddrRule, toAddrItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Lock", tokenRule, senderRule, toAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeLock)
				if err := _Bridge.contract.UnpackLog(event, "Lock", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLock is a log parse operation binding the contract event 0x04751865a034480d02848e192ee5d6a2b5fe77e3440602cc36a28ff559978c6e.
//
// Solidity: event Lock(address indexed token, address indexed sender, bytes32 indexed toAddr, uint256 amount, uint256 lockSequence)
func (_Bridge *BridgeFilterer) ParseLock(log types.Log) (*BridgeLock, error) {
	event := new(BridgeLock)
	if err := _Bridge.contract.UnpackLog(event, "Lock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeUnlockIterator is returned from FilterUnlock and is used to iterate over the raw logs and unpacked data for Unlock events raised by the Bridge contract.
type BridgeUnlockIterator struct {
	Event *BridgeUnlock // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeUnlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeUnlock)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeUnlock)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeUnlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeUnlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeUnlock represents a Unlock event raised by the Bridge contract.
type BridgeUnlock struct {
	Token          common.Address
	ToAddr         common.Address
	Amount         *big.Int
	UnlockSequence *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUnlock is a free log retrieval operation binding the contract event 0xc1640cf787ea538af4a68163da63c6da8b4577194278080f5b040a75df003899.
//
// Solidity: event Unlock(address indexed token, address indexed toAddr, uint256 amount, uint256 unlockSequence)
func (_Bridge *BridgeFilterer) FilterUnlock(opts *bind.FilterOpts, token []common.Address, toAddr []common.Address) (*BridgeUnlockIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toAddrRule []interface{}
	for _, toAddrItem := range toAddr {
		toAddrRule = append(toAddrRule, toAddrItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Unlock", tokenRule, toAddrRule)
	if err != nil {
		return nil, err
	}
	return &BridgeUnlockIterator{contract: _Bridge.contract, event: "Unlock", logs: logs, sub: sub}, nil
}

// WatchUnlock is a free log subscription operation binding the contract event 0xc1640cf787ea538af4a68163da63c6da8b4577194278080f5b040a75df003899.
//
// Solidity: event Unlock(address indexed token, address indexed toAddr, uint256 amount, uint256 unlockSequence)
func (_Bridge *BridgeFilterer) WatchUnlock(opts *bind.WatchOpts, sink chan<- *BridgeUnlock, token []common.Address, toAddr []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toAddrRule []interface{}
	for _, toAddrItem := range toAddr {
		toAddrRule = append(toAddrRule, toAddrItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Unlock", tokenRule, toAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeUnlock)
				if err := _Bridge.contract.UnpackLog(event, "Unlock", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnlock is a log parse operation binding the contract event 0xc1640cf787ea538af4a68163da63c6da8b4577194278080f5b040a75df003899.
//
// Solidity: event Unlock(address indexed token, address indexed toAddr, uint256 amount, uint256 unlockSequence)
func (_Bridge *BridgeFilterer) ParseUnlock(log types.Log) (*BridgeUnlock, error) {
	event := new(BridgeUnlock)
	if err := _Bridge.contract.UnpackLog(event, "Unlock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
