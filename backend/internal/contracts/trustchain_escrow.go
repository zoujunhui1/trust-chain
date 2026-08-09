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

// TrustChainEscrowMilestone is an auto generated low-level Go binding around an user-defined struct.
type TrustChainEscrowMilestone struct {
	Amount      *big.Int
	State       uint8
	ReceiptHash [32]byte
}

// TrustChainEscrowMetaData contains all meta data concerning the TrustChainEscrow contract.
var TrustChainEscrowMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"registry_\",\"type\":\"address\",\"internalType\":\"contractCharityRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"campaignCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createCampaign\",\"inputs\":[{\"name\":\"milestoneAmounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"metadataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"donate\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"donations\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCampaign\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"charity\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"goal\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"raised\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"released\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"milestoneCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMilestones\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structTrustChainEscrow.Milestone[]\",\"components\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumTrustChainEscrow.MilestoneState\"},{\"name\":\"receiptHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractCharityRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"releaseMilestone\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"milestoneIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitReceipt\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"milestoneIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiptHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CampaignCompleted\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CampaignCreated\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"charity\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"goal\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"milestoneCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DonationReceived\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"donor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MilestoneReleased\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"milestoneIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReceiptSubmitted\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"milestoneIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"receiptHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
}

// TrustChainEscrowABI is the input ABI used to generate the binding from.
// Deprecated: Use TrustChainEscrowMetaData.ABI instead.
var TrustChainEscrowABI = TrustChainEscrowMetaData.ABI

// TrustChainEscrow is an auto generated Go binding around an Ethereum contract.
type TrustChainEscrow struct {
	TrustChainEscrowCaller     // Read-only binding to the contract
	TrustChainEscrowTransactor // Write-only binding to the contract
	TrustChainEscrowFilterer   // Log filterer for contract events
}

// TrustChainEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type TrustChainEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustChainEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TrustChainEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustChainEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TrustChainEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustChainEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TrustChainEscrowSession struct {
	Contract     *TrustChainEscrow // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TrustChainEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TrustChainEscrowCallerSession struct {
	Contract *TrustChainEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// TrustChainEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TrustChainEscrowTransactorSession struct {
	Contract     *TrustChainEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TrustChainEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type TrustChainEscrowRaw struct {
	Contract *TrustChainEscrow // Generic contract binding to access the raw methods on
}

// TrustChainEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TrustChainEscrowCallerRaw struct {
	Contract *TrustChainEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// TrustChainEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TrustChainEscrowTransactorRaw struct {
	Contract *TrustChainEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTrustChainEscrow creates a new instance of TrustChainEscrow, bound to a specific deployed contract.
func NewTrustChainEscrow(address common.Address, backend bind.ContractBackend) (*TrustChainEscrow, error) {
	contract, err := bindTrustChainEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrow{TrustChainEscrowCaller: TrustChainEscrowCaller{contract: contract}, TrustChainEscrowTransactor: TrustChainEscrowTransactor{contract: contract}, TrustChainEscrowFilterer: TrustChainEscrowFilterer{contract: contract}}, nil
}

// NewTrustChainEscrowCaller creates a new read-only instance of TrustChainEscrow, bound to a specific deployed contract.
func NewTrustChainEscrowCaller(address common.Address, caller bind.ContractCaller) (*TrustChainEscrowCaller, error) {
	contract, err := bindTrustChainEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrowCaller{contract: contract}, nil
}

// NewTrustChainEscrowTransactor creates a new write-only instance of TrustChainEscrow, bound to a specific deployed contract.
func NewTrustChainEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*TrustChainEscrowTransactor, error) {
	contract, err := bindTrustChainEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrowTransactor{contract: contract}, nil
}

// NewTrustChainEscrowFilterer creates a new log filterer instance of TrustChainEscrow, bound to a specific deployed contract.
func NewTrustChainEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*TrustChainEscrowFilterer, error) {
	contract, err := bindTrustChainEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrowFilterer{contract: contract}, nil
}

// bindTrustChainEscrow binds a generic wrapper to an already deployed contract.
func bindTrustChainEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TrustChainEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TrustChainEscrow *TrustChainEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TrustChainEscrow.Contract.TrustChainEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TrustChainEscrow *TrustChainEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.TrustChainEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TrustChainEscrow *TrustChainEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.TrustChainEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TrustChainEscrow *TrustChainEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TrustChainEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TrustChainEscrow *TrustChainEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TrustChainEscrow *TrustChainEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.contract.Transact(opts, method, params...)
}

// CampaignCount is a free data retrieval call binding the contract method 0x7274e30d.
//
// Solidity: function campaignCount() view returns(uint256)
func (_TrustChainEscrow *TrustChainEscrowCaller) CampaignCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TrustChainEscrow.contract.Call(opts, &out, "campaignCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CampaignCount is a free data retrieval call binding the contract method 0x7274e30d.
//
// Solidity: function campaignCount() view returns(uint256)
func (_TrustChainEscrow *TrustChainEscrowSession) CampaignCount() (*big.Int, error) {
	return _TrustChainEscrow.Contract.CampaignCount(&_TrustChainEscrow.CallOpts)
}

// CampaignCount is a free data retrieval call binding the contract method 0x7274e30d.
//
// Solidity: function campaignCount() view returns(uint256)
func (_TrustChainEscrow *TrustChainEscrowCallerSession) CampaignCount() (*big.Int, error) {
	return _TrustChainEscrow.Contract.CampaignCount(&_TrustChainEscrow.CallOpts)
}

// Donations is a free data retrieval call binding the contract method 0xee1554a3.
//
// Solidity: function donations(uint256 , address ) view returns(uint256)
func (_TrustChainEscrow *TrustChainEscrowCaller) Donations(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TrustChainEscrow.contract.Call(opts, &out, "donations", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Donations is a free data retrieval call binding the contract method 0xee1554a3.
//
// Solidity: function donations(uint256 , address ) view returns(uint256)
func (_TrustChainEscrow *TrustChainEscrowSession) Donations(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _TrustChainEscrow.Contract.Donations(&_TrustChainEscrow.CallOpts, arg0, arg1)
}

// Donations is a free data retrieval call binding the contract method 0xee1554a3.
//
// Solidity: function donations(uint256 , address ) view returns(uint256)
func (_TrustChainEscrow *TrustChainEscrowCallerSession) Donations(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _TrustChainEscrow.Contract.Donations(&_TrustChainEscrow.CallOpts, arg0, arg1)
}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns(address charity, uint256 goal, uint256 raised, uint256 released, bytes32 metadataHash, bool exists, uint256 milestoneCount)
func (_TrustChainEscrow *TrustChainEscrowCaller) GetCampaign(opts *bind.CallOpts, campaignId *big.Int) (struct {
	Charity        common.Address
	Goal           *big.Int
	Raised         *big.Int
	Released       *big.Int
	MetadataHash   [32]byte
	Exists         bool
	MilestoneCount *big.Int
}, error) {
	var out []interface{}
	err := _TrustChainEscrow.contract.Call(opts, &out, "getCampaign", campaignId)

	outstruct := new(struct {
		Charity        common.Address
		Goal           *big.Int
		Raised         *big.Int
		Released       *big.Int
		MetadataHash   [32]byte
		Exists         bool
		MilestoneCount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Charity = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Goal = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Raised = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Released = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.MetadataHash = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.Exists = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.MilestoneCount = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns(address charity, uint256 goal, uint256 raised, uint256 released, bytes32 metadataHash, bool exists, uint256 milestoneCount)
func (_TrustChainEscrow *TrustChainEscrowSession) GetCampaign(campaignId *big.Int) (struct {
	Charity        common.Address
	Goal           *big.Int
	Raised         *big.Int
	Released       *big.Int
	MetadataHash   [32]byte
	Exists         bool
	MilestoneCount *big.Int
}, error) {
	return _TrustChainEscrow.Contract.GetCampaign(&_TrustChainEscrow.CallOpts, campaignId)
}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns(address charity, uint256 goal, uint256 raised, uint256 released, bytes32 metadataHash, bool exists, uint256 milestoneCount)
func (_TrustChainEscrow *TrustChainEscrowCallerSession) GetCampaign(campaignId *big.Int) (struct {
	Charity        common.Address
	Goal           *big.Int
	Raised         *big.Int
	Released       *big.Int
	MetadataHash   [32]byte
	Exists         bool
	MilestoneCount *big.Int
}, error) {
	return _TrustChainEscrow.Contract.GetCampaign(&_TrustChainEscrow.CallOpts, campaignId)
}

// GetMilestones is a free data retrieval call binding the contract method 0x42c549c0.
//
// Solidity: function getMilestones(uint256 campaignId) view returns((uint256,uint8,bytes32)[])
func (_TrustChainEscrow *TrustChainEscrowCaller) GetMilestones(opts *bind.CallOpts, campaignId *big.Int) ([]TrustChainEscrowMilestone, error) {
	var out []interface{}
	err := _TrustChainEscrow.contract.Call(opts, &out, "getMilestones", campaignId)

	if err != nil {
		return *new([]TrustChainEscrowMilestone), err
	}

	out0 := *abi.ConvertType(out[0], new([]TrustChainEscrowMilestone)).(*[]TrustChainEscrowMilestone)

	return out0, err

}

// GetMilestones is a free data retrieval call binding the contract method 0x42c549c0.
//
// Solidity: function getMilestones(uint256 campaignId) view returns((uint256,uint8,bytes32)[])
func (_TrustChainEscrow *TrustChainEscrowSession) GetMilestones(campaignId *big.Int) ([]TrustChainEscrowMilestone, error) {
	return _TrustChainEscrow.Contract.GetMilestones(&_TrustChainEscrow.CallOpts, campaignId)
}

// GetMilestones is a free data retrieval call binding the contract method 0x42c549c0.
//
// Solidity: function getMilestones(uint256 campaignId) view returns((uint256,uint8,bytes32)[])
func (_TrustChainEscrow *TrustChainEscrowCallerSession) GetMilestones(campaignId *big.Int) ([]TrustChainEscrowMilestone, error) {
	return _TrustChainEscrow.Contract.GetMilestones(&_TrustChainEscrow.CallOpts, campaignId)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_TrustChainEscrow *TrustChainEscrowCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TrustChainEscrow.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_TrustChainEscrow *TrustChainEscrowSession) Registry() (common.Address, error) {
	return _TrustChainEscrow.Contract.Registry(&_TrustChainEscrow.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_TrustChainEscrow *TrustChainEscrowCallerSession) Registry() (common.Address, error) {
	return _TrustChainEscrow.Contract.Registry(&_TrustChainEscrow.CallOpts)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0xcbe460fb.
//
// Solidity: function createCampaign(uint256[] milestoneAmounts, bytes32 metadataHash) returns(uint256 campaignId)
func (_TrustChainEscrow *TrustChainEscrowTransactor) CreateCampaign(opts *bind.TransactOpts, milestoneAmounts []*big.Int, metadataHash [32]byte) (*types.Transaction, error) {
	return _TrustChainEscrow.contract.Transact(opts, "createCampaign", milestoneAmounts, metadataHash)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0xcbe460fb.
//
// Solidity: function createCampaign(uint256[] milestoneAmounts, bytes32 metadataHash) returns(uint256 campaignId)
func (_TrustChainEscrow *TrustChainEscrowSession) CreateCampaign(milestoneAmounts []*big.Int, metadataHash [32]byte) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.CreateCampaign(&_TrustChainEscrow.TransactOpts, milestoneAmounts, metadataHash)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0xcbe460fb.
//
// Solidity: function createCampaign(uint256[] milestoneAmounts, bytes32 metadataHash) returns(uint256 campaignId)
func (_TrustChainEscrow *TrustChainEscrowTransactorSession) CreateCampaign(milestoneAmounts []*big.Int, metadataHash [32]byte) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.CreateCampaign(&_TrustChainEscrow.TransactOpts, milestoneAmounts, metadataHash)
}

// Donate is a paid mutator transaction binding the contract method 0xf14faf6f.
//
// Solidity: function donate(uint256 campaignId) payable returns()
func (_TrustChainEscrow *TrustChainEscrowTransactor) Donate(opts *bind.TransactOpts, campaignId *big.Int) (*types.Transaction, error) {
	return _TrustChainEscrow.contract.Transact(opts, "donate", campaignId)
}

// Donate is a paid mutator transaction binding the contract method 0xf14faf6f.
//
// Solidity: function donate(uint256 campaignId) payable returns()
func (_TrustChainEscrow *TrustChainEscrowSession) Donate(campaignId *big.Int) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.Donate(&_TrustChainEscrow.TransactOpts, campaignId)
}

// Donate is a paid mutator transaction binding the contract method 0xf14faf6f.
//
// Solidity: function donate(uint256 campaignId) payable returns()
func (_TrustChainEscrow *TrustChainEscrowTransactorSession) Donate(campaignId *big.Int) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.Donate(&_TrustChainEscrow.TransactOpts, campaignId)
}

// ReleaseMilestone is a paid mutator transaction binding the contract method 0xdbaf0910.
//
// Solidity: function releaseMilestone(uint256 campaignId, uint256 milestoneIndex) returns()
func (_TrustChainEscrow *TrustChainEscrowTransactor) ReleaseMilestone(opts *bind.TransactOpts, campaignId *big.Int, milestoneIndex *big.Int) (*types.Transaction, error) {
	return _TrustChainEscrow.contract.Transact(opts, "releaseMilestone", campaignId, milestoneIndex)
}

// ReleaseMilestone is a paid mutator transaction binding the contract method 0xdbaf0910.
//
// Solidity: function releaseMilestone(uint256 campaignId, uint256 milestoneIndex) returns()
func (_TrustChainEscrow *TrustChainEscrowSession) ReleaseMilestone(campaignId *big.Int, milestoneIndex *big.Int) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.ReleaseMilestone(&_TrustChainEscrow.TransactOpts, campaignId, milestoneIndex)
}

// ReleaseMilestone is a paid mutator transaction binding the contract method 0xdbaf0910.
//
// Solidity: function releaseMilestone(uint256 campaignId, uint256 milestoneIndex) returns()
func (_TrustChainEscrow *TrustChainEscrowTransactorSession) ReleaseMilestone(campaignId *big.Int, milestoneIndex *big.Int) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.ReleaseMilestone(&_TrustChainEscrow.TransactOpts, campaignId, milestoneIndex)
}

// SubmitReceipt is a paid mutator transaction binding the contract method 0x049b6d11.
//
// Solidity: function submitReceipt(uint256 campaignId, uint256 milestoneIndex, bytes32 receiptHash) returns()
func (_TrustChainEscrow *TrustChainEscrowTransactor) SubmitReceipt(opts *bind.TransactOpts, campaignId *big.Int, milestoneIndex *big.Int, receiptHash [32]byte) (*types.Transaction, error) {
	return _TrustChainEscrow.contract.Transact(opts, "submitReceipt", campaignId, milestoneIndex, receiptHash)
}

// SubmitReceipt is a paid mutator transaction binding the contract method 0x049b6d11.
//
// Solidity: function submitReceipt(uint256 campaignId, uint256 milestoneIndex, bytes32 receiptHash) returns()
func (_TrustChainEscrow *TrustChainEscrowSession) SubmitReceipt(campaignId *big.Int, milestoneIndex *big.Int, receiptHash [32]byte) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.SubmitReceipt(&_TrustChainEscrow.TransactOpts, campaignId, milestoneIndex, receiptHash)
}

// SubmitReceipt is a paid mutator transaction binding the contract method 0x049b6d11.
//
// Solidity: function submitReceipt(uint256 campaignId, uint256 milestoneIndex, bytes32 receiptHash) returns()
func (_TrustChainEscrow *TrustChainEscrowTransactorSession) SubmitReceipt(campaignId *big.Int, milestoneIndex *big.Int, receiptHash [32]byte) (*types.Transaction, error) {
	return _TrustChainEscrow.Contract.SubmitReceipt(&_TrustChainEscrow.TransactOpts, campaignId, milestoneIndex, receiptHash)
}

// TrustChainEscrowCampaignCompletedIterator is returned from FilterCampaignCompleted and is used to iterate over the raw logs and unpacked data for CampaignCompleted events raised by the TrustChainEscrow contract.
type TrustChainEscrowCampaignCompletedIterator struct {
	Event *TrustChainEscrowCampaignCompleted // Event containing the contract specifics and raw log

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
func (it *TrustChainEscrowCampaignCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustChainEscrowCampaignCompleted)
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
		it.Event = new(TrustChainEscrowCampaignCompleted)
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
func (it *TrustChainEscrowCampaignCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustChainEscrowCampaignCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustChainEscrowCampaignCompleted represents a CampaignCompleted event raised by the TrustChainEscrow contract.
type TrustChainEscrowCampaignCompleted struct {
	CampaignId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCampaignCompleted is a free log retrieval operation binding the contract event 0xbdef6a2e3aa892961a2dab48272f72f2bc1c7121eb1c0ed7e4a58e081c1e4be0.
//
// Solidity: event CampaignCompleted(uint256 indexed campaignId)
func (_TrustChainEscrow *TrustChainEscrowFilterer) FilterCampaignCompleted(opts *bind.FilterOpts, campaignId []*big.Int) (*TrustChainEscrowCampaignCompletedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.FilterLogs(opts, "CampaignCompleted", campaignIdRule)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrowCampaignCompletedIterator{contract: _TrustChainEscrow.contract, event: "CampaignCompleted", logs: logs, sub: sub}, nil
}

// WatchCampaignCompleted is a free log subscription operation binding the contract event 0xbdef6a2e3aa892961a2dab48272f72f2bc1c7121eb1c0ed7e4a58e081c1e4be0.
//
// Solidity: event CampaignCompleted(uint256 indexed campaignId)
func (_TrustChainEscrow *TrustChainEscrowFilterer) WatchCampaignCompleted(opts *bind.WatchOpts, sink chan<- *TrustChainEscrowCampaignCompleted, campaignId []*big.Int) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.WatchLogs(opts, "CampaignCompleted", campaignIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustChainEscrowCampaignCompleted)
				if err := _TrustChainEscrow.contract.UnpackLog(event, "CampaignCompleted", log); err != nil {
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

// ParseCampaignCompleted is a log parse operation binding the contract event 0xbdef6a2e3aa892961a2dab48272f72f2bc1c7121eb1c0ed7e4a58e081c1e4be0.
//
// Solidity: event CampaignCompleted(uint256 indexed campaignId)
func (_TrustChainEscrow *TrustChainEscrowFilterer) ParseCampaignCompleted(log types.Log) (*TrustChainEscrowCampaignCompleted, error) {
	event := new(TrustChainEscrowCampaignCompleted)
	if err := _TrustChainEscrow.contract.UnpackLog(event, "CampaignCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrustChainEscrowCampaignCreatedIterator is returned from FilterCampaignCreated and is used to iterate over the raw logs and unpacked data for CampaignCreated events raised by the TrustChainEscrow contract.
type TrustChainEscrowCampaignCreatedIterator struct {
	Event *TrustChainEscrowCampaignCreated // Event containing the contract specifics and raw log

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
func (it *TrustChainEscrowCampaignCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustChainEscrowCampaignCreated)
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
		it.Event = new(TrustChainEscrowCampaignCreated)
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
func (it *TrustChainEscrowCampaignCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustChainEscrowCampaignCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustChainEscrowCampaignCreated represents a CampaignCreated event raised by the TrustChainEscrow contract.
type TrustChainEscrowCampaignCreated struct {
	CampaignId     *big.Int
	Charity        common.Address
	Goal           *big.Int
	MilestoneCount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCampaignCreated is a free log retrieval operation binding the contract event 0x91b289a829e71d811b8c69e4a24ba2d40d115d8a236e9a724cb3bb2d43cf7223.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed charity, uint256 goal, uint256 milestoneCount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) FilterCampaignCreated(opts *bind.FilterOpts, campaignId []*big.Int, charity []common.Address) (*TrustChainEscrowCampaignCreatedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var charityRule []interface{}
	for _, charityItem := range charity {
		charityRule = append(charityRule, charityItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.FilterLogs(opts, "CampaignCreated", campaignIdRule, charityRule)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrowCampaignCreatedIterator{contract: _TrustChainEscrow.contract, event: "CampaignCreated", logs: logs, sub: sub}, nil
}

// WatchCampaignCreated is a free log subscription operation binding the contract event 0x91b289a829e71d811b8c69e4a24ba2d40d115d8a236e9a724cb3bb2d43cf7223.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed charity, uint256 goal, uint256 milestoneCount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) WatchCampaignCreated(opts *bind.WatchOpts, sink chan<- *TrustChainEscrowCampaignCreated, campaignId []*big.Int, charity []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var charityRule []interface{}
	for _, charityItem := range charity {
		charityRule = append(charityRule, charityItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.WatchLogs(opts, "CampaignCreated", campaignIdRule, charityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustChainEscrowCampaignCreated)
				if err := _TrustChainEscrow.contract.UnpackLog(event, "CampaignCreated", log); err != nil {
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

// ParseCampaignCreated is a log parse operation binding the contract event 0x91b289a829e71d811b8c69e4a24ba2d40d115d8a236e9a724cb3bb2d43cf7223.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed charity, uint256 goal, uint256 milestoneCount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) ParseCampaignCreated(log types.Log) (*TrustChainEscrowCampaignCreated, error) {
	event := new(TrustChainEscrowCampaignCreated)
	if err := _TrustChainEscrow.contract.UnpackLog(event, "CampaignCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrustChainEscrowDonationReceivedIterator is returned from FilterDonationReceived and is used to iterate over the raw logs and unpacked data for DonationReceived events raised by the TrustChainEscrow contract.
type TrustChainEscrowDonationReceivedIterator struct {
	Event *TrustChainEscrowDonationReceived // Event containing the contract specifics and raw log

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
func (it *TrustChainEscrowDonationReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustChainEscrowDonationReceived)
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
		it.Event = new(TrustChainEscrowDonationReceived)
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
func (it *TrustChainEscrowDonationReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustChainEscrowDonationReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustChainEscrowDonationReceived represents a DonationReceived event raised by the TrustChainEscrow contract.
type TrustChainEscrowDonationReceived struct {
	CampaignId *big.Int
	Donor      common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDonationReceived is a free log retrieval operation binding the contract event 0x0b5b4c52969ff7329ecf7ee536409fda87812b15a8622bc6e8cdeab3aee14a26.
//
// Solidity: event DonationReceived(uint256 indexed campaignId, address indexed donor, uint256 amount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) FilterDonationReceived(opts *bind.FilterOpts, campaignId []*big.Int, donor []common.Address) (*TrustChainEscrowDonationReceivedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var donorRule []interface{}
	for _, donorItem := range donor {
		donorRule = append(donorRule, donorItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.FilterLogs(opts, "DonationReceived", campaignIdRule, donorRule)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrowDonationReceivedIterator{contract: _TrustChainEscrow.contract, event: "DonationReceived", logs: logs, sub: sub}, nil
}

// WatchDonationReceived is a free log subscription operation binding the contract event 0x0b5b4c52969ff7329ecf7ee536409fda87812b15a8622bc6e8cdeab3aee14a26.
//
// Solidity: event DonationReceived(uint256 indexed campaignId, address indexed donor, uint256 amount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) WatchDonationReceived(opts *bind.WatchOpts, sink chan<- *TrustChainEscrowDonationReceived, campaignId []*big.Int, donor []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var donorRule []interface{}
	for _, donorItem := range donor {
		donorRule = append(donorRule, donorItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.WatchLogs(opts, "DonationReceived", campaignIdRule, donorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustChainEscrowDonationReceived)
				if err := _TrustChainEscrow.contract.UnpackLog(event, "DonationReceived", log); err != nil {
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

// ParseDonationReceived is a log parse operation binding the contract event 0x0b5b4c52969ff7329ecf7ee536409fda87812b15a8622bc6e8cdeab3aee14a26.
//
// Solidity: event DonationReceived(uint256 indexed campaignId, address indexed donor, uint256 amount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) ParseDonationReceived(log types.Log) (*TrustChainEscrowDonationReceived, error) {
	event := new(TrustChainEscrowDonationReceived)
	if err := _TrustChainEscrow.contract.UnpackLog(event, "DonationReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrustChainEscrowMilestoneReleasedIterator is returned from FilterMilestoneReleased and is used to iterate over the raw logs and unpacked data for MilestoneReleased events raised by the TrustChainEscrow contract.
type TrustChainEscrowMilestoneReleasedIterator struct {
	Event *TrustChainEscrowMilestoneReleased // Event containing the contract specifics and raw log

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
func (it *TrustChainEscrowMilestoneReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustChainEscrowMilestoneReleased)
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
		it.Event = new(TrustChainEscrowMilestoneReleased)
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
func (it *TrustChainEscrowMilestoneReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustChainEscrowMilestoneReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustChainEscrowMilestoneReleased represents a MilestoneReleased event raised by the TrustChainEscrow contract.
type TrustChainEscrowMilestoneReleased struct {
	CampaignId     *big.Int
	MilestoneIndex *big.Int
	Amount         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMilestoneReleased is a free log retrieval operation binding the contract event 0xd33fcd27cbfb73a2d9058c968b560ccb306500bf536133802048489ecc7355e9.
//
// Solidity: event MilestoneReleased(uint256 indexed campaignId, uint256 indexed milestoneIndex, uint256 amount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) FilterMilestoneReleased(opts *bind.FilterOpts, campaignId []*big.Int, milestoneIndex []*big.Int) (*TrustChainEscrowMilestoneReleasedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var milestoneIndexRule []interface{}
	for _, milestoneIndexItem := range milestoneIndex {
		milestoneIndexRule = append(milestoneIndexRule, milestoneIndexItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.FilterLogs(opts, "MilestoneReleased", campaignIdRule, milestoneIndexRule)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrowMilestoneReleasedIterator{contract: _TrustChainEscrow.contract, event: "MilestoneReleased", logs: logs, sub: sub}, nil
}

// WatchMilestoneReleased is a free log subscription operation binding the contract event 0xd33fcd27cbfb73a2d9058c968b560ccb306500bf536133802048489ecc7355e9.
//
// Solidity: event MilestoneReleased(uint256 indexed campaignId, uint256 indexed milestoneIndex, uint256 amount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) WatchMilestoneReleased(opts *bind.WatchOpts, sink chan<- *TrustChainEscrowMilestoneReleased, campaignId []*big.Int, milestoneIndex []*big.Int) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var milestoneIndexRule []interface{}
	for _, milestoneIndexItem := range milestoneIndex {
		milestoneIndexRule = append(milestoneIndexRule, milestoneIndexItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.WatchLogs(opts, "MilestoneReleased", campaignIdRule, milestoneIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustChainEscrowMilestoneReleased)
				if err := _TrustChainEscrow.contract.UnpackLog(event, "MilestoneReleased", log); err != nil {
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

// ParseMilestoneReleased is a log parse operation binding the contract event 0xd33fcd27cbfb73a2d9058c968b560ccb306500bf536133802048489ecc7355e9.
//
// Solidity: event MilestoneReleased(uint256 indexed campaignId, uint256 indexed milestoneIndex, uint256 amount)
func (_TrustChainEscrow *TrustChainEscrowFilterer) ParseMilestoneReleased(log types.Log) (*TrustChainEscrowMilestoneReleased, error) {
	event := new(TrustChainEscrowMilestoneReleased)
	if err := _TrustChainEscrow.contract.UnpackLog(event, "MilestoneReleased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrustChainEscrowReceiptSubmittedIterator is returned from FilterReceiptSubmitted and is used to iterate over the raw logs and unpacked data for ReceiptSubmitted events raised by the TrustChainEscrow contract.
type TrustChainEscrowReceiptSubmittedIterator struct {
	Event *TrustChainEscrowReceiptSubmitted // Event containing the contract specifics and raw log

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
func (it *TrustChainEscrowReceiptSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustChainEscrowReceiptSubmitted)
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
		it.Event = new(TrustChainEscrowReceiptSubmitted)
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
func (it *TrustChainEscrowReceiptSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustChainEscrowReceiptSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustChainEscrowReceiptSubmitted represents a ReceiptSubmitted event raised by the TrustChainEscrow contract.
type TrustChainEscrowReceiptSubmitted struct {
	CampaignId     *big.Int
	MilestoneIndex *big.Int
	ReceiptHash    [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterReceiptSubmitted is a free log retrieval operation binding the contract event 0x28f27f74dfbf86aa36d83238ef1ebf00ef6fd5911afa22a3d7f62e8a43282656.
//
// Solidity: event ReceiptSubmitted(uint256 indexed campaignId, uint256 indexed milestoneIndex, bytes32 receiptHash)
func (_TrustChainEscrow *TrustChainEscrowFilterer) FilterReceiptSubmitted(opts *bind.FilterOpts, campaignId []*big.Int, milestoneIndex []*big.Int) (*TrustChainEscrowReceiptSubmittedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var milestoneIndexRule []interface{}
	for _, milestoneIndexItem := range milestoneIndex {
		milestoneIndexRule = append(milestoneIndexRule, milestoneIndexItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.FilterLogs(opts, "ReceiptSubmitted", campaignIdRule, milestoneIndexRule)
	if err != nil {
		return nil, err
	}
	return &TrustChainEscrowReceiptSubmittedIterator{contract: _TrustChainEscrow.contract, event: "ReceiptSubmitted", logs: logs, sub: sub}, nil
}

// WatchReceiptSubmitted is a free log subscription operation binding the contract event 0x28f27f74dfbf86aa36d83238ef1ebf00ef6fd5911afa22a3d7f62e8a43282656.
//
// Solidity: event ReceiptSubmitted(uint256 indexed campaignId, uint256 indexed milestoneIndex, bytes32 receiptHash)
func (_TrustChainEscrow *TrustChainEscrowFilterer) WatchReceiptSubmitted(opts *bind.WatchOpts, sink chan<- *TrustChainEscrowReceiptSubmitted, campaignId []*big.Int, milestoneIndex []*big.Int) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var milestoneIndexRule []interface{}
	for _, milestoneIndexItem := range milestoneIndex {
		milestoneIndexRule = append(milestoneIndexRule, milestoneIndexItem)
	}

	logs, sub, err := _TrustChainEscrow.contract.WatchLogs(opts, "ReceiptSubmitted", campaignIdRule, milestoneIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustChainEscrowReceiptSubmitted)
				if err := _TrustChainEscrow.contract.UnpackLog(event, "ReceiptSubmitted", log); err != nil {
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

// ParseReceiptSubmitted is a log parse operation binding the contract event 0x28f27f74dfbf86aa36d83238ef1ebf00ef6fd5911afa22a3d7f62e8a43282656.
//
// Solidity: event ReceiptSubmitted(uint256 indexed campaignId, uint256 indexed milestoneIndex, bytes32 receiptHash)
func (_TrustChainEscrow *TrustChainEscrowFilterer) ParseReceiptSubmitted(log types.Log) (*TrustChainEscrowReceiptSubmitted, error) {
	event := new(TrustChainEscrowReceiptSubmitted)
	if err := _TrustChainEscrow.contract.UnpackLog(event, "ReceiptSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
