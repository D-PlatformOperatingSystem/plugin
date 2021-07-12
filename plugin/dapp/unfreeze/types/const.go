// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//unfreeze action ty
const (
	UnfreezeActionCreate = iota + 1
	UnfreezeActionWithdraw
	UnfreezeActionTerminate

	//log for unfreeze
	TyLogCreateUnfreeze    = 2001 // TODO
	TyLogWithdrawUnfreeze  = 2002
	TyLogTerminateUnfreeze = 2003
)

const (
	// Action_CreateUnfreeze Action
	Action_CreateUnfreeze = "createUnfreeze"
	// Action_WithdrawUnfreeze Action
	Action_WithdrawUnfreeze = "withdrawUnfreeze"
	// Action_TerminateUnfreeze Action
	Action_TerminateUnfreeze = "terminateUnfreeze"
)

const (
	// FuncName_QueryUnfreezeWithdraw
	FuncName_QueryUnfreezeWithdraw = "QueryUnfreezeWithdraw"
)

//
//   github     ，        ,
//      ，
var (
	PackageName    = "dplatformos.unfreeze"
	RPCName        = "DplatformOS.Unfreeze"
	UnfreezeX      = "unfreeze"
	ExecerUnfreeze = []byte(UnfreezeX)

	FixAmountX      = "FixAmount"
	LeftProportionX = "LeftProportion"
	SupportMeans    = []string{"FixAmount", "LeftProportion"}

	ForkTerminatePartX = "ForkTerminatePart"
	ForkUnfreezeIDX    = "ForkUnfreezeIDX"
)
