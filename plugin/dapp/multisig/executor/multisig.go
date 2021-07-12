// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

/*
multiSig          ：
//
//      owner     ：owner add/del/replace
//           ：weight
//
//                ，Addr --->multiSigAddr
//                ，multiSigAddr--->Addr
*/

import (
	"bytes"
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/client"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
)

var multisiglog = log.New("module", "execs.multisig")

var driverName = "multisig"

// Init multisig
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	drivers.Register(cfg, GetName(), newMultiSig, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&MultiSig{}))
}

// GetName multisig  name
func GetName() string {
	return newMultiSig().GetName()
}

// MultiSig multisig
type MultiSig struct {
	drivers.DriverBase
}

func newMultiSig() drivers.Driver {
	m := &MultiSig{}
	m.SetChild(m)
	m.SetExecutorType(types.LoadExecutorType(driverName))
	return m
}

// GetDriverName   multisig  name
func (m *MultiSig) GetDriverName() string {
	return driverName
}

// CheckTx   multisig    ,    amount
func (m *MultiSig) CheckTx(tx *types.Transaction, index int) error {
	ety := m.GetExecutorType()

	//amount check
	amount, err := ety.Amount(tx)
	if err != nil {
		return err
	}
	if amount < 0 {
		return types.ErrAmount
	}

	_, v, err := ety.DecodePayloadValue(tx)
	if err != nil {
		return err
	}
	payload := v.Interface()

	//MultiSigAccCreate
	if ato, ok := payload.(*mty.MultiSigAccCreate); ok {
		return checkAccountCreateTx(ato)
	}

	//MultiSigOwnerOperate
	if ato, ok := payload.(*mty.MultiSigOwnerOperate); ok {
		return checkOwnerOperateTx(ato)
	}
	//MultiSigAccOperate
	if ato, ok := payload.(*mty.MultiSigAccOperate); ok {
		return checkAccountOperateTx(ato)
	}
	//MultiSigConfirmTx
	if ato, ok := payload.(*mty.MultiSigConfirmTx); ok {
		if err := address.CheckMultiSignAddress(ato.GetMultiSigAccAddr()); err != nil {
			return types.ErrInvalidAddress
		}
		return nil
	}

	//MultiSigExecTransferTo
	if ato, ok := payload.(*mty.MultiSigExecTransferTo); ok {
		if err := address.CheckMultiSignAddress(ato.GetTo()); err != nil {
			return types.ErrInvalidAddress
		}
		//assets check
		return mty.IsAssetsInvalid(ato.GetExecname(), ato.GetSymbol())
	}
	//MultiSigExecTransferFrom
	if ato, ok := payload.(*mty.MultiSigExecTransferFrom); ok {
		//from addr check
		if err := address.CheckMultiSignAddress(ato.GetFrom()); err != nil {
			return types.ErrInvalidAddress
		}
		//to addr check
		if err := address.CheckAddress(ato.GetTo()); err != nil {
			return types.ErrInvalidAddress
		}
		//assets check
		return mty.IsAssetsInvalid(ato.GetExecname(), ato.GetSymbol())
	}

	return nil
}
func checkAccountCreateTx(ato *mty.MultiSigAccCreate) error {
	var totalweight uint64
	var ownerCount int

	requiredWeight := ato.GetRequiredWeight()
	if requiredWeight == 0 {
		return mty.ErrInvalidWeight
	}
	owners := ato.GetOwners()
	ownersMap := make(map[string]bool)

	//   requiredweight          owner
	for _, owner := range owners {
		if owner != nil {
			if err := address.CheckAddress(owner.OwnerAddr); err != nil {
				return types.ErrInvalidAddress
			}
			if owner.Weight == 0 {
				return mty.ErrInvalidWeight
			}
			if ownersMap[owner.OwnerAddr] {
				return mty.ErrOwnerExist
			}
			ownersMap[owner.OwnerAddr] = true
			totalweight += owner.Weight
			ownerCount = ownerCount + 1
		}
	}

	if ato.RequiredWeight > totalweight {
		return mty.ErrRequiredweight
	}

	//         owner
	if ownerCount < mty.MinOwnersInit {
		return mty.ErrOwnerLessThanTwo
	}
	//owner
	if ownerCount > mty.MaxOwnersCount {
		return mty.ErrMaxOwnerCount
	}

	dailyLimit := ato.GetDailyLimit()
	//assets check
	return mty.IsAssetsInvalid(dailyLimit.GetExecer(), dailyLimit.GetSymbol())
}

func checkOwnerOperateTx(ato *mty.MultiSigOwnerOperate) error {
	OldOwner := ato.GetOldOwner()
	NewOwner := ato.GetNewOwner()
	NewWeight := ato.GetNewWeight()
	MultiSigAccAddr := ato.GetMultiSigAccAddr()
	if err := address.CheckMultiSignAddress(MultiSigAccAddr); err != nil {
		return types.ErrInvalidAddress
	}

	if ato.OperateFlag == mty.OwnerAdd {
		if err := address.CheckAddress(NewOwner); err != nil {
			return types.ErrInvalidAddress
		}
		if NewWeight <= 0 {
			return mty.ErrInvalidWeight
		}
	}
	if ato.OperateFlag == mty.OwnerDel {
		if err := address.CheckAddress(OldOwner); err != nil {
			return types.ErrInvalidAddress
		}
	}
	if ato.OperateFlag == mty.OwnerModify {
		if err := address.CheckAddress(OldOwner); err != nil {
			return types.ErrInvalidAddress
		}
		if NewWeight <= 0 {
			return mty.ErrInvalidWeight
		}
	}
	if ato.OperateFlag == mty.OwnerReplace {
		if err := address.CheckAddress(OldOwner); err != nil {
			return types.ErrInvalidAddress
		}
		if err := address.CheckAddress(NewOwner); err != nil {
			return types.ErrInvalidAddress
		}
	}
	return nil
}
func checkAccountOperateTx(ato *mty.MultiSigAccOperate) error {
	//MultiSigAccOperate MultiSigAccAddr
	MultiSigAccAddr := ato.GetMultiSigAccAddr()
	if err := address.CheckMultiSignAddress(MultiSigAccAddr); err != nil {
		return types.ErrInvalidAddress
	}

	if ato.OperateFlag == mty.AccWeightOp {
		NewWeight := ato.GetNewRequiredWeight()
		if NewWeight <= 0 {
			return mty.ErrInvalidWeight
		}
	}
	if ato.OperateFlag == mty.AccDailyLimitOp {
		dailyLimit := ato.GetDailyLimit()
		//assets check
		return mty.IsAssetsInvalid(dailyLimit.GetExecer(), dailyLimit.GetSymbol())
	}
	return nil
}

//       Receipt
func (m *MultiSig) execLocalMultiSigReceipt(receiptData *types.ReceiptData, tx *types.Transaction, addOrRollback bool) ([]*types.KeyValue, error) {
	var set []*types.KeyValue
	for _, log := range receiptData.Logs {
		multisiglog.Info("execLocalMultiSigReceipt", "Ty", log.Ty)

		switch log.Ty {
		case mty.TyLogMultiSigAccCreate:
			{
				var receipt mty.MultiSig
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}

				kv, err := m.saveMultiSigAccCreate(&receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv...)
			}
		case mty.TyLogMultiSigOwnerAdd,
			mty.TyLogMultiSigOwnerDel:
			{
				var receipt mty.ReceiptOwnerAddOrDel
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}
				kv, err := m.saveMultiSigOwnerAddOrDel(receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv...)
			}
		case mty.TyLogMultiSigOwnerModify,
			mty.TyLogMultiSigOwnerReplace:
			{
				var receipt mty.ReceiptOwnerModOrRep
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}
				kv, err := m.saveMultiSigOwnerModOrRep(receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv...)
			}
		case mty.TyLogMultiSigAccWeightModify:
			{
				var receipt mty.ReceiptWeightModify
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}
				kv, err := m.saveMultiSigAccWeight(receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv...)
			}
		case mty.TyLogMultiSigAccDailyLimitAdd,
			mty.TyLogMultiSigAccDailyLimitModify:
			{
				var receipt mty.ReceiptDailyLimitOperate
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}
				kv, err := m.saveMultiSigAccDailyLimit(receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv...)
			}
		case mty.TyLogMultiSigConfirmTx, //         ，
			mty.TyLogMultiSigConfirmTxRevoke:
			{
				var receipt mty.ReceiptConfirmTx
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}
				kv, err := m.saveMultiSigConfirmTx(receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv...)
			}
		case mty.TyLogDailyLimitUpdate: //   DailyLimit
			{
				var receipt mty.ReceiptAccDailyLimitUpdate
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}
				kv, err := m.saveDailyLimitUpdate(receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv...)
			}
		case mty.TyLogMultiSigTx: //     owner
			{
				var receipt mty.ReceiptMultiSigTx
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}
				//         tx         owner
				kv1, err := m.saveMultiSigTx(receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv1...)

				//              amount    ,    submit confirm
				if receipt.CurExecuted {
					kv2, err := m.saveMultiSigTransfer(tx, receipt.SubmitOrConfirm, addOrRollback)
					if err != nil {
						return nil, err
					}
					set = append(set, kv2...)
				}
			}
		case mty.TyLogTxCountUpdate:
			{
				var receipt mty.ReceiptTxCountUpdate
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					return nil, err
				}
				kv, err := m.saveMultiSigTxCountUpdate(receipt, addOrRollback)
				if err != nil {
					return nil, err
				}
				set = append(set, kv...)
			}
		default:
			break
		}
	}
	return set, nil
}

//    to        ，Submit    tx。Confirm      txid
func (m *MultiSig) saveMultiSigTransfer(tx *types.Transaction, SubmitOrConfirm, addOrRollback bool) ([]*types.KeyValue, error) {
	var set []*types.KeyValue
	//      GetPayload
	var action mty.MultiSigAction
	err := types.Decode(tx.GetPayload(), &action)
	if err != nil {
		panic(err)
	}
	var to string
	var execname string
	var symbol string
	var amount int64

	//addr-->multiSigAccAddr
	//multiSigAccAddr-->addr
	if SubmitOrConfirm {
		if action.Ty == mty.ActionMultiSigExecTransferTo && action.GetMultiSigExecTransferTo() != nil {
			tx := action.GetMultiSigExecTransferTo()
			to = tx.To
			execname = tx.Execname
			symbol = tx.Symbol
			amount = tx.Amount
		} else if action.Ty == mty.ActionMultiSigExecTransferFrom && action.GetMultiSigExecTransferFrom() != nil {
			tx := action.GetMultiSigExecTransferFrom()
			to = tx.To
			execname = tx.Execname
			symbol = tx.Symbol
			amount = tx.Amount
		} else {
			return set, nil
		}
	} else {
		if action.Ty != mty.ActionMultiSigConfirmTx || action.GetMultiSigConfirmTx() == nil {
			return nil, mty.ErrActionTyNoMatch
		}
		//       txid          multiSigTx  ，    txhash
		multiSigConfirmTx := action.GetMultiSigConfirmTx()
		multiSigTx, err := getMultiSigTx(m.GetLocalDB(), multiSigConfirmTx.MultiSigAccAddr, multiSigConfirmTx.TxId)
		if err != nil {
			return set, err
		}
		tx, err := getTxByHash(m.GetAPI(), multiSigTx.TxHash)
		if err != nil {
			return nil, err
		}
		payload, err := getMultiSigTxPayload(tx)
		if err != nil {
			return nil, err
		}
		if multiSigTx.TxType == mty.TransferOperate {
			tx := payload.GetMultiSigExecTransferFrom()
			to = tx.To
			execname = tx.Execname
			symbol = tx.Symbol
			amount = tx.Amount
		} else {
			return set, nil
		}
	}
	kv, err := updateAddrReciver(m.GetLocalDB(), to, execname, symbol, amount, addOrRollback)
	if err != nil {
		return set, err
	}
	if kv != nil {
		set = append(set, kv)
	}
	return set, nil
}

//localdb Receipt       。        add/Rollback
func (m *MultiSig) saveMultiSigAccCreate(multiSig *mty.MultiSig, addOrRollback bool) ([]*types.KeyValue, error) {
	multiSigAddr := multiSig.MultiSigAddr
	//             localdb ，         local        。
	oldmultiSig, err := getMultiSigAccount(m.GetLocalDB(), multiSigAddr)
	if err != nil {
		return nil, err
	}
	if addOrRollback && oldmultiSig != nil { //
		multisiglog.Error("saveMultiSigAccCreate:getMultiSigAccount", "addOrRollback", addOrRollback, "MultiSigAddr", multiSigAddr, "oldmultiSig", oldmultiSig, "err", err)
		return nil, mty.ErrAccountHasExist

	} else if !addOrRollback && oldmultiSig == nil { //
		multisiglog.Error("saveMultiSigAccCreate:getMultiSigAccount", "addOrRollback", addOrRollback, "MultiSigAddr", multiSigAddr, "err", err)
		return nil, types.ErrAccountNotExist
	}

	err = setMultiSigAccount(m.GetLocalDB(), multiSig, addOrRollback)
	if err != nil {
		return nil, err
	}
	accountkv := getMultiSigAccountKV(multiSig, addOrRollback)

	//
	lastcount, err := getMultiSigAccCount(m.GetLocalDB())
	if err != nil {
		return nil, err
	}

	//      ,       --
	if !addOrRollback && lastcount > 0 {
		lastcount = lastcount - 1
	}
	accCountListkv, err := updateMultiSigAccList(m.GetLocalDB(), multiSig.MultiSigAddr, lastcount, addOrRollback)
	if err != nil {
		return nil, err
	}

	//
	accCountkv, err := updateMultiSigAccCount(m.GetLocalDB(), addOrRollback)
	if err != nil {
		return nil, err
	}
	//  create
	accAddrkv := setMultiSigAddress(m.GetLocalDB(), multiSig.CreateAddr, multiSig.MultiSigAddr, addOrRollback)

	var kvs []*types.KeyValue
	kvs = append(kvs, accCountkv)
	kvs = append(kvs, accountkv)
	kvs = append(kvs, accCountListkv)
	kvs = append(kvs, accAddrkv)

	return kvs, nil
}

//  owner add/del  .    add/del
func (m *MultiSig) saveMultiSigOwnerAddOrDel(ownerOp mty.ReceiptOwnerAddOrDel, addOrRollback bool) ([]*types.KeyValue, error) {
	//             localdb ，         local        。
	multiSig, err := getMultiSigAccount(m.GetLocalDB(), ownerOp.MultiSigAddr)
	multisiglog.Error("saveMultiSigOwnerAddOrDel", "ownerOp", ownerOp)

	if err != nil || multiSig == nil {
		multisiglog.Error("saveMultiSigOwnerAddOrDel", "addOrRollback", addOrRollback, "ownerOp", ownerOp, "err", err)
		return nil, err
	}
	multisiglog.Error("saveMultiSigOwnerAddOrDel", "wonerlen ", len(multiSig.Owners))

	_, index, _, _, find := getOwnerInfoByAddr(multiSig, ownerOp.Owner.OwnerAddr)
	if addOrRollback { //
		if ownerOp.AddOrDel && !find { //add owner
			multiSig.Owners = append(multiSig.Owners, ownerOp.Owner)
		} else if !ownerOp.AddOrDel && find { //dell owner
			multiSig.Owners = delOwner(multiSig.Owners, index)
			//multiSig.Owners = append(multiSig.Owners[0:index], multiSig.Owners[index+1:]...)
		} else {
			multisiglog.Error("saveMultiSigOwnerAddOrDel", "addOrRollback", addOrRollback, "ownerOp", ownerOp, "index", index, "find", find)
			return nil, mty.ErrOwnerNoMatch
		}

	} else { //
		if ownerOp.AddOrDel && find { //  add owner
			multiSig.Owners = delOwner(multiSig.Owners, index)
			//multiSig.Owners = append(multiSig.Owners[0:index], multiSig.Owners[index+1:]...)
		} else if !ownerOp.AddOrDel && !find { //   del owner
			multiSig.Owners = append(multiSig.Owners, ownerOp.Owner)
		} else {
			multisiglog.Error("saveMultiSigOwnerAddOrDel", "addOrRollback", addOrRollback, "ownerOp", ownerOp, "index", index, "find", find)
			return nil, mty.ErrOwnerNoMatch
		}
	}
	multisiglog.Error("saveMultiSigOwnerAddOrDel", "addOrRollback", addOrRollback, "ownerOp", ownerOp, "multiSig", multiSig)

	err = setMultiSigAccount(m.GetLocalDB(), multiSig, true)
	if err != nil {
		return nil, err
	}
	accountkv := getMultiSigAccountKV(multiSig, true)

	var kvs []*types.KeyValue
	kvs = append(kvs, accountkv)
	return kvs, nil
}

//  owner mod/replace
func (m *MultiSig) saveMultiSigOwnerModOrRep(ownerOp mty.ReceiptOwnerModOrRep, addOrRollback bool) ([]*types.KeyValue, error) {
	//           db
	multiSig, err := getMultiSigAccount(m.GetLocalDB(), ownerOp.MultiSigAddr)

	if err != nil || multiSig == nil {
		return nil, err
	}
	if addOrRollback { //
		_, index, _, _, find := getOwnerInfoByAddr(multiSig, ownerOp.PrevOwner.OwnerAddr)
		if ownerOp.ModOrRep && find { //modify owner weight
			multiSig.Owners[index].Weight = ownerOp.CurrentOwner.Weight
		} else if !ownerOp.ModOrRep && find { //replace owner addr
			multiSig.Owners[index].OwnerAddr = ownerOp.CurrentOwner.OwnerAddr
		} else {
			multisiglog.Error("saveMultiSigOwnerModOrRep ", "addOrRollback", addOrRollback, "ownerOp", ownerOp, "index", index, "find", find)
			return nil, mty.ErrOwnerNoMatch
		}

	} else { //
		_, index, _, _, find := getOwnerInfoByAddr(multiSig, ownerOp.CurrentOwner.OwnerAddr)
		if ownerOp.ModOrRep && find { //  modify owner weight
			multiSig.Owners[index].Weight = ownerOp.PrevOwner.Weight
		} else if !ownerOp.ModOrRep && find { //   replace owner addr
			multiSig.Owners[index].OwnerAddr = ownerOp.PrevOwner.OwnerAddr
		} else {
			multisiglog.Error("saveMultiSigOwnerModOrRep ", "addOrRollback", addOrRollback, "ownerOp", ownerOp, "index", index, "find", find)
			return nil, mty.ErrOwnerNoMatch
		}
	}

	err = setMultiSigAccount(m.GetLocalDB(), multiSig, true)
	if err != nil {
		return nil, err
	}
	accountkv := getMultiSigAccountKV(multiSig, true)

	var kvs []*types.KeyValue
	kvs = append(kvs, accountkv)
	return kvs, nil
}

//  weight   mod
func (m *MultiSig) saveMultiSigAccWeight(accountOp mty.ReceiptWeightModify, addOrRollback bool) ([]*types.KeyValue, error) {
	//           db
	multiSig, err := getMultiSigAccount(m.GetLocalDB(), accountOp.MultiSigAddr)

	if err != nil || multiSig == nil {
		return nil, err
	}
	if addOrRollback { //
		multiSig.RequiredWeight = accountOp.CurrentWeight
	} else { //
		multiSig.RequiredWeight = accountOp.PrevWeight
	}

	err = setMultiSigAccount(m.GetLocalDB(), multiSig, true)
	if err != nil {
		return nil, err
	}
	accountkv := getMultiSigAccountKV(multiSig, true)

	var kvs []*types.KeyValue
	kvs = append(kvs, accountkv)
	return kvs, nil
}

//  DailyLimit       add/mod
func (m *MultiSig) saveMultiSigAccDailyLimit(accountOp mty.ReceiptDailyLimitOperate, addOrRollback bool) ([]*types.KeyValue, error) {
	//           db
	multiSig, err := getMultiSigAccount(m.GetLocalDB(), accountOp.MultiSigAddr)

	if err != nil || multiSig == nil {
		return nil, err
	}
	curExecer := accountOp.CurDailyLimit.Execer
	curSymbol := accountOp.CurDailyLimit.Symbol
	curDailyLimit := accountOp.CurDailyLimit
	prevDailyLimit := accountOp.PrevDailyLimit

	//
	index, find := isDailyLimit(multiSig, curExecer, curSymbol)

	if addOrRollback { //
		if accountOp.AddOrModify && !find { //add DailyLimit
			multiSig.DailyLimits = append(multiSig.DailyLimits, curDailyLimit)
		} else if !accountOp.AddOrModify && find { //modifyDailyLimit
			multiSig.DailyLimits[index].DailyLimit = curDailyLimit.DailyLimit
			multiSig.DailyLimits[index].SpentToday = curDailyLimit.SpentToday
			multiSig.DailyLimits[index].LastDay = curDailyLimit.LastDay
		} else {
			multisiglog.Error("saveMultiSigAccDailyLimit", "addOrRollback", addOrRollback, "accountOp", accountOp, "index", index, "find", find)
			return nil, mty.ErrDailyLimitNoMatch
		}
	} else { //
		if accountOp.AddOrModify && find { //     add   DailyLimit
			multiSig.DailyLimits = append(multiSig.DailyLimits[0:index], multiSig.DailyLimits[index+1:]...)
		} else if !accountOp.AddOrModify && find { //         modifyDailyLimit
			multiSig.DailyLimits[index].DailyLimit = prevDailyLimit.DailyLimit
			multiSig.DailyLimits[index].SpentToday = prevDailyLimit.SpentToday
			multiSig.DailyLimits[index].LastDay = prevDailyLimit.LastDay
		} else {
			multisiglog.Error("saveMultiSigAccDailyLimit", "addOrRollback", addOrRollback, "accountOp", accountOp, "index", index, "find", find)
			return nil, mty.ErrDailyLimitNoMatch
		}
	}

	err = setMultiSigAccount(m.GetLocalDB(), multiSig, true)
	if err != nil {
		return nil, err
	}
	accountkv := getMultiSigAccountKV(multiSig, true)

	var kvs []*types.KeyValue
	kvs = append(kvs, accountkv)
	return kvs, nil
}

//         Confirm/Revoke
func (m *MultiSig) saveMultiSigConfirmTx(confirmTx mty.ReceiptConfirmTx, addOrRollback bool) ([]*types.KeyValue, error) {
	multiSigAddr := confirmTx.MultiSigTxOwner.MultiSigAddr
	txid := confirmTx.MultiSigTxOwner.Txid
	owner := confirmTx.MultiSigTxOwner.ConfirmedOwner

	//           db
	multiSigTx, err := getMultiSigTx(m.GetLocalDB(), multiSigAddr, txid)
	if err != nil {
		return nil, err
	}
	if multiSigTx == nil {
		multisiglog.Error("saveMultiSigConfirmTx", "addOrRollback", addOrRollback, "confirmTx", confirmTx)
		return nil, mty.ErrTxidNotExist
	}
	index, exist := isOwnerConfirmedTx(multiSigTx, owner.OwnerAddr)
	if addOrRollback { //
		if confirmTx.ConfirmeOrRevoke && !exist { //add Confirmed Owner
			multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner, owner)
		} else if !confirmTx.ConfirmeOrRevoke && exist { //Revoke Confirmed Owner
			multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner[0:index], multiSigTx.ConfirmedOwner[index+1:]...)
		} else {
			multisiglog.Error("saveMultiSigConfirmTx", "addOrRollback", addOrRollback, "confirmTx", confirmTx, "index", index, "exist", exist)
			return nil, mty.ErrDailyLimitNoMatch
		}
	} else { //
		if confirmTx.ConfirmeOrRevoke && exist { //     add Confirmed Owner
			//multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner, owner)
			multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner[0:index], multiSigTx.ConfirmedOwner[index+1:]...)

		} else if !confirmTx.ConfirmeOrRevoke && !exist { //     Revoke Confirmed Owner
			multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner, owner)
		} else {
			multisiglog.Error("saveMultiSigConfirmTx", "addOrRollback", addOrRollback, "confirmTx", confirmTx, "index", index, "exist", exist)
			return nil, mty.ErrDailyLimitNoMatch
		}
	}

	err = setMultiSigTx(m.GetLocalDB(), multiSigTx, true)
	if err != nil {
		return nil, err
	}
	txkv := getMultiSigTxKV(multiSigTx, true)

	var kvs []*types.KeyValue
	kvs = append(kvs, txkv)
	return kvs, nil
}

//             ,               owner
//
func (m *MultiSig) saveMultiSigTx(execTx mty.ReceiptMultiSigTx, addOrRollback bool) ([]*types.KeyValue, error) {
	multiSigAddr := execTx.MultiSigTxOwner.MultiSigAddr
	txid := execTx.MultiSigTxOwner.Txid
	owner := execTx.MultiSigTxOwner.ConfirmedOwner
	curExecuted := execTx.CurExecuted
	prevExecuted := execTx.PrevExecuted
	submitOrConfirm := execTx.SubmitOrConfirm

	temMultiSigTx := &mty.MultiSigTx{}
	temMultiSigTx.MultiSigAddr = multiSigAddr
	temMultiSigTx.Txid = txid
	temMultiSigTx.TxHash = execTx.TxHash
	temMultiSigTx.TxType = execTx.TxType
	temMultiSigTx.Executed = false
	//           db
	multiSigTx, err := getMultiSigTx(m.GetLocalDB(), multiSigAddr, txid)
	if err != nil {
		multisiglog.Error("saveMultiSigTx getMultiSigTx ", "addOrRollback", addOrRollback, "execTx", execTx, "err", err)
		return nil, err
	}

	//Confirm          txid
	if multiSigTx == nil && !submitOrConfirm {
		multisiglog.Error("saveMultiSigTx", "addOrRollback", addOrRollback, "execTx", execTx)
		return nil, mty.ErrTxidNotExist
	}

	//add submit       txid，
	if submitOrConfirm && addOrRollback {
		if multiSigTx != nil {
			multisiglog.Error("saveMultiSigTx", "addOrRollback", addOrRollback, "execTx", execTx)
			return nil, mty.ErrTxidHasExist
		}
		multiSigTx = temMultiSigTx
	}

	index, exist := isOwnerConfirmedTx(multiSigTx, owner.OwnerAddr)
	if addOrRollback { //
		if !exist { //add Confirmed Owner and modify Executed
			multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner, owner)
			if prevExecuted != multiSigTx.Executed {
				return nil, mty.ErrExecutedNoMatch
			}
			multiSigTx.Executed = curExecuted
		} else {
			multisiglog.Error("saveMultiSigTx", "addOrRollback", addOrRollback, "execTx", execTx, "index", index, "exist", exist)
			return nil, mty.ErrOwnerNoMatch
		}
	} else { //
		if exist { //     add Confirmed Owner and modify Executed
			multiSigTx.ConfirmedOwner = append(multiSigTx.ConfirmedOwner[0:index], multiSigTx.ConfirmedOwner[index+1:]...)
			multiSigTx.Executed = prevExecuted
		} else {
			multisiglog.Error("saveMultiSigTx", "addOrRollback", addOrRollback, "execTx", execTx, "index", index, "exist", exist)
			return nil, mty.ErrOwnerNoMatch
		}
	}
	//submit          txid     nil
	setNil := true
	if !addOrRollback && submitOrConfirm {
		setNil = false
	}

	err = setMultiSigTx(m.GetLocalDB(), multiSigTx, setNil)
	if err != nil {
		return nil, err
	}
	txkv := getMultiSigTxKV(multiSigTx, setNil)

	var kvs []*types.KeyValue
	kvs = append(kvs, txkv)
	return kvs, nil
}

//             ，             ，  txcount
func (m *MultiSig) saveDailyLimitUpdate(execTransfer mty.ReceiptAccDailyLimitUpdate, addOrRollback bool) ([]*types.KeyValue, error) {
	multiSigAddr := execTransfer.MultiSigAddr
	curDailyLimit := execTransfer.CurDailyLimit
	prevDailyLimit := execTransfer.PrevDailyLimit
	execer := execTransfer.CurDailyLimit.Execer
	symbol := execTransfer.CurDailyLimit.Symbol

	//           db
	multiSig, err := getMultiSigAccount(m.GetLocalDB(), multiSigAddr)
	if err != nil {
		return nil, err
	}
	if multiSig == nil {
		multisiglog.Error("saveAccExecTransfer", "addOrRollback", addOrRollback, "execTransfer", execTransfer)
		return nil, types.ErrAccountNotExist
	}
	index, exist := isDailyLimit(multiSig, execer, symbol)
	if !exist {
		return nil, types.ErrAccountNotExist
	}
	if addOrRollback { //
		multiSig.DailyLimits[index].SpentToday = curDailyLimit.SpentToday
		multiSig.DailyLimits[index].LastDay = curDailyLimit.LastDay
	} else { //
		multiSig.DailyLimits[index].SpentToday = prevDailyLimit.SpentToday
		multiSig.DailyLimits[index].LastDay = prevDailyLimit.LastDay
	}

	err = setMultiSigAccount(m.GetLocalDB(), multiSig, true)
	if err != nil {
		return nil, err
	}
	txkv := getMultiSigAccountKV(multiSig, true)

	var kvs []*types.KeyValue
	kvs = append(kvs, txkv)
	return kvs, nil
}

//             ，             ，  txcount
func (m *MultiSig) saveMultiSigTxCountUpdate(accTxCount mty.ReceiptTxCountUpdate, addOrRollback bool) ([]*types.KeyValue, error) {
	multiSigAddr := accTxCount.MultiSigAddr
	curTxCount := accTxCount.CurTxCount

	//           db
	multiSig, err := getMultiSigAccount(m.GetLocalDB(), multiSigAddr)
	if err != nil {
		return nil, err
	}
	if multiSig == nil {
		multisiglog.Error("saveMultiSigTxCountUpdate", "addOrRollback", addOrRollback, "accTxCount", accTxCount)
		return nil, types.ErrAccountNotExist
	}

	if addOrRollback { //
		if multiSig.TxCount+1 == curTxCount {
			multiSig.TxCount = curTxCount
		} else {
			multisiglog.Error("saveMultiSigTxCountUpdate", "addOrRollback", addOrRollback, "accTxCount", accTxCount, "TxCount", multiSig.TxCount)
			return nil, mty.ErrInvalidTxid
		}
	} else { //
		if multiSig.TxCount == curTxCount && curTxCount > 0 {
			multiSig.TxCount = curTxCount - 1
		}
	}

	err = setMultiSigAccount(m.GetLocalDB(), multiSig, true)
	if err != nil {
		return nil, err
	}
	txkv := getMultiSigAccountKV(multiSig, true)

	var kvs []*types.KeyValue
	kvs = append(kvs, txkv)
	return kvs, nil
}

//
func (m *MultiSig) getMultiSigAccAssets(multiSigAddr string, assets *mty.Assets) (*types.Account, error) {
	symbol := getRealSymbol(assets.Symbol)
	cfg := m.GetAPI().GetConfig()
	acc, err := account.NewAccountDB(cfg, assets.Execer, symbol, m.GetStateDB())
	if err != nil {
		return &types.Account{}, err
	}
	var acc1 *types.Account

	execaddress := dapp.ExecAddress(cfg.ExecName(m.GetName()))
	acc1 = acc.LoadExecAccount(multiSigAddr, execaddress)
	return acc1, nil
}

//

//    owner weight  ，owner   index，  owners weight    ，  owner
func getOwnerInfoByAddr(multiSigAcc *mty.MultiSig, oldowner string) (uint64, int, uint64, int, bool) {
	//      owners，     owner    .
	var findindex int
	var totalweight uint64
	var oldweight uint64
	var totalowner int
	flag := false

	for index, owner := range multiSigAcc.Owners {
		if owner.OwnerAddr == oldowner {
			flag = true
			findindex = index
			oldweight = owner.Weight
		}
		totalweight += owner.Weight
		totalowner++
	}
	//owner
	if !flag {
		return 0, 0, totalweight, totalowner, false
	}
	return oldweight, findindex, totalweight, totalowner, true
}

//
func isConfirmed(requiredWeight uint64, multiSigTx *mty.MultiSigTx) bool {
	var totalweight uint64
	for _, owner := range multiSigTx.ConfirmedOwner {
		totalweight += owner.Weight
	}
	return totalweight >= requiredWeight
}

//                 ,      ，    newLastDay
func isUnderLimit(blocktime int64, amount uint64, dailyLimit *mty.DailyLimit) (bool, int64) {

	var lastDay int64
	var newSpentToday uint64

	nowtime := blocktime //types.Now().Unix()
	newSpentToday = dailyLimit.SpentToday

	//        。    LastDay     ，SpentToday    0
	if nowtime > dailyLimit.LastDay+mty.OneDaySecond {
		lastDay = nowtime
		newSpentToday = 0
	}

	if newSpentToday+amount > dailyLimit.DailyLimit || newSpentToday+amount < newSpentToday {
		return false, lastDay
	}
	return true, lastDay
}

//          multiSigAcc       owner,   owner     weight
func isOwner(multiSigAcc *mty.MultiSig, ownerAddr string) (uint64, bool) {
	for _, owner := range multiSigAcc.Owners {
		if owner.OwnerAddr == ownerAddr {
			return owner.Weight, true
		}
	}
	return 0, false
}

//    index owner owners
func delOwner(Owners []*mty.Owner, index int) []*mty.Owner {
	ownerSize := len(Owners)
	multisiglog.Error("delOwner", "ownerSize", ownerSize, "index", index)

	//     owner
	if index == 0 {
		Owners = Owners[1:]
	} else if (ownerSize) == index+1 { //      owner
		multisiglog.Error("delOwner", "ownerSize", ownerSize)
		Owners = Owners[0 : ownerSize-1]
	} else {
		Owners = append(Owners[0:index], Owners[index+1:]...)
	}
	return Owners
}

//
func isDailyLimit(multiSigAcc *mty.MultiSig, execer, symbol string) (int, bool) {
	for index, dailyLimit := range multiSigAcc.DailyLimits {
		if dailyLimit.Execer == execer && dailyLimit.Symbol == symbol {
			return index, true
		}
	}
	return 0, false
}

//owner         txid，        index
func isOwnerConfirmedTx(multiSigTx *mty.MultiSigTx, ownerAddr string) (int, bool) {
	for index, owner := range multiSigTx.ConfirmedOwner {
		if owner.OwnerAddr == ownerAddr {
			return index, true
		}
	}
	return 0, false
}

//  txhash  tx
func getTxByHash(api client.QueueProtocolAPI, txHash string) (*types.TransactionDetail, error) {
	hash, err := hex.DecodeString(txHash)
	if err != nil {
		multisiglog.Error("GetTxByHash DecodeString ", "hash", txHash)
		return nil, err
	}
	txs, err := api.GetTransactionByHash(&types.ReqHashes{Hashes: [][]byte{hash}})
	if err != nil {
		multisiglog.Error("GetTxByHash", "hash", txHash)
		return nil, err
	}
	if len(txs.Txs) != 1 {
		multisiglog.Error("GetTxByHash", "len is not 1", len(txs.Txs))
		return nil, mty.ErrTxHashNoMatch
	}
	if txs.Txs == nil {
		multisiglog.Error("GetTxByHash", "tx hash not found", txHash)
		return nil, mty.ErrTxHashNoMatch
	}
	return txs.Txs[0], nil
}

// tx     payload
func getMultiSigTxPayload(tx *types.TransactionDetail) (*mty.MultiSigAction, error) {
	if !bytes.HasSuffix(tx.Tx.Execer, []byte(mty.MultiSigX)) {
		multisiglog.Error("GetMultiSigTx", "tx.Tx.Execer", string(tx.Tx.Execer), "MultiSigX", mty.MultiSigX)
		return nil, mty.ErrExecerHashNoMatch
	}
	var payload mty.MultiSigAction
	err := types.Decode(tx.Tx.Payload, &payload)
	if err != nil {
		multisiglog.Error("GetMultiSigTx:Decode Payload", "error", err)
		return nil, err
	}
	multisiglog.Error("GetMultiSigTx:Decode Payload", "payload", payload)
	return &payload, nil
}

//dpos      ，   mavl      key
func getRealSymbol(symbol string) string {
	if symbol == types.DOM {
		return "dpos"
	}
	return symbol
}
