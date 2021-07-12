// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	wcom "github.com/D-PlatformOperatingSystem/dpos/wallet/common"
	privacy "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/crypto"
	privacytypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/types"
	"github.com/golang/protobuf/proto"
)

func (policy *privacyPolicy) rescanAllTxAddToUpdateUTXOs() {
	accounts, err := policy.getWalletOperate().GetWalletAccounts()
	if err != nil {
		bizlog.Debug("rescanAllTxToUpdateUTXOs", "walletOperate.GetWalletAccounts error", err)
		return
	}
	bizlog.Debug("rescanAllTxToUpdateUTXOs begin!")
	for _, acc := range accounts {
		// blockchain    Account.Addr
		policy.rescanwg.Add(1)
		go policy.rescanReqTxDetailByAddr(acc.Addr, policy.rescanwg)
	}
	policy.rescanwg.Wait()
	bizlog.Debug("rescanAllTxToUpdateUTXOs success!")
}

// blockchain    addr
func (policy *privacyPolicy) rescanReqTxDetailByAddr(addr string, wg *sync.WaitGroup) {
	defer wg.Done()
	policy.reqTxDetailByAddr(addr)
}

// blockchain    addr
func (policy *privacyPolicy) reqTxDetailByAddr(addr string) {
	if len(addr) == 0 {
		bizlog.Error("reqTxDetailByAddr input addr is nil!")
		return
	}
	var txInfo types.ReplyTxInfo

	i := 0
	operater := policy.getWalletOperate()
	for {
		//   blockchain             hashs  ,
		var ReqAddr types.ReqAddr
		ReqAddr.Addr = addr
		ReqAddr.Flag = 0
		ReqAddr.Direction = 0
		ReqAddr.Count = int32(MaxTxHashsPerTime)
		if i == 0 {
			ReqAddr.Height = -1
			ReqAddr.Index = 0
		} else {
			ReqAddr.Height = txInfo.GetHeight()
			ReqAddr.Index = txInfo.GetIndex()
		}
		i++
		ReplyTxInfos, err := operater.GetAPI().GetTransactionByAddr(&ReqAddr)
		if err != nil {
			bizlog.Error("reqTxDetailByAddr", "GetTransactionByAddr error", err, "addr", addr)
			return
		}
		if ReplyTxInfos == nil {
			bizlog.Info("reqTxDetailByAddr ReplyTxInfos is nil")
			return
		}
		txcount := len(ReplyTxInfos.TxInfos)
		var ReqHashes types.ReqHashes
		ReqHashes.Hashes = make([][]byte, len(ReplyTxInfos.TxInfos))
		for index, ReplyTxInfo := range ReplyTxInfos.TxInfos {
			ReqHashes.Hashes[index] = ReplyTxInfo.GetHash()
			txInfo.Hash = ReplyTxInfo.GetHash()
			txInfo.Height = ReplyTxInfo.GetHeight()
			txInfo.Index = ReplyTxInfo.GetIndex()
		}
		operater.GetTxDetailByHashs(&ReqHashes)
		if txcount < int(MaxTxHashsPerTime) {
			return
		}
	}
}

func (policy *privacyPolicy) isRescanUtxosFlagScaning() (bool, error) {
	if privacytypes.UtxoFlagScaning == policy.GetRescanFlag() {
		return true, privacytypes.ErrRescanFlagScaning
	}
	return false, nil
}

func (policy *privacyPolicy) parseViewSpendPubKeyPair(in string) (viewPubKey, spendPubKey []byte, err error) {
	src, err := common.FromHex(in)
	if err != nil {
		return nil, nil, err
	}
	if 64 != len(src) {
		bizlog.Error("parseViewSpendPubKeyPair", "pair with len", len(src))
		return nil, nil, types.ErrPubKeyLen
	}
	viewPubKey = src[:32]
	spendPubKey = src[32:]
	return
}

func (policy *privacyPolicy) getPrivKeyByAddr(addr string) (crypto.PrivKey, error) {
	//
	Accountstor, err := policy.store.getAccountByAddr(addr)
	if err != nil {
		bizlog.Error("ProcSendToAddress", "GetAccountByAddr err:", err)
		return nil, err
	}

	//  password
	prikeybyte, err := common.FromHex(Accountstor.GetPrivkey())
	if err != nil || len(prikeybyte) == 0 {
		bizlog.Error("ProcSendToAddress", "FromHex err", err)
		return nil, err
	}
	operater := policy.getWalletOperate()
	password := []byte(operater.GetPassword())
	privkey := wcom.CBCDecrypterPrivkey(password, prikeybyte)
	//  privkey    pubkey        addr
	cr, err := crypto.New(types.GetSignName("privacy", operater.GetSignType()))
	if err != nil {
		bizlog.Error("ProcSendToAddress", "err", err)
		return nil, err
	}
	priv, err := cr.PrivKeyFromBytes(privkey)
	if err != nil {
		bizlog.Error("ProcSendToAddress", "PrivKeyFromBytes err", err)
		return nil, err
	}
	return priv, nil
}

func (policy *privacyPolicy) getPrivacykeyPair(addr string) (*privacy.Privacy, error) {
	if accPrivacy, _ := policy.store.getWalletAccountPrivacy(addr); accPrivacy != nil {
		privacyInfo := &privacy.Privacy{}
		password := []byte(policy.getWalletOperate().GetPassword())
		copy(privacyInfo.ViewPubkey[:], accPrivacy.ViewPubkey)
		decrypteredView := wcom.CBCDecrypterPrivkey(password, accPrivacy.ViewPrivKey)
		copy(privacyInfo.ViewPrivKey[:], decrypteredView)
		copy(privacyInfo.SpendPubkey[:], accPrivacy.SpendPubkey)
		decrypteredSpend := wcom.CBCDecrypterPrivkey(password, accPrivacy.SpendPrivKey)
		copy(privacyInfo.SpendPrivKey[:], decrypteredSpend)

		return privacyInfo, nil
	}
	_, err := policy.getPrivKeyByAddr(addr)
	if err != nil {
		return nil, err
	}
	return nil, privacytypes.ErrPrivacyNotEnabled

}

func (policy *privacyPolicy) savePrivacykeyPair(addr string) (*privacy.Privacy, error) {
	priv, err := policy.getPrivKeyByAddr(addr)
	if err != nil {
		return nil, err
	}

	newPrivacy, err := privacy.NewPrivacyWithPrivKey((*[privacy.KeyLen32]byte)(unsafe.Pointer(&priv.Bytes()[0])))
	if err != nil {
		return nil, err
	}

	password := []byte(policy.getWalletOperate().GetPassword())
	encrypteredView := wcom.CBCEncrypterPrivkey(password, newPrivacy.ViewPrivKey.Bytes())
	encrypteredSpend := wcom.CBCEncrypterPrivkey(password, newPrivacy.SpendPrivKey.Bytes())
	walletPrivacy := &privacytypes.WalletAccountPrivacy{
		ViewPubkey:   newPrivacy.ViewPubkey[:],
		ViewPrivKey:  encrypteredView,
		SpendPubkey:  newPrivacy.SpendPubkey[:],
		SpendPrivKey: encrypteredSpend,
	}
	//save the privacy created to wallet db
	policy.store.setWalletAccountPrivacy(addr, walletPrivacy)
	return newPrivacy, nil
}

func (policy *privacyPolicy) enablePrivacy(req *privacytypes.ReqEnablePrivacy) (*privacytypes.RepEnablePrivacy, error) {
	var addrs []string
	if 0 == len(req.Addrs) {
		WalletAccStores, err := policy.store.getAccountByPrefix("Account")
		if err != nil || len(WalletAccStores) == 0 {
			bizlog.Info("enablePrivacy", "GetAccountByPrefix:err", err)
			return nil, types.ErrNotFound
		}
		for _, WalletAccStore := range WalletAccStores {
			addrs = append(addrs, WalletAccStore.Addr)
		}
	} else {
		addrs = append(addrs, req.Addrs...)
	}

	var rep privacytypes.RepEnablePrivacy
	for _, addr := range addrs {
		str := ""
		isOK := true
		_, err := policy.getPrivacykeyPair(addr)
		if err != nil {
			_, err = policy.savePrivacykeyPair(addr)
			if err != nil {
				isOK = false
				str = err.Error()
			}
		}

		priAddrResult := &privacytypes.PriAddrResult{
			Addr: addr,
			IsOK: isOK,
			Msg:  str,
		}

		rep.Results = append(rep.Results, priAddrResult)
	}
	return &rep, nil
}

func (policy *privacyPolicy) showPrivacyKeyPair(reqAddr *types.ReqString) (*privacytypes.ReplyPrivacyPkPair, error) {
	privacyInfo, err := policy.getPrivacykeyPair(reqAddr.GetData())
	if err != nil {
		bizlog.Error("showPrivacyKeyPair", "getPrivacykeyPair error ", err)
		return nil, err
	}

	//pair := privacyInfo.ViewPubkey[:]
	//pair = append(pair, privacyInfo.SpendPubkey[:]...)

	replyPrivacyPkPair := &privacytypes.ReplyPrivacyPkPair{
		ShowSuccessful: true,
		Pubkeypair:     makeViewSpendPubKeyPairToString(privacyInfo.ViewPubkey[:], privacyInfo.SpendPubkey[:]),
	}
	return replyPrivacyPkPair, nil
}

func (policy *privacyPolicy) getPrivacyAccountInfo(req *privacytypes.ReqPrivacyAccount) (*privacytypes.ReplyPrivacyAccount, error) {
	addr := strings.Trim(req.GetAddr(), " ")
	token := req.GetToken()
	reply := &privacytypes.ReplyPrivacyAccount{}
	reply.Displaymode = req.Displaymode
	if len(addr) == 0 {
		return nil, errors.New("Address is empty")
	}

	//
	privacyDBStore, err := policy.store.listAvailableUTXOs(req.GetAssetExec(), token, addr)
	if err != nil {
		bizlog.Error("getPrivacyAccountInfo", "listAvailableUTXOs")
		return nil, err
	}
	utxos := make([]*privacytypes.UTXO, 0)
	for _, ele := range privacyDBStore {
		utxoBasic := &privacytypes.UTXOBasic{
			UtxoGlobalIndex: &privacytypes.UTXOGlobalIndex{
				Outindex: ele.OutIndex,
				Txhash:   ele.Txhash,
			},
			OnetimePubkey: ele.OnetimePublicKey,
		}
		utxo := &privacytypes.UTXO{
			Amount:    ele.Amount,
			UtxoBasic: utxoBasic,
		}
		utxos = append(utxos, utxo)
	}
	reply.Utxos = &privacytypes.UTXOs{Utxos: utxos}

	//
	utxos = make([]*privacytypes.UTXO, 0)
	ftxoslice, err := policy.store.listFrozenUTXOs(req.GetAssetExec(), token, addr)
	if err == nil && ftxoslice != nil {
		for _, ele := range ftxoslice {
			utxos = append(utxos, ele.Utxos...)
		}
	}

	reply.Ftxos = &privacytypes.UTXOs{Utxos: utxos}

	return reply, nil
}

//     UTXO
//     UTXO         12      UTXO
//                12     UTXO
//         UTXO    ，        ，        ，    ，          ，    ，
func (policy *privacyPolicy) selectUTXO(assetExec, token, addr string, amount int64) ([]*txOutputInfo, error) {
	if len(token) == 0 || len(addr) == 0 || amount <= 0 {
		return nil, types.ErrInvalidParam
	}
	wutxos, err := policy.store.getPrivacyTokenUTXOs(assetExec, token, addr)
	if err != nil {
		return nil, types.ErrInsufficientBalance
	}
	operater := policy.getWalletOperate()
	curBlockHeight := operater.GetBlockHeight()
	var confirmUTXOs, unconfirmUTXOs []*walletUTXO
	var balance int64
	for _, wutxo := range wutxos.utxos {
		if curBlockHeight < wutxo.height {
			continue
		}
		if curBlockHeight-wutxo.height > privacytypes.UtxoMaturityDegree {
			balance += wutxo.outinfo.amount
			confirmUTXOs = append(confirmUTXOs, wutxo)
		} else {
			unconfirmUTXOs = append(unconfirmUTXOs, wutxo)
		}
	}
	if balance < amount && len(unconfirmUTXOs) > 0 {
		//      UTXO     ，            ，
		//
		sort.Slice(unconfirmUTXOs, func(i, j int) bool {
			return unconfirmUTXOs[i].height < unconfirmUTXOs[j].height
		})
		for _, wutxo := range unconfirmUTXOs {
			confirmUTXOs = append(confirmUTXOs, wutxo)
			balance += wutxo.outinfo.amount
			if balance >= amount {
				break
			}
		}
	}
	if balance < amount {
		return nil, types.ErrInsufficientBalance
	}
	balance = 0
	var selectedOuts []*txOutputInfo
	for balance < amount {
		index := operater.GetRandom().Intn(len(confirmUTXOs))
		selectedOuts = append(selectedOuts, confirmUTXOs[index].outinfo)
		balance += confirmUTXOs[index].outinfo.amount
		// remove selected utxo
		confirmUTXOs = append(confirmUTXOs[:index], confirmUTXOs[index+1:]...)
	}
	return selectedOuts, nil
}

/*
buildInput

	1.                   UTXO
	2.      (mixcout>0)，   UTXO               UTXO，   UTXO
	3.     x=Hs(aR)+b，       ，   xG = Hs(ar)G+bG = Hs(aR)G+B，
*/
func (policy *privacyPolicy) buildInput(privacykeyParirs *privacy.Privacy, buildInfo *buildInputInfo) (*privacytypes.PrivacyInput, []*privacytypes.UTXOBasics, []*privacytypes.RealKeyInput, []*txOutputInfo, error) {
	operater := policy.getWalletOperate()
	//       utxo
	selectedUtxo, err := policy.selectUTXO(buildInfo.assetExec, buildInfo.assetSymbol, buildInfo.sender, buildInfo.amount)
	if err != nil {
		bizlog.Error("buildInput", "Failed to selectOutput for amount", buildInfo.amount,
			"Due to cause", err)
		return nil, nil, nil, nil, err
	}
	sort.Slice(selectedUtxo, func(i, j int) bool {
		return selectedUtxo[i].amount <= selectedUtxo[j].amount
	})

	reqGetGlobalIndex := privacytypes.ReqUTXOGlobalIndex{
		AssetExec:   buildInfo.assetExec,
		AssetSymbol: buildInfo.assetSymbol,
		MixCount:    0,
	}

	if buildInfo.mixcount > 0 {
		reqGetGlobalIndex.MixCount = common.MinInt32(int32(privacytypes.PrivacyMaxCount), common.MaxInt32(buildInfo.mixcount, 0))
	}
	for _, out := range selectedUtxo {
		reqGetGlobalIndex.Amount = append(reqGetGlobalIndex.Amount, out.amount)
	}
	//      0    blockchain
	var resUTXOGlobalIndex *privacytypes.ResUTXOGlobalIndex
	if buildInfo.mixcount > 0 {
		query := &types.ChainExecutor{
			Driver:   "privacy",
			FuncName: "GetUTXOGlobalIndex",
			Param:    types.Encode(&reqGetGlobalIndex),
		}
		// blockchain         utxo
		data, err := operater.GetAPI().QueryChain(query)
		if err != nil {
			bizlog.Error("buildInput BlockChainQuery", "err", err)
			return nil, nil, nil, nil, err
		}
		resUTXOGlobalIndex = data.(*privacytypes.ResUTXOGlobalIndex)
		if resUTXOGlobalIndex == nil {
			bizlog.Info("buildInput EventBlockChainQuery is nil")
			return nil, nil, nil, nil, err
		}

		sort.Slice(resUTXOGlobalIndex.UtxoIndex4Amount, func(i, j int) bool {
			return resUTXOGlobalIndex.UtxoIndex4Amount[i].Amount <= resUTXOGlobalIndex.UtxoIndex4Amount[j].Amount
		})

		if len(selectedUtxo) != len(resUTXOGlobalIndex.UtxoIndex4Amount) {
			bizlog.Error("buildInput EventBlockChainQuery get not the same count for mix",
				"len(selectedUtxo)", len(selectedUtxo),
				"len(resUTXOGlobalIndex.UtxoIndex4Amount)", len(resUTXOGlobalIndex.UtxoIndex4Amount))
		}
	}

	//    PrivacyInput
	privacyInput := &privacytypes.PrivacyInput{}
	utxosInKeyInput := make([]*privacytypes.UTXOBasics, len(selectedUtxo))
	realkeyInputSlice := make([]*privacytypes.RealKeyInput, len(selectedUtxo))
	for i, utxo2pay := range selectedUtxo {
		var utxoIndex4Amount *privacytypes.UTXOIndex4Amount
		if nil != resUTXOGlobalIndex && i < len(resUTXOGlobalIndex.UtxoIndex4Amount) && utxo2pay.amount == resUTXOGlobalIndex.UtxoIndex4Amount[i].Amount {
			utxoIndex4Amount = resUTXOGlobalIndex.UtxoIndex4Amount[i]
			for j, utxo := range utxoIndex4Amount.Utxos {
				//      UTXO    ，
				if bytes.Equal(utxo.OnetimePubkey, utxo2pay.onetimePublicKey) {
					utxoIndex4Amount.Utxos = append(utxoIndex4Amount.Utxos[:j], utxoIndex4Amount.Utxos[j+1:]...)
					break
				}
			}
		}

		if utxoIndex4Amount == nil {
			utxoIndex4Amount = &privacytypes.UTXOIndex4Amount{}
		}
		if utxoIndex4Amount.Utxos == nil {
			utxoIndex4Amount.Utxos = make([]*privacytypes.UTXOBasic, 0)
		}
		//            utxo        mix   ，      utxo  ，
		if len(utxoIndex4Amount.Utxos) > int(buildInfo.mixcount) {
			utxoIndex4Amount.Utxos = utxoIndex4Amount.Utxos[:len(utxoIndex4Amount.Utxos)-1]
		}

		utxo := &privacytypes.UTXOBasic{
			UtxoGlobalIndex: utxo2pay.utxoGlobalIndex,
			OnetimePubkey:   utxo2pay.onetimePublicKey,
		}
		//    utxo
		utxoIndex4Amount.Utxos = append(utxoIndex4Amount.Utxos, utxo)
		positions := operater.GetRandom().Perm(len(utxoIndex4Amount.Utxos))
		utxos := make([]*privacytypes.UTXOBasic, len(utxoIndex4Amount.Utxos))
		for k, position := range positions {
			utxos[position] = utxoIndex4Amount.Utxos[k]
		}
		utxosInKeyInput[i] = &privacytypes.UTXOBasics{Utxos: utxos}

		//x = Hs(aR) + b
		onetimePriv, err := privacy.RecoverOnetimePriKey(utxo2pay.txPublicKeyR, privacykeyParirs.ViewPrivKey, privacykeyParirs.SpendPrivKey, int64(utxo2pay.utxoGlobalIndex.Outindex))
		if err != nil {
			bizlog.Error("transPri2Pri", "Failed to RecoverOnetimePriKey", err)
			return nil, nil, nil, nil, err
		}

		realkeyInput := &privacytypes.RealKeyInput{
			Realinputkey:   int32(positions[len(positions)-1]),
			Onetimeprivkey: onetimePriv.Bytes(),
		}
		realkeyInputSlice[i] = realkeyInput

		keyImage, err := privacy.GenerateKeyImage(onetimePriv, utxo2pay.onetimePublicKey)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		keyInput := &privacytypes.KeyInput{
			Amount:   utxo2pay.amount,
			KeyImage: keyImage[:],
		}

		for _, utxo := range utxos {
			keyInput.UtxoGlobalIndex = append(keyInput.UtxoGlobalIndex, utxo.UtxoGlobalIndex)
		}
		//    input   ，           ，keyImage   ，
		//       ，            utxo     keyinput              pubkey
		//
		privacyInput.Keyinput = append(privacyInput.Keyinput, keyInput)
	}

	return privacyInput, utxosInKeyInput, realkeyInputSlice, selectedUtxo, nil
}

func (policy *privacyPolicy) createTransaction(req *privacytypes.ReqCreatePrivacyTx) (*types.Transaction, error) {
	switch req.ActionType {
	case privacytypes.ActionPublic2Privacy:
		return policy.createPublic2PrivacyTx(req)
	case privacytypes.ActionPrivacy2Privacy:
		return policy.createPrivacy2PrivacyTx(req)
	case privacytypes.ActionPrivacy2Public:
		return policy.createPrivacy2PublicTx(req)
	}
	return nil, types.ErrInvalidParam
}

func (policy *privacyPolicy) createPublic2PrivacyTx(req *privacytypes.ReqCreatePrivacyTx) (*types.Transaction, error) {
	viewPubSlice, spendPubSlice, err := parseViewSpendPubKeyPair(req.GetPubkeypair())
	if err != nil {
		bizlog.Error("createPublic2PrivacyTx", "parse view spend public key pair failed.  err ", err)
		return nil, err
	}
	amount := req.GetAmount()
	viewPublic := (*[32]byte)(unsafe.Pointer(&viewPubSlice[0]))
	spendPublic := (*[32]byte)(unsafe.Pointer(&spendPubSlice[0]))
	privacyOutput, err := generateOuts(viewPublic, spendPublic, nil, nil, amount, amount, 0)
	if err != nil {
		bizlog.Error("createPublic2PrivacyTx", "generate output failed.  err ", err)
		return nil, err
	}

	value := &privacytypes.Public2Privacy{
		Tokenname: req.Tokenname,
		Amount:    amount,
		Note:      req.GetNote(),
		Output:    privacyOutput,
		AssetExec: req.GetAssetExec(),
	}
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	action := &privacytypes.PrivacyAction{
		Ty:    privacytypes.ActionPublic2Privacy,
		Value: &privacytypes.PrivacyAction_Public2Privacy{Public2Privacy: value},
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(privacytypes.PrivacyX)),
		Payload: types.Encode(action),
		Nonce:   policy.getWalletOperate().Nonce(),
		To:      address.ExecAddress(cfg.ExecName(privacytypes.PrivacyX)),
	}
	tx.SetExpire(cfg, time.Duration(req.Expire))
	tx.Signature = &types.Signature{
		Signature: types.Encode(&privacytypes.PrivacySignatureParam{
			ActionType: action.Ty,
		}),
	}
	tx.Fee, err = tx.GetRealFee(cfg.GetMinTxFeeRate())
	if err != nil {
		bizlog.Error("createPublic2PrivacyTx", "calc fee failed", err)
		return nil, err
	}

	return tx, nil
}

func (policy *privacyPolicy) createPrivacy2PrivacyTx(req *privacytypes.ReqCreatePrivacyTx) (*types.Transaction, error) {

	//     utxo
	var utxoBurnedAmount int64
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	isMainetCoins := !cfg.IsPara() && (req.AssetExec == "coins")
	if isMainetCoins {
		utxoBurnedAmount = privacytypes.PrivacyTxFee
	}
	buildInfo := &buildInputInfo{
		assetExec:   req.GetAssetExec(),
		assetSymbol: req.GetTokenname(),
		sender:      req.GetFrom(),
		amount:      req.GetAmount() + utxoBurnedAmount,
		mixcount:    req.GetMixcount(),
	}
	privacyInfo, err := policy.getPrivacykeyPair(req.GetFrom())
	if err != nil {
		bizlog.Error("createPrivacy2PrivacyTx", "getPrivacykeyPair error", err)
		return nil, err
	}
	//step 1,buildInput
	privacyInput, utxosInKeyInput, realkeyInputSlice, selectedUtxo, err := policy.buildInput(privacyInfo, buildInfo)
	if err != nil {
		return nil, err
	}
	//step 2,generateOuts
	viewPublicSlice, spendPublicSlice, err := parseViewSpendPubKeyPair(req.GetPubkeypair())
	if err != nil {
		bizlog.Error("createPrivacy2PrivacyTx", "parseViewSpendPubKeyPair  ", err)
		return nil, err
	}

	viewPub4change, spendPub4change := privacyInfo.ViewPubkey.Bytes(), privacyInfo.SpendPubkey.Bytes()
	viewPublic := (*[32]byte)(unsafe.Pointer(&viewPublicSlice[0]))
	spendPublic := (*[32]byte)(unsafe.Pointer(&spendPublicSlice[0]))
	viewPub4chgPtr := (*[32]byte)(unsafe.Pointer(&viewPub4change[0]))
	spendPub4chgPtr := (*[32]byte)(unsafe.Pointer(&spendPub4change[0]))

	selectedAmounTotal := int64(0)
	for _, input := range privacyInput.Keyinput {
		selectedAmounTotal += input.Amount
	}
	//    UTXO
	privacyOutput, err := generateOuts(viewPublic, spendPublic, viewPub4chgPtr, spendPub4chgPtr, req.GetAmount(), selectedAmounTotal, utxoBurnedAmount)
	if err != nil {
		return nil, err
	}

	value := &privacytypes.Privacy2Privacy{
		Tokenname: req.GetTokenname(),
		Amount:    req.GetAmount(),
		Note:      req.GetNote(),
		Input:     privacyInput,
		Output:    privacyOutput,
		AssetExec: req.GetAssetExec(),
	}
	action := &privacytypes.PrivacyAction{
		Ty:    privacytypes.ActionPrivacy2Privacy,
		Value: &privacytypes.PrivacyAction_Privacy2Privacy{Privacy2Privacy: value},
	}

	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(privacytypes.PrivacyX)),
		Payload: types.Encode(action),
		Fee:     privacytypes.PrivacyTxFee,
		Nonce:   policy.getWalletOperate().Nonce(),
		To:      address.ExecAddress(cfg.ExecName(privacytypes.PrivacyX)),
	}
	tx.SetExpire(cfg, time.Duration(req.Expire))
	if !isMainetCoins {
		tx.Fee, err = tx.GetRealFee(cfg.GetMinTxFeeRate())
		if err != nil {
			bizlog.Error("createPrivacy2PrivacyTx", "calc fee failed", err)
			return nil, err
		}
	}

	//       ，       UTXO  ，         txHash
	policy.saveFTXOInfo(tx.GetExpire(), req.GetAssetExec(), req.Tokenname, req.GetFrom(), hex.EncodeToString(tx.Hash()), selectedUtxo)
	tx.Signature = &types.Signature{
		Signature: types.Encode(&privacytypes.PrivacySignatureParam{
			ActionType:    action.Ty,
			Utxobasics:    utxosInKeyInput,
			RealKeyInputs: realkeyInputSlice,
		}),
	}
	return tx, nil
}

func (policy *privacyPolicy) createPrivacy2PublicTx(req *privacytypes.ReqCreatePrivacyTx) (*types.Transaction, error) {

	//     utxo
	//     utxo
	var utxoBurnedAmount int64
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	isMainetCoins := !cfg.IsPara() && (req.AssetExec == "coins")
	if isMainetCoins {
		utxoBurnedAmount = privacytypes.PrivacyTxFee
	}
	buildInfo := &buildInputInfo{
		assetExec:   req.GetAssetExec(),
		assetSymbol: req.GetTokenname(),
		sender:      req.GetFrom(),
		amount:      req.GetAmount() + utxoBurnedAmount,
		mixcount:    req.GetMixcount(),
	}
	privacyInfo, err := policy.getPrivacykeyPair(req.GetFrom())
	if err != nil {
		bizlog.Error("createPrivacy2PublicTx failed to getPrivacykeyPair")
		return nil, err
	}
	//step 1,buildInput
	privacyInput, utxosInKeyInput, realkeyInputSlice, selectedUtxo, err := policy.buildInput(privacyInfo, buildInfo)
	if err != nil {
		bizlog.Error("createPrivacy2PublicTx failed to buildInput")
		return nil, err
	}

	viewPub4change, spendPub4change := privacyInfo.ViewPubkey.Bytes(), privacyInfo.SpendPubkey.Bytes()
	viewPub4chgPtr := (*[32]byte)(unsafe.Pointer(&viewPub4change[0]))
	spendPub4chgPtr := (*[32]byte)(unsafe.Pointer(&spendPub4change[0]))

	selectedAmounTotal := int64(0)
	for _, input := range privacyInput.Keyinput {
		if input.Amount <= 0 {
			return nil, errors.New("")
		}
		selectedAmounTotal += input.Amount
	}
	changeAmount := selectedAmounTotal - req.GetAmount()
	//step 2,generateOuts
	//    UTXO,      UTXO
	privacyOutput, err := generateOuts(nil, nil, viewPub4chgPtr, spendPub4chgPtr, 0, changeAmount, utxoBurnedAmount)
	if err != nil {
		return nil, err
	}

	value := &privacytypes.Privacy2Public{
		Tokenname: req.GetTokenname(),
		Amount:    req.GetAmount(),
		Note:      req.GetNote(),
		Input:     privacyInput,
		Output:    privacyOutput,
		To:        req.GetTo(),
		AssetExec: req.GetAssetExec(),
	}
	action := &privacytypes.PrivacyAction{
		Ty:    privacytypes.ActionPrivacy2Public,
		Value: &privacytypes.PrivacyAction_Privacy2Public{Privacy2Public: value},
	}

	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(privacytypes.PrivacyX)),
		Payload: types.Encode(action),
		Fee:     privacytypes.PrivacyTxFee,
		Nonce:   policy.getWalletOperate().Nonce(),
		To:      address.ExecAddress(cfg.ExecName(privacytypes.PrivacyX)),
	}
	tx.SetExpire(cfg, time.Duration(req.Expire))
	if !isMainetCoins {
		tx.Fee, err = tx.GetRealFee(cfg.GetMinTxFeeRate())
		if err != nil {
			bizlog.Error("createPrivacy2PublicTx", "calc fee failed", err)
			return nil, err
		}
	}
	//       ，       UTXO  ，         txHash
	policy.saveFTXOInfo(tx.GetExpire(), req.GetAssetExec(), req.Tokenname, req.GetFrom(), hex.EncodeToString(tx.Hash()), selectedUtxo)
	tx.Signature = &types.Signature{
		Signature: types.Encode(&privacytypes.PrivacySignatureParam{
			ActionType:    action.Ty,
			Utxobasics:    utxosInKeyInput,
			RealKeyInputs: realkeyInputSlice,
		}),
	}
	return tx, nil
}

func (policy *privacyPolicy) saveFTXOInfo(expire int64, assetExec, assetSymbol, sender, txhash string, selectedUtxos []*txOutputInfo) {
	//            utxo    ，
	policy.store.moveUTXO2FTXO(expire, assetExec, assetSymbol, sender, txhash, selectedUtxos)
	//TODO:        ，      txhash       ，                     ，
	//TODO:            ，   FTXO   STXO，added by hezhengjun on 2018.6.5
}

func (policy *privacyPolicy) getPrivacyKeyPairs() ([]addrAndprivacy, error) {
	//  Account
	WalletAccStores, err := policy.store.getAccountByPrefix("Account")
	if err != nil || len(WalletAccStores) == 0 {
		bizlog.Info("getPrivacyKeyPairs", "store getAccountByPrefix error", err)
		return nil, err
	}

	var infoPriRes []addrAndprivacy
	for _, AccStore := range WalletAccStores {
		if len(AccStore.Addr) != 0 {
			if privacyInfo, err := policy.getPrivacykeyPair(AccStore.Addr); err == nil {
				var priInfo addrAndprivacy
				priInfo.Addr = &AccStore.Addr
				priInfo.PrivacyKeyPair = privacyInfo
				infoPriRes = append(infoPriRes, priInfo)
			}
		}
	}

	if 0 == len(infoPriRes) {
		return nil, privacytypes.ErrPrivacyNotEnabled
	}

	return infoPriRes, nil

}

func (policy *privacyPolicy) rescanUTXOs(req *privacytypes.ReqRescanUtxos) (*privacytypes.RepRescanUtxos, error) {
	if req.Flag != 0 {
		return policy.store.getRescanUtxosFlag4Addr(req)
	}
	// Rescan
	var repRescanUtxos privacytypes.RepRescanUtxos
	repRescanUtxos.Flag = req.Flag

	operater := policy.getWalletOperate()
	if operater.IsWalletLocked() {
		return nil, types.ErrWalletIsLocked
	}
	if ok, err := policy.isRescanUtxosFlagScaning(); ok {
		return nil, err
	}
	_, err := policy.getPrivacyKeyPairs()
	if err != nil {
		return nil, err
	}
	policy.SetRescanFlag(privacytypes.UtxoFlagScaning)
	operater.GetWaitGroup().Add(1)
	go policy.rescanReqUtxosByAddr(req.Addrs)
	return &repRescanUtxos, nil
}

// blockchain    addr
func (policy *privacyPolicy) rescanReqUtxosByAddr(addrs []string) {
	defer policy.getWalletOperate().GetWaitGroup().Done()
	bizlog.Debug("RescanAllUTXO begin!")
	policy.reqUtxosByAddr(addrs)
	bizlog.Debug("RescanAllUTXO success!")
}

func (policy *privacyPolicy) reqUtxosByAddr(addrs []string) {
	//
	var storeAddrs []string
	if len(addrs) == 0 {
		WalletAccStores, err := policy.store.getAccountByPrefix("Account")
		if err != nil || len(WalletAccStores) == 0 {
			bizlog.Info("reqUtxosByAddr", "getAccountByPrefix error", err)
			return
		}
		for _, WalletAccStore := range WalletAccStores {
			storeAddrs = append(storeAddrs, WalletAccStore.Addr)
		}
	} else {
		storeAddrs = append(storeAddrs, addrs...)
	}
	policy.store.saveREscanUTXOsAddresses(storeAddrs, privacytypes.UtxoFlagScaning)

	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	reqAddr := address.ExecAddress(cfg.ExecName(privacytypes.PrivacyX))
	var txInfo types.ReplyTxInfo
	i := 0
	operater := policy.getWalletOperate()
	for {
		select {
		case <-operater.GetWalletDone():
			return
		default:
		}

		//   execs           UTXOs,
		// 1
		var ReqAddr types.ReqAddr
		ReqAddr.Addr = reqAddr
		ReqAddr.Flag = 0
		ReqAddr.Direction = 0
		ReqAddr.Count = int32(MaxTxHashsPerTime)
		if i == 0 {
			ReqAddr.Height = -1
			ReqAddr.Index = 0
		} else {
			ReqAddr.Height = txInfo.GetHeight()
			ReqAddr.Index = txInfo.GetIndex()
			if !cfg.IsDappFork(ReqAddr.Height, privacytypes.PrivacyX, "ForkV21Privacy") { //
				break
			}
		}
		i++
		//
		msg, err := operater.GetAPI().Query(privacytypes.PrivacyX, "GetTxsByAddr", &ReqAddr)
		if err != nil {
			bizlog.Error("reqUtxosByAddr", "GetTxsByAddr error", err, "addr", reqAddr)
			break
		}
		ReplyTxInfos := msg.(*types.ReplyTxInfos)
		if ReplyTxInfos == nil {
			bizlog.Info("privacy ReqTxInfosByAddr ReplyTxInfos is nil")
			break
		}
		txcount := len(ReplyTxInfos.TxInfos)

		var ReqHashes types.ReqHashes
		ReqHashes.Hashes = make([][]byte, len(ReplyTxInfos.TxInfos))
		for index, ReplyTxInfo := range ReplyTxInfos.TxInfos {
			ReqHashes.Hashes[index] = ReplyTxInfo.GetHash()
		}

		if txcount > 0 {
			txInfo.Hash = ReplyTxInfos.TxInfos[txcount-1].GetHash()
			txInfo.Height = ReplyTxInfos.TxInfos[txcount-1].GetHeight()
			txInfo.Index = ReplyTxInfos.TxInfos[txcount-1].GetIndex()
		}

		policy.getPrivacyTxDetailByHashs(&ReqHashes, addrs)
		if txcount < int(MaxTxHashsPerTime) {
			break
		}
	}
	//
	policy.SetRescanFlag(privacytypes.UtxoFlagNoScan)
	//   privacyInput
	policy.deleteScanPrivacyInputUtxo()
	policy.store.saveREscanUTXOsAddresses(storeAddrs, privacytypes.UtxoFlagScanEnd)
}

//TODO:input       utxo,          utxo
func (policy *privacyPolicy) deleteScanPrivacyInputUtxo() {
	maxUTXOsPerTime := 1000
	for {
		utxoGlobalIndexs := policy.store.setScanPrivacyInputUTXO(int32(maxUTXOsPerTime))
		policy.store.updateScanInputUTXOs(utxoGlobalIndexs)
		if len(utxoGlobalIndexs) < maxUTXOsPerTime {
			break
		}
	}
}

func (policy *privacyPolicy) getPrivacyTxDetailByHashs(ReqHashes *types.ReqHashes, addrs []string) {
	//  txhashs     txdetail
	TxDetails, err := policy.getWalletOperate().GetAPI().GetTransactionByHash(ReqHashes)
	if err != nil {
		bizlog.Error("getPrivacyTxDetailByHashs", "GetTransactionByHash error", err)
		return
	}
	var privacyInfo []addrAndprivacy
	if len(addrs) > 0 {
		for _, addr := range addrs {
			if privacy, err := policy.getPrivacykeyPair(addr); err == nil {
				priInfo := &addrAndprivacy{
					Addr:           &addr,
					PrivacyKeyPair: privacy,
				}
				privacyInfo = append(privacyInfo, *priInfo)
			}

		}
	} else {
		privacyInfo, _ = policy.getPrivacyKeyPairs()
	}
	policy.store.selectPrivacyTransactionToWallet(TxDetails, privacyInfo)
}

func (policy *privacyPolicy) showPrivacyAccountsSpend(req *privacytypes.ReqPrivBal4AddrToken) (*privacytypes.UTXOHaveTxHashs, error) {
	if ok, err := policy.isRescanUtxosFlagScaning(); ok {
		return nil, err
	}
	utxoHaveTxHashs, err := policy.store.listSpendUTXOs(req.GetAssetExec(), req.GetToken(), req.GetAddr())
	if err != nil {
		return nil, err
	}
	return utxoHaveTxHashs, nil
}

func (policy *privacyPolicy) signatureTx(tx *types.Transaction, privacyInput *privacytypes.PrivacyInput, utxosInKeyInput []*privacytypes.UTXOBasics, realkeyInputSlice []*privacytypes.RealKeyInput) (err error) {
	tx.Signature = nil
	data := types.Encode(tx)
	ringSign := &types.RingSignature{}
	ringSign.Items = make([]*types.RingSignatureItem, len(privacyInput.Keyinput))
	for i, input := range privacyInput.Keyinput {
		utxos := utxosInKeyInput[i]
		h := common.BytesToHash(data)
		item, err := privacy.GenerateRingSignature(h.Bytes(),
			utxos.Utxos,
			realkeyInputSlice[i].Onetimeprivkey,
			int(realkeyInputSlice[i].Realinputkey),
			input.KeyImage)
		if err != nil {
			return err
		}
		ringSign.Items[i] = item
	}
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	ringSignData := types.Encode(ringSign)
	tx.Signature = &types.Signature{
		Ty:        privacytypes.RingBaseonED25519,
		Signature: ringSignData,
		//             ，
		Pubkey: address.ExecPubKey(cfg.ExecName(privacytypes.PrivacyX)),
	}
	return nil
}

func (policy *privacyPolicy) buildAndStoreWalletTxDetail(param *buildStoreWalletTxDetailParam) {
	blockheight := param.block.Block.Height*maxTxNumPerBlock + int64(param.index)
	heightstr := fmt.Sprintf("%018d", blockheight)
	bizlog.Debug("buildAndStoreWalletTxDetail", "heightstr", heightstr, "addDelType", param.addDelType)
	if AddTx == param.addDelType {
		var txdetail types.WalletTxDetail
		key := calcTxKey(heightstr)
		txdetail.Tx = param.tx
		txdetail.Height = param.block.Block.Height
		txdetail.Index = int64(param.index)
		txdetail.Receipt = param.block.Receipts[param.index]
		txdetail.Blocktime = param.block.Block.BlockTime

		txdetail.ActionName = txdetail.Tx.ActionName()
		txdetail.Amount, _ = param.tx.Amount()
		txdetail.Fromaddr = param.senderRecver
		//txdetail.Spendrecv = param.utxos

		txdetailbyte, err := proto.Marshal(&txdetail)
		if err != nil {
			bizlog.Error("buildAndStoreWalletTxDetail err", "Height", param.block.Block.Height, "index", param.index)
			return
		}

		param.newbatch.Set(key, txdetailbyte)
		if param.isprivacy {
			//
			if sendTx == param.sendRecvFlag {
				param.newbatch.Set(calcSendPrivacyTxKey(param.assetExec, param.tokenname, param.senderRecver, heightstr), key)
			} else if recvTx == param.sendRecvFlag {
				param.newbatch.Set(calcRecvPrivacyTxKey(param.assetExec, param.tokenname, param.senderRecver, heightstr), key)
			}
		}
	} else {
		param.newbatch.Delete(calcTxKey(heightstr))
		if param.isprivacy {
			if sendTx == param.sendRecvFlag {
				param.newbatch.Delete(calcSendPrivacyTxKey(param.assetExec, param.tokenname, param.senderRecver, heightstr))
			} else if recvTx == param.sendRecvFlag {
				param.newbatch.Delete(calcRecvPrivacyTxKey(param.assetExec, param.tokenname, param.senderRecver, heightstr))
			}
		}
	}
}

func (policy *privacyPolicy) checkExpireFTXOOnTimer() {
	operater := policy.getWalletOperate()

	header := operater.GetLastHeader()
	if header == nil {
		bizlog.Error("checkExpireFTXOOnTimer Can not get last header.")
		return
	}
	policy.store.moveFTXO2UTXOWhenFTXOExpire(header.Height, header.BlockTime)
}

func (policy *privacyPolicy) checkWalletStoreData() {
	operater := policy.getWalletOperate()
	defer operater.GetWaitGroup().Done()
	timecount := 10
	checkTicker := time.NewTicker(time.Duration(timecount) * time.Second)
	for {
		select {
		case <-checkTicker.C:
			policy.checkExpireFTXOOnTimer()

			//newbatch := wallet.walletStore.NewBatch(true)
			//err := wallet.procInvalidTxOnTimer(newbatch)
			//if err != nil && err != dbm.ErrNotFoundInDb {
			//	walletlog.Error("checkWalletStoreData", "procInvalidTxOnTimer error ", err)
			//	return
			//}
			//newbatch.Write()
		case <-operater.GetWalletDone():
			return
		}
	}
}

func (policy *privacyPolicy) addDelPrivacyTxsFromBlock(tx *types.Transaction, index int32, block *types.BlockDetail, newbatch db.Batch, addDelType int32) {
	txhash := tx.Hash()
	txhashstr := hex.EncodeToString(txhash)
	_, err := tx.Amount()
	if err != nil {
		bizlog.Error("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "addDelType", addDelType, "index", index, "tx.Amount() error", err)
		return
	}

	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	txExecRes := block.Receipts[index].Ty
	var privateAction privacytypes.PrivacyAction
	if err := types.Decode(tx.GetPayload(), &privateAction); err != nil {
		bizlog.Error("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "addDelType", addDelType, "index", index, "Decode tx.GetPayload() error", err)
		return
	}
	bizlog.Info("addDelPrivacyTxsFromBlock", "Enter addDelPrivacyTxsFromBlock txhash", txhashstr, "index", index, "addDelType", addDelType)

	privacyOutput := privateAction.GetOutput()
	if privacyOutput == nil {
		bizlog.Error("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "addDelType", addDelType, "index", index, "privacyOutput is", privacyOutput)
		return
	}
	assetExec, tokenname := privateAction.GetAssetExecSymbol()
	if assetExec == "" {
		assetExec = "coins"
	}
	RpubKey := privacyOutput.GetRpubKeytx()

	totalUtxosLeft := len(privacyOutput.Keyoutput)
	//  output
	if privacyInfo, err := policy.getPrivacyKeyPairs(); err == nil {
		matchedCount := 0
		utxoProcessed := make([]bool, len(privacyOutput.Keyoutput))
		for _, info := range privacyInfo {
			privacykeyParirs := info.PrivacyKeyPair
			matched4addr := false
			var utxos []*privacytypes.UTXO
			for indexoutput, output := range privacyOutput.Keyoutput {
				if utxoProcessed[indexoutput] {
					continue
				}
				priv, err := privacy.RecoverOnetimePriKey(RpubKey, privacykeyParirs.ViewPrivKey, privacykeyParirs.SpendPrivKey, int64(indexoutput))
				if err == nil {
					recoverPub := priv.PubKey().Bytes()[:]
					if bytes.Equal(recoverPub, output.Onetimepubkey) {
						//                  ，
						//               ，
						//1.     ，      ，            ，
						//2.                ，       change，   2
						matched4addr = true
						totalUtxosLeft--
						utxoProcessed[indexoutput] = true
						//                UTXO
						if types.ExecOk == txExecRes {
							if AddTx == addDelType {
								info2store := &privacytypes.PrivacyDBStore{
									AssetExec:        assetExec,
									Txhash:           txhash,
									Tokenname:        tokenname,
									Amount:           output.Amount,
									OutIndex:         int32(indexoutput),
									TxPublicKeyR:     RpubKey,
									OnetimePublicKey: output.Onetimepubkey,
									Owner:            *info.Addr,
									Height:           block.Block.Height,
									Txindex:          index,
									Blockhash:        block.Block.Hash(cfg),
								}

								utxoGlobalIndex := &privacytypes.UTXOGlobalIndex{
									Outindex: int32(indexoutput),
									Txhash:   txhash,
								}

								utxoCreated := &privacytypes.UTXO{
									Amount: output.Amount,
									UtxoBasic: &privacytypes.UTXOBasic{
										UtxoGlobalIndex: utxoGlobalIndex,
										OnetimePubkey:   output.Onetimepubkey,
									},
								}

								utxos = append(utxos, utxoCreated)
								policy.store.setUTXO(info.Addr, &txhashstr, indexoutput, info2store, newbatch)
								bizlog.Info("addDelPrivacyTxsFromBlock", "add tx txhash", txhashstr, "setUTXO addr ", *info.Addr, "indexoutput", indexoutput)
							} else {
								policy.store.unsetUTXO(assetExec, tokenname, *info.Addr, txhashstr, indexoutput, newbatch)
								bizlog.Info("addDelPrivacyTxsFromBlock", "delete tx txhash", txhashstr, "unsetUTXO addr ", *info.Addr, "indexoutput", indexoutput)
							}
						} else {
							//         ，
							bizlog.Error("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "txExecRes", txExecRes)
							break
						}
					}
				} else {
					bizlog.Error("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "RecoverOnetimePriKey error", err)
				}
			}
			if matched4addr {
				matchedCount++
				//      2 ，
				bizlog.Debug("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "address", *info.Addr, "totalUtxosLeft", totalUtxosLeft, "matchedCount", matchedCount)
				param := &buildStoreWalletTxDetailParam{
					assetExec:    assetExec,
					tokenname:    tokenname,
					block:        block,
					tx:           tx,
					index:        int(index),
					newbatch:     newbatch,
					senderRecver: *info.Addr,
					isprivacy:    true,
					addDelType:   addDelType,
					sendRecvFlag: recvTx,
					utxos:        utxos,
				}
				policy.buildAndStoreWalletTxDetail(param)
				if 2 == matchedCount || 0 == totalUtxosLeft || types.ExecOk != txExecRes {
					bizlog.Info("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "Get matched privacy transfer for address address", *info.Addr, "totalUtxosLeft", totalUtxosLeft, "matchedCount", matchedCount)
					break
				}
			}
		}
	}

	//  input,          ，     output
	//                    ，       utxo     TODO:              input(    keyimage)
	if AddTx == addDelType {
		ftxos, keys := policy.store.getFTXOlist()
		for i, ftxo := range ftxos {
			//
			if ftxo.Txhash != txhashstr {
				continue
			}
			if types.ExecOk == txExecRes && privacytypes.ActionPublic2Privacy != privateAction.Ty {
				bizlog.Info("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "addDelType", addDelType, "moveFTXO2STXO, key", string(keys[i]), "txExecRes", txExecRes)
				policy.store.moveFTXO2STXO(keys[i], txhashstr, newbatch)
			} else if types.ExecOk != txExecRes && privacytypes.ActionPublic2Privacy != privateAction.Ty {
				//
				bizlog.Info("PrivacyTrading AddDelPrivacyTxsFromBlock", "txhash", txhashstr, "addDelType", addDelType, "moveFTXO2UTXO, key", string(keys[i]), "txExecRes", txExecRes)
				policy.store.moveFTXO2UTXO(keys[i], newbatch)
			}
			//         ，
			param := &buildStoreWalletTxDetailParam{
				assetExec:    assetExec,
				tokenname:    tokenname,
				block:        block,
				tx:           tx,
				index:        int(index),
				newbatch:     newbatch,
				senderRecver: ftxo.Sender,
				isprivacy:    true,
				addDelType:   addDelType,
				sendRecvFlag: sendTx,
				utxos:        nil,
			}
			policy.buildAndStoreWalletTxDetail(param)
		}
	} else {
		//        ，    STXO        ，      FTXO，
		stxosInOneTx, _, _ := policy.store.getWalletFtxoStxo(STXOs4Tx)
		for _, ftxo := range stxosInOneTx {
			if ftxo.Txhash == txhashstr {
				param := &buildStoreWalletTxDetailParam{
					assetExec:    assetExec,
					tokenname:    tokenname,
					block:        block,
					tx:           tx,
					index:        int(index),
					newbatch:     newbatch,
					senderRecver: "",
					isprivacy:    true,
					addDelType:   addDelType,
					sendRecvFlag: sendTx,
					utxos:        nil,
				}

				if types.ExecOk == txExecRes && privacytypes.ActionPublic2Privacy != privateAction.Ty {
					bizlog.Info("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "addDelType", addDelType, "moveSTXO2FTXO txExecRes", txExecRes)
					policy.store.moveSTXO2FTXO(tx, txhashstr, newbatch)
					policy.buildAndStoreWalletTxDetail(param)
				} else if types.ExecOk != txExecRes && privacytypes.ActionPublic2Privacy != privateAction.Ty {
					bizlog.Info("addDelPrivacyTxsFromBlock", "txhash", txhashstr, "addDelType", addDelType)
					policy.buildAndStoreWalletTxDetail(param)
				}
			}
		}
	}
}
