// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var multisiglog = log15.New("module", "execs.multisig")

// OwnerAdd :
var (
	OwnerAdd     uint64 = 1
	OwnerDel     uint64 = 2
	OwnerModify  uint64 = 3
	OwnerReplace uint64 = 4
	//AccWeightOp
	AccWeightOp     = true
	AccDailyLimitOp = false
	//OwnerOperate         ：  ，owner  ，account
	OwnerOperate    uint64 = 1
	AccountOperate  uint64 = 2
	TransferOperate uint64 = 3
	//IsSubmit ：
	IsSubmit  = true
	IsConfirm = false

	MultiSigX            = "multisig"
	OneDaySecond   int64 = 24 * 3600
	MinOwnersInit        = 2
	MinOwnersCount       = 1  //                owner
	MaxOwnersCount       = 20 //             20 owner

	Multisiglog = log15.New("module", MultiSigX)
)

// MultiSig    actionid
const (
	ActionMultiSigAccCreate        = 10000
	ActionMultiSigOwnerOperate     = 10001
	ActionMultiSigAccOperate       = 10002
	ActionMultiSigConfirmTx        = 10003
	ActionMultiSigExecTransferTo   = 10004
	ActionMultiSigExecTransferFrom = 10005
)

//           logid
const (
	TyLogMultiSigAccCreate = 10000 //

	TyLogMultiSigOwnerAdd     = 10001 //  add owner：addr weight
	TyLogMultiSigOwnerDel     = 10002 //  del owner：addr weight
	TyLogMultiSigOwnerModify  = 10003 //  modify owner：preweight  currentweight
	TyLogMultiSigOwnerReplace = 10004 //  old owner   ：     owner  ：addr+weight

	TyLogMultiSigAccWeightModify     = 10005 //            ：preReqWeight curReqWeight
	TyLogMultiSigAccDailyLimitAdd    = 10006 //  add DailyLimit：Symbol DailyLimit
	TyLogMultiSigAccDailyLimitModify = 10007 //  modify DailyLimit：preDailyLimit  currentDailyLimit

	TyLogMultiSigConfirmTx       = 10008 //
	TyLogMultiSigConfirmTxRevoke = 10009 //

	TyLogDailyLimitUpdate = 10010 //DailyLimit  ，DailyLimit Submit Confirm
	TyLogMultiSigTx       = 10011 // Submit
	TyLogTxCountUpdate    = 10012 //txcount   Submit

)

//AccAssetsResult     cli   ，   amount
type AccAssetsResult struct {
	Execer   string `json:"execer,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Currency int32  `json:"currency,omitempty"`
	Balance  string `json:"balance,omitempty"`
	Frozen   string `json:"frozen,omitempty"`
	Receiver string `json:"receiver,omitempty"`
	Addr     string `json:"addr,omitempty"`
}

//DailyLimitResult          cli
type DailyLimitResult struct {
	Symbol     string `json:"symbol,omitempty"`
	Execer     string `json:"execer,omitempty"`
	DailyLimit string `json:"dailyLimit,omitempty"`
	SpentToday string `json:"spent,omitempty"`
	LastDay    string `json:"lastday,omitempty"`
}

//MultiSigResult            cli
type MultiSigResult struct {
	CreateAddr     string              `json:"createAddr,omitempty"`
	MultiSigAddr   string              `json:"multiSigAddr,omitempty"`
	Owners         []*Owner            `json:"owners,omitempty"`
	DailyLimits    []*DailyLimitResult `json:"dailyLimits,omitempty"`
	TxCount        uint64              `json:"txCount,omitempty"`
	RequiredWeight uint64              `json:"requiredWeight,omitempty"`
}

//UnSpentAssetsResult               cli
type UnSpentAssetsResult struct {
	Symbol  string `json:"symbol,omitempty"`
	Execer  string `json:"execer,omitempty"`
	UnSpent string `json:"unspent,omitempty"`
}

//IsAssetsInvalid         ，Symbol：      ，  ：DOM,coins.DOM。exec：   types.AllowUserExec
func IsAssetsInvalid(exec, symbol string) error {

	//exec
	allowExeName := types.AllowUserExec
	nameLen := len(allowExeName)
	execValid := false
	for i := 0; i < nameLen; i++ {
		if exec == string(allowExeName[i]) {
			execValid = true
			break
		}
	}
	if !execValid {
		multisiglog.Error("IsAssetsInvalid", "exec", exec)
		return ErrInvalidExec
	}
	//Symbol
	return nil
}
