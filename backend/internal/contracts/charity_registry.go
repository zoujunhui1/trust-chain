// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// CharityRegistryMetaData contains all meta data concerning the CharityRegistry contract.
var CharityRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isVerified\",\"inputs\":[{\"name\":\"charity\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeCharity\",\"inputs\":[{\"name\":\"charity\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyCharity\",\"inputs\":[{\"name\":\"charity\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CharityRevoked\",\"inputs\":[{\"name\":\"charity\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CharityVerified\",\"inputs\":[{\"name\":\"charity\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// CharityRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use CharityRegistryMetaData.ABI instead.
var CharityRegistryABI = CharityRegistryMetaData.ABI

// CharityRegistry is an auto generated Go binding around an Ethereum contract.
type CharityRegistry struct {
	CharityRegistryCaller     // Read-only binding to the contract
	CharityRegistryTransactor // Write-only binding to the contract
	CharityRegistryFilterer   // Log filterer for contract events
}

// CharityRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type CharityRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CharityRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CharityRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CharityRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CharityRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CharityRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CharityRegistrySession struct {
	Contract     *CharityRegistry  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CharityRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CharityRegistryCallerSession struct {
	Contract *CharityRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CharityRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CharityRegistryTransactorSession struct {
	Contract     *CharityRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CharityRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type CharityRegistryRaw struct {
	Contract *CharityRegistry // Generic contract binding to access the raw methods on
}

// CharityRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CharityRegistryCallerRaw struct {
	Contract *CharityRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// CharityRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CharityRegistryTransactorRaw struct {
	Contract *CharityRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCharityRegistry creates a new instance of CharityRegistry, bound to a specific deployed contract.
func NewCharityRegistry(address common.Address, backend bind.ContractBackend) (*CharityRegistry, error) {
	contract, err := bindCharityRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CharityRegistry{CharityRegistryCaller: CharityRegistryCaller{contract: contract}, CharityRegistryTransactor: CharityRegistryTransactor{contract: contract}, CharityRegistryFilterer: CharityRegistryFilterer{contract: contract}}, nil
}

// NewCharityRegistryCaller creates a new read-only instance of CharityRegistry, bound to a specific deployed contract.
func NewCharityRegistryCaller(address common.Address, caller bind.ContractCaller) (*CharityRegistryCaller, error) {
	contract, err := bindCharityRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CharityRegistryCaller{contract: contract}, nil
}

// NewCharityRegistryTransactor creates a new write-only instance of CharityRegistry, bound to a specific deployed contract.
func NewCharityRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*CharityRegistryTransactor, error) {
	contract, err := bindCharityRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CharityRegistryTransactor{contract: contract}, nil
}

// NewCharityRegistryFilterer creates a new log filterer instance of CharityRegistry, bound to a specific deployed contract.
func NewCharityRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*CharityRegistryFilterer, error) {
	contract, err := bindCharityRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CharityRegistryFilterer{contract: contract}, nil
}

// bindCharityRegistry binds a generic wrapper to an already deployed contract.
func bindCharityRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CharityRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CharityRegistry *CharityRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CharityRegistry.Contract.CharityRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CharityRegistry *CharityRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CharityRegistry.Contract.CharityRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CharityRegistry *CharityRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CharityRegistry.Contract.CharityRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CharityRegistry *CharityRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CharityRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CharityRegistry *CharityRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CharityRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CharityRegistry *CharityRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CharityRegistry.Contract.contract.Transact(opts, method, params...)
}

// IsVerified is a free data retrieval call binding the contract method 0xb9209e33.
//
// Solidity: function isVerified(address charity) view returns(bool)
func (_CharityRegistry *CharityRegistryCaller) IsVerified(opts *bind.CallOpts, charity common.Address) (bool, error) {
	var out []interface{}
	err := _CharityRegistry.contract.Call(opts, &out, "isVerified", charity)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsVerified is a free data retrieval call binding the contract method 0xb9209e33.
//
// Solidity: function isVerified(address charity) view returns(bool)
func (_CharityRegistry *CharityRegistrySession) IsVerified(charity common.Address) (bool, error) {
	return _CharityRegistry.Contract.IsVerified(&_CharityRegistry.CallOpts, charity)
}

// IsVerified is a free data retrieval call binding the contract method 0xb9209e33.
//
// Solidity: function isVerified(address charity) view returns(bool)
func (_CharityRegistry *CharityRegistryCallerSession) IsVerified(charity common.Address) (bool, error) {
	return _CharityRegistry.Contract.IsVerified(&_CharityRegistry.CallOpts, charity)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CharityRegistry *CharityRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CharityRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CharityRegistry *CharityRegistrySession) Owner() (common.Address, error) {
	return _CharityRegistry.Contract.Owner(&_CharityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CharityRegistry *CharityRegistryCallerSession) Owner() (common.Address, error) {
	return _CharityRegistry.Contract.Owner(&_CharityRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CharityRegistry *CharityRegistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CharityRegistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CharityRegistry *CharityRegistrySession) PendingOwner() (common.Address, error) {
	return _CharityRegistry.Contract.PendingOwner(&_CharityRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CharityRegistry *CharityRegistryCallerSession) PendingOwner() (common.Address, error) {
	return _CharityRegistry.Contract.PendingOwner(&_CharityRegistry.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CharityRegistry *CharityRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CharityRegistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CharityRegistry *CharityRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _CharityRegistry.Contract.AcceptOwnership(&_CharityRegistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CharityRegistry *CharityRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CharityRegistry.Contract.AcceptOwnership(&_CharityRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CharityRegistry *CharityRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CharityRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CharityRegistry *CharityRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _CharityRegistry.Contract.RenounceOwnership(&_CharityRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CharityRegistry *CharityRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CharityRegistry.Contract.RenounceOwnership(&_CharityRegistry.TransactOpts)
}

// RevokeCharity is a paid mutator transaction binding the contract method 0x18c9c1dc.
//
// Solidity: function revokeCharity(address charity) returns()
func (_CharityRegistry *CharityRegistryTransactor) RevokeCharity(opts *bind.TransactOpts, charity common.Address) (*types.Transaction, error) {
	return _CharityRegistry.contract.Transact(opts, "revokeCharity", charity)
}

// RevokeCharity is a paid mutator transaction binding the contract method 0x18c9c1dc.
//
// Solidity: function revokeCharity(address charity) returns()
func (_CharityRegistry *CharityRegistrySession) RevokeCharity(charity common.Address) (*types.Transaction, error) {
	return _CharityRegistry.Contract.RevokeCharity(&_CharityRegistry.TransactOpts, charity)
}

// RevokeCharity is a paid mutator transaction binding the contract method 0x18c9c1dc.
//
// Solidity: function revokeCharity(address charity) returns()
func (_CharityRegistry *CharityRegistryTransactorSession) RevokeCharity(charity common.Address) (*types.Transaction, error) {
	return _CharityRegistry.Contract.RevokeCharity(&_CharityRegistry.TransactOpts, charity)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CharityRegistry *CharityRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CharityRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CharityRegistry *CharityRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CharityRegistry.Contract.TransferOwnership(&_CharityRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CharityRegistry *CharityRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CharityRegistry.Contract.TransferOwnership(&_CharityRegistry.TransactOpts, newOwner)
}

// VerifyCharity is a paid mutator transaction binding the contract method 0x72097548.
//
// Solidity: function verifyCharity(address charity) returns()
func (_CharityRegistry *CharityRegistryTransactor) VerifyCharity(opts *bind.TransactOpts, charity common.Address) (*types.Transaction, error) {
	return _CharityRegistry.contract.Transact(opts, "verifyCharity", charity)
}

// VerifyCharity is a paid mutator transaction binding the contract method 0x72097548.
//
// Solidity: function verifyCharity(address charity) returns()
func (_CharityRegistry *CharityRegistrySession) VerifyCharity(charity common.Address) (*types.Transaction, error) {
	return _CharityRegistry.Contract.VerifyCharity(&_CharityRegistry.TransactOpts, charity)
}

// VerifyCharity is a paid mutator transaction binding the contract method 0x72097548.
//
// Solidity: function verifyCharity(address charity) returns()
func (_CharityRegistry *CharityRegistryTransactorSession) VerifyCharity(charity common.Address) (*types.Transaction, error) {
	return _CharityRegistry.Contract.VerifyCharity(&_CharityRegistry.TransactOpts, charity)
}

// CharityRegistryCharityRevokedIterator is returned from FilterCharityRevoked and is used to iterate over the raw logs and unpacked data for CharityRevoked events raised by the CharityRegistry contract.
type CharityRegistryCharityRevokedIterator struct {
	Event *CharityRegistryCharityRevoked // Event containing the contract specifics and raw log

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
func (it *CharityRegistryCharityRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharityRegistryCharityRevoked)
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
		it.Event = new(CharityRegistryCharityRevoked)
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
func (it *CharityRegistryCharityRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharityRegistryCharityRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharityRegistryCharityRevoked represents a CharityRevoked event raised by the CharityRegistry contract.
type CharityRegistryCharityRevoked struct {
	Charity common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterCharityRevoked is a free log retrieval operation binding the contract event 0x5c2c8d962dfdf7cf138ba30b3641ae6d187b74b7b4fd6506e381e344ce09ad83.
//
// Solidity: event CharityRevoked(address indexed charity)
func (_CharityRegistry *CharityRegistryFilterer) FilterCharityRevoked(opts *bind.FilterOpts, charity []common.Address) (*CharityRegistryCharityRevokedIterator, error) {

	var charityRule []interface{}
	for _, charityItem := range charity {
		charityRule = append(charityRule, charityItem)
	}

	logs, sub, err := _CharityRegistry.contract.FilterLogs(opts, "CharityRevoked", charityRule)
	if err != nil {
		return nil, err
	}
	return &CharityRegistryCharityRevokedIterator{contract: _CharityRegistry.contract, event: "CharityRevoked", logs: logs, sub: sub}, nil
}

// WatchCharityRevoked is a free log subscription operation binding the contract event 0x5c2c8d962dfdf7cf138ba30b3641ae6d187b74b7b4fd6506e381e344ce09ad83.
//
// Solidity: event CharityRevoked(address indexed charity)
func (_CharityRegistry *CharityRegistryFilterer) WatchCharityRevoked(opts *bind.WatchOpts, sink chan<- *CharityRegistryCharityRevoked, charity []common.Address) (event.Subscription, error) {

	var charityRule []interface{}
	for _, charityItem := range charity {
		charityRule = append(charityRule, charityItem)
	}

	logs, sub, err := _CharityRegistry.contract.WatchLogs(opts, "CharityRevoked", charityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharityRegistryCharityRevoked)
				if err := _CharityRegistry.contract.UnpackLog(event, "CharityRevoked", log); err != nil {
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

// ParseCharityRevoked is a log parse operation binding the contract event 0x5c2c8d962dfdf7cf138ba30b3641ae6d187b74b7b4fd6506e381e344ce09ad83.
//
// Solidity: event CharityRevoked(address indexed charity)
func (_CharityRegistry *CharityRegistryFilterer) ParseCharityRevoked(log types.Log) (*CharityRegistryCharityRevoked, error) {
	event := new(CharityRegistryCharityRevoked)
	if err := _CharityRegistry.contract.UnpackLog(event, "CharityRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharityRegistryCharityVerifiedIterator is returned from FilterCharityVerified and is used to iterate over the raw logs and unpacked data for CharityVerified events raised by the CharityRegistry contract.
type CharityRegistryCharityVerifiedIterator struct {
	Event *CharityRegistryCharityVerified // Event containing the contract specifics and raw log

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
func (it *CharityRegistryCharityVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharityRegistryCharityVerified)
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
		it.Event = new(CharityRegistryCharityVerified)
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
func (it *CharityRegistryCharityVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharityRegistryCharityVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharityRegistryCharityVerified represents a CharityVerified event raised by the CharityRegistry contract.
type CharityRegistryCharityVerified struct {
	Charity common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterCharityVerified is a free log retrieval operation binding the contract event 0x27ecb4f695edff6a3981d86a210339a309eba3476cebdde8ed999783024dd594.
//
// Solidity: event CharityVerified(address indexed charity)
func (_CharityRegistry *CharityRegistryFilterer) FilterCharityVerified(opts *bind.FilterOpts, charity []common.Address) (*CharityRegistryCharityVerifiedIterator, error) {

	var charityRule []interface{}
	for _, charityItem := range charity {
		charityRule = append(charityRule, charityItem)
	}

	logs, sub, err := _CharityRegistry.contract.FilterLogs(opts, "CharityVerified", charityRule)
	if err != nil {
		return nil, err
	}
	return &CharityRegistryCharityVerifiedIterator{contract: _CharityRegistry.contract, event: "CharityVerified", logs: logs, sub: sub}, nil
}

// WatchCharityVerified is a free log subscription operation binding the contract event 0x27ecb4f695edff6a3981d86a210339a309eba3476cebdde8ed999783024dd594.
//
// Solidity: event CharityVerified(address indexed charity)
func (_CharityRegistry *CharityRegistryFilterer) WatchCharityVerified(opts *bind.WatchOpts, sink chan<- *CharityRegistryCharityVerified, charity []common.Address) (event.Subscription, error) {

	var charityRule []interface{}
	for _, charityItem := range charity {
		charityRule = append(charityRule, charityItem)
	}

	logs, sub, err := _CharityRegistry.contract.WatchLogs(opts, "CharityVerified", charityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharityRegistryCharityVerified)
				if err := _CharityRegistry.contract.UnpackLog(event, "CharityVerified", log); err != nil {
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

// ParseCharityVerified is a log parse operation binding the contract event 0x27ecb4f695edff6a3981d86a210339a309eba3476cebdde8ed999783024dd594.
//
// Solidity: event CharityVerified(address indexed charity)
func (_CharityRegistry *CharityRegistryFilterer) ParseCharityVerified(log types.Log) (*CharityRegistryCharityVerified, error) {
	event := new(CharityRegistryCharityVerified)
	if err := _CharityRegistry.contract.UnpackLog(event, "CharityVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharityRegistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the CharityRegistry contract.
type CharityRegistryOwnershipTransferStartedIterator struct {
	Event *CharityRegistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *CharityRegistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharityRegistryOwnershipTransferStarted)
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
		it.Event = new(CharityRegistryOwnershipTransferStarted)
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
func (it *CharityRegistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharityRegistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharityRegistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the CharityRegistry contract.
type CharityRegistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CharityRegistry *CharityRegistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CharityRegistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CharityRegistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CharityRegistryOwnershipTransferStartedIterator{contract: _CharityRegistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CharityRegistry *CharityRegistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *CharityRegistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CharityRegistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharityRegistryOwnershipTransferStarted)
				if err := _CharityRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CharityRegistry *CharityRegistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*CharityRegistryOwnershipTransferStarted, error) {
	event := new(CharityRegistryOwnershipTransferStarted)
	if err := _CharityRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CharityRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CharityRegistry contract.
type CharityRegistryOwnershipTransferredIterator struct {
	Event *CharityRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CharityRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CharityRegistryOwnershipTransferred)
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
		it.Event = new(CharityRegistryOwnershipTransferred)
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
func (it *CharityRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CharityRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CharityRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the CharityRegistry contract.
type CharityRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CharityRegistry *CharityRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CharityRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CharityRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CharityRegistryOwnershipTransferredIterator{contract: _CharityRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CharityRegistry *CharityRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CharityRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CharityRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CharityRegistryOwnershipTransferred)
				if err := _CharityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CharityRegistry *CharityRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*CharityRegistryOwnershipTransferred, error) {
	event := new(CharityRegistryOwnershipTransferred)
	if err := _CharityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
