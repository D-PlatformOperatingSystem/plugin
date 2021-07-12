// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var name string

func init() {
	name = AutonomyX
	types.AllowUserExec = append(types.AllowUserExec, []byte(name))
	types.RegFork(name, InitFork)
	types.RegExec(name, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(AutonomyX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(AutonomyX, NewType(cfg))
}

// NewType         
func NewType(cfg *types.DplatformOSConfig) *AutonomyType {
	c := &AutonomyType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// AutonomyType        
type AutonomyType struct {
	types.ExecTypeBase
}

// GetName        
func (a *AutonomyType) GetName() string {
	return AutonomyX
}

// GetLogMap         
func (a *AutonomyType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogPropBoard:      {Ty: reflect.TypeOf(ReceiptProposalBoard{}), Name: "LogPropBoard"},
		TyLogRvkPropBoard:   {Ty: reflect.TypeOf(ReceiptProposalBoard{}), Name: "LogRvkPropBoard"},
		TyLogVotePropBoard:  {Ty: reflect.TypeOf(ReceiptProposalBoard{}), Name: "LogVotePropBoard"},
		TyLogTmintPropBoard: {Ty: reflect.TypeOf(ReceiptProposalBoard{}), Name: "LogTmintPropBoard"},

		TyLogPropProject:        {Ty: reflect.TypeOf(ReceiptProposalProject{}), Name: "LogPropProject"},
		TyLogRvkPropProject:     {Ty: reflect.TypeOf(ReceiptProposalProject{}), Name: "LogRvkPropProject"},
		TyLogVotePropProject:    {Ty: reflect.TypeOf(ReceiptProposalProject{}), Name: "LogVotePropProject"},
		TyLogPubVotePropProject: {Ty: reflect.TypeOf(ReceiptProposalProject{}), Name: "LogPubVotePropProject"},
		TyLogTmintPropProject:   {Ty: reflect.TypeOf(ReceiptProposalProject{}), Name: "LogTmintPropProject"},

		TyLogPropRule:      {Ty: reflect.TypeOf(ReceiptProposalRule{}), Name: "LogPropRule"},
		TyLogRvkPropRule:   {Ty: reflect.TypeOf(ReceiptProposalRule{}), Name: "LogRvkPropRule"},
		TyLogVotePropRule:  {Ty: reflect.TypeOf(ReceiptProposalRule{}), Name: "LogVotePropRule"},
		TyLogTmintPropRule: {Ty: reflect.TypeOf(ReceiptProposalRule{}), Name: "LogTmintPropRule"},

		TyLogCommentProp: {Ty: reflect.TypeOf(ReceiptProposalComment{}), Name: "LogCommentProp"},

		TyLogPropChange:      {Ty: reflect.TypeOf(ReceiptProposalChange{}), Name: "LogPropChange"},
		TyLogRvkPropChange:   {Ty: reflect.TypeOf(ReceiptProposalChange{}), Name: "LogRvkPropChange"},
		TyLogVotePropChange:  {Ty: reflect.TypeOf(ReceiptProposalChange{}), Name: "LogVotePropChange"},
		TyLogTmintPropChange: {Ty: reflect.TypeOf(ReceiptProposalChange{}), Name: "LogTmintPropChange"},
	}
}

// GetPayload     Unfreeze   Payload
func (a *AutonomyType) GetPayload() types.Message {
	return &AutonomyAction{}
}

// GetTypeMap   Action     
func (a *AutonomyType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"PropBoard":      AutonomyActionPropBoard,
		"RvkPropBoard":   AutonomyActionRvkPropBoard,
		"VotePropBoard":  AutonomyActionVotePropBoard,
		"TmintPropBoard": AutonomyActionTmintPropBoard,

		"PropProject":        AutonomyActionPropProject,
		"RvkPropProject":     AutonomyActionRvkPropProject,
		"VotePropProject":    AutonomyActionVotePropProject,
		"PubVotePropProject": AutonomyActionPubVotePropProject,
		"TmintPropProject":   AutonomyActionTmintPropProject,

		"PropRule":      AutonomyActionPropRule,
		"RvkPropRule":   AutonomyActionRvkPropRule,
		"VotePropRule":  AutonomyActionVotePropRule,
		"TmintPropRule": AutonomyActionTmintPropRule,

		"Transfer":    AutonomyActionTransfer,
		"CommentProp": AutonomyActionCommentProp,

		"PropChange":      AutonomyActionPropChange,
		"RvkPropChange":   AutonomyActionRvkPropChange,
		"VotePropChange":  AutonomyActionVotePropChange,
		"TmintPropChange": AutonomyActionTmintPropChange,
	}
}
