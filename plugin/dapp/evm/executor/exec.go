// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"encoding/hex"
	"fmt"
	"strings"

	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/abi"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/model"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/runtime"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/state"
	evmtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"
)

// Exec
func (evm *EVMExecutor) Exec(tx *types.Transaction, index int) (*types.Receipt, error) {
	evm.CheckInit()
	//
	msg, err := evm.GetMessage(tx, index)
	if err != nil {
		return nil, err
	}

	return evm.innerExec(msg, tx.Hash(), index, evm.GetTxFee(tx, index), false)
}

//    EVM
// readOnly       ，   evm abi    true
func (evm *EVMExecutor) innerExec(msg *common.Message, txHash []byte, index int, txFee int64, readOnly bool) (receipt *types.Receipt, err error) {
	//               EVM
	context := evm.NewEVMContext(msg)
	cfg := evm.GetAPI().GetConfig()
	//   EVM
	env := runtime.NewEVM(context, evm.mStateDB, *evm.vmCfg, cfg)
	isCreate := strings.Compare(msg.To().String(), EvmAddress) == 0
	var (
		ret          []byte
		vmerr        error
		leftOverGas  uint64
		contractAddr common.Address
		snapshot     int
		execName     string
		methodName   string
	)

	//       ，        ，
	if isCreate {
		//                （                   ，        ）
		contractAddr = evm.getNewAddr(txHash)
		if !env.StateDB.Empty(contractAddr.String()) {
			return receipt, model.ErrContractAddressCollision
		}
		//
		execName = fmt.Sprintf("%s%s", cfg.ExecName(evmtypes.EvmPrefix), common.BytesToHash(txHash).Hex())
	} else {
		contractAddr = *msg.To()
	}

	//
	evm.mStateDB.Prepare(common.BytesToHash(txHash), index)

	if isCreate {
		//     ABI  ，
		if len(msg.ABI()) > 0 && cfg.IsDappFork(evm.GetHeight(), "evm", evmtypes.ForkEVMABI) {
			_, err = abi.JSON(strings.NewReader(msg.ABI()))
			if err != nil {
				return receipt, err
			}
		}
		ret, snapshot, leftOverGas, vmerr = env.Create(runtime.AccountRef(msg.From()), contractAddr, msg.Data(), context.GasLimit, execName, msg.Alias(), msg.ABI())
	} else {
		inData := msg.Data()
		//      ABI
		if len(msg.ABI()) > 0 && cfg.IsDappFork(evm.GetHeight(), "evm", evmtypes.ForkEVMABI) {
			funcName, packData, err := abi.Pack(msg.ABI(), evm.mStateDB.GetAbi(msg.To().String()), readOnly)
			if err != nil {
				return receipt, err
			}
			inData = packData
			methodName = funcName
			log.Debug("call contract ", "abi funcName", funcName, "packData", common.Bytes2Hex(inData))
		}
		ret, snapshot, leftOverGas, vmerr = env.Call(runtime.AccountRef(msg.From()), *msg.To(), inData, context.GasLimit, msg.Value())
		log.Debug("call(create) contract ", "input", common.Bytes2Hex(inData))
	}
	usedGas := msg.GasLimit() - leftOverGas
	logMsg := "call contract details:"
	if isCreate {
		logMsg = "create contract details:"
	}
	log.Debug(logMsg, "caller address", msg.From().String(), "contract address", contractAddr.String(), "exec name", execName, "alias name", msg.Alias(), "usedGas", usedGas, "return data", common.Bytes2Hex(ret))
	curVer := evm.mStateDB.GetLastSnapshot()
	if vmerr != nil {
		log.Error("evm contract exec error", "error info", vmerr)
		return receipt, vmerr
	}

	//          （       ）
	usedFee, overflow := common.SafeMul(usedGas, uint64(msg.GasPrice()))
	//       ，
	if overflow || usedFee > uint64(txFee) {
		//         ，
		if curVer != nil && snapshot >= curVer.GetID() && curVer.GetID() > -1 {
			evm.mStateDB.RevertToSnapshot(snapshot)
		}
		return receipt, model.ErrOutOfGas
	}

	//
	evm.mStateDB.PrintLogs()

	//
	if curVer == nil {
		return receipt, nil
	}

	//
	kvSet, logs := evm.mStateDB.GetChangedData(curVer.GetID())
	contractReceipt := &evmtypes.ReceiptEVMContract{Caller: msg.From().String(), ContractName: execName, ContractAddr: contractAddr.String(), UsedGas: usedGas, Ret: ret}
	//     ABI
	if len(methodName) > 0 && len(msg.ABI()) > 0 && cfg.IsDappFork(evm.GetHeight(), "evm", evmtypes.ForkEVMABI) {
		jsonRet, err := abi.Unpack(ret, methodName, evm.mStateDB.GetAbi(msg.To().String()))
		if err != nil {
			//            ，
			log.Error("unpack evm return error", "error", err)
		}
		contractReceipt.JsonRet = jsonRet
	}
	logs = append(logs, &types.ReceiptLog{Ty: evmtypes.TyLogCallContract, Log: types.Encode(contractReceipt)})
	logs = append(logs, evm.mStateDB.GetReceiptLogs(contractAddr.String())...)

	if cfg.IsDappFork(evm.GetHeight(), "evm", evmtypes.ForkEVMKVHash) {
		//
		hashKV := evm.calcKVHash(contractAddr, logs)
		if hashKV != nil {
			kvSet = append(kvSet, hashKV)
		}
	}

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kvSet, Logs: logs}

	//     ，
	if evm.mStateDB != nil {
		evm.mStateDB.WritePreimages(evm.GetHeight())
	}

	//
	state.ProcessFork(cfg, evm.GetHeight(), txHash, receipt)

	evm.collectEvmTxLog(txHash, contractReceipt, receipt)

	return receipt, nil
}

// CheckInit
func (evm *EVMExecutor) CheckInit() {
	if evm.mStateDB == nil {
		evm.mStateDB = state.NewMemoryStateDB(evm.GetStateDB(), evm.GetLocalDB(), evm.GetCoinsAccount(), evm.GetHeight(), evm.GetAPI())
	}
}

// GetMessage       ，   coins  ，     payload ，      ，    Transaction
func (evm *EVMExecutor) GetMessage(tx *types.Transaction, index int) (msg *common.Message, err error) {
	var action evmtypes.EVMContractAction
	err = types.Decode(tx.Payload, &action)
	if err != nil {
		return msg, err
	}
	//                 ，dplatformos mempool
	from := getCaller(tx)
	to := getReceiver(tx)
	if to == nil {
		return msg, types.ErrInvalidAddress
	}

	gasLimit := action.GasLimit
	gasPrice := action.GasPrice
	if gasLimit == 0 {
		gasLimit = uint64(evm.GetTxFee(tx, index))
	}
	if gasPrice == 0 {
		gasPrice = uint32(1)
	}

	//    GasLimit
	msg = common.NewMessage(from, to, tx.Nonce, action.Amount, gasLimit, gasPrice, action.Code, action.GetAlias(), action.Abi)
	return msg, err
}

func (evm *EVMExecutor) collectEvmTxLog(txHash []byte, cr *evmtypes.ReceiptEVMContract, receipt *types.Receipt) {
	log.Debug("evm collect begin")
	log.Debug("Tx info", "txHash", common.Bytes2Hex(txHash), "height", evm.GetHeight())
	log.Debug("ReceiptEVMContract", "data", fmt.Sprintf("caller=%v, name=%v, addr=%v, usedGas=%v, ret=%v", cr.Caller, cr.ContractName, cr.ContractAddr, cr.UsedGas, common.Bytes2Hex(cr.Ret)))
	log.Debug("receipt data", "type", receipt.Ty)
	for _, kv := range receipt.KV {
		log.Debug("KeyValue", "key", common.Bytes2Hex(kv.Key), "value", common.Bytes2Hex(kv.Value))
	}
	for _, kv := range receipt.Logs {
		log.Debug("ReceiptLog", "Type", kv.Ty, "log", common.Bytes2Hex(kv.Log))
	}
	log.Debug("evm collect end")
}

func (evm *EVMExecutor) calcKVHash(addr common.Address, logs []*types.ReceiptLog) (kv *types.KeyValue) {
	hashes := []byte{}
	//                ，     KV
	for _, logItem := range logs {
		if evmtypes.TyLogEVMStateChangeItem == logItem.Ty {
			data := logItem.Log
			hashes = append(hashes, common.ToHash(data).Bytes()...)
		}
	}

	if len(hashes) > 0 {
		hash := common.ToHash(hashes)
		return &types.KeyValue{Key: getDataHashKey(addr), Value: hash.Bytes()}
	}
	return nil
}

// GetTxFee        ，
func (evm *EVMExecutor) GetTxFee(tx *types.Transaction, index int) int64 {
	fee := tx.Fee
	cfg := evm.GetAPI().GetConfig()
	if fee == 0 && cfg.IsDappFork(evm.GetHeight(), "evm", evmtypes.ForkEVMTxGroup) {
		if tx.GroupCount >= 2 {
			txs, err := evm.GetTxGroup(index)
			if err != nil {
				log.Error("evm GetTxFee", "get tx group fail", err, "hash", hex.EncodeToString(tx.Hash()))
				return 0
			}
			fee = txs[0].Fee
		}
	}
	return fee
}

func getDataHashKey(addr common.Address) []byte {
	return []byte(fmt.Sprintf("mavl-%v-data-hash:%v", evmtypes.ExecutorName, addr))
}

//
func getCaller(tx *types.Transaction) common.Address {
	return *common.StringToAddress(tx.From())
}

//               ，        ，
func getReceiver(tx *types.Transaction) *common.Address {
	if tx.To == "" {
		return nil
	}
	return common.StringToAddress(tx.To)
}
