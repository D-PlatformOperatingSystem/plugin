// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"fmt"
	"sync/atomic"

	evmtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"

	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common/math"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/gas"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/mm"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/model"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/params"
)

// Config
type Config struct {
	// Debug
	Debug bool
	// Tracer
	Tracer Tracer
	// NoRecursion      Call, CallCode, DelegateCall
	NoRecursion bool
	// EnablePreimageRecording SHA3/keccak
	EnablePreimageRecording bool
	// JumpTable
	JumpTable [256]Operation
}

// Interpreter
type Interpreter struct {
	evm      *EVM
	cfg      Config
	gasTable gas.Table
	//
	readOnly bool
	//
	ReturnData []byte
}

// NewInterpreter
func NewInterpreter(evm *EVM, cfg Config) *Interpreter {
	//          STOP    jump table
	//     ，        ，          ，
	if !cfg.JumpTable[STOP].Valid {
		cfg.JumpTable = ConstantinopleInstructionSet
		if evm.cfg.IsDappFork(evm.StateDB.GetBlockHeight(), "evm", evmtypes.ForkEVMYoloV1) {
			//
			cfg.JumpTable = YoloV1InstructionSet
		}
	}

	return &Interpreter{
		evm:      evm,
		cfg:      cfg,
		gasTable: evm.GasTable(evm.BlockNumber),
	}
}

func (in *Interpreter) enforceRestrictions(op OpCode, operation Operation, stack *mm.Stack) error {
	if in.readOnly {
		//               ，
		//           （           ）
		if operation.Writes || (op == CALL && stack.Back(2).BitLen() > 0) {
			return model.ErrWriteProtection
		}
	}
	return nil
}

// Run
//       ，        ，        Gas
// （      ErrExecutionReverted，           Gas）
func (in *Interpreter) Run(contract *Contract, input []byte) (ret []byte, err error) {
	//TODO           ,        ？
	//       ，   1
	in.evm.depth++
	defer func() { in.evm.depth-- }()

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This makes also sure that the readOnly flag isn't removed for child calls.
	//if !in.readOnly {
	//	in.readOnly = true
	//	defer func() { in.readOnly = false }()
	//}
	//
	in.ReturnData = nil

	//
	if len(contract.Code) == 0 {
		return nil, nil
	}

	var (
		//
		op OpCode
		//
		mem = mm.NewMemory()
		//
		stack = mm.NewStack()
		//
		returns     = mm.NewReturnStack() // local returns stack
		callContext = &callCtx{
			memory:   mem,
			stack:    stack,
			rstack:   returns,
			contract: contract,
		}
		//
		pc = uint64(0)
		//      Gas
		cost uint64
		//    tracer       ，
		pcCopy  uint64
		gasCopy uint64
		logged  bool
		//
		res []byte
	)
	contract.Input = input

	//      ，
	defer func() {
		mm.Returnstack(stack)
		mm.ReturnRStack(returns)
	}()

	if in.cfg.Debug {
		defer func() {
			if err != nil {
				if !logged {
					in.cfg.Tracer.CaptureState(in.evm, pcCopy, op, gasCopy, cost, mem, stack, returns, in.ReturnData, contract, in.evm.depth, err)
				} else {
					in.cfg.Tracer.CaptureFault(in.evm, pcCopy, op, gasCopy, cost, mem, stack, returns, contract, in.evm.depth, err)
				}
			}
		}()
	}
	//             ，        （  、  、  、  、  ）
	steps := 0
	for {
		steps++
		if steps%1000 == 0 && atomic.LoadInt32(&in.evm.abort) != 0 {
			break
		}
		if in.cfg.Debug {
			//
			logged, pcCopy, gasCopy = false, pc, contract.Gas
		}

		//
		op = contract.GetOp(pc)
		operation := in.cfg.JumpTable[op]
		if !operation.Valid {
			return nil, fmt.Errorf("invalid OpCode 0x%x", int(op))
		}
		if err := operation.ValidateStack(stack); err != nil {
			return nil, err
		}
		//
		if err := in.enforceRestrictions(op, operation, stack); err != nil {
			return nil, err
		}
		var memorySize uint64
		//
		if operation.MemorySize != nil {
			memSize, overflow := operation.MemorySize(stack)
			if overflow {
				return nil, model.ErrGasUintOverflow
			}
			// memory is expanded in words of 32 bytes. Gas
			// is also calculated in words.
			if memorySize, overflow = math.SafeMul(common.ToWordSize(memSize), 32); overflow {
				return nil, model.ErrGasUintOverflow
			}
		}
		//             Gas
		evmParam := buildEVMParam(in.evm)
		gasParam := buildGasParam(contract)
		cost, err = operation.GasCost(in.gasTable, evmParam, gasParam, stack, mem, memorySize)
		fillEVM(evmParam, in.evm)

		if err != nil || !contract.UseGas(cost) {
			return nil, model.ErrOutOfGas
		}
		if memorySize > 0 {
			//
			mem.Resize(memorySize)
		}

		if in.cfg.Debug {
			in.cfg.Tracer.CaptureState(in.evm, pc, op, gasCopy, cost, mem, stack, returns, in.ReturnData, contract, in.evm.depth, err)
			logged = true
		}

		//            （       ）
		res, err = operation.Execute(&pc, in.evm, callContext)
		//          ，
		if operation.Returns {
			in.ReturnData = common.CopyBytes(res)
		}

		switch {
		case err != nil:
			return nil, err
		case operation.Reverts:
			return res, model.ErrExecutionReverted
		case operation.Halts:
			return res, nil
		case !operation.Jumps:
			pc++
		}
	}
	return nil, nil
}

//  Contract       GasFunc
//             ，   GasFun  Gas       Contract
//        GasParam
func buildGasParam(contract *Contract) *params.GasParam {
	return &params.GasParam{Gas: contract.Gas, Address: contract.Address()}
}

//  EVM       GasFunc
//             ，   GasFun  Gas       EVM
//        EVMParam
func buildEVMParam(evm *EVM) *params.EVMParam {
	return &params.EVMParam{
		StateDB:     evm.StateDB,
		CallGasTemp: evm.CallGasTemp,
		BlockNumber: evm.BlockNumber,
	}
}

//           EVM
//       CallGasTemp，             ，         EVM
func fillEVM(param *params.EVMParam, evm *EVM) {
	evm.CallGasTemp = param.CallGasTemp
}

// callCtx contains the things that are per-call, such as stack and memory,
// but not transients like pc and gas
type callCtx struct {
	memory   *mm.Memory
	stack    *mm.Stack
	rstack   *mm.ReturnStack
	contract *Contract
}
