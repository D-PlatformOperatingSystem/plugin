// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of policy source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"sync"
	"sync/atomic"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	wcom "github.com/D-PlatformOperatingSystem/dpos/wallet/common"
	privacytypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/types"
)

var (
	bizlog = log15.New("module", "wallet.privacy")
	// MaxTxHashsPerTime
	MaxTxHashsPerTime int64 = 100
	// maxTxNumPerBlock
	maxTxNumPerBlock int64 = types.MaxTxsPerBlock
)

func init() {
	wcom.RegisterPolicy(privacytypes.PrivacyX, New())
}

// New
func New() wcom.WalletBizPolicy {
	return &privacyPolicy{
		mtx:            &sync.Mutex{},
		rescanwg:       &sync.WaitGroup{},
		rescanUTXOflag: privacytypes.UtxoFlagNoScan,
	}
}

type privacyPolicy struct {
	mtx            *sync.Mutex
	store          *privacyStore
	walletOperate  wcom.WalletOperate
	rescanwg       *sync.WaitGroup
	rescanUTXOflag int32
}

func (policy *privacyPolicy) setWalletOperate(walletBiz wcom.WalletOperate) {
	policy.mtx.Lock()
	defer policy.mtx.Unlock()
	policy.walletOperate = walletBiz
}

func (policy *privacyPolicy) getWalletOperate() wcom.WalletOperate {
	policy.mtx.Lock()
	defer policy.mtx.Unlock()
	return policy.walletOperate
}

// Init
func (policy *privacyPolicy) Init(walletOperate wcom.WalletOperate, sub []byte) {
	policy.setWalletOperate(walletOperate)
	policy.store = newStore(walletOperate.GetDBStore())
	//         FTXO
	walletOperate.GetWaitGroup().Add(1)
	go policy.checkWalletStoreData()
}

// OnCreateNewAccount
func (policy *privacyPolicy) OnCreateNewAccount(acc *types.Account) {
	wg := policy.getWalletOperate().GetWaitGroup()
	wg.Add(1)
	go policy.rescanReqTxDetailByAddr(acc.Addr, wg)
}

// OnImportPrivateKey
func (policy *privacyPolicy) OnImportPrivateKey(acc *types.Account) {
	wg := policy.getWalletOperate().GetWaitGroup()
	wg.Add(1)
	go policy.rescanReqTxDetailByAddr(acc.Addr, wg)
}

// OnAddBlockFinish
func (policy *privacyPolicy) OnAddBlockFinish(block *types.BlockDetail) {

}

// OnDeleteBlockFinish
func (policy *privacyPolicy) OnDeleteBlockFinish(block *types.BlockDetail) {

}

// OnClose
func (policy *privacyPolicy) OnClose() {

}

// OnSetQueueClient
func (policy *privacyPolicy) OnSetQueueClient() {
	version := policy.store.getVersion()
	if version < PRIVACYDBVERSION {
		policy.rescanAllTxAddToUpdateUTXOs()
		policy.store.setVersion()
	}
}

// OnWalletLocked
func (policy *privacyPolicy) OnWalletLocked() {
}

// OnWalletUnlocked
func (policy *privacyPolicy) OnWalletUnlocked(WalletUnLock *types.WalletUnLock) {
}

// Call
func (policy *privacyPolicy) Call(funName string, in types.Message) (ret types.Message, err error) {
	switch funName {
	case "GetUTXOScaningFlag":
		isok := policy.GetRescanFlag() == privacytypes.UtxoFlagScaning
		ret = &types.Reply{IsOk: isok}
	default:
		err = types.ErrNotSupport
	}
	return
}

// SignTransaction
func (policy *privacyPolicy) SignTransaction(key crypto.PrivKey, req *types.ReqSignRawTx) (needSysSign bool, signtxhex string, err error) {
	needSysSign = false
	bytes, err := common.FromHex(req.GetTxHex())
	if err != nil {
		bizlog.Error("SignTransaction", "common.FromHex error", err)
		return
	}
	tx := new(types.Transaction)
	if err = types.Decode(bytes, tx); err != nil {
		bizlog.Error("SignTransaction", "Decode Transaction error", err)
		return
	}
	signParam := &privacytypes.PrivacySignatureParam{}
	if err = types.Decode(tx.Signature.Signature, signParam); err != nil {
		bizlog.Error("SignTransaction", "Decode PrivacySignatureParam error", err)
		return
	}
	action := new(privacytypes.PrivacyAction)
	if err = types.Decode(tx.Payload, action); err != nil {
		bizlog.Error("SignTransaction", "Decode PrivacyAction error", err)
		return
	}
	if action.Ty != signParam.ActionType {
		bizlog.Error("SignTransaction", "action type ", action.Ty, "signature action type ", signParam.ActionType)
		return
	}
	switch action.Ty {
	case privacytypes.ActionPublic2Privacy:
		//           ，
		tx.Sign(int32(policy.getWalletOperate().GetSignType()), key)

	case privacytypes.ActionPrivacy2Privacy, privacytypes.ActionPrivacy2Public:
		//         、
		if err = policy.signatureTx(tx, action.GetInput(), signParam.GetUtxobasics(), signParam.GetRealKeyInputs()); err != nil {
			return
		}
	default:
		bizlog.Error("SignTransaction", "Invalid action type ", action.Ty)
		err = types.ErrInvalidParam
	}
	signtxhex = common.ToHex(types.Encode(tx))
	return
}

type buildStoreWalletTxDetailParam struct {
	assetExec    string
	tokenname    string
	block        *types.BlockDetail
	tx           *types.Transaction
	index        int
	newbatch     db.Batch
	senderRecver string
	isprivacy    bool
	addDelType   int32
	sendRecvFlag int32
	utxos        []*privacytypes.UTXO
}

// OnAddBlockTx
func (policy *privacyPolicy) OnAddBlockTx(block *types.BlockDetail, tx *types.Transaction, index int32, dbbatch db.Batch) *types.WalletTxDetail {
	policy.addDelPrivacyTxsFromBlock(tx, index, block, dbbatch, AddTx)
	//          ，
	return nil
}

// OnDeleteBlockTx
func (policy *privacyPolicy) OnDeleteBlockTx(block *types.BlockDetail, tx *types.Transaction, index int32, dbbatch db.Batch) *types.WalletTxDetail {
	policy.addDelPrivacyTxsFromBlock(tx, index, block, dbbatch, DelTx)
	//          ，
	return nil
}

// GetRescanFlag get rescan utxo flag
func (policy *privacyPolicy) GetRescanFlag() int32 {
	return atomic.LoadInt32(&policy.rescanUTXOflag)
}

// SetRescanFlag set rescan utxos flag
func (policy *privacyPolicy) SetRescanFlag(flag int32) {
	atomic.StoreInt32(&policy.rescanUTXOflag, flag)
}
