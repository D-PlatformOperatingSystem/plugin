// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var clog = log.New("module", "execs.hashlock")

const minLockTime = 60

var driverName = "hashlock"

// Init hashlock
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	drivers.Register(cfg, GetName(), newHashlock, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Hashlock{}))
}

// GetName for hashlock
func GetName() string {
	return newHashlock().GetName()
}

// Hashlock driver
type Hashlock struct {
	drivers.DriverBase
}

func newHashlock() drivers.Driver {
	h := &Hashlock{}
	h.SetChild(h)
	h.SetExecutorType(types.LoadExecutorType(driverName))
	return h
}

// GetDriverName driverName
func (h *Hashlock) GetDriverName() string {
	return driverName
}

// CheckTx nil
func (h *Hashlock) CheckTx(tx *types.Transaction, index int) error {
	return nil
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (h *Hashlock) CheckReceiptExecOk() bool {
	return true
}
