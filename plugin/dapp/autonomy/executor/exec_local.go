// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	auty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/autonomy/types"
)

//

// ExecLocal_PropBoard
func (a *Autonomy) ExecLocal_PropBoard(payload *auty.ProposalBoard, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalBoard(tx, receiptData)
}

// ExecLocal_RvkPropBoard
func (a *Autonomy) ExecLocal_RvkPropBoard(payload *auty.RevokeProposalBoard, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalBoard(tx, receiptData)
}

// ExecLocal_VotePropBoard
func (a *Autonomy) ExecLocal_VotePropBoard(payload *auty.VoteProposalBoard, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalBoard(tx, receiptData)
}

// ExecLocal_TmintPropBoard
func (a *Autonomy) ExecLocal_TmintPropBoard(payload *auty.TerminateProposalBoard, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalBoard(tx, receiptData)
}

//

// ExecLocal_PropProject
func (a *Autonomy) ExecLocal_PropProject(payload *auty.ProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

// ExecLocal_RvkPropProject
func (a *Autonomy) ExecLocal_RvkPropProject(payload *auty.RevokeProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

// ExecLocal_VotePropProject
func (a *Autonomy) ExecLocal_VotePropProject(payload *auty.VoteProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

// ExecLocal_PubVotePropProject
func (a *Autonomy) ExecLocal_PubVotePropProject(payload *auty.PubVoteProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

// ExecLocal_TmintPropProject
func (a *Autonomy) ExecLocal_TmintPropProject(payload *auty.TerminateProposalProject, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalProject(tx, receiptData)
}

//

// ExecLocal_PropRule
func (a *Autonomy) ExecLocal_PropRule(payload *auty.ProposalRule, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalRule(tx, receiptData)
}

// ExecLocal_RvkPropRule
func (a *Autonomy) ExecLocal_RvkPropRule(payload *auty.RevokeProposalRule, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalRule(tx, receiptData)
}

// ExecLocal_VotePropRule
func (a *Autonomy) ExecLocal_VotePropRule(payload *auty.VoteProposalRule, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalRule(tx, receiptData)
}

// ExecLocal_TmintPropRule
func (a *Autonomy) ExecLocal_TmintPropRule(payload *auty.TerminateProposalRule, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalRule(tx, receiptData)
}

// ExecLocal_CommentProp
func (a *Autonomy) ExecLocal_CommentProp(payload *auty.Comment, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalCommentProp(tx, receiptData)
}

//

// ExecLocal_PropChange
func (a *Autonomy) ExecLocal_PropChange(payload *auty.ProposalChange, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalChange(tx, receiptData)
}

// ExecLocal_RvkPropChange
func (a *Autonomy) ExecLocal_RvkPropChange(payload *auty.RevokeProposalChange, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalChange(tx, receiptData)
}

// ExecLocal_VotePropChange
func (a *Autonomy) ExecLocal_VotePropChange(payload *auty.VoteProposalChange, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalChange(tx, receiptData)
}

// ExecLocal_TmintPropChange
func (a *Autonomy) ExecLocal_TmintPropChange(payload *auty.TerminateProposalChange, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return a.execAutoLocalChange(tx, receiptData)
}
