// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
)

//ExecLocal_MultiSigAccCreate         ,  payload receiptData            db
func (m *MultiSig) ExecLocal_MultiSigAccCreate(payload *mty.MultiSigAccCreate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, true)
	if err != nil {
		multisiglog.Error("ExecLocal_MultiSigAccCreate", "err", err)
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

//ExecLocal_MultiSigOwnerOperate       owner     ：owner add/del/replace/modify
func (m *MultiSig) ExecLocal_MultiSigOwnerOperate(payload *mty.MultiSigOwnerOperate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, true)
	if err != nil {
		multisiglog.Error("ExecLocal_MultiSigOwnerOperate", "err", err)
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

//ExecLocal_MultiSigAccOperate            ：weight
func (m *MultiSig) ExecLocal_MultiSigAccOperate(payload *mty.MultiSigAccOperate, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, true)
	if err != nil {
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

//ExecLocal_MultiSigConfirmTx
func (m *MultiSig) ExecLocal_MultiSigConfirmTx(payload *mty.MultiSigConfirmTx, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, true)
	if err != nil {
		multisiglog.Error("ExecLocal_MultiSigConfirmTx", "err", err)
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

//ExecLocal_MultiSigExecTransferTo                 ，Addr --->multiSigAddr
func (m *MultiSig) ExecLocal_MultiSigExecTransferTo(payload *mty.MultiSigExecTransferTo, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.saveMultiSigTransfer(tx, mty.IsSubmit, true)
	if err != nil {
		multisiglog.Error("ExecLocal_MultiSigExecTransferTo", "err", err)
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}

//ExecLocal_MultiSigExecTransferFrom                 ，multiSigAddr--->Addr
func (m *MultiSig) ExecLocal_MultiSigExecTransferFrom(payload *mty.MultiSigExecTransferFrom, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	if receiptData.GetTy() != types.ExecOk {
		return &types.LocalDBSet{}, nil
	}

	kv, err := m.execLocalMultiSigReceipt(receiptData, tx, true)
	if err != nil {
		multisiglog.Error("ExecLocal_MultiSigExecTransferFrom", "err", err)
		return nil, err
	}
	return &types.LocalDBSet{KV: kv}, nil
}
