// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// autonomy action ty
const (
	AutonomyActionPropBoard = iota + 1
	AutonomyActionRvkPropBoard
	AutonomyActionVotePropBoard
	AutonomyActionTmintPropBoard

	AutonomyActionPropProject
	AutonomyActionRvkPropProject
	AutonomyActionVotePropProject
	AutonomyActionPubVotePropProject
	AutonomyActionTmintPropProject

	AutonomyActionPropRule
	AutonomyActionRvkPropRule
	AutonomyActionVotePropRule
	AutonomyActionTmintPropRule

	AutonomyActionTransfer
	AutonomyActionCommentProp

	AutonomyActionPropChange
	AutonomyActionRvkPropChange
	AutonomyActionVotePropChange
	AutonomyActionTmintPropChange

	//log for autonomy
	TyLogPropBoard      = 2101
	TyLogRvkPropBoard   = 2102
	TyLogVotePropBoard  = 2103
	TyLogTmintPropBoard = 2104

	TyLogPropProject        = 2111
	TyLogRvkPropProject     = 2112
	TyLogVotePropProject    = 2113
	TyLogPubVotePropProject = 2114
	TyLogTmintPropProject   = 2115

	TyLogPropRule      = 2121
	TyLogRvkPropRule   = 2122
	TyLogVotePropRule  = 2123
	TyLogTmintPropRule = 2124

	TyLogCommentProp = 2131

	TyLogPropChange      = 2141
	TyLogRvkPropChange   = 2142
	TyLogVotePropChange  = 2143
	TyLogTmintPropChange = 2144
)

// Board status
const (
	AutonomyStatusProposalBoard = iota + 1
	AutonomyStatusRvkPropBoard
	AutonomyStatusVotePropBoard
	AutonomyStatusTmintPropBoard
)

// Project status
const (
	AutonomyStatusProposalProject = iota + 1
	AutonomyStatusRvkPropProject
	AutonomyStatusVotePropProject
	AutonomyStatusPubVotePropProject
	AutonomyStatusTmintPropProject
)

// Rule status
const (
	AutonomyStatusProposalRule = iota + 1
	AutonomyStatusRvkPropRule
	AutonomyStatusVotePropRule
	AutonomyStatusTmintPropRule
)

// Change status
const (
	AutonomyStatusProposalChange = iota + 1
	AutonomyStatusRvkPropChange
	AutonomyStatusVotePropChange
	AutonomyStatusTmintPropChange
)

const (
	// GetProposalBoard    cmd
	GetProposalBoard = "GetProposalBoard"
	// ListProposalBoard
	ListProposalBoard = "ListProposalBoard"
	// GetActiveBoard
	GetActiveBoard = "GetActiveBoard"
	// GetProposalProject    cmd
	GetProposalProject = "GetProposalProject"
	// ListProposalProject
	ListProposalProject = "ListProposalProject"
	// GetProposalRule    cmd
	GetProposalRule = "GetProposalRule"
	// ListProposalRule
	ListProposalRule = "ListProposalRule"
	// GetActiveRule
	GetActiveRule = "GetActiveRule"
	// ListProposalComment
	ListProposalComment = "ListProposalComment"
	// GetProposalChange    cmd
	GetProposalChange = "GetProposalChange"
	// ListProposalChange
	ListProposalChange = "ListProposalChange"
)

//
//   github     ，        ,
//      ，
var (
	AutonomyX      = "autonomy"
	ExecerAutonomy = []byte(AutonomyX)
	// TicketX        ticket
	TicketX = "ticket"
)
