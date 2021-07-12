/*
 * Copyright D-Platform Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package types

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/types"
)

func init() {
	// init executor type
	types.AllowUserExec = append(types.AllowUserExec, []byte(OracleX))
	types.RegFork(OracleX, InitFork)
	types.RegExec(OracleX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(OracleX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(OracleX, NewType(cfg))
}

// OracleType
type OracleType struct {
	types.ExecTypeBase
}

// NewType
func NewType(cfg *types.DplatformOSConfig) *OracleType {
	c := &OracleType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName
func (o *OracleType) GetName() string {
	return OracleX
}

// GetPayload   oracle action
func (o *OracleType) GetPayload() types.Message {
	return &OracleAction{}
}

// GetTypeMap     map
func (o *OracleType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"EventPublish":     ActionEventPublish,
		"EventAbort":       ActionEventAbort,
		"ResultPrePublish": ActionResultPrePublish,
		"ResultAbort":      ActionResultAbort,
		"ResultPublish":    ActionResultPublish,
	}
}

// GetLogMap     map
func (o *OracleType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogEventPublish:     {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogEventPublish"},
		TyLogEventAbort:       {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogEventAbort"},
		TyLogResultPrePublish: {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogResultPrePublish"},
		TyLogResultAbort:      {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogResultAbort"},
		TyLogResultPublish:    {Ty: reflect.TypeOf(ReceiptOracle{}), Name: "LogResultPublish"},
	}
}
