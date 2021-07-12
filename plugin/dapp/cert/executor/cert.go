// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/cert/authority"
	ct "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/cert/types"
)

var clog = log.New("module", "execs.cert")
var driverName = ct.CertX

// Init
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	driverName = name
	var scfg ct.Authority
	if sub != nil {
		types.MustDecode(sub, &scfg)
	}
	err := authority.Author.Init(&scfg)
	if err != nil {
		clog.Error("error to initialize authority", err)
		return
	}
	drivers.Register(cfg, driverName, newCert, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

// InitExecType Init Exec Type
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Cert{}))
}

// GetName   cert
func GetName() string {
	return newCert().GetName()
}

// Cert cert
type Cert struct {
	drivers.DriverBase
}

func newCert() drivers.Driver {
	c := &Cert{}
	c.SetChild(c)
	c.SetIsFree(true)
	c.SetExecutorType(types.LoadExecutorType(driverName))
	return c
}

// GetDriverName   cert
func (c *Cert) GetDriverName() string {
	return driverName
}

// CheckTx cert   tx
func (c *Cert) CheckTx(tx *types.Transaction, index int) error {
	//
	err := c.DriverBase.CheckTx(tx, index)
	if err != nil {
		return err
	}

	// auth
	if !authority.IsAuthEnable {
		clog.Error("Authority is not available. Please check the authority config or authority initialize error logs.")
		return ct.ErrInitializeAuthority
	}

	//
	if authority.Author.HistoryCertCache.CurHeight == -1 {
		err := c.loadHistoryByPrefix()
		if err != nil {
			return err
		}
	}

	//     <        ，cert
	if c.GetHeight() <= authority.Author.HistoryCertCache.CurHeight {
		err := c.loadHistoryByPrefix()
		if err != nil {
			return err
		}
	}

	//     >        ，      -1，          ，  cert
	nxtHeight := authority.Author.HistoryCertCache.NxtHeight
	if nxtHeight != -1 && c.GetHeight() > nxtHeight {
		err := c.loadHistoryByHeight()
		if err != nil {
			return err
		}
	}

	// auth
	return authority.Author.Validate(tx.GetSignature())
}

/**
            ，cert  、  、
*/
func (c *Cert) loadHistoryByPrefix() error {
	parm := &types.LocalDBList{
		Prefix:    []byte("LODB-cert-"),
		Key:       nil,
		Direction: 0,
		Count:     0,
	}
	result, err := c.DriverBase.GetAPI().LocalList(parm)
	if err != nil {
		return err
	}

	//          ，        cert
	if len(result.Values) == 0 {
		authority.Author.HistoryCertCache.CurHeight = 0
		return nil
	}

	//
	var historyData types.HistoryCertStore
	for _, v := range result.Values {
		err := types.Decode(v, &historyData)
		if err != nil {
			return err
		}
		if historyData.CurHeigth < c.GetHeight() && (historyData.NxtHeight >= c.GetHeight() || historyData.NxtHeight == -1) {
			return authority.Author.ReloadCert(&historyData)
		}
	}

	return ct.ErrGetHistoryCertData
}

/**
            ，cert
*/
func (c *Cert) loadHistoryByHeight() error {
	key := calcCertHeightKey(c.GetHeight())
	parm := &types.LocalDBGet{Keys: [][]byte{key}}
	result, err := c.DriverBase.GetAPI().LocalGet(parm)
	if err != nil {
		return err
	}
	var historyData types.HistoryCertStore
	for _, v := range result.Values {
		err := types.Decode(v, &historyData)
		if err != nil {
			return err
		}
		if historyData.CurHeigth < c.GetHeight() && historyData.NxtHeight >= c.GetHeight() {
			return authority.Author.ReloadCert(&historyData)
		}
	}
	return ct.ErrGetHistoryCertData
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (c *Cert) CheckReceiptExecOk() bool {
	return true
}
