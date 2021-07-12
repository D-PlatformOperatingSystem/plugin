// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

// Message
//  EVM         ，   Tx
type Message struct {
	to       *Address
	from     Address
	alias    string
	nonce    int64
	amount   uint64
	gasLimit uint64
	gasPrice uint32
	data     []byte
	abi      string
}

// NewMessage
func NewMessage(from Address, to *Address, nonce int64, amount uint64, gasLimit uint64, gasPrice uint32, data []byte, alias, abi string) *Message {
	return &Message{
		from:     from,
		to:       to,
		nonce:    nonce,
		amount:   amount,
		gasLimit: gasLimit,
		gasPrice: gasPrice,
		data:     data,
		alias:    alias,
		abi:      abi,
	}
}

// From
func (m Message) From() Address { return m.from }

// To
func (m Message) To() *Address { return m.to }

// GasPrice Gas
func (m Message) GasPrice() uint32 { return m.gasPrice }

// Value
func (m Message) Value() uint64 { return m.amount }

// Nonce  nonce
func (m Message) Nonce() int64 { return m.nonce }

// Data
func (m Message) Data() []byte { return m.data }

// GasLimit Gas
func (m Message) GasLimit() uint64 { return m.gasLimit }

// Alias
func (m Message) Alias() string { return m.alias }

// ABI   ABI
func (m Message) ABI() string { return m.abi }
