// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"bytes"
	"math/big"

	"os"

	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/runtime"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/state"
	evmtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"
)

var (
	evmDebug = false

	// EvmAddress
	EvmAddress = ""
)

var driverName = evmtypes.ExecutorName

// Init
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	driverName = name
	drivers.Register(cfg, driverName, newEVMDriver, cfg.GetDappFork(driverName, evmtypes.EVMEnable))
	EvmAddress = address.ExecAddress(cfg.ExecName(name))
	//
	state.InitForkData()
	InitExecType()
}

// InitExecType Init Exec Type
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&EVMExecutor{}))
}

// GetName
func GetName() string {
	return newEVMDriver().GetName()
}

func newEVMDriver() drivers.Driver {
	evm := NewEVMExecutor()
	evm.vmCfg.Debug = evmDebug
	return evm
}

// EVMExecutor EVM
type EVMExecutor struct {
	drivers.DriverBase
	vmCfg    *runtime.Config
	mStateDB *state.MemoryStateDB
}

// NewEVMExecutor
func NewEVMExecutor() *EVMExecutor {
	exec := &EVMExecutor{}

	exec.vmCfg = &runtime.Config{}
	exec.vmCfg.Tracer = runtime.NewJSONLogger(os.Stdout)

	exec.SetChild(exec)
	return exec
}

// GetFuncMap
func (evm *EVMExecutor) GetFuncMap() map[string]reflect.Method {
	ety := types.LoadExecutorType(driverName)
	return ety.GetExecFuncMap()
}

// GetDriverName
func (evm *EVMExecutor) GetDriverName() string {
	return evmtypes.ExecutorName
}

// ExecutorOrder   localdb EnableRead
func (evm *EVMExecutor) ExecutorOrder() int64 {
	cfg := evm.GetAPI().GetConfig()
	if cfg.IsFork(evm.GetHeight(), "ForkLocalDBAccess") {
		return drivers.ExecLocalSameTime
	}
	return evm.DriverBase.ExecutorOrder()
}

// Allow
func (evm *EVMExecutor) Allow(tx *types.Transaction, index int) error {
	err := evm.DriverBase.Allow(tx, index)
	if err == nil {
		return nil
	}
	//      :
	//  : user.evm.xxx     evm
	//   : user.p.guodun.user.evm.xxx    evm
	cfg := evm.GetAPI().GetConfig()
	exec := cfg.GetParaExec(tx.Execer)
	if evm.AllowIsUserDot2(exec) {
		return nil
	}
	return types.ErrNotAllow
}

// IsFriend        KEY
func (evm *EVMExecutor) IsFriend(myexec, writekey []byte, othertx *types.Transaction) bool {
	if othertx == nil {
		return false
	}
	cfg := evm.GetAPI().GetConfig()
	exec := cfg.GetParaExec(othertx.Execer)
	if exec == nil || len(bytes.TrimSpace(exec)) == 0 {
		return false
	}
	if bytes.HasPrefix(exec, evmtypes.UserPrefix) || bytes.Equal(exec, evmtypes.ExecerEvm) {
		if bytes.HasPrefix(writekey, []byte("mavl-evm-")) {
			return true
		}
	}
	return false
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (evm *EVMExecutor) CheckReceiptExecOk() bool {
	return true
}

//
func (evm *EVMExecutor) getNewAddr(txHash []byte) common.Address {
	cfg := evm.GetAPI().GetConfig()
	return common.NewAddress(cfg, txHash)
}

// CheckTx
func (evm *EVMExecutor) CheckTx(tx *types.Transaction, index int) error {
	return nil
}

// GetActionName
func (evm *EVMExecutor) GetActionName(tx *types.Transaction) string {
	cfg := evm.GetAPI().GetConfig()
	if bytes.Equal(tx.Execer, []byte(cfg.ExecName(evmtypes.ExecutorName))) {
		return cfg.ExecName(evmtypes.ExecutorName)
	}
	return tx.ActionName()
}

// GetMStateDB
func (evm *EVMExecutor) GetMStateDB() *state.MemoryStateDB {
	return evm.mStateDB
}

// GetVMConfig   VM
func (evm *EVMExecutor) GetVMConfig() *runtime.Config {
	return evm.vmCfg
}

// NewEVMContext       EVM
func (evm *EVMExecutor) NewEVMContext(msg *common.Message) runtime.Context {
	return runtime.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(evm.GetAPI()),
		Origin:      msg.From(),
		Coinbase:    nil,
		BlockNumber: new(big.Int).SetInt64(evm.GetHeight()),
		Time:        new(big.Int).SetInt64(evm.GetBlockTime()),
		Difficulty:  new(big.Int).SetUint64(evm.GetDifficulty()),
		GasLimit:    msg.GasLimit(),
		GasPrice:    msg.GasPrice(),
	}
}
