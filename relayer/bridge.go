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
	_ = abi.ConvertType
)

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"relayer_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"toKavaAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lockSequence\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"toAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockSequence\",\"type\":\"uint256\"}],\"name\":\"Unlock\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toKavaAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockSequence\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516107a23803806107a283398101604081905261002f9161005d565b600160008181559055600280546001600160a01b0319166001600160a01b039290921691909117905561008d565b60006020828403121561006f57600080fd5b81516001600160a01b038116811461008657600080fd5b9392505050565b6107068061009c6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634202e907146100465780637750c9f01461005b5780638406c0791461006e575b600080fd5b6100596100543660046105b5565b61008d565b005b6100596100693660046105f7565b6101b9565b600254604080516001600160a01b039092168252519081900360200190f35b600260005414156100e55760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c0060448201526064015b60405180910390fd5b60026000819055546001600160a01b031633146101445760405162461bcd60e51b815260206004820152601960248201527f4272696467653a20756e7472757374656420616464726573730000000000000060448201526064016100dc565b6101586001600160a01b038516848461029e565b826001600160a01b0316846001600160a01b03167fc1640cf787ea538af4a68163da63c6da8b4577194278080f5b040a75df00389984846040516101a6929190918252602082015260400190565b60405180910390a3505060016000555050565b6002600054141561020c5760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c0060448201526064016100dc565b600260005561021e6001805481019055565b6102336001600160a01b038416333084610306565b816001600160a01b0316336001600160a01b0316846001600160a01b03167f749e347e95185169edffb86e003abcbf08f3510641663b00600acf707f098f478461027c60015490565b6040805192835260208301919091520160405180910390a45050600160005550565b6040516001600160a01b03831660248201526044810182905261030190849063a9059cbb60e01b906064015b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152610344565b505050565b6040516001600160a01b038085166024830152831660448201526064810182905261033e9085906323b872dd60e01b906084016102ca565b50505050565b6000610399826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166104169092919063ffffffff16565b80519091501561030157808060200190518101906103b79190610633565b6103015760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084016100dc565b6060610425848460008561042f565b90505b9392505050565b6060824710156104905760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b60648201526084016100dc565b6001600160a01b0385163b6104e75760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016100dc565b600080866001600160a01b031685876040516105039190610681565b60006040518083038185875af1925050503d8060008114610540576040519150601f19603f3d011682016040523d82523d6000602084013e610545565b606091505b5091509150610555828286610560565b979650505050505050565b6060831561056f575081610428565b82511561057f5782518084602001fd5b8160405162461bcd60e51b81526004016100dc919061069d565b80356001600160a01b03811681146105b057600080fd5b919050565b600080600080608085870312156105cb57600080fd5b6105d485610599565b93506105e260208601610599565b93969395505050506040820135916060013590565b60008060006060848603121561060c57600080fd5b61061584610599565b925061062360208501610599565b9150604084013590509250925092565b60006020828403121561064557600080fd5b8151801515811461042857600080fd5b60005b83811015610670578181015183820152602001610658565b8381111561033e5750506000910152565b60008251610693818460208701610655565b9190910192915050565b60208152600082518060208401526106bc816040850160208701610655565b601f01601f1916919091016040019291505056fea2646970667358221220b5376e1fd0eccc27196aa613da565619399b75112f6b6ea195dc52e3685e7c2864736f6c63430008090033",
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
	parsed, err := BridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// Lock is a paid mutator transaction binding the contract method 0x7750c9f0.
//
// Solidity: function lock(address token, address toKavaAddr, uint256 amount) returns()
func (_Bridge *BridgeTransactor) Lock(opts *bind.TransactOpts, token common.Address, toKavaAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "lock", token, toKavaAddr, amount)
}

// Lock is a paid mutator transaction binding the contract method 0x7750c9f0.
//
// Solidity: function lock(address token, address toKavaAddr, uint256 amount) returns()
func (_Bridge *BridgeSession) Lock(token common.Address, toKavaAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Lock(&_Bridge.TransactOpts, token, toKavaAddr, amount)
}

// Lock is a paid mutator transaction binding the contract method 0x7750c9f0.
//
// Solidity: function lock(address token, address toKavaAddr, uint256 amount) returns()
func (_Bridge *BridgeTransactorSession) Lock(token common.Address, toKavaAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Lock(&_Bridge.TransactOpts, token, toKavaAddr, amount)
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
	ToKavaAddr   common.Address
	Amount       *big.Int
	LockSequence *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLock is a free log retrieval operation binding the contract event 0x749e347e95185169edffb86e003abcbf08f3510641663b00600acf707f098f47.
//
// Solidity: event Lock(address indexed token, address indexed sender, address indexed toKavaAddr, uint256 amount, uint256 lockSequence)
func (_Bridge *BridgeFilterer) FilterLock(opts *bind.FilterOpts, token []common.Address, sender []common.Address, toKavaAddr []common.Address) (*BridgeLockIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var toKavaAddrRule []interface{}
	for _, toKavaAddrItem := range toKavaAddr {
		toKavaAddrRule = append(toKavaAddrRule, toKavaAddrItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Lock", tokenRule, senderRule, toKavaAddrRule)
	if err != nil {
		return nil, err
	}
	return &BridgeLockIterator{contract: _Bridge.contract, event: "Lock", logs: logs, sub: sub}, nil
}

// WatchLock is a free log subscription operation binding the contract event 0x749e347e95185169edffb86e003abcbf08f3510641663b00600acf707f098f47.
//
// Solidity: event Lock(address indexed token, address indexed sender, address indexed toKavaAddr, uint256 amount, uint256 lockSequence)
func (_Bridge *BridgeFilterer) WatchLock(opts *bind.WatchOpts, sink chan<- *BridgeLock, token []common.Address, sender []common.Address, toKavaAddr []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var toKavaAddrRule []interface{}
	for _, toKavaAddrItem := range toKavaAddr {
		toKavaAddrRule = append(toKavaAddrRule, toKavaAddrItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Lock", tokenRule, senderRule, toKavaAddrRule)
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

// ParseLock is a log parse operation binding the contract event 0x749e347e95185169edffb86e003abcbf08f3510641663b00600acf707f098f47.
//
// Solidity: event Lock(address indexed token, address indexed sender, address indexed toKavaAddr, uint256 amount, uint256 lockSequence)
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
