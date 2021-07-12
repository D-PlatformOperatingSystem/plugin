// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/holiman/uint256"
)

// ContractRef
type ContractRef interface {
	Address() common.Address
}

// AccountRef        （         ContractRef  ）
//           ，           ，         ，
type AccountRef common.Address

// Address
func (ar AccountRef) Address() common.Address { return (common.Address)(ar) }

// Contract     ，
//
type Contract struct {
	// CallerAddress      ，
	//              ，
	CallerAddress common.Address

	//         ，        （     ），         （     ）
	caller ContractRef

	//
	//   ，     （      CallCode         ，              ，   caller  ）
	self ContractRef

	// Jumpdests       ， JUMP JUMPI
	Jumpdests Destinations

	// Locally cached result of JUMPDEST analysis
	analysis bitvec

	// Code
	Code []byte
	// CodeHash
	CodeHash common.Hash

	// CodeAddr
	CodeAddr *common.Address
	// Input
	Input []byte

	// Gas         Gas（            ）
	Gas uint64

	// value        ，        ，
	value uint64

	// DelegateCall      ，        true
	DelegateCall bool
}

// NewContract
//         ，                       ，
func NewContract(caller ContractRef, object ContractRef, value uint64, gas uint64) *Contract {

	c := &Contract{CallerAddress: caller.Address(), caller: caller, self: object}

	//             ，         jumpdests
	//   ，      jumpdests
	if parent, ok := caller.(*Contract); ok {
		c.Jumpdests = parent.Jumpdests
	} else {
		c.Jumpdests = make(Destinations)
	}

	//   gas  ，         gas
	c.Gas = gas
	c.value = value

	return c
}

func (c *Contract) validJumpdest(dest *uint256.Int) bool {
	udest, overflow := dest.Uint64WithOverflow()
	// PC cannot go beyond len(code) and certainly can't be bigger than 63bits.
	// Don't bother checking for JUMPDEST in that case.
	if overflow || udest >= uint64(len(c.Code)) {
		return false
	}
	// Only JUMPDESTs allowed for destinations
	if OpCode(c.Code[udest]) != JUMPDEST {
		return false
	}
	return c.isCode(udest)
}

func (c *Contract) validJumpSubdest(udest uint64) bool {
	// PC cannot go beyond len(code) and certainly can't be bigger than 63 bits.
	// Don't bother checking for BEGINSUB in that case.
	if int64(udest) < 0 || udest >= uint64(len(c.Code)) {
		return false
	}
	// Only BEGINSUBs allowed for destinations
	if OpCode(c.Code[udest]) != BEGINSUB {
		return false
	}
	return c.isCode(udest)
}

//      PC         ，   PUSHN       ，isCode  true
func (c *Contract) isCode(udest uint64) bool {
	// Do we have a contract hash already?
	if c.CodeHash != (common.Hash{}) {
		// Does parent context have the analysis?
		analysis, exist := c.Jumpdests[c.CodeHash]
		if !exist {
			// Do the analysis and save in parent context
			// We do not need to store it in c.analysis
			analysis = codeBitmap(c.Code)
			c.Jumpdests[c.CodeHash] = analysis
		}
		// Also stash it in current contract for faster access
		c.analysis = analysis
		return analysis.codeSegment(udest)
	}
	// We don't have the code hash, most likely a piece of initcode not already
	// in state trie. In that case, we do an analysis, and save it locally, so
	// we don't have to recalculate it for every JUMP instruction in the execution
	// However, we don't save it within the parent context
	if c.analysis == nil {
		c.analysis = codeBitmap(c.Code)
	}
	return c.analysis.codeSegment(udest)
}

// AsDelegate
//            ，
func (c *Contract) AsDelegate() *Contract {
	c.DelegateCall = true

	//         ，          ，
	parent := c.caller.(*Contract)

	//             ，
	c.CallerAddress = parent.CallerAddress

	//
	c.value = parent.value
	return c
}

// GetOp
func (c *Contract) GetOp(n uint64) OpCode {
	return OpCode(c.GetByte(n))
}

// GetByte
func (c *Contract) GetByte(n uint64) byte {
	if n < uint64(len(c.Code)) {
		return c.Code[n]
	}

	return 0
}

// Caller
//            ，           ，      ，     caller          caller
//             ，     caller
func (c *Contract) Caller() common.Address {
	return c.CallerAddress
}

// UseGas       gas   gas
func (c *Contract) UseGas(gas uint64) (ok bool) {
	if c.Gas < gas {
		return false
	}
	c.Gas -= gas
	return true
}

// Address
//   ，     CallCode   ，                  ，
func (c *Contract) Address() common.Address {
	return c.self.Address()
}

// Value          ，
func (c *Contract) Value() uint64 {
	return c.value
}

// SetCode
func (c *Contract) SetCode(hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
}

// SetCallCode
func (c *Contract) SetCallCode(addr *common.Address, hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
	c.CodeAddr = addr
}
