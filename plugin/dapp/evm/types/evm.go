// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/golang/protobuf/proto"
)

var (
	elog = log.New("module", "exectype.evm")

	actionName = map[string]int32{
		"EvmCreate": EvmCreateAction,
		"EvmCall":   EvmCallAction,
	}
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, ExecerEvm)
	types.RegFork(ExecutorName, InitFork)
	types.RegExec(ExecutorName, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(ExecutorName, EVMEnable, 500000)
	// EVM          ，
	cfg.RegisterDappFork(ExecutorName, ForkEVMState, 650000)
	// EVM          ，      StateDB
	cfg.RegisterDappFork(ExecutorName, ForkEVMKVHash, 1000000)
	// EVM    ABI
	cfg.RegisterDappFork(ExecutorName, ForkEVMABI, 1250000)
	// EEVM
	cfg.RegisterDappFork(ExecutorName, ForkEVMFrozen, 1300000)
	// EEVM   v1
	cfg.RegisterDappFork(ExecutorName, ForkEVMYoloV1, 9500000)
	// EVM
	cfg.RegisterDappFork(ExecutorName, ForkEVMTxGroup, 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(ExecutorName, NewType(cfg))
}

// EvmType EVM
type EvmType struct {
	types.ExecTypeBase
}

// NewType   EVM
func NewType(cfg *types.DplatformOSConfig) *EvmType {
	c := &EvmType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName
func (evm *EvmType) GetName() string {
	return ExecutorName
}

// GetPayload
func (evm *EvmType) GetPayload() types.Message {
	return &EVMContractAction{}
}

// ActionName   ActionName
func (evm EvmType) ActionName(tx *types.Transaction) string {
	//                  Action
	//         ，  evm       ，      ，
	cfg := evm.GetConfig()
	if strings.EqualFold(tx.To, address.ExecAddress(cfg.ExecName(ExecutorName))) {
		return "createEvmContract"
	}
	return "callEvmContract"
}

// GetTypeMap
func (evm *EvmType) GetTypeMap() map[string]int32 {
	return actionName
}

// GetRealToAddr
func (evm EvmType) GetRealToAddr(tx *types.Transaction) string {
	if string(tx.Execer) == ExecutorName {
		return tx.To
	}
	var action EVMContractAction
	err := types.Decode(tx.Payload, &action)
	if err != nil {
		return tx.To
	}
	return tx.To
}

// Amount
func (evm EvmType) Amount(tx *types.Transaction) (int64, error) {
	return 0, nil
}

// CreateTx
func (evm EvmType) CreateTx(action string, message json.RawMessage) (*types.Transaction, error) {
	elog.Debug("evm.CreateTx", "action", action)
	if action == "CreateCall" {
		var param CreateCallTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			elog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return createEvmTx(evm.GetConfig(), &param)
	}
	return nil, types.ErrNotSupport
}

// GetLogMap
func (evm *EvmType) GetLogMap() map[int64]*types.LogInfo {
	return logInfo
}

func createEvmTx(cfg *types.DplatformOSConfig, param *CreateCallTx) (*types.Transaction, error) {
	if param == nil {
		elog.Error("createEvmTx", "param", param)
		return nil, types.ErrInvalidParam
	}

	//         ：
	//                ，    ，  ABI
	//       ， ABI    0x00000000

	action := &EVMContractAction{
		Amount:   param.Amount,
		GasLimit: param.GasLimit,
		GasPrice: param.GasPrice,
		Note:     param.Note,
		Alias:    param.Alias,
	}
	// Abi              ，    ABI
	if len(param.Abi) > 0 {
		action.Abi = strings.TrimSpace(param.Abi)
	}
	if len(param.Code) > 0 {
		bCode, err := common.FromHex(param.Code)
		if err != nil {
			elog.Error("create evm Tx error, code is invalid", "param.Code", param.Code)
			return nil, err
		}
		action.Code = bCode
	}

	if param.IsCreate {
		if len(action.Abi) > 0 && len(action.Code) == 0 {
			elog.Error("create evm Tx error, code is empty")
			return nil, errors.New("code must be set in create tx")
		}

		return createRawTx(cfg, action, "", param.Fee)
	}
	return createRawTx(cfg, action, param.Name, param.Fee)
}

func createRawTx(cfg *types.DplatformOSConfig, action proto.Message, name string, fee int64) (*types.Transaction, error) {
	tx := &types.Transaction{}
	if len(name) == 0 {
		tx = &types.Transaction{
			Execer:  []byte(cfg.ExecName(ExecutorName)),
			Payload: types.Encode(action),
			To:      address.ExecAddress(cfg.ExecName(ExecutorName)),
		}
	} else {
		tx = &types.Transaction{
			Execer:  []byte(cfg.ExecName(name)),
			Payload: types.Encode(action),
			To:      address.ExecAddress(cfg.ExecName(name)),
		}
	}
	tx, err := types.FormatTx(cfg, string(tx.Execer), tx)
	if err != nil {
		return nil, err
	}

	if tx.Fee < fee {
		tx.Fee = fee
	}

	return tx, nil
}
