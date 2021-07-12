// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	dty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/dposvote/types"
)

//Exec_Regist DPos
func (d *DPos) Exec_Regist(payload *dty.DposCandidatorRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.Regist(payload)
}

//Exec_CancelRegist DPos
func (d *DPos) Exec_CancelRegist(payload *dty.DposCandidatorCancelRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.CancelRegist(payload)
}

//Exec_ReRegist DPos
func (d *DPos) Exec_ReRegist(payload *dty.DposCandidatorRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.ReRegist(payload)
}

//Exec_Vote DPos
func (d *DPos) Exec_Vote(payload *dty.DposVote, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.Vote(payload)
}

//Exec_CancelVote DPos
func (d *DPos) Exec_CancelVote(payload *dty.DposCancelVote, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.CancelVote(payload)
}

//Exec_RegistVrfM DPos            Vrf M
func (d *DPos) Exec_RegistVrfM(payload *dty.DposVrfMRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.RegistVrfM(payload)
}

//Exec_RegistVrfRP DPos            Vrf R/P
func (d *DPos) Exec_RegistVrfRP(payload *dty.DposVrfRPRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.RegistVrfRP(payload)
}

//Exec_RecordCB DPos     CycleBoundary
func (d *DPos) Exec_RecordCB(payload *dty.DposCBInfo, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.RecordCB(payload)
}

//Exec_RegistTopN DPos       cycle  TOPN
func (d *DPos) Exec_RegistTopN(payload *dty.TopNCandidatorRegist, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(d, tx, index)
	return action.RegistTopN(payload)
}
