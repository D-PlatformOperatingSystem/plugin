// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	auty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/autonomy/types"
)

// Query_GetProposalBoard        
func (a *Autonomy) Query_GetProposalBoard(in *types.ReqString) (types.Message, error) {
	return a.getProposalBoard(in)
}

// Query_ListProposalBoard     
func (a *Autonomy) Query_ListProposalBoard(in *auty.ReqQueryProposalBoard) (types.Message, error) {
	return a.listProposalBoard(in)
}

// Query_GetActiveBoard     board
func (a *Autonomy) Query_GetActiveBoard(in *types.ReqString) (types.Message, error) {
	return a.getActiveBoard()
}

// Query_GetProposalProject       
func (a *Autonomy) Query_GetProposalProject(in *types.ReqString) (types.Message, error) {
	return a.getProposalProject(in)
}

// Query_ListProposalProject     
func (a *Autonomy) Query_ListProposalProject(in *auty.ReqQueryProposalProject) (types.Message, error) {
	return a.listProposalProject(in)
}

// Query_GetProposalRule       
func (a *Autonomy) Query_GetProposalRule(in *types.ReqString) (types.Message, error) {
	return a.getProposalRule(in)
}

// Query_ListProposalRule     
func (a *Autonomy) Query_ListProposalRule(in *auty.ReqQueryProposalRule) (types.Message, error) {
	return a.listProposalRule(in)
}

// Query_GetActiveRule     rule
func (a *Autonomy) Query_GetActiveRule(in *types.ReqString) (types.Message, error) {
	return a.getActiveRule()
}

// Query_ListProposalComment         
func (a *Autonomy) Query_ListProposalComment(in *auty.ReqQueryProposalComment) (types.Message, error) {
	return a.listProposalComment(in)
}

// Query_GetProposalChange            
func (a *Autonomy) Query_GetProposalChange(in *types.ReqString) (types.Message, error) {
	return a.getProposalChange(in)
}

// Query_ListProposalChange     
func (a *Autonomy) Query_ListProposalChange(in *auty.ReqQueryProposalChange) (types.Message, error) {
	return a.listProposalChange(in)
}
