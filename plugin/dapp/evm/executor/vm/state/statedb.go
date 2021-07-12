// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package state

import (
	"fmt"
	"strings"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/client"
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/model"
	evmtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"
)

// MemoryStateDB        ，
//            ，
//         ，             （          blockchain statedb   ）
//         ，          ，        ，         ，   blockchain
//     Exec     ：    、    （      、    、      ）
//     ExecLocal     ：
type MemoryStateDB struct {
	// StateDB   DB，
	StateDB db.KV

	// LocalDB   DB，
	LocalDB db.KVDB

	// CoinsAccount Coins      ，
	CoinsAccount *account.DB

	//
	accounts map[string]*ContractAccount

	//
	refund uint64

	//   makeLogN
	logs    map[common.Hash][]*model.ContractLog
	logSize uint

	//    ，
	snapshots  []*Snapshot
	currentVer *Snapshot
	versionID  int

	//   sha3       ，   debug
	preimages map[common.Hash][]byte

	//
	txHash  common.Hash
	txIndex int

	//
	blockHeight int64

	//
	stateDirty map[string]interface{}
	dataDirty  map[string]interface{}
	api        client.QueueProtocolAPI
}

// NewMemoryStateDB           DB
//               ，                    DB
//           （       setEnv            ），      DB
func NewMemoryStateDB(StateDB db.KV, LocalDB db.KVDB, CoinsAccount *account.DB, blockHeight int64, api client.QueueProtocolAPI) *MemoryStateDB {
	mdb := &MemoryStateDB{
		StateDB:      StateDB,
		LocalDB:      LocalDB,
		CoinsAccount: CoinsAccount,
		accounts:     make(map[string]*ContractAccount),
		logs:         make(map[common.Hash][]*model.ContractLog),
		preimages:    make(map[common.Hash][]byte),
		stateDirty:   make(map[string]interface{}),
		dataDirty:    make(map[string]interface{}),
		blockHeight:  blockHeight,
		refund:       0,
		txIndex:      0,
		api:          api,
	}
	return mdb
}

// Prepare               ，
//
func (mdb *MemoryStateDB) Prepare(txHash common.Hash, txIndex int) {
	mdb.txHash = txHash
	mdb.txIndex = txIndex
}

// CreateAccount
func (mdb *MemoryStateDB) CreateAccount(addr, creator string, execName, alias string) {
	acc := mdb.GetAccount(addr)
	if acc == nil {
		//
		acc := NewContractAccount(addr, mdb)
		acc.SetCreator(creator)
		acc.SetExecName(execName)
		acc.SetAliasName(alias)
		mdb.accounts[addr] = acc
		mdb.addChange(createAccountChange{baseChange: baseChange{}, account: addr})
	}
}

func (mdb *MemoryStateDB) addChange(entry DataChange) {
	if mdb.currentVer != nil {
		mdb.currentVer.append(entry)
	}
}

// SubBalance          （            ）
func (mdb *MemoryStateDB) SubBalance(addr, caddr string, value uint64) {
	res := mdb.Transfer(addr, caddr, value)
	log15.Debug("transfer result", "from", addr, "to", caddr, "amount", value, "result", res)
}

// AddBalance          （                  ）
func (mdb *MemoryStateDB) AddBalance(addr, caddr string, value uint64) {
	res := mdb.Transfer(caddr, addr, value)
	log15.Debug("transfer result", "from", addr, "to", caddr, "amount", value, "result", res)
}

// GetBalance         ，       ，                    ；
//        ，
func (mdb *MemoryStateDB) GetBalance(addr string) uint64 {
	if mdb.CoinsAccount == nil {
		return 0
	}
	isExec := mdb.Exist(addr)
	var ac *types.Account
	if isExec {
		cfg := mdb.api.GetConfig()
		if cfg.IsDappFork(mdb.GetBlockHeight(), "evm", evmtypes.ForkEVMFrozen) {
			ac = mdb.CoinsAccount.LoadExecAccount(addr, addr)
		} else {
			contract := mdb.GetAccount(addr)
			if contract == nil {
				return 0
			}
			creator := contract.GetCreator()
			if len(creator) == 0 {
				return 0
			}
			ac = mdb.CoinsAccount.LoadExecAccount(creator, addr)
		}
	} else {
		ac = mdb.CoinsAccount.LoadAccount(addr)
	}
	if ac != nil {
		return uint64(ac.Balance)
	}
	return 0
}

// GetNonce   dplatformos        nonce  ，            ；
//   ，         nonce
func (mdb *MemoryStateDB) GetNonce(addr string) uint64 {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.GetNonce()
	}
	return 0
}

// SetNonce   nonce
func (mdb *MemoryStateDB) SetNonce(addr string, nonce uint64) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		acc.SetNonce(nonce)
	}
}

// GetCodeHash
func (mdb *MemoryStateDB) GetCodeHash(addr string) common.Hash {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return common.BytesToHash(acc.Data.GetCodeHash())
	}
	return common.Hash{}
}

// GetCode
func (mdb *MemoryStateDB) GetCode(addr string) []byte {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.Data.GetCode()
	}
	return nil
}

// SetCode
func (mdb *MemoryStateDB) SetCode(addr string, code []byte) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		mdb.dataDirty[addr] = true
		acc.SetCode(code)
	}
}

// SetAbi   ABI
func (mdb *MemoryStateDB) SetAbi(addr, abi string) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		mdb.dataDirty[addr] = true
		acc.SetAbi(abi)
	}
}

// GetAbi   ABI
func (mdb *MemoryStateDB) GetAbi(addr string) string {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.Data.GetAbi()
	}
	return ""
}

// GetCodeSize
//    EXTCODESIZE
func (mdb *MemoryStateDB) GetCodeSize(addr string) int {
	code := mdb.GetCode(addr)
	if code != nil {
		return len(code)
	}
	return 0
}

// AddRefund      SSTORE   ，  Gas
func (mdb *MemoryStateDB) AddRefund(gas uint64) {
	mdb.addChange(refundChange{baseChange: baseChange{}, prev: mdb.refund})
	mdb.refund += gas
}

// GetRefund
func (mdb *MemoryStateDB) GetRefund() uint64 {
	return mdb.refund
}

// GetAccount
func (mdb *MemoryStateDB) GetAccount(addr string) *ContractAccount {
	if acc, ok := mdb.accounts[addr]; ok {
		return acc
	}
	//         ，
	contract := NewContractAccount(addr, mdb)
	contract.LoadContract(mdb.StateDB)
	if contract.Empty() {
		return nil
	}
	mdb.accounts[addr] = contract
	return contract
}

// GetState SLOAD
func (mdb *MemoryStateDB) GetState(addr string, key common.Hash) common.Hash {
	//
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.GetState(key)
	}
	return common.Hash{}
}

// SetState SSTORE
func (mdb *MemoryStateDB) SetState(addr string, key common.Hash, value common.Hash) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		acc.SetState(key, value)
		//
		cfg := mdb.api.GetConfig()
		if !cfg.IsDappFork(mdb.blockHeight, "evm", evmtypes.ForkEVMState) {
			mdb.stateDirty[addr] = true
		}
	}
}

// TransferStateData
func (mdb *MemoryStateDB) TransferStateData(addr string) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		acc.TransferState()
	}
}

// UpdateState                 ，
func (mdb *MemoryStateDB) UpdateState(addr string) {
	mdb.stateDirty[addr] = true
}

// Suicide SELFDESTRUCT
//      ，        ，       ，
func (mdb *MemoryStateDB) Suicide(addr string) bool {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		mdb.addChange(suicideChange{
			baseChange: baseChange{},
			account:    addr,
			prev:       acc.State.GetSuicided(),
		})
		mdb.stateDirty[addr] = true
		return acc.Suicide()
	}
	return false
}

// HasSuicided
//
func (mdb *MemoryStateDB) HasSuicided(addr string) bool {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.HasSuicided()
	}
	return false
}

// Exist
func (mdb *MemoryStateDB) Exist(addr string) bool {
	return mdb.GetAccount(addr) != nil
}

// Empty
func (mdb *MemoryStateDB) Empty(addr string) bool {
	acc := mdb.GetAccount(addr)

	//         ，
	if acc != nil && !acc.Empty() {
		return false
	}

	//      ，
	if mdb.GetBalance(addr) != 0 {
		return false
	}
	return true
}

// RevertToSnapshot               （            ）
func (mdb *MemoryStateDB) RevertToSnapshot(version int) {
	if version >= len(mdb.snapshots) {
		return
	}

	ver := mdb.snapshots[version]

	//        ，
	if ver == nil || ver.id != version {
		log15.Crit(fmt.Errorf("Snapshot id %v cannot be reverted", version).Error())
		return
	}

	//
	for index := len(mdb.snapshots) - 1; index >= version; index-- {
		mdb.snapshots[index].revert()
	}

	//
	mdb.snapshots = mdb.snapshots[:version]
	mdb.versionID = version
	if version == 0 {
		mdb.currentVer = nil
	} else {
		mdb.currentVer = mdb.snapshots[version-1]
	}

}

// Snapshot            ，        ，
func (mdb *MemoryStateDB) Snapshot() int {
	id := mdb.versionID
	mdb.versionID++
	mdb.currentVer = &Snapshot{id: id, statedb: mdb}
	mdb.snapshots = append(mdb.snapshots, mdb.currentVer)
	return id
}

// GetLastSnapshot
func (mdb *MemoryStateDB) GetLastSnapshot() *Snapshot {
	if mdb.versionID == 0 {
		return nil
	}
	return mdb.snapshots[mdb.versionID-1]
}

// GetReceiptLogs
func (mdb *MemoryStateDB) GetReceiptLogs(addr string) (logs []*types.ReceiptLog) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		if mdb.stateDirty[addr] != nil {
			stateLog := acc.BuildStateLog()
			if stateLog != nil {
				logs = append(logs, stateLog)
			}
		}

		if mdb.dataDirty[addr] != nil {
			logs = append(logs, acc.BuildDataLog())
		}
		return
	}
	return
}

// GetChangedData
//                  MemoryStateDB，  ，        0   ，
//          0          ；
//   ，             ，        ，         ，          。
func (mdb *MemoryStateDB) GetChangedData(version int) (kvSet []*types.KeyValue, logs []*types.ReceiptLog) {
	if version < 0 {
		return
	}

	for _, snapshot := range mdb.snapshots {
		kv, log := snapshot.getData()
		if kv != nil {
			kvSet = append(kvSet, kv...)
		}

		if log != nil {
			logs = append(logs, log...)
		}
	}
	return
}

// CanTransfer   coins
func (mdb *MemoryStateDB) CanTransfer(sender, recipient string, amount uint64) bool {

	log15.Debug("check CanTransfer", "sender", sender, "recipient", recipient, "amount", amount)

	tType, errInfo := mdb.checkTransfer(sender, recipient, amount)

	if errInfo != nil {
		log15.Error("check transfer error", "sender", sender, "recipient", recipient, "amount", amount, "err info", errInfo)
		return false
	}

	value := int64(amount)
	if value < 0 {
		return false
	}

	switch tType {
	case NoNeed:
		return true
	case ToExec:
		//                   ，
		accFrom := mdb.CoinsAccount.LoadExecAccount(sender, recipient)
		b := accFrom.GetBalance() - value
		if b < 0 {
			log15.Error("check transfer error", "error info", types.ErrNoBalance)
			return false
		}
		return true
	case FromExec:
		return mdb.checkExecAccount(sender, value)
	default:
		return false
	}
}
func (mdb *MemoryStateDB) checkExecAccount(execAddr string, value int64) bool {
	var err error
	defer func() {
		if err != nil {
			log15.Error("checkExecAccount error", "error info", err)
		}
	}()
	//        ，
	if !types.CheckAmount(value) {
		err = types.ErrAmount
		return false
	}
	contract := mdb.GetAccount(execAddr)
	if contract == nil {
		err = model.ErrAddrNotExists
		return false
	}
	creator := contract.GetCreator()
	if len(creator) == 0 {
		err = model.ErrNoCreator
		return false
	}

	var accFrom *types.Account
	cfg := mdb.api.GetConfig()
	if cfg.IsDappFork(mdb.GetBlockHeight(), "evm", evmtypes.ForkEVMFrozen) {
		//    ，
		accFrom = mdb.CoinsAccount.LoadExecAccount(execAddr, execAddr)
	} else {
		accFrom = mdb.CoinsAccount.LoadExecAccount(creator, execAddr)
	}
	balance := accFrom.GetBalance()
	remain := balance - value
	if remain < 0 {
		err = types.ErrNoBalance
		return false
	}
	return true
}

// TransferType
type TransferType int

const (
	_ TransferType = iota
	// NoNeed
	NoNeed
	// ToExec
	ToExec
	// FromExec
	FromExec
	// Error
	Error
)

func (mdb *MemoryStateDB) checkTransfer(sender, recipient string, amount uint64) (tType TransferType, err error) {
	if amount == 0 {
		return NoNeed, nil
	}
	if mdb.CoinsAccount == nil {
		log15.Error("no coinsaccount exists", "sender", sender, "recipient", recipient, "amount", amount)
		return Error, model.ErrNoCoinsAccount
	}

	//              ，
	execSender := mdb.Exist(sender)
	execRecipient := mdb.Exist(recipient)

	if execRecipient && execSender {
		//         ，
		err = model.ErrTransferBetweenContracts
		tType = Error
	} else if execSender {
		//              （                 ）
		tType = FromExec
		err = nil
	} else if execRecipient {
		//
		tType = ToExec
		err = nil
	} else {
		//         ，
		err = model.ErrTransferBetweenEOA
		tType = Error
	}

	return tType, err
}

// Transfer   coins
//              ，
func (mdb *MemoryStateDB) Transfer(sender, recipient string, amount uint64) bool {
	log15.Debug("transfer from contract to external(contract)", "sender", sender, "recipient", recipient, "amount", amount)

	tType, errInfo := mdb.checkTransfer(sender, recipient, amount)

	if errInfo != nil {
		log15.Error("transfer error", "sender", sender, "recipient", recipient, "amount", amount, "err info", errInfo)
		return false
	}

	var (
		ret *types.Receipt
		err error
	)

	value := int64(amount)
	if value < 0 {
		return false
	}

	switch tType {
	case NoNeed:
		return true
	case ToExec:
		ret, err = mdb.transfer2Contract(sender, recipient, value)
	case FromExec:
		ret, err = mdb.transfer2External(sender, recipient, value)
	default:
		return false
	}

	//                ，    sender    ，
	if err != nil {
		log15.Error("transfer error", "sender", sender, "recipient", recipient, "amount", amount, "err info", err)
		return false
	}
	if ret != nil {
		mdb.addChange(transferChange{
			baseChange: baseChange{},
			amount:     value,
			data:       ret.KV,
			logs:       ret.Logs,
		})
	}
	return true
}

//   dplatformos   ，                  ：
// A   X   <-> B   X  ；
//        ，      EVM                  ，  A  B   X    ，       ：
// A -> A:X -> B:X；  (             )
//             ;
func (mdb *MemoryStateDB) transfer2Contract(sender, recipient string, amount int64) (ret *types.Receipt, err error) {
	//
	contract := mdb.GetAccount(recipient)
	if contract == nil {
		return nil, model.ErrAddrNotExists
	}
	creator := contract.GetCreator()
	if len(creator) == 0 {
		return nil, model.ErrNoCreator
	}
	execAddr := recipient

	ret = &types.Receipt{}

	cfg := mdb.api.GetConfig()
	if cfg.IsDappFork(mdb.GetBlockHeight(), "evm", evmtypes.ForkEVMFrozen) {
		//         ，         execAddr:execAddr
		rs, err := mdb.CoinsAccount.ExecTransfer(sender, execAddr, execAddr, amount)
		if err != nil {
			return nil, err
		}

		ret.KV = append(ret.KV, rs.KV...)
		ret.Logs = append(ret.Logs, rs.Logs...)
	} else {
		if strings.Compare(sender, creator) != 0 {
			//         ，
			rs, err := mdb.CoinsAccount.ExecTransfer(sender, creator, execAddr, amount)
			if err != nil {
				return nil, err
			}

			ret.KV = append(ret.KV, rs.KV...)
			ret.Logs = append(ret.Logs, rs.Logs...)
		}
	}

	return ret, nil
}

// dplatformos          Transfer2Contract ；
//                     ；
func (mdb *MemoryStateDB) transfer2External(sender, recipient string, amount int64) (ret *types.Receipt, err error) {
	//
	contract := mdb.GetAccount(sender)
	if contract == nil {
		return nil, model.ErrAddrNotExists
	}
	creator := contract.GetCreator()
	if len(creator) == 0 {
		return nil, model.ErrNoCreator
	}

	execAddr := sender

	cfg := mdb.api.GetConfig()
	if cfg.IsDappFork(mdb.GetBlockHeight(), "evm", evmtypes.ForkEVMFrozen) {
		//           ，
		ret, err = mdb.CoinsAccount.ExecTransfer(execAddr, recipient, execAddr, amount)
		if err != nil {
			return nil, err
		}
	} else {
		//
		//               ，
		if strings.Compare(creator, recipient) != 0 {
			ret, err = mdb.CoinsAccount.ExecTransfer(creator, recipient, execAddr, amount)
			if err != nil {
				return nil, err
			}
		}
	}
	return ret, nil
}

func (mdb *MemoryStateDB) mergeResult(one, two *types.Receipt) (ret *types.Receipt) {
	ret = one
	if ret == nil {
		ret = two
	} else if two != nil {
		ret.KV = append(ret.KV, two.KV...)
		ret.Logs = append(ret.Logs, two.Logs...)
	}
	return
}

// AddLog LOG0-4
//          ，
func (mdb *MemoryStateDB) AddLog(log *model.ContractLog) {
	mdb.addChange(addLogChange{txhash: mdb.txHash})
	log.TxHash = mdb.txHash
	log.Index = int(mdb.logSize)
	mdb.logs[mdb.txHash] = append(mdb.logs[mdb.txHash], log)
	mdb.logSize++
}

// AddPreimage   sha3
func (mdb *MemoryStateDB) AddPreimage(hash common.Hash, data []byte) {
	//
	if _, ok := mdb.preimages[hash]; !ok {
		mdb.addChange(addPreimageChange{hash: hash})
		pi := make([]byte, len(data))
		copy(pi, data)
		mdb.preimages[hash] = pi
	}
}

// PrintLogs                   （   ）
//                ，            ，        ，
func (mdb *MemoryStateDB) PrintLogs() {
	items := mdb.logs[mdb.txHash]
	for _, item := range items {
		item.PrintLog()
	}
}

// WritePreimages          preimages
func (mdb *MemoryStateDB) WritePreimages(number int64) {
	for k, v := range mdb.preimages {
		log15.Debug("Contract preimages ", "key:", k.Str(), "value:", common.Bytes2Hex(v), "block height:", number)
	}
}

// ResetDatas    ，
func (mdb *MemoryStateDB) ResetDatas() {
	mdb.currentVer = nil
	mdb.snapshots = mdb.snapshots[:0]
}

// GetBlockHeight
func (mdb *MemoryStateDB) GetBlockHeight() int64 {
	return mdb.blockHeight
}

// GetConfig
func (mdb *MemoryStateDB) GetConfig() *types.DplatformOSConfig {
	return mdb.api.GetConfig()
}
