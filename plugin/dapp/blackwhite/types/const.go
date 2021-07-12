// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

// status
const (
	BlackwhiteStatusCreate = iota + 1
	BlackwhiteStatusPlay
	BlackwhiteStatusShow
	BlackwhiteStatusTimeout
	BlackwhiteStatusDone
)

const (
	// TyLogBlackwhiteCreate log for blackwhite create game
	TyLogBlackwhiteCreate = 750
	// TyLogBlackwhitePlay log for blackwhite play game
	TyLogBlackwhitePlay = 751
	// TyLogBlackwhiteShow log for blackwhite show game
	TyLogBlackwhiteShow = 752
	// TyLogBlackwhiteTimeout log for blackwhite timeout game
	TyLogBlackwhiteTimeout = 753
	// TyLogBlackwhiteDone log for blackwhite down game
	TyLogBlackwhiteDone = 754
	// TyLogBlackwhiteLoopInfo log for blackwhite LoopInfo game
	TyLogBlackwhiteLoopInfo = 755
)

const (
	// GetBlackwhiteRoundInfo    cmd
	GetBlackwhiteRoundInfo = "GetBlackwhiteRoundInfo"
	// GetBlackwhiteByStatusAndAddr    cmd
	GetBlackwhiteByStatusAndAddr = "GetBlackwhiteByStatusAndAddr"
	// GetBlackwhiteloopResult    cmd
	GetBlackwhiteloopResult = "GetBlackwhiteloopResult"
)

var (
	// BlackwhiteX
	BlackwhiteX = "blackwhite"
	glog        = log15.New("module", BlackwhiteX)
	// JRPCName json RPC name
	JRPCName = "Blackwhite"
	// ExecerBlackwhite      byte
	ExecerBlackwhite = []byte(BlackwhiteX)
	actionName       = map[string]int32{
		"Create":      BlackwhiteActionCreate,
		"Play":        BlackwhiteActionPlay,
		"Show":        BlackwhiteActionShow,
		"TimeoutDone": BlackwhiteActionTimeoutDone,
	}
	logInfo = map[int64]*types.LogInfo{
		TyLogBlackwhiteCreate:   {Ty: reflect.TypeOf(ReceiptBlackwhite{}), Name: "LogBlackwhiteCreate"},
		TyLogBlackwhitePlay:     {Ty: reflect.TypeOf(ReceiptBlackwhite{}), Name: "LogBlackwhitePlay"},
		TyLogBlackwhiteShow:     {Ty: reflect.TypeOf(ReceiptBlackwhite{}), Name: "LogBlackwhiteShow"},
		TyLogBlackwhiteTimeout:  {Ty: reflect.TypeOf(ReceiptBlackwhite{}), Name: "LogBlackwhiteTimeout"},
		TyLogBlackwhiteDone:     {Ty: reflect.TypeOf(ReceiptBlackwhite{}), Name: "LogBlackwhiteDone"},
		TyLogBlackwhiteLoopInfo: {Ty: reflect.TypeOf(ReplyLoopResults{}), Name: "LogBlackwhiteLoopInfo"},
	}
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, ExecerBlackwhite)
}
