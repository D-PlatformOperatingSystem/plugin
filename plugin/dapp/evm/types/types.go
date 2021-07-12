// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/types"
)

const (
	// EvmCreateAction
	EvmCreateAction = 1
	// EvmCallAction
	EvmCallAction = 2

	// TyLogContractData
	TyLogContractData = 601
	// TyLogContractState
	TyLogContractState = 602
	// TyLogCallContract
	TyLogCallContract = 603
	// TyLogEVMStateChangeItem
	TyLogEVMStateChangeItem = 604

	// MaxGasLimit    Gas
	MaxGasLimit = 10000000
)

const (
	// EVMEnable   EVM
	EVMEnable = "Enable"
	// ForkEVMState EVM          ，
	ForkEVMState = "ForkEVMState"
	// ForkEVMKVHash EVM          ，      StateDB
	ForkEVMKVHash = "ForkEVMKVHash"
	// ForkEVMABI EVM    ABI
	ForkEVMABI = "ForkEVMABI"
	// ForkEVMFrozen EVM
	ForkEVMFrozen = "ForkEVMFrozen"
	// ForkEVMYoloV1 YoloV1
	ForkEVMYoloV1 = "ForkEVMYoloV1"
	//ForkEVMTxGroup          GAS
	ForkEVMTxGroup = "ForkEVMTxGroup"
)

var (
	// EvmPrefix
	EvmPrefix = "user.evm."
	// ExecutorName
	ExecutorName = "evm"

	// ExecerEvm EVM
	ExecerEvm = []byte(ExecutorName)
	// UserPrefix
	UserPrefix = []byte(EvmPrefix)

	logInfo = map[int64]*types.LogInfo{
		TyLogCallContract:       {Ty: reflect.TypeOf(ReceiptEVMContract{}), Name: "LogCallContract"},
		TyLogContractData:       {Ty: reflect.TypeOf(EVMContractData{}), Name: "LogContractData"},
		TyLogContractState:      {Ty: reflect.TypeOf(EVMContractState{}), Name: "LogContractState"},
		TyLogEVMStateChangeItem: {Ty: reflect.TypeOf(EVMStateChangeItem{}), Name: "LogEVMStateChangeItem"},
	}
)
