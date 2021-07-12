// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/client"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
)

//action
type action struct {
	coinsAccount *account.DB
	db           dbm.KV
	localdb      dbm.KVDB
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	index        int32
	execaddr     string
	api          client.QueueProtocolAPI
}

func newAction(t *MultiSig, tx *types.Transaction, index int32) *action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &action{t.GetCoinsAccount(), t.GetStateDB(), t.GetLocalDB(), hash, fromaddr,
		t.GetBlockTime(), t.GetHeight(), index, dapp.ExecAddress(string(tx.Execer)), t.GetAPI()}
}

//MultiSigAccCreate
func (a *action) MultiSigAccCreate(accountCreate *mty.MultiSigAccCreate) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var totalweight uint64
	var ownerCount int
	var dailyLimit mty.DailyLimit

	//
	if accountCreate == nil {
		return nil, types.ErrInvalidParam
	}
	//   requiredweight          owner
	for _, owner := range accountCreate.Owners {
		if owner != nil {
			totalweight += owner.Weight
			ownerCount = ownerCount + 1
		}
	}

	if accountCreate.RequiredWeight > totalweight {
		return nil, mty.ErrRequiredweight
	}

	//         owner
	if ownerCount < mty.MinOwnersInit {
		return nil, mty.ErrOwnerLessThanTwo
	}
	//owner
	if ownerCount > mty.MaxOwnersCount {
		return nil, mty.ErrMaxOwnerCount
	}

	multiSigAccount := &mty.MultiSig{}
	multiSigAccount.CreateAddr = a.fromaddr
	multiSigAccount.Owners = accountCreate.Owners
	multiSigAccount.TxCount = 0
	multiSigAccount.RequiredWeight = accountCreate.RequiredWeight

	//
	if accountCreate.DailyLimit != nil {
		symbol := accountCreate.DailyLimit.Symbol
		execer := accountCreate.DailyLimit.Execer
		err := mty.IsAssetsInvalid(execer, symbol)
		if err != nil {
			return nil, err
		}
		dailyLimit.Symbol = symbol
		dailyLimit.Execer = execer
		dailyLimit.DailyLimit = accountCreate.DailyLimit.DailyLimit
		dailyLimit.SpentToday = 0
		dailyLimit.LastDay = a.blocktime //types.Now().Unix()
		multiSigAccount.DailyLimits = append(multiSigAccount.DailyLimits, &dailyLimit)
	}
	//       txhash              NewAddrFromString
	addr := address.MultiSignAddress(a.txhash)
	//
	multiSig, err := getMultiSigAccFromDb(a.db, addr)
	if err == nil && multiSig != nil {
		return nil, mty.ErrAccountHasExist
	}

	multiSigAccount.MultiSigAddr = addr
	receiptLog := &types.ReceiptLog{}
	receiptLog.Ty = mty.TyLogMultiSigAccCreate
	receiptLog.Log = types.Encode(multiSigAccount)
	logs = append(logs, receiptLog)

	key, value := setMultiSigAccToDb(a.db, multiSigAccount)
	kv = append(kv, &types.KeyValue{Key: key, Value: value})

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//MultiSigAccOperate            ：weight
//  requiredweight           owner
func (a *action) MultiSigAccOperate(AccountOperate *mty.MultiSigAccOperate) (*types.Receipt, error) {

	if AccountOperate == nil {
		return nil, types.ErrInvalidParam
	}
	//   statedb   MultiSigAccAddr
	multiSigAccAddr := AccountOperate.MultiSigAccAddr
	multiSigAccount, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("MultiSigAccountOperate", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, err
	}

	if multiSigAccount == nil {
		multisiglog.Error("MultiSigAccountOperate:getMultiSigAccFromDb is nil", "MultiSigAccAddr", multiSigAccAddr)
		return nil, types.ErrAccountNotExist
	}

	//              owner
	owneraddr := a.fromaddr
	ownerWeight, isowner := isOwner(multiSigAccount, owneraddr)
	if !isowner {
		return nil, mty.ErrIsNotOwner
	}

	//dailylimit             assets
	if !AccountOperate.OperateFlag {
		execer := AccountOperate.DailyLimit.Execer
		symbol := AccountOperate.DailyLimit.Symbol
		err := mty.IsAssetsInvalid(execer, symbol)
		if err != nil {
			return nil, err
		}
	}
	//    txid,          Txs
	txID := multiSigAccount.TxCount
	newMultiSigTx := &mty.MultiSigTx{}
	newMultiSigTx.Txid = txID
	newMultiSigTx.TxHash = hex.EncodeToString(a.txhash)
	newMultiSigTx.Executed = false
	newMultiSigTx.TxType = mty.AccountOperate
	newMultiSigTx.MultiSigAddr = multiSigAccAddr
	confirmOwner := &mty.Owner{OwnerAddr: owneraddr, Weight: ownerWeight}
	newMultiSigTx.ConfirmedOwner = append(newMultiSigTx.ConfirmedOwner, confirmOwner)

	return a.executeAccOperateTx(multiSigAccount, newMultiSigTx, AccountOperate, confirmOwner, true)
}

//MultiSigOwnerOperate       owner     ：owner add/del/replace
// del replace owner          owner       requiredweight
func (a *action) MultiSigOwnerOperate(AccOwnerOperate *mty.MultiSigOwnerOperate) (*types.Receipt, error) {
	multiSigAccAddr := AccOwnerOperate.MultiSigAccAddr

	//   statedb   MultiSigAccAddr
	multiSigAccount, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("MultiSigAccountOperate", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, err
	}
	if multiSigAccount == nil {
		multisiglog.Error("MultiSigAccountOwnerOperate:getMultiSigAccFromDb is nil", "MultiSigAccAddr", multiSigAccAddr)
		return nil, types.ErrAccountNotExist
	}

	//              owner
	owneraddr := a.fromaddr
	ownerWeight, isowner := isOwner(multiSigAccount, owneraddr)
	if !isowner {
		return nil, mty.ErrIsNotOwner
	}

	//    txid,          Txs
	txID := multiSigAccount.TxCount
	newMultiSigTx := &mty.MultiSigTx{}
	newMultiSigTx.Txid = txID
	newMultiSigTx.TxHash = hex.EncodeToString(a.txhash)
	newMultiSigTx.Executed = false
	newMultiSigTx.TxType = mty.OwnerOperate
	newMultiSigTx.MultiSigAddr = multiSigAccAddr
	confirmOwner := &mty.Owner{OwnerAddr: owneraddr, Weight: ownerWeight}
	newMultiSigTx.ConfirmedOwner = append(newMultiSigTx.ConfirmedOwner, confirmOwner)

	return a.executeOwnerOperateTx(multiSigAccount, newMultiSigTx, AccOwnerOperate, confirmOwner, true)
}

//MultiSigExecTransferFrom                  ，         ，  ExecTransferFrozen
//
//                ，multiSigAddr--->Addr
func (a *action) MultiSigExecTransferFrom(multiSigAccTransfer *mty.MultiSigExecTransferFrom) (*types.Receipt, error) {

	//   statedb   MultiSigAccAddr
	multiSigAccAddr := multiSigAccTransfer.From
	multiSigAcc, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("MultiSigAccExecTransfer", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, err
	}

	// to
	multiSigAccTo, err := getMultiSigAccFromDb(a.db, multiSigAccTransfer.To)
	if multiSigAccTo != nil && err == nil {
		multisiglog.Error("MultiSigExecTransferFrom", "multiSigAccTo", multiSigAccTo, "ToAddr", multiSigAccTransfer.To)
		return nil, mty.ErrAddrNotSupport
	}

	//              owner
	owneraddr := a.fromaddr
	ownerWeight, isowner := isOwner(multiSigAcc, owneraddr)
	if !isowner {
		return nil, mty.ErrIsNotOwner
	}

	//assete
	err = mty.IsAssetsInvalid(multiSigAccTransfer.Execname, multiSigAccTransfer.Symbol)
	if err != nil {
		return nil, err
	}
	//    txid,          Txs
	txID := multiSigAcc.TxCount
	newMultiSigTx := &mty.MultiSigTx{}
	newMultiSigTx.Txid = txID
	newMultiSigTx.TxHash = hex.EncodeToString(a.txhash)
	newMultiSigTx.Executed = false
	newMultiSigTx.TxType = mty.TransferOperate
	newMultiSigTx.MultiSigAddr = multiSigAccAddr
	confirmOwner := &mty.Owner{OwnerAddr: owneraddr, Weight: ownerWeight}
	newMultiSigTx.ConfirmedOwner = append(newMultiSigTx.ConfirmedOwner, confirmOwner)

	//
	return a.executeTransferTx(multiSigAcc, newMultiSigTx, multiSigAccTransfer, confirmOwner, mty.IsSubmit)
}

//MultiSigExecTransferTo             Execname.Symbol           ，from:Addr --->to:multiSigAddr
// from    tx       ，payload from       TransferTo
func (a *action) MultiSigExecTransferTo(execTransfer *mty.MultiSigExecTransferTo) (*types.Receipt, error) {

	//from
	multiSigAccFrom, err := getMultiSigAccFromDb(a.db, a.fromaddr)
	if multiSigAccFrom != nil && err == nil {
		multisiglog.Error("MultiSigExecTransferTo", "multiSigAccFrom", multiSigAccFrom, "From", a.fromaddr)
		return nil, mty.ErrAddrNotSupport
	}
	// to
	multiSigAccTo, err := getMultiSigAccFromDb(a.db, execTransfer.To)
	if multiSigAccTo == nil || err != nil {
		multisiglog.Error("MultiSigExecTransferTo", "ToAddr", execTransfer.To)
		return nil, mty.ErrAddrNotSupport
	}
	//assete
	err = mty.IsAssetsInvalid(execTransfer.Execname, execTransfer.Symbol)
	if err != nil {
		return nil, err
	}

	//          balance          balance
	symbol := getRealSymbol(execTransfer.Symbol)
	cfg := a.api.GetConfig()
	newAccountDB, err := account.NewAccountDB(cfg, execTransfer.Execname, symbol, a.db)
	if err != nil {
		return nil, err
	}
	receiptExecTransfer, err := newAccountDB.ExecTransfer(a.fromaddr, execTransfer.To, a.execaddr, execTransfer.Amount)
	if err != nil {
		multisiglog.Error("MultiSigExecTransfer:ExecTransfer", "From", a.fromaddr,
			"To", execTransfer.To, "execaddr", a.execaddr,
			"amount", execTransfer.Amount, "Execer", execTransfer.Execname, "Symbol", execTransfer.Symbol, "error", err)
		return nil, err
	}
	//         balance
	receiptExecFrozen, err := newAccountDB.ExecFrozen(execTransfer.To, a.execaddr, execTransfer.Amount)
	if err != nil {
		multisiglog.Error("MultiSigExecTransfer:ExecFrozen", "addr", execTransfer.To, "execaddr", a.execaddr,
			"amount", execTransfer.Amount, "Execer", execTransfer.Execname, "Symbol", execTransfer.Symbol, "error", err)
		return nil, err
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	logs = append(logs, receiptExecTransfer.Logs...)
	logs = append(logs, receiptExecFrozen.Logs...)
	kv = append(kv, receiptExecTransfer.KV...)
	kv = append(kv, receiptExecFrozen.KV...)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//MultiSigConfirmTx        MultiSigAcc  Transfer
//              ，         ，  ExecTransferFrozen
//             owner
//      ，              ，         owner
func (a *action) MultiSigConfirmTx(ConfirmTx *mty.MultiSigConfirmTx) (*types.Receipt, error) {

	//   statedb   MultiSigAccAddr
	multiSigAccAddr := ConfirmTx.MultiSigAccAddr
	multiSigAcc, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("MultiSigConfirmTx:getMultiSigAccFromDb", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, err
	}
	//              owner
	owneraddr := a.fromaddr
	ownerWeight, isowner := isOwner(multiSigAcc, owneraddr)
	if !isowner {
		multisiglog.Error("MultiSigConfirmTx: is not Owner", "MultiSigAccAddr", multiSigAccAddr, "txFrom", owneraddr, "err", err)
		return nil, mty.ErrIsNotOwner
	}
	//TxId
	if ConfirmTx.TxId > multiSigAcc.TxCount {
		multisiglog.Error("MultiSigConfirmTx: Invalid Txid", "MultiSigAccTxCount", multiSigAcc.TxCount, "Confirm TxId", ConfirmTx.TxId, "err", err)
		return nil, mty.ErrInvalidTxid
	}
	//         txid
	multiSigTx, err := getMultiSigAccTxFromDb(a.db, multiSigAccAddr, ConfirmTx.TxId)
	if err != nil {
		multisiglog.Error("MultiSigConfirmTx:getMultiSigAccTxFromDb", "multiSigAccAddr", multiSigAccAddr, "Confirm TxId", ConfirmTx.TxId, "err", err)
		return nil, mty.ErrTxidNotExist
	}
	//              /
	if multiSigTx.Executed {
		return nil, mty.ErrTxHasExecuted
	}
	// owneraddr        txid
	findindex, exist := isOwnerConfirmedTx(multiSigTx, owneraddr)

	//
	if exist && ConfirmTx.ConfirmOrRevoke {
		return nil, mty.ErrDupConfirmed
	}
	//
	if !exist && !ConfirmTx.ConfirmOrRevoke {
		return nil, mty.ErrConfirmNotExist
	}

	owner := &mty.Owner{OwnerAddr: owneraddr, Weight: ownerWeight}

	//          ， owneraddr
	if exist && !ConfirmTx.ConfirmOrRevoke {
		multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner[0:findindex], multiSigTx.ConfirmedOwner[findindex+1:]...)
	} else if !exist && ConfirmTx.ConfirmOrRevoke {
		//   owner      multiSigTx
		multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner, owner)
	}

	multiSigTxOwner := &mty.MultiSigTxOwner{MultiSigAddr: multiSigAccAddr, Txid: ConfirmTx.TxId, ConfirmedOwner: owner}
	isConfirm := isConfirmed(multiSigAcc.RequiredWeight, multiSigTx)

	//               ，  MultiSigConfirmTx receiptLog
	if !isConfirm || !ConfirmTx.ConfirmOrRevoke {
		return a.confirmTransaction(multiSigTx, multiSigTxOwner, ConfirmTx.ConfirmOrRevoke)
	}
	//  txhash
	tx, err := getTxByHash(a.api, multiSigTx.TxHash)
	if err != nil {
		return nil, err
	}
	payload, err := getMultiSigTxPayload(tx)
	if err != nil {
		return nil, err
	}

	//                  ，     owner/account
	if multiSigTx.TxType == mty.OwnerOperate && payload != nil {
		transfer := payload.GetMultiSigOwnerOperate()
		return a.executeOwnerOperateTx(multiSigAcc, multiSigTx, transfer, owner, false)
	} else if multiSigTx.TxType == mty.AccountOperate {
		transfer := payload.GetMultiSigAccOperate()
		return a.executeAccOperateTx(multiSigAcc, multiSigTx, transfer, owner, false)
	} else if multiSigTx.TxType == mty.TransferOperate {
		transfer := payload.GetMultiSigExecTransferFrom()
		return a.executeTransferTx(multiSigAcc, multiSigTx, transfer, owner, mty.IsConfirm)
	}
	multisiglog.Error("MultiSigConfirmTx:GetMultiSigTx", "multiSigAccAddr", multiSigAccAddr, "Confirm TxId", ConfirmTx.TxId, "TxType unknown", multiSigTx.TxType)
	return nil, mty.ErrTxTypeNoMatch
}

//             ,    KeyValue  ReceiptLog
func (a *action) multiSigWeightModify(multiSigAccAddr string, newRequiredWeight uint64) (*types.KeyValue, *types.ReceiptLog, error) {

	multiSigAccount, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("multiSigWeightModify", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, nil, err
	}

	if multiSigAccount == nil {
		multisiglog.Error("multiSigWeightModify:getMultiSigAccFromDb is nil", "MultiSigAccAddr", multiSigAccAddr)
		return nil, nil, types.ErrAccountNotExist
	}

	//      owner     ，    newRequiredWeight    owner
	var totalweight uint64
	receiptLog := &types.ReceiptLog{}
	for _, owner := range multiSigAccount.Owners {
		if owner != nil {
			totalweight += owner.Weight
		}
	}
	if newRequiredWeight > totalweight {
		return nil, nil, mty.ErrRequiredweight
	}
	//  RequiredWeight
	prevWeight := multiSigAccount.RequiredWeight
	multiSigAccount.RequiredWeight = newRequiredWeight

	//  receiptLog
	receiptWeight := &mty.ReceiptWeightModify{}
	receiptWeight.MultiSigAddr = multiSigAccount.MultiSigAddr
	receiptWeight.PrevWeight = prevWeight
	receiptWeight.CurrentWeight = multiSigAccount.RequiredWeight
	receiptLog.Ty = mty.TyLogMultiSigAccWeightModify
	receiptLog.Log = types.Encode(receiptWeight)

	key, value := setMultiSigAccToDb(a.db, multiSigAccount)
	kv := &types.KeyValue{Key: key, Value: value}
	return kv, receiptLog, nil
}

//                   ,
func (a *action) multiSigDailyLimitOperate(multiSigAccAddr string, dailylimit *mty.SymbolDailyLimit) (*types.KeyValue, *types.ReceiptLog, error) {

	multiSigAccount, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("multiSigDailyLimitOperate", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, nil, err
	}

	if multiSigAccount == nil {
		multisiglog.Error("multiSigDailyLimitOperate:getMultiSigAccFromDb is nil", "MultiSigAccAddr", multiSigAccAddr)
		return nil, nil, types.ErrAccountNotExist
	}

	flag := false
	var addOrModify bool
	var findindex int
	var curDailyLimit *mty.DailyLimit

	receiptLog := &types.ReceiptLog{}

	newSymbol := dailylimit.Symbol
	newExecer := dailylimit.Execer
	newDailyLimit := dailylimit.DailyLimit

	prevDailyLimit := &mty.DailyLimit{Symbol: newSymbol, Execer: newExecer, DailyLimit: 0, SpentToday: 0, LastDay: 0}
	//           symbol    ,
	for index, dailyLimit := range multiSigAccount.DailyLimits {
		if dailyLimit.Symbol == newSymbol && dailyLimit.Execer == newExecer {
			prevDailyLimit.DailyLimit = dailyLimit.DailyLimit
			prevDailyLimit.SpentToday = dailyLimit.SpentToday
			prevDailyLimit.LastDay = dailyLimit.LastDay
			flag = true
			findindex = index
			break
		}
	}
	if flag { //modify old DailyLimit
		multiSigAccount.DailyLimits[findindex].DailyLimit = newDailyLimit
		curDailyLimit = multiSigAccount.DailyLimits[findindex]
		addOrModify = false
	} else { //add new DailyLimit
		temDailyLimit := &mty.DailyLimit{}
		temDailyLimit.Symbol = newSymbol
		temDailyLimit.Execer = newExecer
		temDailyLimit.DailyLimit = newDailyLimit
		temDailyLimit.SpentToday = 0
		temDailyLimit.LastDay = a.blocktime //types.Now().Unix()
		multiSigAccount.DailyLimits = append(multiSigAccount.DailyLimits, temDailyLimit)

		curDailyLimit = temDailyLimit
		addOrModify = true
	}
	receiptDailyLimit := &mty.ReceiptDailyLimitOperate{
		MultiSigAddr:   multiSigAccount.MultiSigAddr,
		PrevDailyLimit: prevDailyLimit,
		CurDailyLimit:  curDailyLimit,
		AddOrModify:    addOrModify,
	}
	receiptLog.Ty = mty.TyLogMultiSigAccDailyLimitModify
	receiptLog.Log = types.Encode(receiptDailyLimit)

	key, value := setMultiSigAccToDb(a.db, multiSigAccount)
	kv := &types.KeyValue{Key: key, Value: value}
	return kv, receiptLog, nil
}

//         ,    KeyValue  ReceiptLog
func (a *action) multiSigOwnerAdd(multiSigAccAddr string, AccOwnerOperate *mty.MultiSigOwnerOperate) (*types.KeyValue, *types.ReceiptLog, error) {

	//  newowner    owner
	var newOwner mty.Owner
	newOwner.OwnerAddr = AccOwnerOperate.NewOwner
	newOwner.Weight = AccOwnerOperate.NewWeight
	return a.receiptOwnerAddOrDel(multiSigAccAddr, &newOwner, true)
}

//         ,    KeyValue  ReceiptLog
func (a *action) multiSigOwnerDel(multiSigAccAddr string, AccOwnerOperate *mty.MultiSigOwnerOperate) (*types.KeyValue, *types.ReceiptLog, error) {
	var owner mty.Owner
	owner.OwnerAddr = AccOwnerOperate.OldOwner
	owner.Weight = 0
	return a.receiptOwnerAddOrDel(multiSigAccAddr, &owner, false)
}

//  add/del owner receipt
func (a *action) receiptOwnerAddOrDel(multiSigAccAddr string, owner *mty.Owner, addOrDel bool) (*types.KeyValue, *types.ReceiptLog, error) {
	receiptLog := &types.ReceiptLog{}

	multiSigAcc, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("receiptOwnerAddOrDel", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, nil, err
	}

	if multiSigAcc == nil {
		multisiglog.Error("receiptOwnerAddOrDel:getMultiSigAccFromDb is nil", "MultiSigAccAddr", multiSigAccAddr)
		return nil, nil, types.ErrAccountNotExist
	}

	oldweight, index, totalWeight, totalowner, exist := getOwnerInfoByAddr(multiSigAcc, owner.OwnerAddr)

	if addOrDel {
		if exist {
			return nil, nil, mty.ErrOwnerExist
		}
		if totalowner >= mty.MaxOwnersCount {
			return nil, nil, mty.ErrMaxOwnerCount
		}
		multiSigAcc.Owners = append(multiSigAcc.Owners, owner)
		receiptLog.Ty = mty.TyLogMultiSigOwnerAdd
	} else {
		if !exist {
			return nil, nil, mty.ErrOwnerNotExist
		}
		//            owners         reqweight
		if totalWeight-oldweight < multiSigAcc.RequiredWeight {
			return nil, nil, mty.ErrTotalWeightNotEnough
		}
		//       owner
		if totalowner <= 1 {
			return nil, nil, mty.ErrOnlyOneOwner
		}
		owner.Weight = oldweight
		receiptLog.Ty = mty.TyLogMultiSigOwnerDel
		multiSigAcc.Owners = delOwner(multiSigAcc.Owners, index)
	}

	//  receiptLog
	receiptOwner := &mty.ReceiptOwnerAddOrDel{}
	receiptOwner.MultiSigAddr = multiSigAcc.MultiSigAddr
	receiptOwner.Owner = owner
	receiptOwner.AddOrDel = addOrDel
	receiptLog.Log = types.Encode(receiptOwner)

	key, value := setMultiSigAccToDb(a.db, multiSigAcc)
	keyValue := &types.KeyValue{Key: key, Value: value}
	return keyValue, receiptLog, nil
}

//      owner   ,    KeyValue  ReceiptLog
func (a *action) multiSigOwnerModify(multiSigAccAddr string, AccOwnerOperate *mty.MultiSigOwnerOperate) (*types.KeyValue, *types.ReceiptLog, error) {

	prev := &mty.Owner{OwnerAddr: AccOwnerOperate.OldOwner, Weight: 0}
	cur := &mty.Owner{OwnerAddr: AccOwnerOperate.OldOwner, Weight: AccOwnerOperate.NewWeight}
	return a.receiptOwnerModOrRep(multiSigAccAddr, prev, cur, true)
}

//      owner   ,    KeyValue  ReceiptLog
func (a *action) multiSigOwnerReplace(multiSigAccAddr string, AccOwnerOperate *mty.MultiSigOwnerOperate) (*types.KeyValue, *types.ReceiptLog, error) {

	prev := &mty.Owner{OwnerAddr: AccOwnerOperate.OldOwner, Weight: 0}
	cur := &mty.Owner{OwnerAddr: AccOwnerOperate.NewOwner, Weight: 0}
	return a.receiptOwnerModOrRep(multiSigAccAddr, prev, cur, false)
}

//    /  owner receipt
func (a *action) receiptOwnerModOrRep(multiSigAccAddr string, prev *mty.Owner, cur *mty.Owner, modOrRep bool) (*types.KeyValue, *types.ReceiptLog, error) {
	receiptLog := &types.ReceiptLog{}

	multiSigAcc, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("receiptOwnerModOrRep", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, nil, err
	}

	if multiSigAcc == nil {
		multisiglog.Error("receiptOwnerModOrRep:getMultiSigAccFromDb is nil", "MultiSigAccAddr", multiSigAccAddr)
		return nil, nil, types.ErrAccountNotExist
	}
	oldweight, index, totalWeight, _, exist := getOwnerInfoByAddr(multiSigAcc, prev.OwnerAddr)
	if modOrRep {
		if !exist {
			return nil, nil, mty.ErrOwnerNotExist
		}
		//            owners         reqweight
		if totalWeight-oldweight+cur.Weight < multiSigAcc.RequiredWeight {
			return nil, nil, mty.ErrTotalWeightNotEnough
		}
		prev.Weight = oldweight
		multiSigAcc.Owners[index].Weight = cur.Weight
		receiptLog.Ty = mty.TyLogMultiSigOwnerModify
	} else {
		if !exist {
			return nil, nil, mty.ErrOwnerNotExist
		}
		//   newowner
		_, _, _, _, find := getOwnerInfoByAddr(multiSigAcc, cur.OwnerAddr)
		if find {
			return nil, nil, mty.ErrNewOwnerExist
		}
		prev.Weight = oldweight
		cur.Weight = oldweight
		multiSigAcc.Owners[index].OwnerAddr = cur.OwnerAddr
		receiptLog.Ty = mty.TyLogMultiSigOwnerReplace
	}
	//  receiptLog
	receiptAddOwner := &mty.ReceiptOwnerModOrRep{}
	receiptAddOwner.MultiSigAddr = multiSigAcc.MultiSigAddr
	receiptAddOwner.PrevOwner = prev
	receiptAddOwner.CurrentOwner = cur
	receiptAddOwner.ModOrRep = modOrRep
	receiptLog.Log = types.Encode(receiptAddOwner)

	key, value := setMultiSigAccToDb(a.db, multiSigAcc)
	keyValue := &types.KeyValue{Key: key, Value: value}
	return keyValue, receiptLog, nil
}

//  AccExecTransfer receipt  ,              ，
func (a *action) receiptDailyLimitUpdate(multiSigAccAddr string, findindex int, curdailyLimit *mty.DailyLimit) (*types.KeyValue, *types.ReceiptLog, error) {

	multiSigAcc, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("receiptDailyLimitUpdate", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, nil, err
	}

	if multiSigAcc == nil {
		multisiglog.Error("receiptDailyLimitUpdate:getMultiSigAccFromDb is nil", "MultiSigAccAddr", multiSigAccAddr)
		return nil, nil, types.ErrAccountNotExist
	}
	receiptLog := &types.ReceiptLog{}

	//  receiptLog
	receipt := &mty.ReceiptAccDailyLimitUpdate{}
	receipt.MultiSigAddr = multiSigAcc.MultiSigAddr
	receipt.PrevDailyLimit = multiSigAcc.DailyLimits[findindex]
	receipt.CurDailyLimit = curdailyLimit
	receiptLog.Ty = mty.TyLogDailyLimitUpdate
	receiptLog.Log = types.Encode(receipt)

	//  DailyLimit
	multiSigAcc.DailyLimits[findindex].SpentToday = curdailyLimit.SpentToday
	multiSigAcc.DailyLimits[findindex].LastDay = curdailyLimit.LastDay

	key, value := setMultiSigAccToDb(a.db, multiSigAcc)
	keyValue := &types.KeyValue{Key: key, Value: value}
	return keyValue, receiptLog, nil
}

//                  receipt
func (a *action) receiptTxCountUpdate(multiSigAccAddr string) (*types.KeyValue, *types.ReceiptLog, error) {

	multiSigAcc, err := getMultiSigAccFromDb(a.db, multiSigAccAddr)
	if err != nil {
		multisiglog.Error("receiptTxCountUpdate", "MultiSigAccAddr", multiSigAccAddr, "err", err)
		return nil, nil, err
	}

	if multiSigAcc == nil {
		multisiglog.Error("receiptTxCountUpdate:getMultiSigAccFromDb is nil", "MultiSigAccAddr", multiSigAccAddr)
		return nil, nil, types.ErrAccountNotExist
	}

	receiptLog := &types.ReceiptLog{}

	//  receiptLog
	multiSigAcc.TxCount++

	receiptLogTxCount := &mty.ReceiptTxCountUpdate{
		MultiSigAddr: multiSigAcc.MultiSigAddr,
		CurTxCount:   multiSigAcc.TxCount,
	}

	receiptLog.Ty = mty.TyLogTxCountUpdate
	receiptLog.Log = types.Encode(receiptLogTxCount)

	key, value := setMultiSigAccToDb(a.db, multiSigAcc)
	keyValue := &types.KeyValue{Key: key, Value: value}
	return keyValue, receiptLog, nil
}

//  MultiSigAccTx receipt
func (a *action) receiptMultiSigTx(multiSigTx *mty.MultiSigTx, owner *mty.Owner, prevExecutes, subOrConfirm bool) (*types.KeyValue, *types.ReceiptLog) {
	receiptLog := &types.ReceiptLog{}

	//  receiptLog
	receiptLogTx := &mty.ReceiptMultiSigTx{}
	multiSigTxOwner := &mty.MultiSigTxOwner{MultiSigAddr: multiSigTx.MultiSigAddr, Txid: multiSigTx.Txid, ConfirmedOwner: owner}

	receiptLogTx.MultiSigTxOwner = multiSigTxOwner
	receiptLogTx.PrevExecuted = prevExecutes
	receiptLogTx.CurExecuted = multiSigTx.Executed
	receiptLogTx.SubmitOrConfirm = subOrConfirm
	if subOrConfirm {
		receiptLogTx.TxHash = multiSigTx.TxHash
		receiptLogTx.TxType = multiSigTx.TxType
	}

	receiptLog.Ty = mty.TyLogMultiSigTx
	receiptLog.Log = types.Encode(receiptLogTx)

	key, value := setMultiSigAccTxToDb(a.db, multiSigTx)
	keyValue := &types.KeyValue{Key: key, Value: value}
	return keyValue, receiptLog
}

//         ：  submitTx confirmtx  。
func (a *action) executeTransferTx(multiSigAcc *mty.MultiSig, newMultiSigTx *mty.MultiSigTx, transfer *mty.MultiSigExecTransferFrom, confOwner *mty.Owner, subOrConfirm bool) (*types.Receipt, error) {

	//
	var findindex int
	curDailyLimit := &mty.DailyLimit{Symbol: transfer.Symbol, Execer: transfer.Execname, DailyLimit: 0, SpentToday: 0, LastDay: 0}
	for Index, dailyLimit := range multiSigAcc.DailyLimits {
		if dailyLimit.Symbol == transfer.Symbol && dailyLimit.Execer == transfer.Execname {
			curDailyLimit.DailyLimit = dailyLimit.DailyLimit
			curDailyLimit.SpentToday = dailyLimit.SpentToday
			curDailyLimit.LastDay = dailyLimit.LastDay
			findindex = Index
			break
		}
	}
	//     0
	if curDailyLimit.DailyLimit == 0 {
		return nil, mty.ErrDailyLimitIsZero
	}

	//                ，
	amount := transfer.Amount
	confirmed := isConfirmed(multiSigAcc.RequiredWeight, newMultiSigTx)
	underLimit, newlastday := isUnderLimit(a.blocktime, uint64(amount), curDailyLimit)

	//      lastday spenttoday
	if newlastday != 0 {
		curDailyLimit.LastDay = newlastday
		curDailyLimit.SpentToday = 0
	}

	prevExecuted := newMultiSigTx.Executed

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	//            ，       ，          ，       ，             tx
	if confirmed || underLimit {

		//     ，              ， multiSig
		symbol := getRealSymbol(transfer.Symbol)
		cfg := a.api.GetConfig()
		execerAccDB, err := account.NewAccountDB(cfg, transfer.Execname, symbol, a.db)
		if err != nil {
			multisiglog.Error("executeTransaction:NewAccountDB", "From", transfer.From, "To", transfer.To,
				"execaddr", a.execaddr, "amount", amount, "Execer", transfer.Execname, "Symbol", transfer.Symbol, "error", err)
			return nil, err
		}
		receiptFromMultiSigAcc, err := execerAccDB.ExecTransferFrozen(transfer.From, transfer.To, a.execaddr, amount)
		if err != nil {
			multisiglog.Error("executeTransaction:ExecTransferFrozen", "From", transfer.From, "To", transfer.To,
				"execaddr", a.execaddr, "amount", amount, "Execer", transfer.Execname, "Symbol", transfer.Symbol, "error", err)
			return nil, err
		}
		logs = append(logs, receiptFromMultiSigAcc.Logs...)
		kv = append(kv, receiptFromMultiSigAcc.KV...)

		//
		newMultiSigTx.Executed = true

		//        ,
		if !confirmed && subOrConfirm {
			curDailyLimit.SpentToday += uint64(amount)
		}
	}

	//  multiSigAcc  :txcount    submit
	if subOrConfirm {
		keyvalue, receiptlog, err := a.receiptTxCountUpdate(multiSigAcc.MultiSigAddr)
		if err != nil {
			multisiglog.Error("executeTransaction:receiptTxCountUpdate", "error", err)
		}
		kv = append(kv, keyvalue)
		logs = append(logs, receiptlog)
	}

	//  multiSigAcc  :
	keyvalue, receiptlog, err := a.receiptDailyLimitUpdate(multiSigAcc.MultiSigAddr, findindex, curDailyLimit)
	if err != nil {
		multisiglog.Error("executeTransaction:receiptDailyLimitUpdate", "error", err)
	}
	//  newMultiSigTx   ：MultiSigTx      owner，
	keyvaluetx, receiptlogtx := a.receiptMultiSigTx(newMultiSigTx, confOwner, prevExecuted, subOrConfirm)

	logs = append(logs, receiptlog)
	logs = append(logs, receiptlogtx)
	kv = append(kv, keyvalue)
	kv = append(kv, keyvaluetx)

	//test
	multisiglog.Error("executeTransferTx", "multiSigAcc", multiSigAcc, "newMultiSigTx", newMultiSigTx)

	return &types.Receipt{
		Ty:   types.ExecOk,
		KV:   kv,
		Logs: logs,
	}, nil
}

//              ：  submitTx confirmtx  。
func (a *action) executeAccOperateTx(multiSigAcc *mty.MultiSig, newMultiSigTx *mty.MultiSigTx, accountOperate *mty.MultiSigAccOperate, confOwner *mty.Owner, subOrConfirm bool) (*types.Receipt, error) {

	//
	confirmed := isConfirmed(multiSigAcc.RequiredWeight, newMultiSigTx)
	prevExecuted := newMultiSigTx.Executed

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	var accAttrkv *types.KeyValue
	var accAttrReceiptLog *types.ReceiptLog
	var err error

	//           ，             tx
	if confirmed {
		//    RequiredWeight
		if accountOperate.OperateFlag {
			accAttrkv, accAttrReceiptLog, err = a.multiSigWeightModify(multiSigAcc.MultiSigAddr, accountOperate.NewRequiredWeight)
			if err != nil {
				multisiglog.Error("executeAccOperateTx", "multiSigWeightModify", err)
				return nil, err
			}
		} else { //
			accAttrkv, accAttrReceiptLog, err = a.multiSigDailyLimitOperate(multiSigAcc.MultiSigAddr, accountOperate.DailyLimit)
			if err != nil {
				multisiglog.Error("executeAccOperateTx", "multiSigDailyLimitOperate", err)
				return nil, err
			}
		}
		logs = append(logs, accAttrReceiptLog)
		kv = append(kv, accAttrkv)
		//
		newMultiSigTx.Executed = true
	}

	//  multiSigAcc  :txcount    submit
	if subOrConfirm {
		keyvalue, receiptlog, err := a.receiptTxCountUpdate(multiSigAcc.MultiSigAddr)
		if err != nil {
			multisiglog.Error("executeAccOperateTx:receiptTxCountUpdate", "error", err)
		}
		kv = append(kv, keyvalue)
		logs = append(logs, receiptlog)
	}
	//  newMultiSigTx   ：MultiSigTx      owner，
	keyvaluetx, receiptlogtx := a.receiptMultiSigTx(newMultiSigTx, confOwner, prevExecuted, subOrConfirm)
	logs = append(logs, receiptlogtx)
	kv = append(kv, keyvaluetx)
	return &types.Receipt{
		Ty:   types.ExecOk,
		KV:   kv,
		Logs: logs,
	}, nil
}

//       owner     ：  submitTx confirmtx  。
func (a *action) executeOwnerOperateTx(multiSigAccount *mty.MultiSig, newMultiSigTx *mty.MultiSigTx, accountOperate *mty.MultiSigOwnerOperate, confOwner *mty.Owner, subOrConfirm bool) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	var multiSigkv *types.KeyValue
	var receiptLog *types.ReceiptLog
	var err error
	//
	confirmed := isConfirmed(multiSigAccount.RequiredWeight, newMultiSigTx)
	prevExecuted := newMultiSigTx.Executed

	flag := accountOperate.OperateFlag

	//           ，             tx
	if confirmed {
		//add
		if mty.OwnerAdd == flag {
			multiSigkv, receiptLog, err = a.multiSigOwnerAdd(multiSigAccount.MultiSigAddr, accountOperate)
			if err != nil {
				multisiglog.Error("MultiSigAccountOwnerOperate", "multiSigOwnerAdd err", err)
				return nil, err
			}

		} else if mty.OwnerDel == flag {
			multiSigkv, receiptLog, err = a.multiSigOwnerDel(multiSigAccount.MultiSigAddr, accountOperate)
			if err != nil {
				multisiglog.Error("MultiSigAccountOwnerOperate", "multiSigOwnerAdd err", err)
				return nil, err
			}
		} else if mty.OwnerModify == flag { //modify owner
			multiSigkv, receiptLog, err = a.multiSigOwnerModify(multiSigAccount.MultiSigAddr, accountOperate)
			if err != nil {
				multisiglog.Error("MultiSigAccountOwnerOperate", "multiSigOwnerModify err", err)
				return nil, err
			}
		} else if mty.OwnerReplace == flag { //replace owner
			multiSigkv, receiptLog, err = a.multiSigOwnerReplace(multiSigAccount.MultiSigAddr, accountOperate)
			if err != nil {
				multisiglog.Error("MultiSigAccountOwnerOperate", "multiSigOwnerReplace err", err)
				return nil, err
			}
		} else {
			multisiglog.Error("MultiSigAccountOwnerOperate", "OperateFlag", flag)
			return nil, mty.ErrOperateType
		}
		logs = append(logs, receiptLog)
		kv = append(kv, multiSigkv)

		//
		newMultiSigTx.Executed = true
	}

	//  multiSigAcc  :txcount    submit
	if subOrConfirm {
		keyvalue, receiptlog, err := a.receiptTxCountUpdate(multiSigAccount.MultiSigAddr)
		if err != nil {
			multisiglog.Error("executeOwnerOperateTx:receiptTxCountUpdate", "error", err)
		}
		kv = append(kv, keyvalue)
		logs = append(logs, receiptlog)
	}
	//  newMultiSigTx   ：MultiSigTx      owner，
	keyvaluetx, receiptlogtx := a.receiptMultiSigTx(newMultiSigTx, confOwner, prevExecuted, subOrConfirm)
	logs = append(logs, receiptlogtx)
	kv = append(kv, keyvaluetx)

	//test
	multisiglog.Error("executeOwnerOperateTx", "multiSigAccount", multiSigAccount, "newMultiSigTx", newMultiSigTx)

	return &types.Receipt{
		Ty:   types.ExecOk,
		KV:   kv,
		Logs: logs,
	}, nil
}

//       receiptLog
func (a *action) confirmTransaction(multiSigTx *mty.MultiSigTx, multiSigTxOwner *mty.MultiSigTxOwner, ConfirmOrRevoke bool) (*types.Receipt, error) {
	receiptLog := &types.ReceiptLog{}

	receiptLogUnConfirmTx := &mty.ReceiptConfirmTx{MultiSigTxOwner: multiSigTxOwner, ConfirmeOrRevoke: ConfirmOrRevoke}
	if ConfirmOrRevoke {
		receiptLog.Ty = mty.TyLogMultiSigConfirmTx
	} else {
		receiptLog.Ty = mty.TyLogMultiSigConfirmTxRevoke
	}
	receiptLog.Log = types.Encode(receiptLogUnConfirmTx)

	//  MultiSigAccTx
	key, value := setMultiSigAccTxToDb(a.db, multiSigTx)
	kv := &types.KeyValue{Key: key, Value: value}

	return &types.Receipt{
		Ty:   types.ExecOk,
		KV:   []*types.KeyValue{kv},
		Logs: []*types.ReceiptLog{receiptLog},
	}, nil
}
