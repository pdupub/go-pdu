// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package poster

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ContractABI is the input ABI used to generate the binding from.
const ContractABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_recordID\",\"type\":\"uint256\"}],\"name\":\"getRecord\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_minDepositVal\",\"type\":\"uint256\"},{\"name\":\"_minDepositLen\",\"type\":\"uint256\"}],\"name\":\"setMinDeposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_recordID\",\"type\":\"uint256\"},{\"name\":\"_info\",\"type\":\"string\"}],\"name\":\"updateRecord\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"records\",\"outputs\":[{\"name\":\"info\",\"type\":\"string\"},{\"name\":\"author\",\"type\":\"address\"},{\"name\":\"deposit\",\"type\":\"uint256\"},{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"update\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_notice\",\"type\":\"string\"}],\"name\":\"setNotice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"termination\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minDepositVal\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_recordID\",\"type\":\"uint256\"},{\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"getStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_recordID\",\"type\":\"uint256\"},{\"name\":\"_val\",\"type\":\"uint256\"}],\"name\":\"fadeRecord\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextRecordID\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"terminateSys\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"notice\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_recordID\",\"type\":\"uint256\"}],\"name\":\"propRecord\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_recordID\",\"type\":\"uint256\"}],\"name\":\"delRecord\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_info\",\"type\":\"string\"}],\"name\":\"addRecord\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxDepositLen\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minDepositLen\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"Terminate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"author\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"info\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"}],\"name\":\"AddRecord\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnerSet\",\"type\":\"event\"}]"

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_Contract *ContractCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "getOwner")
	return *ret0, err
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_Contract *ContractSession) GetOwner() (common.Address, error) {
	return _Contract.Contract.GetOwner(&_Contract.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_Contract *ContractCallerSession) GetOwner() (common.Address, error) {
	return _Contract.Contract.GetOwner(&_Contract.CallOpts)
}

// GetRecord is a free data retrieval call binding the contract method 0x03e9e609.
//
// Solidity: function getRecord(uint256 _recordID) constant returns(string, address, uint256, uint256, uint256, uint256, address[])
func (_Contract *ContractCaller) GetRecord(opts *bind.CallOpts, _recordID *big.Int) (string, common.Address, *big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	var (
		ret0 = new(string)
		ret1 = new(common.Address)
		ret2 = new(*big.Int)
		ret3 = new(*big.Int)
		ret4 = new(*big.Int)
		ret5 = new(*big.Int)
		ret6 = new([]common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
		ret5,
		ret6,
	}
	err := _Contract.contract.Call(opts, out, "getRecord", _recordID)
	return *ret0, *ret1, *ret2, *ret3, *ret4, *ret5, *ret6, err
}

// GetRecord is a free data retrieval call binding the contract method 0x03e9e609.
//
// Solidity: function getRecord(uint256 _recordID) constant returns(string, address, uint256, uint256, uint256, uint256, address[])
func (_Contract *ContractSession) GetRecord(_recordID *big.Int) (string, common.Address, *big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _Contract.Contract.GetRecord(&_Contract.CallOpts, _recordID)
}

// GetRecord is a free data retrieval call binding the contract method 0x03e9e609.
//
// Solidity: function getRecord(uint256 _recordID) constant returns(string, address, uint256, uint256, uint256, uint256, address[])
func (_Contract *ContractCallerSession) GetRecord(_recordID *big.Int) (string, common.Address, *big.Int, *big.Int, *big.Int, *big.Int, []common.Address, error) {
	return _Contract.Contract.GetRecord(&_Contract.CallOpts, _recordID)
}

// GetStake is a free data retrieval call binding the contract method 0x68c5805e.
//
// Solidity: function getStake(uint256 _recordID, address _user) constant returns(uint256, uint256)
func (_Contract *ContractCaller) GetStake(opts *bind.CallOpts, _recordID *big.Int, _user common.Address) (*big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Contract.contract.Call(opts, out, "getStake", _recordID, _user)
	return *ret0, *ret1, err
}

// GetStake is a free data retrieval call binding the contract method 0x68c5805e.
//
// Solidity: function getStake(uint256 _recordID, address _user) constant returns(uint256, uint256)
func (_Contract *ContractSession) GetStake(_recordID *big.Int, _user common.Address) (*big.Int, *big.Int, error) {
	return _Contract.Contract.GetStake(&_Contract.CallOpts, _recordID, _user)
}

// GetStake is a free data retrieval call binding the contract method 0x68c5805e.
//
// Solidity: function getStake(uint256 _recordID, address _user) constant returns(uint256, uint256)
func (_Contract *ContractCallerSession) GetStake(_recordID *big.Int, _user common.Address) (*big.Int, *big.Int, error) {
	return _Contract.Contract.GetStake(&_Contract.CallOpts, _recordID, _user)
}

// MaxDepositLen is a free data retrieval call binding the contract method 0xdc05a1af.
//
// Solidity: function maxDepositLen() constant returns(uint256)
func (_Contract *ContractCaller) MaxDepositLen(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "maxDepositLen")
	return *ret0, err
}

// MaxDepositLen is a free data retrieval call binding the contract method 0xdc05a1af.
//
// Solidity: function maxDepositLen() constant returns(uint256)
func (_Contract *ContractSession) MaxDepositLen() (*big.Int, error) {
	return _Contract.Contract.MaxDepositLen(&_Contract.CallOpts)
}

// MaxDepositLen is a free data retrieval call binding the contract method 0xdc05a1af.
//
// Solidity: function maxDepositLen() constant returns(uint256)
func (_Contract *ContractCallerSession) MaxDepositLen() (*big.Int, error) {
	return _Contract.Contract.MaxDepositLen(&_Contract.CallOpts)
}

// MinDepositLen is a free data retrieval call binding the contract method 0xf36f76ba.
//
// Solidity: function minDepositLen() constant returns(uint256)
func (_Contract *ContractCaller) MinDepositLen(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minDepositLen")
	return *ret0, err
}

// MinDepositLen is a free data retrieval call binding the contract method 0xf36f76ba.
//
// Solidity: function minDepositLen() constant returns(uint256)
func (_Contract *ContractSession) MinDepositLen() (*big.Int, error) {
	return _Contract.Contract.MinDepositLen(&_Contract.CallOpts)
}

// MinDepositLen is a free data retrieval call binding the contract method 0xf36f76ba.
//
// Solidity: function minDepositLen() constant returns(uint256)
func (_Contract *ContractCallerSession) MinDepositLen() (*big.Int, error) {
	return _Contract.Contract.MinDepositLen(&_Contract.CallOpts)
}

// MinDepositVal is a free data retrieval call binding the contract method 0x6890a153.
//
// Solidity: function minDepositVal() constant returns(uint256)
func (_Contract *ContractCaller) MinDepositVal(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minDepositVal")
	return *ret0, err
}

// MinDepositVal is a free data retrieval call binding the contract method 0x6890a153.
//
// Solidity: function minDepositVal() constant returns(uint256)
func (_Contract *ContractSession) MinDepositVal() (*big.Int, error) {
	return _Contract.Contract.MinDepositVal(&_Contract.CallOpts)
}

// MinDepositVal is a free data retrieval call binding the contract method 0x6890a153.
//
// Solidity: function minDepositVal() constant returns(uint256)
func (_Contract *ContractCallerSession) MinDepositVal() (*big.Int, error) {
	return _Contract.Contract.MinDepositVal(&_Contract.CallOpts)
}

// NextRecordID is a free data retrieval call binding the contract method 0x858a5a08.
//
// Solidity: function nextRecordID() constant returns(uint256)
func (_Contract *ContractCaller) NextRecordID(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "nextRecordID")
	return *ret0, err
}

// NextRecordID is a free data retrieval call binding the contract method 0x858a5a08.
//
// Solidity: function nextRecordID() constant returns(uint256)
func (_Contract *ContractSession) NextRecordID() (*big.Int, error) {
	return _Contract.Contract.NextRecordID(&_Contract.CallOpts)
}

// NextRecordID is a free data retrieval call binding the contract method 0x858a5a08.
//
// Solidity: function nextRecordID() constant returns(uint256)
func (_Contract *ContractCallerSession) NextRecordID() (*big.Int, error) {
	return _Contract.Contract.NextRecordID(&_Contract.CallOpts)
}

// Notice is a free data retrieval call binding the contract method 0x9c94e6c6.
//
// Solidity: function notice() constant returns(string)
func (_Contract *ContractCaller) Notice(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "notice")
	return *ret0, err
}

// Notice is a free data retrieval call binding the contract method 0x9c94e6c6.
//
// Solidity: function notice() constant returns(string)
func (_Contract *ContractSession) Notice() (string, error) {
	return _Contract.Contract.Notice(&_Contract.CallOpts)
}

// Notice is a free data retrieval call binding the contract method 0x9c94e6c6.
//
// Solidity: function notice() constant returns(string)
func (_Contract *ContractCallerSession) Notice() (string, error) {
	return _Contract.Contract.Notice(&_Contract.CallOpts)
}

// Records is a free data retrieval call binding the contract method 0x34461067.
//
// Solidity: function records(uint256 ) constant returns(string info, address author, uint256 deposit, uint256 start, uint256 update)
func (_Contract *ContractCaller) Records(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Info    string
	Author  common.Address
	Deposit *big.Int
	Start   *big.Int
	Update  *big.Int
}, error) {
	ret := new(struct {
		Info    string
		Author  common.Address
		Deposit *big.Int
		Start   *big.Int
		Update  *big.Int
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "records", arg0)
	return *ret, err
}

// Records is a free data retrieval call binding the contract method 0x34461067.
//
// Solidity: function records(uint256 ) constant returns(string info, address author, uint256 deposit, uint256 start, uint256 update)
func (_Contract *ContractSession) Records(arg0 *big.Int) (struct {
	Info    string
	Author  common.Address
	Deposit *big.Int
	Start   *big.Int
	Update  *big.Int
}, error) {
	return _Contract.Contract.Records(&_Contract.CallOpts, arg0)
}

// Records is a free data retrieval call binding the contract method 0x34461067.
//
// Solidity: function records(uint256 ) constant returns(string info, address author, uint256 deposit, uint256 start, uint256 update)
func (_Contract *ContractCallerSession) Records(arg0 *big.Int) (struct {
	Info    string
	Author  common.Address
	Deposit *big.Int
	Start   *big.Int
	Update  *big.Int
}, error) {
	return _Contract.Contract.Records(&_Contract.CallOpts, arg0)
}

// Termination is a free data retrieval call binding the contract method 0x40734387.
//
// Solidity: function termination() constant returns(uint256)
func (_Contract *ContractCaller) Termination(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "termination")
	return *ret0, err
}

// Termination is a free data retrieval call binding the contract method 0x40734387.
//
// Solidity: function termination() constant returns(uint256)
func (_Contract *ContractSession) Termination() (*big.Int, error) {
	return _Contract.Contract.Termination(&_Contract.CallOpts)
}

// Termination is a free data retrieval call binding the contract method 0x40734387.
//
// Solidity: function termination() constant returns(uint256)
func (_Contract *ContractCallerSession) Termination() (*big.Int, error) {
	return _Contract.Contract.Termination(&_Contract.CallOpts)
}

// AddRecord is a paid mutator transaction binding the contract method 0xd81aa8c4.
//
// Solidity: function addRecord(string _info) returns()
func (_Contract *ContractTransactor) AddRecord(opts *bind.TransactOpts, _info string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addRecord", _info)
}

// AddRecord is a paid mutator transaction binding the contract method 0xd81aa8c4.
//
// Solidity: function addRecord(string _info) returns()
func (_Contract *ContractSession) AddRecord(_info string) (*types.Transaction, error) {
	return _Contract.Contract.AddRecord(&_Contract.TransactOpts, _info)
}

// AddRecord is a paid mutator transaction binding the contract method 0xd81aa8c4.
//
// Solidity: function addRecord(string _info) returns()
func (_Contract *ContractTransactorSession) AddRecord(_info string) (*types.Transaction, error) {
	return _Contract.Contract.AddRecord(&_Contract.TransactOpts, _info)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Contract *ContractTransactor) ChangeOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "changeOwner", newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Contract *ContractSession) ChangeOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.ChangeOwner(&_Contract.TransactOpts, newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Contract *ContractTransactorSession) ChangeOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.ChangeOwner(&_Contract.TransactOpts, newOwner)
}

// DelRecord is a paid mutator transaction binding the contract method 0xc20c6d9d.
//
// Solidity: function delRecord(uint256 _recordID) returns()
func (_Contract *ContractTransactor) DelRecord(opts *bind.TransactOpts, _recordID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "delRecord", _recordID)
}

// DelRecord is a paid mutator transaction binding the contract method 0xc20c6d9d.
//
// Solidity: function delRecord(uint256 _recordID) returns()
func (_Contract *ContractSession) DelRecord(_recordID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DelRecord(&_Contract.TransactOpts, _recordID)
}

// DelRecord is a paid mutator transaction binding the contract method 0xc20c6d9d.
//
// Solidity: function delRecord(uint256 _recordID) returns()
func (_Contract *ContractTransactorSession) DelRecord(_recordID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DelRecord(&_Contract.TransactOpts, _recordID)
}

// FadeRecord is a paid mutator transaction binding the contract method 0x69896183.
//
// Solidity: function fadeRecord(uint256 _recordID, uint256 _val) returns()
func (_Contract *ContractTransactor) FadeRecord(opts *bind.TransactOpts, _recordID *big.Int, _val *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "fadeRecord", _recordID, _val)
}

// FadeRecord is a paid mutator transaction binding the contract method 0x69896183.
//
// Solidity: function fadeRecord(uint256 _recordID, uint256 _val) returns()
func (_Contract *ContractSession) FadeRecord(_recordID *big.Int, _val *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.FadeRecord(&_Contract.TransactOpts, _recordID, _val)
}

// FadeRecord is a paid mutator transaction binding the contract method 0x69896183.
//
// Solidity: function fadeRecord(uint256 _recordID, uint256 _val) returns()
func (_Contract *ContractTransactorSession) FadeRecord(_recordID *big.Int, _val *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.FadeRecord(&_Contract.TransactOpts, _recordID, _val)
}

// PropRecord is a paid mutator transaction binding the contract method 0xae86a2d2.
//
// Solidity: function propRecord(uint256 _recordID) returns()
func (_Contract *ContractTransactor) PropRecord(opts *bind.TransactOpts, _recordID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "propRecord", _recordID)
}

// PropRecord is a paid mutator transaction binding the contract method 0xae86a2d2.
//
// Solidity: function propRecord(uint256 _recordID) returns()
func (_Contract *ContractSession) PropRecord(_recordID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PropRecord(&_Contract.TransactOpts, _recordID)
}

// PropRecord is a paid mutator transaction binding the contract method 0xae86a2d2.
//
// Solidity: function propRecord(uint256 _recordID) returns()
func (_Contract *ContractTransactorSession) PropRecord(_recordID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PropRecord(&_Contract.TransactOpts, _recordID)
}

// SetMinDeposit is a paid mutator transaction binding the contract method 0x141028b5.
//
// Solidity: function setMinDeposit(uint256 _minDepositVal, uint256 _minDepositLen) returns()
func (_Contract *ContractTransactor) SetMinDeposit(opts *bind.TransactOpts, _minDepositVal *big.Int, _minDepositLen *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setMinDeposit", _minDepositVal, _minDepositLen)
}

// SetMinDeposit is a paid mutator transaction binding the contract method 0x141028b5.
//
// Solidity: function setMinDeposit(uint256 _minDepositVal, uint256 _minDepositLen) returns()
func (_Contract *ContractSession) SetMinDeposit(_minDepositVal *big.Int, _minDepositLen *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetMinDeposit(&_Contract.TransactOpts, _minDepositVal, _minDepositLen)
}

// SetMinDeposit is a paid mutator transaction binding the contract method 0x141028b5.
//
// Solidity: function setMinDeposit(uint256 _minDepositVal, uint256 _minDepositLen) returns()
func (_Contract *ContractTransactorSession) SetMinDeposit(_minDepositVal *big.Int, _minDepositLen *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetMinDeposit(&_Contract.TransactOpts, _minDepositVal, _minDepositLen)
}

// SetNotice is a paid mutator transaction binding the contract method 0x3cf572a7.
//
// Solidity: function setNotice(string _notice) returns()
func (_Contract *ContractTransactor) SetNotice(opts *bind.TransactOpts, _notice string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setNotice", _notice)
}

// SetNotice is a paid mutator transaction binding the contract method 0x3cf572a7.
//
// Solidity: function setNotice(string _notice) returns()
func (_Contract *ContractSession) SetNotice(_notice string) (*types.Transaction, error) {
	return _Contract.Contract.SetNotice(&_Contract.TransactOpts, _notice)
}

// SetNotice is a paid mutator transaction binding the contract method 0x3cf572a7.
//
// Solidity: function setNotice(string _notice) returns()
func (_Contract *ContractTransactorSession) SetNotice(_notice string) (*types.Transaction, error) {
	return _Contract.Contract.SetNotice(&_Contract.TransactOpts, _notice)
}

// TerminateSys is a paid mutator transaction binding the contract method 0x91e0cdb3.
//
// Solidity: function terminateSys() returns()
func (_Contract *ContractTransactor) TerminateSys(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "terminateSys")
}

// TerminateSys is a paid mutator transaction binding the contract method 0x91e0cdb3.
//
// Solidity: function terminateSys() returns()
func (_Contract *ContractSession) TerminateSys() (*types.Transaction, error) {
	return _Contract.Contract.TerminateSys(&_Contract.TransactOpts)
}

// TerminateSys is a paid mutator transaction binding the contract method 0x91e0cdb3.
//
// Solidity: function terminateSys() returns()
func (_Contract *ContractTransactorSession) TerminateSys() (*types.Transaction, error) {
	return _Contract.Contract.TerminateSys(&_Contract.TransactOpts)
}

// UpdateRecord is a paid mutator transaction binding the contract method 0x1c631466.
//
// Solidity: function updateRecord(uint256 _recordID, string _info) returns()
func (_Contract *ContractTransactor) UpdateRecord(opts *bind.TransactOpts, _recordID *big.Int, _info string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateRecord", _recordID, _info)
}

// UpdateRecord is a paid mutator transaction binding the contract method 0x1c631466.
//
// Solidity: function updateRecord(uint256 _recordID, string _info) returns()
func (_Contract *ContractSession) UpdateRecord(_recordID *big.Int, _info string) (*types.Transaction, error) {
	return _Contract.Contract.UpdateRecord(&_Contract.TransactOpts, _recordID, _info)
}

// UpdateRecord is a paid mutator transaction binding the contract method 0x1c631466.
//
// Solidity: function updateRecord(uint256 _recordID, string _info) returns()
func (_Contract *ContractTransactorSession) UpdateRecord(_recordID *big.Int, _info string) (*types.Transaction, error) {
	return _Contract.Contract.UpdateRecord(&_Contract.TransactOpts, _recordID, _info)
}

// ContractAddRecordIterator is returned from FilterAddRecord and is used to iterate over the raw logs and unpacked data for AddRecord events raised by the Contract contract.
type ContractAddRecordIterator struct {
	Event *ContractAddRecord // Event containing the contract specifics and raw log

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
func (it *ContractAddRecordIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAddRecord)
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
		it.Event = new(ContractAddRecord)
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
func (it *ContractAddRecordIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAddRecordIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAddRecord represents a AddRecord event raised by the Contract contract.
type ContractAddRecord struct {
	Author  common.Address
	Info    string
	Deposit *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddRecord is a free log retrieval operation binding the contract event 0x7807102e7be173ae03e3bf3bfa7cac55a6d8fcca979058fef5a9240537b76bfa.
//
// Solidity: event AddRecord(address indexed author, string info, uint256 deposit)
func (_Contract *ContractFilterer) FilterAddRecord(opts *bind.FilterOpts, author []common.Address) (*ContractAddRecordIterator, error) {

	var authorRule []interface{}
	for _, authorItem := range author {
		authorRule = append(authorRule, authorItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "AddRecord", authorRule)
	if err != nil {
		return nil, err
	}
	return &ContractAddRecordIterator{contract: _Contract.contract, event: "AddRecord", logs: logs, sub: sub}, nil
}

// WatchAddRecord is a free log subscription operation binding the contract event 0x7807102e7be173ae03e3bf3bfa7cac55a6d8fcca979058fef5a9240537b76bfa.
//
// Solidity: event AddRecord(address indexed author, string info, uint256 deposit)
func (_Contract *ContractFilterer) WatchAddRecord(opts *bind.WatchOpts, sink chan<- *ContractAddRecord, author []common.Address) (event.Subscription, error) {

	var authorRule []interface{}
	for _, authorItem := range author {
		authorRule = append(authorRule, authorItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "AddRecord", authorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAddRecord)
				if err := _Contract.contract.UnpackLog(event, "AddRecord", log); err != nil {
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

// ParseAddRecord is a log parse operation binding the contract event 0x7807102e7be173ae03e3bf3bfa7cac55a6d8fcca979058fef5a9240537b76bfa.
//
// Solidity: event AddRecord(address indexed author, string info, uint256 deposit)
func (_Contract *ContractFilterer) ParseAddRecord(log types.Log) (*ContractAddRecord, error) {
	event := new(ContractAddRecord)
	if err := _Contract.contract.UnpackLog(event, "AddRecord", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractOwnerSetIterator is returned from FilterOwnerSet and is used to iterate over the raw logs and unpacked data for OwnerSet events raised by the Contract contract.
type ContractOwnerSetIterator struct {
	Event *ContractOwnerSet // Event containing the contract specifics and raw log

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
func (it *ContractOwnerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOwnerSet)
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
		it.Event = new(ContractOwnerSet)
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
func (it *ContractOwnerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOwnerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOwnerSet represents a OwnerSet event raised by the Contract contract.
type ContractOwnerSet struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnerSet is a free log retrieval operation binding the contract event 0x342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a735.
//
// Solidity: event OwnerSet(address indexed oldOwner, address indexed newOwner)
func (_Contract *ContractFilterer) FilterOwnerSet(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*ContractOwnerSetIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OwnerSet", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractOwnerSetIterator{contract: _Contract.contract, event: "OwnerSet", logs: logs, sub: sub}, nil
}

// WatchOwnerSet is a free log subscription operation binding the contract event 0x342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a735.
//
// Solidity: event OwnerSet(address indexed oldOwner, address indexed newOwner)
func (_Contract *ContractFilterer) WatchOwnerSet(opts *bind.WatchOpts, sink chan<- *ContractOwnerSet, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OwnerSet", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOwnerSet)
				if err := _Contract.contract.UnpackLog(event, "OwnerSet", log); err != nil {
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

// ParseOwnerSet is a log parse operation binding the contract event 0x342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a735.
//
// Solidity: event OwnerSet(address indexed oldOwner, address indexed newOwner)
func (_Contract *ContractFilterer) ParseOwnerSet(log types.Log) (*ContractOwnerSet, error) {
	event := new(ContractOwnerSet)
	if err := _Contract.contract.UnpackLog(event, "OwnerSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractTerminateIterator is returned from FilterTerminate and is used to iterate over the raw logs and unpacked data for Terminate events raised by the Contract contract.
type ContractTerminateIterator struct {
	Event *ContractTerminate // Event containing the contract specifics and raw log

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
func (it *ContractTerminateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractTerminate)
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
		it.Event = new(ContractTerminate)
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
func (it *ContractTerminateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractTerminateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractTerminate represents a Terminate event raised by the Contract contract.
type ContractTerminate struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTerminate is a free log retrieval operation binding the contract event 0x22680ec819c18f669abb4fa84dfd5fe059ea6d707972b79ba5db1d9ef1e531cb.
//
// Solidity: event Terminate(address indexed owner)
func (_Contract *ContractFilterer) FilterTerminate(opts *bind.FilterOpts, owner []common.Address) (*ContractTerminateIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "Terminate", ownerRule)
	if err != nil {
		return nil, err
	}
	return &ContractTerminateIterator{contract: _Contract.contract, event: "Terminate", logs: logs, sub: sub}, nil
}

// WatchTerminate is a free log subscription operation binding the contract event 0x22680ec819c18f669abb4fa84dfd5fe059ea6d707972b79ba5db1d9ef1e531cb.
//
// Solidity: event Terminate(address indexed owner)
func (_Contract *ContractFilterer) WatchTerminate(opts *bind.WatchOpts, sink chan<- *ContractTerminate, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "Terminate", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractTerminate)
				if err := _Contract.contract.UnpackLog(event, "Terminate", log); err != nil {
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

// ParseTerminate is a log parse operation binding the contract event 0x22680ec819c18f669abb4fa84dfd5fe059ea6d707972b79ba5db1d9ef1e531cb.
//
// Solidity: event Terminate(address indexed owner)
func (_Contract *ContractFilterer) ParseTerminate(log types.Log) (*ContractTerminate, error) {
	event := new(ContractTerminate)
	if err := _Contract.contract.UnpackLog(event, "Terminate", log); err != nil {
		return nil, err
	}
	return event, nil
}
