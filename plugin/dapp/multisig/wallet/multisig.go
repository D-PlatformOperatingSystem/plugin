// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of policy source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"sync"

	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	wcom "github.com/D-PlatformOperatingSystem/dpos/wallet/common"
	mtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
)

var (
	bizlog = log15.New("module", "wallet.multisig")
)

func init() {
	wcom.RegisterPolicy(mtypes.MultiSigX, New())
}

// New
func New() wcom.WalletBizPolicy {
	return &multisigPolicy{
		mtx:      &sync.Mutex{},
		rescanwg: &sync.WaitGroup{},
	}
}

type multisigPolicy struct {
	mtx           *sync.Mutex
	store         *multisigStore
	walletOperate wcom.WalletOperate
	rescanwg      *sync.WaitGroup
	cfg           *subConfig
}
type subConfig struct {
	RescanMultisigAddr bool `json:"rescanMultisigAddr"`
}

func (policy *multisigPolicy) setWalletOperate(walletBiz wcom.WalletOperate) {
	policy.mtx.Lock()
	defer policy.mtx.Unlock()
	policy.walletOperate = walletBiz
}

func (policy *multisigPolicy) getWalletOperate() wcom.WalletOperate {
	policy.mtx.Lock()
	defer policy.mtx.Unlock()
	return policy.walletOperate
}

func (policy *multisigPolicy) getRescanMultisigAddr() bool {
	policy.mtx.Lock()
	defer policy.mtx.Unlock()
	return policy.cfg.RescanMultisigAddr
}

// Init
func (policy *multisigPolicy) Init(walletOperate wcom.WalletOperate, sub []byte) {
	policy.setWalletOperate(walletOperate)
	policy.store = newStore(walletOperate.GetDBStore())
	var subcfg subConfig
	if sub != nil {
		types.MustDecode(sub, &subcfg)
	}
	policy.cfg = &subcfg
}

// OnCreateNewAccount
func (policy *multisigPolicy) OnCreateNewAccount(acc *types.Account) {
	if policy.getRescanMultisigAddr() {
		policy.rescanwg.Add(1)
		go policy.rescanOwnerAttrByAddr(acc.Addr)
	}
}

// OnImportPrivateKey
func (policy *multisigPolicy) OnImportPrivateKey(acc *types.Account) {
	if policy.getRescanMultisigAddr() {
		policy.rescanwg.Add(1)
		go policy.rescanOwnerAttrByAddr(acc.Addr)
	}
}

// OnAddBlockTx
func (policy *multisigPolicy) OnAddBlockTx(block *types.BlockDetail, tx *types.Transaction, index int32, dbbatch db.Batch) *types.WalletTxDetail {
	policy.filterMultisigTxsFromBlock(tx, index, block, dbbatch, true)
	return policy.proceWalletTxDetail(block, tx, index)
}

// OnDeleteBlockTx
func (policy *multisigPolicy) OnDeleteBlockTx(block *types.BlockDetail, tx *types.Transaction, index int32, dbbatch db.Batch) *types.WalletTxDetail {
	policy.filterMultisigTxsFromBlock(tx, index, block, dbbatch, false)
	return policy.proceWalletTxDetail(block, tx, index)
}

// OnAddBlockFinish
func (policy *multisigPolicy) OnAddBlockFinish(block *types.BlockDetail) {

}

// OnDeleteBlockFinish
func (policy *multisigPolicy) OnDeleteBlockFinish(block *types.BlockDetail) {

}

// OnClose
func (policy *multisigPolicy) OnClose() {

}

// OnSetQueueClient
func (policy *multisigPolicy) OnSetQueueClient() {
}

// OnWalletLocked
func (policy *multisigPolicy) OnWalletLocked() {
}

// OnWalletUnlocked
func (policy *multisigPolicy) OnWalletUnlocked(WalletUnLock *types.WalletUnLock) {
}

// Call
func (policy *multisigPolicy) Call(funName string, in types.Message) (ret types.Message, err error) {
	err = types.ErrNotSupport
	return
}

// SignTransaction :
func (policy *multisigPolicy) SignTransaction(key crypto.PrivKey, req *types.ReqSignRawTx) (needSysSign bool, signtxhex string, err error) {
	needSysSign = true
	return
}

func (policy *multisigPolicy) filterMultisigTxsFromBlock(tx *types.Transaction, index int32, block *types.BlockDetail, newbatch db.Batch, addOrRollback bool) {

	receiptData := block.Receipts[index]
	if receiptData.GetTy() != types.ExecOk {
		return
	}

	for _, log := range receiptData.Logs {
		bizlog.Debug("filterMultisigTxsFromBlock", "Ty", log.Ty)

		switch log.Ty {
		case mtypes.TyLogMultiSigAccCreate:
			{
				var receipt mtypes.MultiSig
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					bizlog.Error("filterMultisigTxsFromBlock Decode err", "Ty", log.Ty, "err", err)
					return
				}
				policy.saveMultiSigAccCreate(&receipt, newbatch, addOrRollback)
			}
		case mtypes.TyLogMultiSigOwnerAdd,
			mtypes.TyLogMultiSigOwnerDel:
			{
				var receipt mtypes.ReceiptOwnerAddOrDel
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					bizlog.Error("filterMultisigTxsFromBlock Decode err", "Ty", log.Ty, "err", err)
					return
				}
				policy.saveMultiSigOwnerAddOrDel(&receipt, newbatch, addOrRollback)
			}
		case mtypes.TyLogMultiSigOwnerModify,
			mtypes.TyLogMultiSigOwnerReplace:
			{
				var receipt mtypes.ReceiptOwnerModOrRep
				err := types.Decode(log.Log, &receipt)
				if err != nil {
					bizlog.Error("filterMultisigTxsFromBlock Decode err", "Ty", log.Ty, "err", err)
					return
				}
				policy.saveMultiSigOwnerModOrRep(&receipt, newbatch, addOrRollback)
			}
		default:
			continue
		}
	}
}

//        add/Rollback
func (policy *multisigPolicy) saveMultiSigAccCreate(multiSig *mtypes.MultiSig, newbatch db.Batch, addOrRollback bool) {
	wallet := policy.getWalletOperate()

	for _, owner := range multiSig.Owners {
		if wallet.AddrInWallet(owner.OwnerAddr) {
			ownerAttrs, err := policy.store.listOwnerAttrsByAddr(owner.OwnerAddr)
			if err != nil && err != types.ErrNotFound {
				bizlog.Error("saveMultiSigAccCreate ", "owner.OwnerAddr", owner.OwnerAddr, "err", err)
				continue
			}
			//add tx
			if addOrRollback {
				ownerAttr := &mtypes.OwnerAttr{
					MultiSigAddr: multiSig.MultiSigAddr,
					OwnerAddr:    owner.OwnerAddr,
					Weight:       owner.Weight,
				}
				//     owner
				if err == types.ErrNotFound {
					AddOwnerAttr(true, nil, ownerAttr, newbatch)
				} else if ownerAttrs != nil {
					AddOwnerAttr(false, ownerAttrs, ownerAttr, newbatch)
				}
			} else if ownerAttrs != nil && !addOrRollback {
				DelOwnerAttr(ownerAttrs, owner.OwnerAddr, multiSig.MultiSigAddr, newbatch)
			}
		}
	}
}

//  owner add/del  .    add/del
func (policy *multisigPolicy) saveMultiSigOwnerAddOrDel(ownerOp *mtypes.ReceiptOwnerAddOrDel, newbatch db.Batch, addOrRollback bool) {
	wallet := policy.getWalletOperate()
	owner := ownerOp.Owner

	if wallet.AddrInWallet(owner.OwnerAddr) {
		ownerAttrs, err := policy.store.listOwnerAttrsByAddr(owner.OwnerAddr)
		if err != nil && err != types.ErrNotFound {
			bizlog.Error("saveMultiSigOwnerAddOrDel ", "owner.OwnerAddr", owner.OwnerAddr, "err", err)
			return
		}
		ownerAttr := &mtypes.OwnerAttr{
			MultiSigAddr: ownerOp.MultiSigAddr,
			OwnerAddr:    owner.OwnerAddr,
			Weight:       owner.Weight,
		}
		//add tx
		if addOrRollback {
			if ownerOp.AddOrDel {
				if err == types.ErrNotFound {
					AddOwnerAttr(true, nil, ownerAttr, newbatch)
				} else {
					AddOwnerAttr(false, ownerAttrs, ownerAttr, newbatch)
				}
			} else if ownerAttrs != nil {
				DelOwnerAttr(ownerAttrs, owner.OwnerAddr, ownerOp.MultiSigAddr, newbatch)
			}
		} else {
			//  add owner
			if ownerOp.AddOrDel && ownerAttrs != nil {
				DelOwnerAttr(ownerAttrs, owner.OwnerAddr, ownerOp.MultiSigAddr, newbatch)
			} else if !ownerOp.AddOrDel { //   del owner
				if err == types.ErrNotFound {
					AddOwnerAttr(true, nil, ownerAttr, newbatch)
				} else if ownerAttrs != nil {
					AddOwnerAttr(false, ownerAttrs, ownerAttr, newbatch)
				}
			}
		}
	}
}

//  owner mod/replace
func (policy *multisigPolicy) saveMultiSigOwnerModOrRep(ownerOp *mtypes.ReceiptOwnerModOrRep, newbatch db.Batch, addOrRollback bool) {
	wallet := policy.getWalletOperate()
	prevOwner := ownerOp.PrevOwner
	curOwner := ownerOp.CurrentOwner
	multiSigAddr := ownerOp.MultiSigAddr

	//    prevOwner
	if wallet.AddrInWallet(prevOwner.OwnerAddr) {
		ownerAttrs, err := policy.store.listOwnerAttrsByAddr(prevOwner.OwnerAddr)
		if err != nil && err != types.ErrNotFound {
			bizlog.Error("saveMultiSigOwnerModOrRep ", "prevOwner.OwnerAddr", prevOwner.OwnerAddr, "err", err)
			return
		}
		ownerAttr := &mtypes.OwnerAttr{
			MultiSigAddr: multiSigAddr,
			OwnerAddr:    prevOwner.OwnerAddr,
			Weight:       prevOwner.Weight,
		}
		//add tx
		if addOrRollback && ownerAttrs != nil {
			if ownerOp.ModOrRep {
				ModOwnerAttr(ownerAttrs, prevOwner.OwnerAddr, multiSigAddr, curOwner.Weight, newbatch)
			} else {
				DelOwnerAttr(ownerAttrs, prevOwner.OwnerAddr, multiSigAddr, newbatch)
			}
		} else if !addOrRollback { //  Rollback
			if ownerOp.ModOrRep && ownerAttrs != nil {
				ModOwnerAttr(ownerAttrs, prevOwner.OwnerAddr, multiSigAddr, prevOwner.Weight, newbatch)
			} else {
				if err == types.ErrNotFound {
					AddOwnerAttr(true, nil, ownerAttr, newbatch)
				} else if ownerAttrs != nil {
					AddOwnerAttr(false, ownerAttrs, ownerAttr, newbatch)
				}
			}
		}
	}
	//    curOwner,replace
	if wallet.AddrInWallet(curOwner.OwnerAddr) && !ownerOp.ModOrRep {
		ownerAttrs, err := policy.store.listOwnerAttrsByAddr(curOwner.OwnerAddr)
		if err != nil && err != types.ErrNotFound {
			bizlog.Error("saveMultiSigOwnerModOrRep ", "curOwner.OwnerAddr", curOwner.OwnerAddr, "err", err)
			return
		}
		ownerAttr := &mtypes.OwnerAttr{
			MultiSigAddr: multiSigAddr,
			OwnerAddr:    curOwner.OwnerAddr,
			Weight:       curOwner.Weight,
		}
		if addOrRollback {
			if err == types.ErrNotFound {
				AddOwnerAttr(true, nil, ownerAttr, newbatch)
			} else if ownerAttrs != nil {
				AddOwnerAttr(false, ownerAttrs, ownerAttr, newbatch)
			}
		} else if ownerAttrs != nil {
			DelOwnerAttr(ownerAttrs, curOwner.OwnerAddr, multiSigAddr, newbatch)
		}
	}
}

//AddOwnerAttr :   owmer
func AddOwnerAttr(firstAdd bool, ownerAttrs *mtypes.OwnerAttrs, ownerAttr *mtypes.OwnerAttr, newbatch db.Batch) {
	if firstAdd {
		var firstownerAttrs mtypes.OwnerAttrs
		addOwnerAttr(&firstownerAttrs, ownerAttr)
		batchSet(&firstownerAttrs, ownerAttr.OwnerAddr, newbatch)
		return
	}
	addOwnerAttr(ownerAttrs, ownerAttr)
	batchSet(ownerAttrs, ownerAttr.OwnerAddr, newbatch)
}

//DelOwnerAttr ：  owner
func DelOwnerAttr(ownerAttrs *mtypes.OwnerAttrs, ownerAddr string, multiSigAddr string, newbatch db.Batch) {
	index, find := getOwnerAttr(ownerAttrs, multiSigAddr)
	if find {
		//         value
		if len(ownerAttrs.Items) == 1 && index == 0 {
			newbatch.Delete(calcMultisigAddr(ownerAddr))
		} else {
			ownerAttrs = delOwnerAttr(ownerAttrs, index)
			batchSet(ownerAttrs, ownerAddr, newbatch)
		}
	}
}

//ModOwnerAttr ：  owner weight
func ModOwnerAttr(ownerAttrs *mtypes.OwnerAttrs, ownerAddr string, multiSigAddr string, weight uint64, newbatch db.Batch) {
	index, find := getOwnerAttr(ownerAttrs, multiSigAddr)
	if find {
		ownerAttrs.Items[index].Weight = weight
		batchSet(ownerAttrs, ownerAddr, newbatch)
	}
}

//batchSet :
func batchSet(ownerAttrs *mtypes.OwnerAttrs, addr string, newbatch db.Batch) {
	v := *ownerAttrs
	ownerAttrsbyte := types.Encode(&v)
	newbatch.Set(calcMultisigAddr(addr), ownerAttrsbyte)
}

//delOwnerAttr :
func delOwnerAttr(ownerAttrs *mtypes.OwnerAttrs, index int) *mtypes.OwnerAttrs {
	ownerSize := len(ownerAttrs.Items)
	//     owner
	if index == 0 {
		ownerAttrs.Items = ownerAttrs.Items[1:]
	} else if (ownerSize) == index+1 { //      owner
		ownerAttrs.Items = ownerAttrs.Items[0 : ownerSize-1]
	} else {
		ownerAttrs.Items = append(ownerAttrs.Items[0:index], ownerAttrs.Items[index+1:]...)
	}
	return ownerAttrs
}

//addOwnerAttr ：
func addOwnerAttr(ownerAttrs *mtypes.OwnerAttrs, ownerAttr *mtypes.OwnerAttr) *mtypes.OwnerAttrs {
	ownerAttrs.Items = append(ownerAttrs.Items, ownerAttr)
	return ownerAttrs
}

// getOwnerAttr :
func getOwnerAttr(ownerAttrs *mtypes.OwnerAttrs, multiSigAddr string) (int, bool) {
	for index, owner := range ownerAttrs.Items {
		if owner.MultiSigAddr == multiSigAddr {
			return index, true
		}
	}
	return 0, false
}

//            ， blockchain            ，
func (policy *multisigPolicy) rescanOwnerAttrByAddr(addr string) {
	beg := types.Now()
	defer func() {
		bizlog.Info("rescanOwnerAttrByAddr", "addr", addr, "cost", types.Since(beg))
	}()

	defer policy.rescanwg.Done()
	if len(addr) == 0 {
		bizlog.Error("rescanOwnerAttrByAddr input addr is nil!")
		return
	}

	operater := policy.getWalletOperate()
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	//
	msg, err := operater.GetAPI().Query(cfg.ExecName(mtypes.MultiSigX), "MultiSigAccCount", &types.ReqNil{})
	if err != nil {
		bizlog.Error("rescanOwnerAttrByAddr Query MultiSigAccCount err", "MultiSigX", mtypes.MultiSigX, "addr", addr, "err", err)
		return
	}
	replay := msg.(*types.Int64)
	if replay == nil {
		bizlog.Error("rescanOwnerAttrByAddr Query MultiSigAccCount is nil")
		return
	}
	totalCount := replay.Data
	bizlog.Info("rescanOwnerAttrByAddr MultiSigAccCount ", "totalCount", totalCount, "addr", addr)
	if totalCount <= 0 {
		return
	}
	var curCount int64
	for {
		var req mtypes.ReqMultiSigAccs
		if totalCount <= MaxCountPerTime || (curCount+MaxCountPerTime) >= totalCount {
			req.Start = curCount
			req.End = totalCount - 1
			curCount = req.End
		} else if curCount+MaxCountPerTime < totalCount {
			req.Start = curCount
			req.End = req.Start + MaxCountPerTime
			curCount = req.End
		}
		msg, err := operater.GetAPI().Query(cfg.ExecName(mtypes.MultiSigX), "MultiSigAccounts", &req)
		if err != nil {
			bizlog.Error("rescanOwnerAttrByAddr", "MultiSigAccounts error", err, "addr", addr)
			return
		}

		replay := msg.(*mtypes.ReplyMultiSigAccs)
		if replay == nil {
			bizlog.Error("rescanOwnerAttrByAddr Query MultiSigAccounts is nil")
			return
		}
		policy.proceMultiSigAcc(replay, addr)
		if curCount >= totalCount-1 {
			return
		}
		curCount = curCount + 1
	}
}
func (policy *multisigPolicy) proceMultiSigAcc(multiSigAccs *mtypes.ReplyMultiSigAccs, owneraddr string) {
	operater := policy.getWalletOperate()
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	for _, multiSigaddr := range multiSigAccs.Address {
		req := mtypes.ReqMultiSigAccInfo{
			MultiSigAccAddr: multiSigaddr,
		}
		msg, err := operater.GetAPI().Query(cfg.ExecName(mtypes.MultiSigX), "MultiSigAccountInfo", &req)
		if err != nil {
			bizlog.Error("ProceMultiSigAcc", "MultiSigAccountInfo error", err, "multiSigaddr", multiSigaddr)
			continue
		}
		replay := msg.(*mtypes.MultiSig)
		if replay == nil {
			bizlog.Error("ProceMultiSigAcc Query MultiSigAccountInfo is nil", "multiSigaddr", multiSigaddr)
			continue
		}

		for _, owner := range replay.Owners {
			if owner.OwnerAddr == owneraddr {
				ownerAttrs, err := policy.store.listOwnerAttrsByAddr(owneraddr)
				if err != nil && err != types.ErrNotFound {
					bizlog.Error("ProceMultiSigAcc ", "owneraddr", owneraddr, "err", err)
					break
				}
				ownerAttr := &mtypes.OwnerAttr{
					MultiSigAddr: multiSigaddr,
					OwnerAddr:    owner.OwnerAddr,
					Weight:       owner.Weight,
				}
				newbatch := policy.store.NewBatch(true)
				if err == types.ErrNotFound {
					AddOwnerAttr(true, nil, ownerAttr, newbatch)
				} else if ownerAttrs != nil {
					AddOwnerAttr(false, ownerAttrs, ownerAttr, newbatch)
				}
				err = newbatch.Write()
				if err != nil {
					bizlog.Error("ProceMultiSigAcc Write", "owneraddr", owneraddr, "err", err)
				}
				break
			}
		}
	}
}

func (policy *multisigPolicy) proceWalletTxDetail(block *types.BlockDetail, tx *types.Transaction, index int32) *types.WalletTxDetail {
	receipt := block.Receipts[index]
	amount, err := tx.Amount()
	if err != nil {
		bizlog.Error("proceWalletTxDetail:tx.Amount()", "err", err)
	}
	wtxdetail := &types.WalletTxDetail{
		Tx:         tx,
		Height:     block.Block.Height,
		Index:      int64(index),
		Receipt:    receipt,
		Blocktime:  block.Block.BlockTime,
		ActionName: tx.ActionName(),
		Amount:     amount,
		Payload:    nil,
	}
	if len(wtxdetail.Fromaddr) <= 0 {
		pubkey := tx.Signature.GetPubkey()
		address := address.PubKeyToAddress(pubkey)
		//from addr
		fromaddress := address.String()
		if len(fromaddress) != 0 && policy.walletOperate.AddrInWallet(fromaddress) {
			wtxdetail.Fromaddr = fromaddress
		}
	}
	if len(wtxdetail.Fromaddr) <= 0 {
		toaddr := tx.GetTo()
		if len(toaddr) != 0 && policy.walletOperate.AddrInWallet(toaddr) {
			wtxdetail.Fromaddr = toaddr
		}
	}
	return wtxdetail
}
