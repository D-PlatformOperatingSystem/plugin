// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"encoding/json"
	"io"
	"math/big"
	"time"

	"github.com/holiman/uint256"

	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/mm"
)

// Tracer                   。
// CaptureState   EVM         。
//       ，            ，     EVM    ；           ，      。
type Tracer interface {
	// CaptureStart
	CaptureStart(from common.Address, to common.Address, call bool, input []byte, gas uint64, value uint64) error
	// CaptureState
	CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *mm.Memory, stack *mm.Stack, rStack *mm.ReturnStack, rData []byte, contract *Contract, depth int, err error) error
	// CaptureFault
	CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *mm.Memory, stack *mm.Stack, rStack *mm.ReturnStack, contract *Contract, depth int, err error) error
	// CaptureEnd
	CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error
}

// JSONLogger   json
type JSONLogger struct {
	encoder *json.Encoder
}

// Storage represents a contract's storage.
type Storage map[common.Hash]common.Hash

// Copy duplicates the current storage.
func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}
	return cpy
}

// LogConfig are the configuration options for structured logger the EVM
type LogConfig struct {
	DisableMemory     bool // disable memory capture
	DisableStack      bool // disable stack capture
	DisableStorage    bool // disable storage capture
	DisableReturnData bool // disable return data capture
	Debug             bool // print output during capture end
	Limit             int  // maximum length of output, but zero means unlimited
}

// StructLog
type StructLog struct {
	// Pc pc
	Pc uint64 `json:"pc"`
	// Op
	Op OpCode `json:"op"`
	// Gas gas
	Gas uint64 `json:"gas"`
	// GasCost
	GasCost uint64 `json:"gasCost"`
	// Memory
	Memory []string `json:"memory"`
	// MemorySize
	MemorySize int `json:"memSize"`
	// Stack
	Stack []*big.Int `json:"stack"`
	// ReturnStack
	ReturnStack []uint32 `json:"returnStack"`
	// ReturnData
	ReturnData []byte `json:"returnData"`
	// Storage
	Storage map[common.Hash]common.Hash `json:"-"`
	// Depth
	Depth int `json:"depth"`
	// RefundCounter
	RefundCounter uint64 `json:"refund"`
	// Err
	Err error `json:"-"`
}

// NewJSONLogger
func NewJSONLogger(writer io.Writer) *JSONLogger {
	return &JSONLogger{json.NewEncoder(writer)}
}

// CaptureStart
func (logger *JSONLogger) CaptureStart(from common.Address, to common.Address, create bool, input []byte, gas uint64, value uint64) error {
	return nil
}

// CaptureState
func (logger *JSONLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *mm.Memory, stack *mm.Stack, rStack *mm.ReturnStack, rData []byte, contract *Contract, depth int, err error) error {
	log := StructLog{
		Pc:         pc,
		Op:         op,
		Gas:        gas,
		GasCost:    cost,
		MemorySize: memory.Len(),
		Storage:    nil,
		Depth:      depth,
		Err:        err,
	}
	log.Memory = formatMemory(memory.Data())
	log.Stack = formatStack(stack.Data())
	log.ReturnStack = rStack.Data()
	log.ReturnData = rData
	return logger.encoder.Encode(log)
}

func formatStack(data []uint256.Int) (res []*big.Int) {
	for _, v := range data {
		res = append(res, v.ToBig())
	}
	return
}

func formatMemory(data []byte) (res []string) {
	for idx := 0; idx < len(data); idx += 32 {
		res = append(res, common.Bytes2HexTrim(data[idx:idx+32]))
	}
	return
}

//CaptureFault
func (logger *JSONLogger) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *mm.Memory, stack *mm.Stack, rStack *mm.ReturnStack, contract *Contract, depth int, err error) error {
	return nil
}

// CaptureEnd
func (logger *JSONLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	type endLog struct {
		Output  string        `json:"output"`
		GasUsed int64         `json:"gasUsed"`
		Time    time.Duration `json:"time"`
		Err     string        `json:"error,omitempty"`
	}

	if err != nil {
		return logger.encoder.Encode(endLog{common.Bytes2Hex(output), int64(gasUsed), t, err.Error()})
	}
	return logger.encoder.Encode(endLog{common.Bytes2Hex(output), int64(gasUsed), t, ""})
}
