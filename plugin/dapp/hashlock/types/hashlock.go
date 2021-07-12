// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var (
	hlog = log.New("module", "exectype.hashlock")
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(HashlockX))
	types.RegFork(HashlockX, InitFork)
	types.RegExec(HashlockX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(HashlockX, "Enable", 0)
	cfg.RegisterDappFork(HashlockX, ForkBadRepeatSecretX, 2715575)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(HashlockX, NewType(cfg))
}

// HashlockType def
type HashlockType struct {
	types.ExecTypeBase
}

// NewType method
func NewType(cfg *types.DplatformOSConfig) *HashlockType {
	c := &HashlockType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName
func (hashlock *HashlockType) GetName() string {
	return HashlockX
}

// GetPayload method
func (hashlock *HashlockType) GetPayload() types.Message {
	return &HashlockAction{}
}

// CreateTx method
func (hashlock *HashlockType) CreateTx(action string, message json.RawMessage) (*types.Transaction, error) {
	hlog.Debug("hashlock.CreateTx", "action", action)

	cfg := hashlock.GetConfig()
	if action == "HashlockLock" {
		var param HashlockLockTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			hlog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return CreateRawHashlockLockTx(cfg, &param)
	} else if action == "HashlockUnlock" {
		var param HashlockUnlockTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			hlog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return CreateRawHashlockUnlockTx(cfg, &param)
	} else if action == "HashlockSend" {
		var param HashlockSendTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			hlog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return CreateRawHashlockSendTx(cfg, &param)
	}
	return nil, types.ErrNotSupport

}

// GetTypeMap method
func (hashlock *HashlockType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Hlock":   HashlockActionLock,
		"Hsend":   HashlockActionSend,
		"Hunlock": HashlockActionUnlock,
	}
}

// GetLogMap method
func (hashlock *HashlockType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{}
}

// CreateRawHashlockLockTx method
func CreateRawHashlockLockTx(cfg *types.DplatformOSConfig, parm *HashlockLockTx) (*types.Transaction, error) {
	if parm == nil {
		hlog.Error("CreateRawHashlockLockTx", "parm", parm)
		return nil, types.ErrInvalidParam
	}

	v := &HashlockLock{
		Amount:        parm.Amount,
		Time:          parm.Time,
		Hash:          common.Sha256([]byte(parm.Secret)),
		ToAddress:     parm.ToAddr,
		ReturnAddress: parm.ReturnAddr,
	}
	lock := &HashlockAction{
		Ty:    HashlockActionLock,
		Value: &HashlockAction_Hlock{v},
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(HashlockX)),
		Payload: types.Encode(lock),
		Fee:     parm.Fee,
		To:      address.ExecAddress(cfg.ExecName(HashlockX)),
	}
	tx, err := types.FormatTx(cfg, cfg.ExecName(HashlockX), tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// CreateRawHashlockUnlockTx method
func CreateRawHashlockUnlockTx(cfg *types.DplatformOSConfig, parm *HashlockUnlockTx) (*types.Transaction, error) {
	if parm == nil {
		hlog.Error("CreateRawHashlockUnlockTx", "parm", parm)
		return nil, types.ErrInvalidParam
	}

	v := &HashlockUnlock{
		Secret: []byte(parm.Secret),
	}
	unlock := &HashlockAction{
		Ty:    HashlockActionUnlock,
		Value: &HashlockAction_Hunlock{v},
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(HashlockX)),
		Payload: types.Encode(unlock),
		Fee:     parm.Fee,
		To:      address.ExecAddress(HashlockX),
	}

	tx, err := types.FormatTx(cfg, cfg.ExecName(HashlockX), tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// CreateRawHashlockSendTx method
func CreateRawHashlockSendTx(cfg *types.DplatformOSConfig, parm *HashlockSendTx) (*types.Transaction, error) {
	if parm == nil {
		hlog.Error("CreateRawHashlockSendTx", "parm", parm)
		return nil, types.ErrInvalidParam
	}

	v := &HashlockSend{
		Secret: []byte(parm.Secret),
	}
	send := &HashlockAction{
		Ty:    HashlockActionSend,
		Value: &HashlockAction_Hsend{v},
	}
	tx := &types.Transaction{
		Execer:  []byte(HashlockX),
		Payload: types.Encode(send),
		Fee:     parm.Fee,
		To:      address.ExecAddress(HashlockX),
	}
	tx, err := types.FormatTx(cfg, cfg.ExecName(HashlockX), tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}