// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "github.com/D-PlatformOperatingSystem/dpos/types"

//cert
const (
	CertActionNew    = 1
	CertActionUpdate = 2
	CertActionNormal = 3

	AuthECDSA = 257
	AuthSM2   = 258
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, ExecerCert)
	// init executor type
	types.RegFork(CertX, InitFork)
	types.RegExec(CertX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(CertX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(CertX, NewType(cfg))
}

// CertType cert       
type CertType struct {
	types.ExecTypeBase
}

// NewType   cert    
func NewType(cfg *types.DplatformOSConfig) *CertType {
	c := &CertType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetPayload   payload
func (b *CertType) GetPayload() types.Message {
	return &CertAction{}
}

// GetName       
func (b *CertType) GetName() string {
	return CertX
}

// GetLogMap   logmap
func (b *CertType) GetLogMap() map[int64]*types.LogInfo {
	return nil
}

// GetTypeMap     map
func (b *CertType) GetTypeMap() map[string]int32 {
	return actionName
}
