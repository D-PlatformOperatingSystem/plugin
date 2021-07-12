// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package state

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	evmtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"
	"github.com/golang/protobuf/proto"
)

// ContractAccount
type ContractAccount struct {
	mdb *MemoryStateDB

	// Addr
	Addr string

	// Data
	Data evmtypes.EVMContractData

	// State
	State evmtypes.EVMContractState

	//
	stateCache map[string]common.Hash
}

// NewContractAccount
//   ，
//     LoadContract
func NewContractAccount(addr string, db *MemoryStateDB) *ContractAccount {
	if len(addr) == 0 || db == nil {
		log15.Error("NewContractAccount error, something is missing", "contract addr", addr, "db", db)
		return nil
	}
	ca := &ContractAccount{Addr: addr, mdb: db}
	ca.State.Storage = make(map[string][]byte)
	ca.stateCache = make(map[string]common.Hash)
	return ca
}

// GetState       ；
//         ，            ，      ，  localdb
func (ca *ContractAccount) GetState(key common.Hash) common.Hash {
	//  ForkV19  ，         KEY
	cfg := ca.mdb.api.GetConfig()
	if cfg.IsDappFork(ca.mdb.blockHeight, "evm", evmtypes.ForkEVMState) {
		if val, ok := ca.stateCache[key.Hex()]; ok {
			return val
		}
		keyStr := getStateItemKey(ca.Addr, key.Hex())
		//           ，
		val, err := ca.mdb.LocalDB.Get([]byte(keyStr))
		if err != nil {
			log15.Debug("GetState error!", "key", key, "error", err)
			return common.Hash{}
		}
		valHash := common.BytesToHash(val)
		ca.stateCache[key.Hex()] = valHash
		return valHash
	}
	return common.BytesToHash(ca.State.GetStorage()[key.Hex()])
}

// SetState
func (ca *ContractAccount) SetState(key, value common.Hash) {
	ca.mdb.addChange(storageChange{
		baseChange: baseChange{},
		account:    ca.Addr,
		key:        key,
		prevalue:   ca.GetState(key),
	})
	cfg := ca.mdb.api.GetConfig()
	if cfg.IsDappFork(ca.mdb.blockHeight, "evm", evmtypes.ForkEVMState) {
		ca.stateCache[key.Hex()] = value
		//     localdb ，                  ，
		keyStr := getStateItemKey(ca.Addr, key.Hex())
		ca.mdb.LocalDB.Set([]byte(keyStr), value.Bytes())
	} else {
		ca.State.GetStorage()[key.Hex()] = value.Bytes()
		ca.updateStorageHash()
	}
}

// TransferState            ，            KEY，
func (ca *ContractAccount) TransferState() {
	if len(ca.State.Storage) > 0 {
		storage := ca.State.Storage
		//          、        ，
		ca.State.Storage = make(map[string][]byte)
		ca.State.StorageHash = common.ToHash([]byte{}).Bytes()

		//                 ，
		for key, value := range storage {
			ca.SetState(common.BytesToHash(common.FromHex(key)), common.BytesToHash(value))
		}
		//           （    map    ）
		ca.mdb.UpdateState(ca.Addr)
		return
	}
}

func (ca *ContractAccount) updateStorageHash() {
	//  ForkV20  ，        KEY
	cfg := ca.mdb.api.GetConfig()
	if cfg.IsDappFork(ca.mdb.blockHeight, "evm", evmtypes.ForkEVMState) {
		return
	}
	var state = &evmtypes.EVMContractState{Suicided: ca.State.Suicided, Nonce: ca.State.Nonce}
	state.Storage = make(map[string][]byte)
	for k, v := range ca.State.GetStorage() {
		state.Storage[k] = v
	}
	ret, err := proto.Marshal(state)
	if err != nil {
		log15.Error("marshal contract state data error", "error", err)
		return
	}

	ca.State.StorageHash = common.ToHash(ret).Bytes()
}

//
func (ca *ContractAccount) resotreData(data []byte) {
	var content evmtypes.EVMContractData
	err := proto.Unmarshal(data, &content)
	if err != nil {
		log15.Error("read contract data error", ca.Addr)
		return
	}

	ca.Data = content
}

//
func (ca *ContractAccount) resotreState(data []byte) {
	var content evmtypes.EVMContractState
	err := proto.Unmarshal(data, &content)
	if err != nil {
		log15.Error("read contract state error", ca.Addr)
		return
	}
	ca.State = content
	if ca.State.Storage == nil {
		ca.State.Storage = make(map[string][]byte)
	}
}

// LoadContract            （           ）
func (ca *ContractAccount) LoadContract(db db.KV) {
	//
	data, err := db.Get(ca.GetDataKey())
	if err != nil {
		return
	}
	ca.resotreData(data)

	//
	data, err = db.Get(ca.GetStateKey())
	if err != nil {
		return
	}
	ca.resotreState(data)
}

// SetCode
//
func (ca *ContractAccount) SetCode(code []byte) {
	prevcode := ca.Data.GetCode()
	ca.mdb.addChange(codeChange{
		baseChange: baseChange{},
		account:    ca.Addr,
		prevhash:   ca.Data.GetCodeHash(),
		prevcode:   prevcode,
	})
	ca.Data.Code = code
	ca.Data.CodeHash = common.ToHash(code).Bytes()
}

// SetAbi        ABI
func (ca *ContractAccount) SetAbi(abi string) {
	cfg := ca.mdb.api.GetConfig()
	if cfg.IsDappFork(ca.mdb.GetBlockHeight(), "evm", evmtypes.ForkEVMABI) {
		ca.mdb.addChange(abiChange{
			baseChange: baseChange{},
			account:    ca.Addr,
			prevabi:    ca.Data.Abi,
		})
		ca.Data.Abi = abi
	}
}

// SetCreator
func (ca *ContractAccount) SetCreator(creator string) {
	if len(creator) == 0 {
		log15.Error("SetCreator error", "creator", creator)
		return
	}
	ca.Data.Creator = creator
}

// SetExecName
func (ca *ContractAccount) SetExecName(execName string) {
	if len(execName) == 0 {
		log15.Error("SetExecName error", "execName", execName)
		return
	}
	ca.Data.Name = execName
}

// SetAliasName
func (ca *ContractAccount) SetAliasName(alias string) {
	if len(alias) == 0 {
		log15.Error("SetAliasName error", "aliasName", alias)
		return
	}
	ca.Data.Alias = alias
}

// GetAliasName
func (ca *ContractAccount) GetAliasName() string {
	return ca.Data.Alias
}

// GetCreator
func (ca *ContractAccount) GetCreator() string {
	return ca.Data.Creator
}

// GetExecName
func (ca *ContractAccount) GetExecName() string {
	return ca.Data.Name
}

// GetDataKV       ，      ，
func (ca *ContractAccount) GetDataKV() (kvSet []*types.KeyValue) {
	ca.Data.Addr = ca.Addr
	datas, err := proto.Marshal(&ca.Data)
	if err != nil {
		log15.Error("marshal contract data error!", "addr", ca.Addr, "error", err)
		return
	}
	kvSet = append(kvSet, &types.KeyValue{Key: ca.GetDataKey(), Value: datas})
	return
}

// GetStateKV         ，  nonce、    、    、
func (ca *ContractAccount) GetStateKV() (kvSet []*types.KeyValue) {
	datas, err := proto.Marshal(&ca.State)
	if err != nil {
		log15.Error("marshal contract state error!", "addr", ca.Addr, "error", err)
		return
	}
	kvSet = append(kvSet, &types.KeyValue{Key: ca.GetStateKey(), Value: datas})
	return
}

// BuildDataLog
func (ca *ContractAccount) BuildDataLog() (log *types.ReceiptLog) {
	datas, err := proto.Marshal(&ca.Data)
	if err != nil {
		log15.Error("marshal contract data error!", "addr", ca.Addr, "error", err)
		return
	}
	return &types.ReceiptLog{Ty: evmtypes.TyLogContractData, Log: datas}
}

// BuildStateLog
func (ca *ContractAccount) BuildStateLog() (log *types.ReceiptLog) {
	datas, err := proto.Marshal(&ca.State)
	if err != nil {
		log15.Error("marshal contract state log error!", "addr", ca.Addr, "error", err)
		return
	}

	return &types.ReceiptLog{Ty: evmtypes.TyLogContractState, Log: datas}
}

// GetDataKey     KEY
func (ca *ContractAccount) GetDataKey() []byte {
	return []byte("mavl-" + evmtypes.ExecutorName + "-data: " + ca.Addr)
}

// GetStateKey     key
func (ca *ContractAccount) GetStateKey() []byte {
	return []byte("mavl-" + evmtypes.ExecutorName + "-state: " + ca.Addr)
}

//        LocalDB
func getStateItemKey(addr, key string) string {
	return fmt.Sprintf("LODB-"+evmtypes.ExecutorName+"-state:%v:%v", addr, key)
}

// Suicide
func (ca *ContractAccount) Suicide() bool {
	ca.State.Suicided = true
	return true
}

// HasSuicided
func (ca *ContractAccount) HasSuicided() bool {
	return ca.State.GetSuicided()
}

// Empty
func (ca *ContractAccount) Empty() bool {
	return ca.Data.GetCodeHash() == nil || len(ca.Data.GetCodeHash()) == 0
}

// SetNonce   nonce
func (ca *ContractAccount) SetNonce(nonce uint64) {
	ca.mdb.addChange(nonceChange{
		baseChange: baseChange{},
		account:    ca.Addr,
		prev:       ca.State.GetNonce(),
	})
	ca.State.Nonce = nonce
}

// GetNonce   nonce
func (ca *ContractAccount) GetNonce() uint64 {
	return ca.State.GetNonce()
}
