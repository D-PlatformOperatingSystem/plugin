// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
)

//Query_MultiSigAccCount            ，
//  ReplyMultiSigAccounts
func (m *MultiSig) Query_MultiSigAccCount(in *types.ReqNil) (types.Message, error) {
	db := m.GetLocalDB()
	count, err := getMultiSigAccCount(db)
	if err != nil {
		return nil, err
	}

	return &types.Int64{Data: count}, nil
}

//Query_MultiSigAccounts
//  ：
//message ReqMultiSigAccs {
//	int64	start	= 1;
//	int64	end		= 2;
//  ：
//message ReplyMultiSigAccs {
//    repeated string address = 1;
func (m *MultiSig) Query_MultiSigAccounts(in *mty.ReqMultiSigAccs) (types.Message, error) {
	accountAddrs := &mty.ReplyMultiSigAccs{}

	if in.Start > in.End || in.Start < 0 {
		return nil, types.ErrInvalidParam
	}

	db := m.GetLocalDB()
	totalcount, err := getMultiSigAccCount(db)
	if err != nil {
		return nil, err
	}
	if totalcount == 0 {
		return accountAddrs, nil
	}
	if in.End >= totalcount {
		return nil, types.ErrInvalidParam
	}
	for index := in.Start; index <= in.End; index++ {
		addr, err := getMultiSigAccList(db, index)
		if err == nil {
			accountAddrs.Address = append(accountAddrs.Address, addr)
		}
	}
	return accountAddrs, nil
}

//Query_MultiSigAccountInfo
//  ：
//message ReqMultiSigAccountInfo {
//	string MultiSigAccAddr = 1;
//  ：
//message MultiSig {
//    string 							createAddr        	= 1;
//    string 							multiSigAddr      	= 2;
//    repeated Owner           			owners				= 3;
//    repeated DailyLimit          		dailyLimits   		= 4;
//    uint64           					txCount				= 5;
//	  uint64           					requiredWeight		= 6;
func (m *MultiSig) Query_MultiSigAccountInfo(in *mty.ReqMultiSigAccInfo) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	db := m.GetLocalDB()
	addr := in.MultiSigAccAddr

	if err := address.CheckMultiSignAddress(addr); err != nil {
		return nil, types.ErrInvalidAddress
	}
	multiSigAcc, err := getMultiSigAccount(db, addr)
	if err != nil {
		return nil, err
	}
	if multiSigAcc == nil {
		multiSigAcc = &mty.MultiSig{}
	}
	return multiSigAcc, nil
}

//Query_MultiSigAccTxCount             tx
//  ：
//message ReqMultiSigAccountInfo {
//	string MultiSigAccAddr = 1;
//  ：
//uint64
func (m *MultiSig) Query_MultiSigAccTxCount(in *mty.ReqMultiSigAccInfo) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	db := m.GetLocalDB()
	addr := in.MultiSigAccAddr

	if err := address.CheckMultiSignAddress(addr); err != nil {
		return nil, types.ErrInvalidAddress
	}

	multiSigAcc, err := getMultiSigAccount(db, addr)
	if err != nil {
		return nil, err
	}
	if multiSigAcc == nil {
		return nil, mty.ErrAccountHasExist
	}
	return &mty.Uint64{Data: multiSigAcc.TxCount}, nil
}

//Query_MultiSigTxids   txids            ，pending, executed
//  ：
//message ReqMultiSigTxids {
//  string multisigaddr = 1;
//	uint64 fromtxid = 2;
//	uint64 totxid = 3;
//	bool   pending = 4;
//	bool   executed	= 5;
//   :
//message ReplyMultiSigTxids {
//  string 			multisigaddr = 1;
//	repeated uint64	txids		 = 2;
func (m *MultiSig) Query_MultiSigTxids(in *mty.ReqMultiSigTxids) (types.Message, error) {
	if in == nil || in.FromTxId > in.ToTxId || in.FromTxId < 0 {
		return nil, types.ErrInvalidParam
	}

	db := m.GetLocalDB()
	addr := in.MultiSigAddr

	if err := address.CheckMultiSignAddress(addr); err != nil {
		return nil, types.ErrInvalidAddress
	}
	multiSigAcc, err := getMultiSigAccount(db, addr)
	if err != nil {
		return nil, err
	}
	if multiSigAcc == nil || multiSigAcc.TxCount <= in.ToTxId {
		return nil, types.ErrInvalidParam
	}

	multiSigTxids := &mty.ReplyMultiSigTxids{}
	multiSigTxids.MultiSigAddr = addr
	for txid := in.FromTxId; txid <= in.ToTxId; txid++ {
		multiSigTx, err := getMultiSigTx(db, addr, txid)
		if err != nil || multiSigTx == nil {
			multisiglog.Error("Query_MultiSigTxids:getMultiSigTx", "addr", addr, "txid", txid, "err", err)
			continue
		}
		findTxid := txid
		//  Pending/Executed   txid
		if in.Pending && !multiSigTx.Executed || in.Executed && multiSigTx.Executed {
			multiSigTxids.Txids = append(multiSigTxids.Txids, findTxid)
		}
	}
	return multiSigTxids, nil

}

//Query_MultiSigTxInfo   txid     ，       owner
//  :
//message ReqMultiSigTxInfo {
//  string multisigaddr = 1;
//	uint64 txid = 2;
//  :
//message ReplyMultiSigTxInfo {
//    MultiSigTransaction multisigtxinfo = 1;
//    repeated Owner confirmowners = 3;
func (m *MultiSig) Query_MultiSigTxInfo(in *mty.ReqMultiSigTxInfo) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	db := m.GetLocalDB()
	addr := in.MultiSigAddr
	txid := in.TxId

	if err := address.CheckMultiSignAddress(addr); err != nil {
		return nil, types.ErrInvalidAddress
	}

	multiSigTx, err := getMultiSigTx(db, addr, txid)
	if err != nil {
		return nil, err
	}
	if multiSigTx == nil {
		multiSigTx = &mty.MultiSigTx{}
	} else { //       hex.EncodeToString()     ，   0x，                 0x
		multiSigTx.TxHash = "0x" + multiSigTx.TxHash
	}
	return multiSigTx, nil
}

//Query_MultiSigTxConfirmedWeight   txid
//  :
//message ReqMultiSigTxInfo {
//  string multisigaddr = 1;
//	uint64 txid = 2;
//  :
//message Int64
func (m *MultiSig) Query_MultiSigTxConfirmedWeight(in *mty.ReqMultiSigTxInfo) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	db := m.GetLocalDB()
	addr := in.MultiSigAddr
	txid := in.TxId

	if err := address.CheckMultiSignAddress(addr); err != nil {
		return nil, types.ErrInvalidAddress
	}

	multiSigTx, err := getMultiSigTx(db, addr, txid)
	if err != nil {
		return nil, err
	}
	if multiSigTx == nil {
		return nil, mty.ErrTxidNotExist
	}
	var totalWeight uint64
	for _, owner := range multiSigTx.ConfirmedOwner {
		totalWeight += owner.Weight
	}

	return &mty.Uint64{Data: totalWeight}, nil
}

//Query_MultiSigAccUnSpentToday
//  :
//message ReqMultiSigAccUnSpentToday {
//	string multiSigAddr = 1;
//	string execer 		= 2;
//	string symbol 		= 3;
//  :
//message ReplyMultiSigAccUnSpentToday {
//	uint64 	amount = 1;
func (m *MultiSig) Query_MultiSigAccUnSpentToday(in *mty.ReqAccAssets) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	db := m.GetLocalDB()
	addr := in.MultiSigAddr
	isAll := in.IsAll

	if err := address.CheckMultiSignAddress(addr); err != nil {
		return nil, types.ErrInvalidAddress
	}
	multiSigAcc, err := getMultiSigAccount(db, addr)
	if err != nil {
		return nil, err
	}

	replyUnSpentAssets := &mty.ReplyUnSpentAssets{}
	if multiSigAcc == nil {
		return replyUnSpentAssets, nil
	}
	if isAll {
		for _, dailyLimit := range multiSigAcc.DailyLimits {
			var unSpentAssets mty.UnSpentAssets
			assets := &mty.Assets{
				Execer: dailyLimit.Execer,
				Symbol: dailyLimit.Symbol,
			}
			unSpentAssets.Assets = assets
			unSpentAssets.Amount = 0
			if dailyLimit.DailyLimit > dailyLimit.SpentToday {
				unSpentAssets.Amount = dailyLimit.DailyLimit - dailyLimit.SpentToday
			}
			replyUnSpentAssets.UnSpentAssets = append(replyUnSpentAssets.UnSpentAssets, &unSpentAssets)
		}
	} else {
		//assets
		err := mty.IsAssetsInvalid(in.Assets.Execer, in.Assets.Symbol)
		if err != nil {
			return nil, err
		}

		for _, dailyLimit := range multiSigAcc.DailyLimits {
			var unSpentAssets mty.UnSpentAssets

			if dailyLimit.Execer == in.Assets.Execer && dailyLimit.Symbol == in.Assets.Symbol {
				assets := &mty.Assets{
					Execer: dailyLimit.Execer,
					Symbol: dailyLimit.Symbol,
				}
				unSpentAssets.Assets = assets
				unSpentAssets.Amount = 0
				if dailyLimit.DailyLimit > dailyLimit.SpentToday {
					unSpentAssets.Amount = dailyLimit.DailyLimit - dailyLimit.SpentToday
				}
				replyUnSpentAssets.UnSpentAssets = append(replyUnSpentAssets.UnSpentAssets, &unSpentAssets)
				break
			}
		}
	}
	return replyUnSpentAssets, nil
}

//Query_MultiSigAccAssets                ，
//  :
//message ReqAccAssets {
//	string multiSigAddr = 1;
//	Assets assets 		= 2;
//	bool   isAll 		= 3;
//  :
//message MultiSigAccAssets {
//	Assets 		assets 		= 1;
//	int64   	recvAmount 	= 2;
//   Account 	account 	= 3;
func (m *MultiSig) Query_MultiSigAccAssets(in *mty.ReqAccAssets) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	//
	if err := address.CheckMultiSignAddress(in.MultiSigAddr); err != nil {
		if err = address.CheckAddress(in.MultiSigAddr); err != nil {
			return nil, types.ErrInvalidAddress
		}
	}

	replyAccAssets := &mty.ReplyAccAssets{}
	//
	if in.IsAll {
		values, err := getMultiSigAccAllAssets(m.GetLocalDB(), in.MultiSigAddr)
		if err != nil {
			return nil, err
		}
		if len(values) != 0 {
			for _, value := range values {
				reciver := mty.AccountAssets{}
				err = types.Decode(value, &reciver)
				if err != nil {
					continue
				}
				accAssets := &mty.AccAssets{}
				account, err := m.getMultiSigAccAssets(reciver.MultiSigAddr, reciver.Assets)
				if err != nil {
					multisiglog.Error("Query_MultiSigAccAssets:getMultiSigAccAssets", "MultiSigAddr", reciver.MultiSigAddr, "err", err)
				}
				accAssets.Account = account
				accAssets.Assets = reciver.Assets
				accAssets.RecvAmount = reciver.Amount

				replyAccAssets.AccAssets = append(replyAccAssets.AccAssets, accAssets)
			}
		}
	} else { //
		accAssets := &mty.AccAssets{}
		//assets
		err := mty.IsAssetsInvalid(in.Assets.Execer, in.Assets.Symbol)
		if err != nil {
			return nil, err
		}
		account, err := m.getMultiSigAccAssets(in.MultiSigAddr, in.Assets)
		if err != nil {
			multisiglog.Error("Query_MultiSigAccAssets:getMultiSigAccAssets", "MultiSigAddr", in.MultiSigAddr, "err", err)
		}
		accAssets.Account = account
		accAssets.Assets = in.Assets

		amount, err := getAddrReciver(m.GetLocalDB(), in.MultiSigAddr, in.Assets.Execer, in.Assets.Symbol)
		if err != nil {
			multisiglog.Error("Query_MultiSigAccAssets:getAddrReciver", "MultiSigAddr", in.MultiSigAddr, "err", err)
		}
		accAssets.RecvAmount = amount

		replyAccAssets.AccAssets = append(replyAccAssets.AccAssets, accAssets)
	}

	return replyAccAssets, nil
}

//Query_MultiSigAccAllAddress
//  :
//createaddr
//  :
//[]string
func (m *MultiSig) Query_MultiSigAccAllAddress(in *mty.ReqMultiSigAccInfo) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	if err := address.CheckAddress(in.MultiSigAccAddr); err != nil {
		return nil, types.ErrInvalidAddress
	}
	return getMultiSigAccAllAddress(m.GetLocalDB(), in.MultiSigAccAddr)
}
