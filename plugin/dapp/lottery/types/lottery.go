// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var (
	llog = log.New("module", "exectype."+LotteryX)
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(LotteryX))
	types.RegFork(LotteryX, InitFork)
	types.RegExec(LotteryX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(LotteryX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(LotteryX, NewType(cfg))
}

// LotteryType def
type LotteryType struct {
	types.ExecTypeBase
}

// NewType method
func NewType(cfg *types.DplatformOSConfig) *LotteryType {
	c := &LotteryType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName
func (lottery *LotteryType) GetName() string {
	return LotteryX
}

// GetLogMap method
func (lottery *LotteryType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogLotteryCreate: {Ty: reflect.TypeOf(ReceiptLottery{}), Name: "LogLotteryCreate"},
		TyLogLotteryBuy:    {Ty: reflect.TypeOf(ReceiptLottery{}), Name: "LogLotteryBuy"},
		TyLogLotteryDraw:   {Ty: reflect.TypeOf(ReceiptLottery{}), Name: "LogLotteryDraw"},
		TyLogLotteryClose:  {Ty: reflect.TypeOf(ReceiptLottery{}), Name: "LogLotteryClose"},
	}
}

// GetPayload method
func (lottery *LotteryType) GetPayload() types.Message {
	return &LotteryAction{}
}

// CreateTx method
func (lottery LotteryType) CreateTx(action string, message json.RawMessage) (*types.Transaction, error) {
	llog.Debug("lottery.CreateTx", "action", action)
	cfg := lottery.GetConfig()
	if action == "LotteryCreate" {
		var param LotteryCreateTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			llog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return CreateRawLotteryCreateTx(cfg, &param)
	} else if action == "LotteryBuy" {
		var param LotteryBuyTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			llog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return CreateRawLotteryBuyTx(cfg, &param)
	} else if action == "LotteryDraw" {
		var param LotteryDrawTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			llog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return CreateRawLotteryDrawTx(cfg, &param)
	} else if action == "LotteryClose" {
		var param LotteryCloseTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			llog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return CreateRawLotteryCloseTx(cfg, &param)
	} else {
		return nil, types.ErrNotSupport
	}
}

// GetTypeMap method
func (lottery LotteryType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Create": LotteryActionCreate,
		"Buy":    LotteryActionBuy,
		"Draw":   LotteryActionDraw,
		"Close":  LotteryActionClose,
	}
}

// CreateRawLotteryCreateTx method
func CreateRawLotteryCreateTx(cfg *types.DplatformOSConfig, parm *LotteryCreateTx) (*types.Transaction, error) {
	if parm == nil {
		llog.Error("CreateRawLotteryCreateTx", "parm", parm)
		return nil, types.ErrInvalidParam
	}

	v := &LotteryCreate{
		PurBlockNum:    parm.PurBlockNum,
		DrawBlockNum:   parm.DrawBlockNum,
		OpRewardRatio:  parm.OpRewardRatio,
		DevRewardRatio: parm.DevRewardRatio,
	}
	create := &LotteryAction{
		Ty:    LotteryActionCreate,
		Value: &LotteryAction_Create{v},
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(LotteryX)),
		Payload: types.Encode(create),
		Fee:     parm.Fee,
		To:      address.ExecAddress(cfg.ExecName(LotteryX)),
	}
	name := cfg.ExecName(LotteryX)
	tx, err := types.FormatTx(cfg, name, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// CreateRawLotteryBuyTx method
func CreateRawLotteryBuyTx(cfg *types.DplatformOSConfig, parm *LotteryBuyTx) (*types.Transaction, error) {
	if parm == nil {
		llog.Error("CreateRawLotteryBuyTx", "parm", parm)
		return nil, types.ErrInvalidParam
	}

	v := &LotteryBuy{
		LotteryId: parm.LotteryID,
		Amount:    parm.Amount,
		Number:    parm.Number,
		Way:       parm.Way,
	}
	buy := &LotteryAction{
		Ty:    LotteryActionBuy,
		Value: &LotteryAction_Buy{v},
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(LotteryX)),
		Payload: types.Encode(buy),
		Fee:     parm.Fee,
		To:      address.ExecAddress(cfg.ExecName(LotteryX)),
	}
	name := cfg.ExecName(LotteryX)
	tx, err := types.FormatTx(cfg, name, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// CreateRawLotteryDrawTx method
func CreateRawLotteryDrawTx(cfg *types.DplatformOSConfig, parm *LotteryDrawTx) (*types.Transaction, error) {
	if parm == nil {
		llog.Error("CreateRawLotteryDrawTx", "parm", parm)
		return nil, types.ErrInvalidParam
	}

	v := &LotteryDraw{
		LotteryId: parm.LotteryID,
	}
	draw := &LotteryAction{
		Ty:    LotteryActionDraw,
		Value: &LotteryAction_Draw{v},
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(LotteryX)),
		Payload: types.Encode(draw),
		Fee:     parm.Fee,
		To:      address.ExecAddress(cfg.ExecName(LotteryX)),
	}
	name := cfg.ExecName(LotteryX)
	tx, err := types.FormatTx(cfg, name, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// CreateRawLotteryCloseTx method
func CreateRawLotteryCloseTx(cfg *types.DplatformOSConfig, parm *LotteryCloseTx) (*types.Transaction, error) {
	if parm == nil {
		llog.Error("CreateRawLotteryCloseTx", "parm", parm)
		return nil, types.ErrInvalidParam
	}

	v := &LotteryClose{
		LotteryId: parm.LotteryID,
	}
	close := &LotteryAction{
		Ty:    LotteryActionClose,
		Value: &LotteryAction_Close{v},
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(LotteryX)),
		Payload: types.Encode(close),
		Fee:     parm.Fee,
		To:      address.ExecAddress(cfg.ExecName(LotteryX)),
	}

	name := cfg.ExecName(LotteryX)
	tx, err := types.FormatTx(cfg, name, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
