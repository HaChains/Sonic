// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package coinbase

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

// CoinbaseMetaData contains all meta data concerning the Coinbase contract.
var CoinbaseMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"LogAddress\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"logCoinBaseAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"touchAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"touchCoinbase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b5061020a8061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c8063713ee28514610043578063874395f81461004d578063d847a22b14610057575b5f5ffd5b61004b610073565b005b6100556100ac565b005b610071600480360381019061006c9190610135565b6100b7565b005b7fb123f68b8ba02b447d91a6629e121111b7dd6061ff418a60139c8bf00522a284416040516100a291906101bb565b60405180910390a1565b6100b5416100b7565b565b8073ffffffffffffffffffffffffffffffffffffffff16315f8190555050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610104826100db565b9050919050565b610114816100fa565b811461011e575f5ffd5b50565b5f8135905061012f8161010b565b92915050565b5f6020828403121561014a576101496100d7565b5b5f61015784828501610121565b91505092915050565b5f819050919050565b5f61018361017e610179846100db565b610160565b6100db565b9050919050565b5f61019482610169565b9050919050565b5f6101a58261018a565b9050919050565b6101b58161019b565b82525050565b5f6020820190506101ce5f8301846101ac565b9291505056fea26469706673582212209ae6955b9fbfcc111b9360ab4a553ab92eec2e564767d625efd66a4475e9e79564736f6c634300081c0033",
}

// CoinbaseABI is the input ABI used to generate the binding from.
// Deprecated: Use CoinbaseMetaData.ABI instead.
var CoinbaseABI = CoinbaseMetaData.ABI

// CoinbaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CoinbaseMetaData.Bin instead.
var CoinbaseBin = CoinbaseMetaData.Bin

// DeployCoinbase deploys a new Ethereum contract, binding an instance of Coinbase to it.
func DeployCoinbase(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Coinbase, error) {
	parsed, err := CoinbaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CoinbaseBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Coinbase{CoinbaseCaller: CoinbaseCaller{contract: contract}, CoinbaseTransactor: CoinbaseTransactor{contract: contract}, CoinbaseFilterer: CoinbaseFilterer{contract: contract}}, nil
}

// Coinbase is an auto generated Go binding around an Ethereum contract.
type Coinbase struct {
	CoinbaseCaller     // Read-only binding to the contract
	CoinbaseTransactor // Write-only binding to the contract
	CoinbaseFilterer   // Log filterer for contract events
}

// CoinbaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type CoinbaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoinbaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CoinbaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoinbaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CoinbaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoinbaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CoinbaseSession struct {
	Contract     *Coinbase         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CoinbaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CoinbaseCallerSession struct {
	Contract *CoinbaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// CoinbaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CoinbaseTransactorSession struct {
	Contract     *CoinbaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// CoinbaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type CoinbaseRaw struct {
	Contract *Coinbase // Generic contract binding to access the raw methods on
}

// CoinbaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CoinbaseCallerRaw struct {
	Contract *CoinbaseCaller // Generic read-only contract binding to access the raw methods on
}

// CoinbaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CoinbaseTransactorRaw struct {
	Contract *CoinbaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCoinbase creates a new instance of Coinbase, bound to a specific deployed contract.
func NewCoinbase(address common.Address, backend bind.ContractBackend) (*Coinbase, error) {
	contract, err := bindCoinbase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Coinbase{CoinbaseCaller: CoinbaseCaller{contract: contract}, CoinbaseTransactor: CoinbaseTransactor{contract: contract}, CoinbaseFilterer: CoinbaseFilterer{contract: contract}}, nil
}

// NewCoinbaseCaller creates a new read-only instance of Coinbase, bound to a specific deployed contract.
func NewCoinbaseCaller(address common.Address, caller bind.ContractCaller) (*CoinbaseCaller, error) {
	contract, err := bindCoinbase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CoinbaseCaller{contract: contract}, nil
}

// NewCoinbaseTransactor creates a new write-only instance of Coinbase, bound to a specific deployed contract.
func NewCoinbaseTransactor(address common.Address, transactor bind.ContractTransactor) (*CoinbaseTransactor, error) {
	contract, err := bindCoinbase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CoinbaseTransactor{contract: contract}, nil
}

// NewCoinbaseFilterer creates a new log filterer instance of Coinbase, bound to a specific deployed contract.
func NewCoinbaseFilterer(address common.Address, filterer bind.ContractFilterer) (*CoinbaseFilterer, error) {
	contract, err := bindCoinbase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CoinbaseFilterer{contract: contract}, nil
}

// bindCoinbase binds a generic wrapper to an already deployed contract.
func bindCoinbase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CoinbaseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Coinbase *CoinbaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Coinbase.Contract.CoinbaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Coinbase *CoinbaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coinbase.Contract.CoinbaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Coinbase *CoinbaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Coinbase.Contract.CoinbaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Coinbase *CoinbaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Coinbase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Coinbase *CoinbaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coinbase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Coinbase *CoinbaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Coinbase.Contract.contract.Transact(opts, method, params...)
}

// LogCoinBaseAddress is a paid mutator transaction binding the contract method 0x713ee285.
//
// Solidity: function logCoinBaseAddress() returns()
func (_Coinbase *CoinbaseTransactor) LogCoinBaseAddress(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coinbase.contract.Transact(opts, "logCoinBaseAddress")
}

// LogCoinBaseAddress is a paid mutator transaction binding the contract method 0x713ee285.
//
// Solidity: function logCoinBaseAddress() returns()
func (_Coinbase *CoinbaseSession) LogCoinBaseAddress() (*types.Transaction, error) {
	return _Coinbase.Contract.LogCoinBaseAddress(&_Coinbase.TransactOpts)
}

// LogCoinBaseAddress is a paid mutator transaction binding the contract method 0x713ee285.
//
// Solidity: function logCoinBaseAddress() returns()
func (_Coinbase *CoinbaseTransactorSession) LogCoinBaseAddress() (*types.Transaction, error) {
	return _Coinbase.Contract.LogCoinBaseAddress(&_Coinbase.TransactOpts)
}

// TouchAddress is a paid mutator transaction binding the contract method 0xd847a22b.
//
// Solidity: function touchAddress(address addr) returns()
func (_Coinbase *CoinbaseTransactor) TouchAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _Coinbase.contract.Transact(opts, "touchAddress", addr)
}

// TouchAddress is a paid mutator transaction binding the contract method 0xd847a22b.
//
// Solidity: function touchAddress(address addr) returns()
func (_Coinbase *CoinbaseSession) TouchAddress(addr common.Address) (*types.Transaction, error) {
	return _Coinbase.Contract.TouchAddress(&_Coinbase.TransactOpts, addr)
}

// TouchAddress is a paid mutator transaction binding the contract method 0xd847a22b.
//
// Solidity: function touchAddress(address addr) returns()
func (_Coinbase *CoinbaseTransactorSession) TouchAddress(addr common.Address) (*types.Transaction, error) {
	return _Coinbase.Contract.TouchAddress(&_Coinbase.TransactOpts, addr)
}

// TouchCoinbase is a paid mutator transaction binding the contract method 0x874395f8.
//
// Solidity: function touchCoinbase() returns()
func (_Coinbase *CoinbaseTransactor) TouchCoinbase(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coinbase.contract.Transact(opts, "touchCoinbase")
}

// TouchCoinbase is a paid mutator transaction binding the contract method 0x874395f8.
//
// Solidity: function touchCoinbase() returns()
func (_Coinbase *CoinbaseSession) TouchCoinbase() (*types.Transaction, error) {
	return _Coinbase.Contract.TouchCoinbase(&_Coinbase.TransactOpts)
}

// TouchCoinbase is a paid mutator transaction binding the contract method 0x874395f8.
//
// Solidity: function touchCoinbase() returns()
func (_Coinbase *CoinbaseTransactorSession) TouchCoinbase() (*types.Transaction, error) {
	return _Coinbase.Contract.TouchCoinbase(&_Coinbase.TransactOpts)
}

// CoinbaseLogAddressIterator is returned from FilterLogAddress and is used to iterate over the raw logs and unpacked data for LogAddress events raised by the Coinbase contract.
type CoinbaseLogAddressIterator struct {
	Event *CoinbaseLogAddress // Event containing the contract specifics and raw log

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
func (it *CoinbaseLogAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoinbaseLogAddress)
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
		it.Event = new(CoinbaseLogAddress)
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
func (it *CoinbaseLogAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoinbaseLogAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoinbaseLogAddress represents a LogAddress event raised by the Coinbase contract.
type CoinbaseLogAddress struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogAddress is a free log retrieval operation binding the contract event 0xb123f68b8ba02b447d91a6629e121111b7dd6061ff418a60139c8bf00522a284.
//
// Solidity: event LogAddress(address addr)
func (_Coinbase *CoinbaseFilterer) FilterLogAddress(opts *bind.FilterOpts) (*CoinbaseLogAddressIterator, error) {

	logs, sub, err := _Coinbase.contract.FilterLogs(opts, "LogAddress")
	if err != nil {
		return nil, err
	}
	return &CoinbaseLogAddressIterator{contract: _Coinbase.contract, event: "LogAddress", logs: logs, sub: sub}, nil
}

// WatchLogAddress is a free log subscription operation binding the contract event 0xb123f68b8ba02b447d91a6629e121111b7dd6061ff418a60139c8bf00522a284.
//
// Solidity: event LogAddress(address addr)
func (_Coinbase *CoinbaseFilterer) WatchLogAddress(opts *bind.WatchOpts, sink chan<- *CoinbaseLogAddress) (event.Subscription, error) {

	logs, sub, err := _Coinbase.contract.WatchLogs(opts, "LogAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoinbaseLogAddress)
				if err := _Coinbase.contract.UnpackLog(event, "LogAddress", log); err != nil {
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

// ParseLogAddress is a log parse operation binding the contract event 0xb123f68b8ba02b447d91a6629e121111b7dd6061ff418a60139c8bf00522a284.
//
// Solidity: event LogAddress(address addr)
func (_Coinbase *CoinbaseFilterer) ParseLogAddress(log types.Log) (*CoinbaseLogAddress, error) {
	event := new(CoinbaseLogAddress)
	if err := _Coinbase.contract.UnpackLog(event, "LogAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
