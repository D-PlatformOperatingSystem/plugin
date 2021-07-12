// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/types"
)

func init() {
	// init executor type
	types.AllowUserExec = append(types.AllowUserExec, ExecerPokerBull)
	types.RegFork(PokerBullX, InitFork)
	types.RegExec(PokerBullX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(PokerBullX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(PokerBullX, NewType(cfg))
}

// PokerBullType
type PokerBullType struct {
	types.ExecTypeBase
}

// NewType   pokerbull
func NewType(cfg *types.DplatformOSConfig) *PokerBullType {
	c := &PokerBullType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName
func (t *PokerBullType) GetName() string {
	return PokerBullX
}

// GetPayload   payload
func (t *PokerBullType) GetPayload() types.Message {
	return &PBGameAction{}
}

// GetTypeMap     map
func (t *PokerBullType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Start":    PBGameActionStart,
		"Continue": PBGameActionContinue,
		"Quit":     PBGameActionQuit,
		"Query":    PBGameActionQuery,
		"Play":     PBGameActionPlay,
	}
}

// GetLogMap     map
func (t *PokerBullType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogPBGameStart:    {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGameStart"},
		TyLogPBGameContinue: {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGameContinue"},
		TyLogPBGameQuit:     {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGameQuit"},
		TyLogPBGameQuery:    {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGameQuery"},
		TyLogPBGamePlay:     {Ty: reflect.TypeOf(ReceiptPBGame{}), Name: "TyLogPBGamePlay"},
	}
}
