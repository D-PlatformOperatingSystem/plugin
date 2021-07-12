// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package state

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/model"
)

// EVMStateDB        ，  EVM      ；
//          ，                     StateDB  ；
// StateDB             （     ），         ，                  ；
// StateDB        ，                  ，                      ；
type EVMStateDB interface {
	// CreateAccount         
	CreateAccount(string, string, string, string)

	// SubBalance          
	SubBalance(string, string, uint64)
	// AddBalance          
	AddBalance(string, string, uint64)
	// GetBalance          
	GetBalance(string) uint64

	// GetNonce   nonce （       ，     0）
	GetNonce(string) uint64
	// SetNonce   nonce （       ，     0）
	SetNonce(string, uint64)

	// GetCodeHash              
	GetCodeHash(string) common.Hash
	// GetCode           
	GetCode(string) []byte
	// SetCode           
	SetCode(string, []byte)
	// GetCodeSize             
	GetCodeSize(string) int
	// SetAbi   ABI  
	SetAbi(addr, abi string)
	// GetAbi   ABI
	GetAbi(addr string) string

	// AddRefund   Gas    
	AddRefund(uint64)
	// GetRefund     Gas  
	GetRefund() uint64

	// GetState         
	GetState(string, common.Hash) common.Hash
	// SetState         
	SetState(string, common.Hash, common.Hash)

	// Suicide      
	Suicide(string) bool
	// HasSuicided         
	HasSuicided(string) bool

	// Exist             （               ）
	Exist(string) bool
	// Empty             （       、          ）
	Empty(string) bool

	// RevertToSnapshot        （                     ）
	RevertToSnapshot(int)
	// Snapshot         （    ）
	Snapshot() int
	// TransferStateData           
	TransferStateData(addr string)

	// AddLog         
	AddLog(*model.ContractLog)
	// AddPreimage   sha3  
	AddPreimage(common.Hash, []byte)

	// CanTransfer             
	CanTransfer(sender, recipient string, amount uint64) bool
	// Transfer     
	Transfer(sender, recipient string, amount uint64) bool

	// GetBlockHeight         
	GetBlockHeight() int64

	// GetConfig       
	GetConfig() *types.DplatformOSConfig
}
