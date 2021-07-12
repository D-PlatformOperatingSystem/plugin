// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
)

// ContractLog      ，  EVM  Log  ，
//                     ，
type ContractLog struct {
	// Address
	Address common.Address

	// TxHash
	TxHash common.Hash

	// Index
	Index int

	// Topics
	Topics []common.Hash

	// Data
	Data []byte
	//
	BlockNumber uint64 `json:"blockNumber"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex"`
	// hash of the block in which the transaction was included
	BlockHash common.Hash `json:"blockHash"`
}

// PrintLog
func (log *ContractLog) PrintLog() {
	log15.Debug("!Contract Log!", "Contract address", log.Address.String(), "TxHash", log.TxHash.Hex(), "Log Index", log.Index, "Log Topics", log.Topics, "Log Data", common.Bytes2Hex(log.Data))
}
