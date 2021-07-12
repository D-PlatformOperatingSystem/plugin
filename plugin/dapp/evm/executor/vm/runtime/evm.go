// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"math/big"
	"sync/atomic"

	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/gas"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/model"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/params"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/state"
	evmtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"
)

type (
	// CanTransferFunc
	CanTransferFunc func(state.EVMStateDB, common.Address, common.Address, uint64) bool

	// TransferFunc
	TransferFunc func(state.EVMStateDB, common.Address, common.Address, uint64) bool

	// GetHashFunc
	//   BLOCKHASH
	GetHashFunc func(uint64) common.Hash
)

//                 ，    ，
func run(evm *EVM, contract *Contract, input []byte) (ret []byte, err error) {
	if contract.CodeAddr != nil {
		//                 ，      ，
		precompiles := PrecompiledContractsByzantium
		//       ： dplatformos                v1  （         ）
		if evm.cfg.IsDappFork(evm.StateDB.GetBlockHeight(), "evm", evmtypes.ForkEVMYoloV1) {
			precompiles = PrecompiledContractsYoloV1
		}
		if p := precompiles[*contract.CodeAddr]; p != nil {
			return RunPrecompiledContract(p, input, contract)
		}
	}
	//
	ret, err = evm.Interpreter.Run(contract, input)
	if err != nil {
		log.Error("error occurs while run evm contract", "error info", err)
	}

	return ret, err
}

// Context EVM
//      EVM     ，EVM
type Context struct {

	// CanTransfer
	CanTransfer CanTransferFunc
	// Transfer
	Transfer TransferFunc
	// GetHash
	GetHash GetHashFunc

	// Origin       ，
	Origin common.Address
	// GasPrice
	GasPrice uint32

	// Coinbase   ，
	Coinbase *common.Address
	// GasLimit   ，     GasLimit
	GasLimit uint64

	// BlockNumber NUMBER   ，
	BlockNumber *big.Int
	// Time   ，
	Time *big.Int
	// Difficulty   ，
	Difficulty *big.Int
}

// EVM              ，         EVM
//                 （          ）
//                  ，                ，          Gas
//               ，
type EVM struct {
	// Context
	Context
	// EVMStateDB
	StateDB state.EVMStateDB
	//
	depth int

	// VMConfig
	VMConfig Config

	// Interpreter EVM     ，     EVM
	Interpreter *Interpreter

	// EVM
	//
	abort int32

	// CallGasTemp               Gas
	//       ，      gasCost  ，           Gas，
	//      opCall ，         Gas
	CallGasTemp uint64

	//
	maxCodeSize int

	// dplatformos
	cfg *types.DplatformOSConfig
}

// NewEVM       EVM
//        ，  EVM
func NewEVM(ctx Context, statedb state.EVMStateDB, vmConfig Config, cfg *types.DplatformOSConfig) *EVM {
	evm := &EVM{
		Context:     ctx,
		StateDB:     statedb,
		VMConfig:    vmConfig,
		maxCodeSize: params.MaxCodeSize,
		cfg:         cfg,
	}

	evm.Interpreter = NewInterpreter(evm, vmConfig)
	return evm
}

// GasTable          Gas
//           ，
func (evm *EVM) GasTable(num *big.Int) gas.Table {
	return gas.TableHomestead
}

// Cancel               EVM       ，
func (evm *EVM) Cancel() {
	atomic.StoreInt32(&evm.abort, 1)
}

// SetMaxCodeSize
func (evm *EVM) SetMaxCodeSize(maxCodeSize int) {
	if maxCodeSize < 1 || maxCodeSize > params.MaxCodeSize {
		return
	}

	evm.maxCodeSize = maxCodeSize
}

//
func (evm *EVM) preCheck(caller ContractRef, recipient common.Address, value uint64) (pass bool, err error) {
	//
	if evm.VMConfig.NoRecursion && evm.depth > 0 {
		return false, nil
	}

	//     ，
	if evm.depth > int(params.CallCreateDepth) {
		return false, model.ErrDepth
	}

	//     ，
	if value > 0 {
		if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), recipient, value) {
			return false, model.ErrInsufficientBalance
		}
	}

	//        ，
	return true, nil
}

// Call
//                ，input
//
func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value uint64) (ret []byte, snapshot int, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, addr, value)
	if !pass {
		return nil, -1, gas, err
	}

	if !evm.StateDB.Exist(addr.String()) {
		//       ： dplatformos                v1  （         ）
		precompiles := PrecompiledContractsByzantium
		//       v1
		if evm.cfg.IsDappFork(evm.StateDB.GetBlockHeight(), "evm", evmtypes.ForkEVMYoloV1) {
			precompiles = PrecompiledContractsYoloV1
		}
		//                       ，
		if precompiles[addr] == nil {
			//             ，
			if len(input) > 0 || value == 0 {
				//             ，
				if evm.VMConfig.Debug && evm.depth == 0 {
					evm.VMConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)
					evm.VMConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
				}
				return nil, -1, gas, model.ErrAddrNotExists
			}
		} else {
			//   ，      ，
			//       ，                      ，                 ，
			// evm.EVMStateDB.CreateAccount(addr, caller.Address())
		}
	}

	//
	if evm.StateDB.HasSuicided(addr.String()) {
		return nil, -1, gas, model.ErrDestruct
	}

	//    ，
	snapshot = evm.StateDB.Snapshot()
	to := AccountRef(addr)

	//
	evm.Transfer(evm.StateDB, caller.Address(), to.Address(), value)
	log.Info("evm call", "caller address", caller.Address().String(), "contract address", to.Address().String(), "value", value)
	//         ，            ，  Gas
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr.String()), evm.StateDB.GetCode(addr.String()))

	start := types.Now()

	//
	if evm.VMConfig.Debug && evm.depth == 0 {
		evm.VMConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)

		defer func() {
			evm.VMConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, types.Since(start), err)
		}()
	}

	//  ForkV20EVMState  ，          ，
	cfg := evm.StateDB.GetConfig()
	if cfg.IsDappFork(evm.BlockNumber.Int64(), "evm", evmtypes.ForkEVMState) {
		evm.StateDB.TransferStateData(addr.String())
	}

	ret, err = run(evm, contract, input)

	//         ，      （            ），         gas
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != model.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, snapshot, contract.Gas, err
}

// CallCode
//      Call  ，         ：
//         ，          （     self  ）    caller
func (evm *EVM) CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64, value uint64) (ret []byte, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, addr, value)
	if !pass {
		return nil, gas, err
	}

	//
	if evm.StateDB.HasSuicided(addr.String()) {
		return nil, gas, model.ErrDestruct
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	//        ，
	contract := NewContract(caller, to, value, gas)
	//
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr.String()), evm.StateDB.GetCode(addr.String()))

	ret, err = run(evm, contract, input)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != model.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// DelegateCall
//
//  CallCode    ，               caller caller
func (evm *EVM) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, addr, 0)
	if !pass {
		return nil, gas, err
	}

	//
	if evm.StateDB.HasSuicided(addr.String()) {
		return nil, gas, model.ErrDestruct
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	//              ，      ，
	//     ，      ，             （         ）
	contract := NewContract(caller, to, 0, gas).AsDelegate()
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr.String()), evm.StateDB.GetCode(addr.String()))

	//      StaticCall
	ret, err = run(evm, contract, input)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != model.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// StaticCall
//
//       ，                       ，  ，         MemoryStateDB      ，
func (evm *EVM) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, addr, 0)
	if !pass {
		return nil, gas, err
	}

	//
	if evm.StateDB.HasSuicided(addr.String()) {
		return nil, gas, model.ErrDestruct
	}

	//               ，         ，
	//                 ，      ，
	if !evm.Interpreter.readOnly {
		evm.Interpreter.readOnly = true
		defer func() { evm.Interpreter.readOnly = false }()
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)

	//              ，      ，
	contract := NewContract(caller, to, 0, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr.String()), evm.StateDB.GetCode(addr.String()))

	//            ，      ，       Gas
	ret, err = run(evm, contract, input)
	if err != nil {
		//
		//   ，               ，               ，
		evm.StateDB.RevertToSnapshot(snapshot)

		//          ，    ，      Gas
		if err != model.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// Create              ；
//                ；
//   dplatformos        ，                  ，
//   ，
func (evm *EVM) Create(caller ContractRef, contractAddr common.Address, code []byte, gas uint64, execName, alias, abi string) (ret []byte, snapshot int, leftOverGas uint64, err error) {
	pass, err := evm.preCheck(caller, contractAddr, 0)
	if !pass {
		return nil, -1, gas, err
	}

	//         ，            ，  Gas
	contract := NewContract(caller, AccountRef(contractAddr), 0, gas)
	contract.SetCallCode(&contractAddr, common.ToHash(code), code)

	//           （    ）
	snapshot = evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(contractAddr.String(), contract.CallerAddress.String(), execName, alias)

	if evm.VMConfig.Debug && evm.depth == 0 {
		evm.VMConfig.Tracer.CaptureStart(caller.Address(), contractAddr, true, code, gas, 0)
	}
	start := types.Now()

	//
	ret, err = run(evm, contract, nil)

	//
	maxCodeSizeExceeded := len(ret) > evm.maxCodeSize

	cfg := evm.StateDB.GetConfig()
	//       ，             Gas
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			evm.StateDB.SetCode(contractAddr.String(), ret)
			//    ABI (     )，
			if len(abi) > 0 && cfg.IsDappFork(evm.StateDB.GetBlockHeight(), "evm", evmtypes.ForkEVMABI) {
				evm.StateDB.SetAbi(contractAddr.String(), abi)
			}
		} else {
			//   Gas  ，      ，
			err = model.ErrCodeStoreOutOfGas
		}
	}

	//         ，     Gas
	//
	if maxCodeSizeExceeded || (err != nil && err != model.ErrCodeStoreOutOfGas) {
		evm.StateDB.RevertToSnapshot(snapshot)

		//         ，      ，   Gas
		if err != model.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	//             ，          ，            ，
	if maxCodeSizeExceeded && err == nil {
		err = model.ErrMaxCodeSizeExceeded
	}

	if evm.VMConfig.Debug && evm.depth == 0 {
		evm.VMConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, types.Since(start), err)
	}

	return ret, snapshot, contract.Gas, err
}
