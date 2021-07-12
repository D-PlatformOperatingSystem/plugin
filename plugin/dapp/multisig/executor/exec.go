// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
)

//Exec_MultiSigAccCreate
func (m *MultiSig) Exec_MultiSigAccCreate(payload *mty.MultiSigAccCreate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigAccCreate(payload)
}

//Exec_MultiSigOwnerOperate       owner     ：owner add/del/replace
func (m *MultiSig) Exec_MultiSigOwnerOperate(payload *mty.MultiSigOwnerOperate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigOwnerOperate(payload)
}

//Exec_MultiSigAccOperate            ：weight
func (m *MultiSig) Exec_MultiSigAccOperate(payload *mty.MultiSigAccOperate, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigAccOperate(payload)
}

//Exec_MultiSigConfirmTx
func (m *MultiSig) Exec_MultiSigConfirmTx(payload *mty.MultiSigConfirmTx, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigConfirmTx(payload)
}

//Exec_MultiSigExecTransferTo                 ，Addr --->multiSigAddr
func (m *MultiSig) Exec_MultiSigExecTransferTo(payload *mty.MultiSigExecTransferTo, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigExecTransferTo(payload)
}

//Exec_MultiSigExecTransferFrom                 ，multiSigAddr--->Addr
func (m *MultiSig) Exec_MultiSigExecTransferFrom(payload *mty.MultiSigExecTransferFrom, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(m, tx, int32(index))
	return action.MultiSigExecTransferFrom(payload)
}
