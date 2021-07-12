// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	auty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/autonomy/types"
)

type subConfig struct {
	Total      string `json:"total"`
	UseBalance bool   `json:"useBalance"`
}

var (
	alog         = log.New("module", "execs.autonomy")
	driverName   = auty.AutonomyX
	autonomyAddr string
	subcfg       subConfig
	ticketName   = auty.TicketX
)

// Init
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	if sub != nil {
		types.MustDecode(sub, &subcfg)
	}
	autonomyAddr = address.ExecAddress(cfg.ExecName(auty.AutonomyX))
	ticketName = cfg.ExecName(auty.TicketX)
	drivers.Register(cfg, GetName(), newAutonomy, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Autonomy{}))
}

// Autonomy
type Autonomy struct {
	drivers.DriverBase
}

func newAutonomy() drivers.Driver {
	t := &Autonomy{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetName
func GetName() string {
	return newAutonomy().GetName()
}

// GetDriverName
func (u *Autonomy) GetDriverName() string {
	return driverName
}
