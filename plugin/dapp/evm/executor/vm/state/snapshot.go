// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package state

import (
	"sort"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	evmtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"
)

// DataChange
//                 ，
//       （   Tx   ），
//         ，         ，                   ，
//         ，         ，                  ，
type DataChange interface {
	revert(mdb *MemoryStateDB)
	getData(mdb *MemoryStateDB) []*types.KeyValue
	getLog(mdb *MemoryStateDB) []*types.ReceiptLog
}

// Snapshot     ，
type Snapshot struct {
	id      int
	entries []DataChange
	statedb *MemoryStateDB
}

// GetID   ID
func (ver *Snapshot) GetID() int {
	return ver.id
}

//
func (ver *Snapshot) revert() bool {
	if ver.entries == nil {
		return true
	}
	for _, entry := range ver.entries {
		entry.revert(ver.statedb)
	}
	return true
}

//
func (ver *Snapshot) append(entry DataChange) {
	ver.entries = append(ver.entries, entry)
}

//
func (ver *Snapshot) getData() (kvSet []*types.KeyValue, logs []*types.ReceiptLog) {
	//
	dataMap := make(map[string]*types.KeyValue)

	for _, entry := range ver.entries {
		items := entry.getData(ver.statedb)
		logEntry := entry.getLog(ver.statedb)
		if logEntry != nil {
			logs = append(logs, entry.getLog(ver.statedb)...)
		}

		//
		for _, kv := range items {
			dataMap[string(kv.Key)] = kv
		}
	}

	//                   ，    （   KV           ，           ）
	names := make([]string, 0, len(dataMap))
	for name := range dataMap {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		kvSet = append(kvSet, dataMap[name])
	}

	return kvSet, logs
}

type (

	//       ，
	baseChange struct {
	}

	//
	createAccountChange struct {
		baseChange
		account string
	}

	//
	suicideChange struct {
		baseChange
		account string
		prev    bool // whether account had already suicided
	}

	// nonce
	nonceChange struct {
		baseChange
		account string
		prev    uint64
	}

	//
	storageChange struct {
		baseChange
		account       string
		key, prevalue common.Hash
	}

	//
	codeChange struct {
		baseChange
		account            string
		prevcode, prevhash []byte
	}

	//   ABI
	abiChange struct {
		baseChange
		account string
		prevabi string
	}

	//
	refundChange struct {
		baseChange
		prev uint64
	}

	//
	//            ，
	transferChange struct {
		baseChange
		amount int64
		data   []*types.KeyValue
		logs   []*types.ReceiptLog
	}

	//
	addLogChange struct {
		baseChange
		txhash common.Hash
	}

	//     sha3
	addPreimageChange struct {
		baseChange
		hash common.Hash
	}
)

//  baseChang         ，
func (ch baseChange) revert(s *MemoryStateDB) {
}

func (ch baseChange) getData(s *MemoryStateDB) (kvset []*types.KeyValue) {
	return nil
}

func (ch baseChange) getLog(s *MemoryStateDB) (logs []*types.ReceiptLog) {
	return nil
}

//          ，
func (ch createAccountChange) revert(s *MemoryStateDB) {
	delete(s.accounts, ch.account)
}

//
func (ch createAccountChange) getData(s *MemoryStateDB) (kvset []*types.KeyValue) {
	acc := s.accounts[ch.account]
	if acc != nil {
		kvset = append(kvset, acc.GetDataKV()...)
		kvset = append(kvset, acc.GetStateKV()...)
		return kvset
	}
	return nil
}

func (ch suicideChange) revert(mdb *MemoryStateDB) {
	//         ，
	if ch.prev {
		return
	}
	acc := mdb.accounts[ch.account]
	if acc != nil {
		acc.State.Suicided = ch.prev
	}
}

func (ch suicideChange) getData(mdb *MemoryStateDB) []*types.KeyValue {
	//         ，
	if ch.prev {
		return nil
	}
	acc := mdb.accounts[ch.account]
	if acc != nil {
		return acc.GetStateKV()
	}
	return nil
}

func (ch nonceChange) revert(mdb *MemoryStateDB) {
	acc := mdb.accounts[ch.account]
	if acc != nil {
		acc.State.Nonce = ch.prev
	}
}

func (ch nonceChange) getData(mdb *MemoryStateDB) []*types.KeyValue {
	// nonce        ，          ，
	//acc := mdb.accounts[ch.account]
	//if acc != nil {
	//	return acc.GetStateKV()
	//}
	return nil
}

func (ch codeChange) revert(mdb *MemoryStateDB) {
	acc := mdb.accounts[ch.account]
	if acc != nil {
		acc.Data.Code = ch.prevcode
		acc.Data.CodeHash = ch.prevhash
	}
}

func (ch codeChange) getData(mdb *MemoryStateDB) (kvset []*types.KeyValue) {
	acc := mdb.accounts[ch.account]
	if acc != nil {
		kvset = append(kvset, acc.GetDataKV()...)
		kvset = append(kvset, acc.GetStateKV()...)
		return kvset
	}
	return nil
}

func (ch abiChange) revert(mdb *MemoryStateDB) {
	acc := mdb.accounts[ch.account]
	if acc != nil {
		acc.Data.Abi = ch.prevabi
	}
}

func (ch abiChange) getData(mdb *MemoryStateDB) (kvset []*types.KeyValue) {
	acc := mdb.accounts[ch.account]
	if acc != nil {
		kvset = append(kvset, acc.GetDataKV()...)
		return kvset
	}
	return nil
}

func (ch storageChange) revert(mdb *MemoryStateDB) {
	acc := mdb.accounts[ch.account]
	if acc != nil {
		acc.SetState(ch.key, ch.prevalue)
	}
}

func (ch storageChange) getData(mdb *MemoryStateDB) []*types.KeyValue {
	acc := mdb.accounts[ch.account]
	if _, ok := mdb.stateDirty[ch.account]; ok && acc != nil {
		return acc.GetStateKV()
	}
	return nil
}

func (ch storageChange) getLog(mdb *MemoryStateDB) []*types.ReceiptLog {
	cfg := mdb.api.GetConfig()
	if cfg.IsDappFork(mdb.blockHeight, "evm", evmtypes.ForkEVMState) {
		acc := mdb.accounts[ch.account]
		if acc != nil {
			currentVal := acc.GetState(ch.key)
			receipt := &evmtypes.EVMStateChangeItem{Key: getStateItemKey(ch.account, ch.key.Hex()), PreValue: ch.prevalue.Bytes(), CurrentValue: currentVal.Bytes()}
			return []*types.ReceiptLog{{Ty: evmtypes.TyLogEVMStateChangeItem, Log: types.Encode(receipt)}}
		}
	}
	return nil
}

func (ch refundChange) revert(mdb *MemoryStateDB) {
	mdb.refund = ch.prev
}

func (ch addLogChange) revert(mdb *MemoryStateDB) {
	logs := mdb.logs[ch.txhash]
	if len(logs) == 1 {
		delete(mdb.logs, ch.txhash)
	} else {
		mdb.logs[ch.txhash] = logs[:len(logs)-1]
	}
	mdb.logSize--
}

func (ch addPreimageChange) revert(mdb *MemoryStateDB) {
	delete(mdb.preimages, ch.hash)
}

func (ch transferChange) getData(mdb *MemoryStateDB) []*types.KeyValue {
	return ch.data
}
func (ch transferChange) getLog(mdb *MemoryStateDB) []*types.ReceiptLog {
	return ch.logs
}
