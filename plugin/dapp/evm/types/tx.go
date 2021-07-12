// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// CreateCallTx
type CreateCallTx struct {
	// Amount
	Amount uint64 `json:"amount"`
	// Code
	Code string `json:"code"`
	// GasLimit gas
	GasLimit uint64 `json:"gasLimit"`
	// GasPrice gas
	GasPrice uint32 `json:"gasPrice"`
	// Note
	Note string `json:"note"`
	// Alias
	Alias string `json:"alias"`
	// Fee
	Fee int64 `json:"fee"`
	// Name
	Name string `json:"name"`
	// IsCreate
	IsCreate bool `json:"isCreate"`
	// Abi              ABI
	Abi string `json:"abi"`
}
