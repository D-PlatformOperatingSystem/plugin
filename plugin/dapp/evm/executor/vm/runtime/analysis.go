// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/holiman/uint256"
)

// Destinations
// Key     ，Value
// JUMPDEST
type Destinations map[common.Hash]bitvec

// Has   PC         JUMPDEST  ，
func (d Destinations) Has(codehash common.Hash, code []byte, dest *uint256.Int) bool {
	//       PC（    ），           ，      63
	//   ，     dest  PC
	udest, overflow := dest.Uint64WithOverflow()
	// PC cannot go beyond len(code) and certainly can't be bigger than 63bits.
	// Don't bother checking for JUMPDEST in that case.
	if overflow || udest >= uint64(len(code)) {
		return false
	}

	//             ，     ，
	m, analysed := d[codehash]
	if !analysed {
		m = codeBitmap(code)
		d[codehash] = m
	}

	// PC    ，
	return OpCode(code[udest]) == JUMPDEST && m.codeSegment(udest)
}

//    ，          ，
type bitvec []byte

func (bits *bitvec) set(pos uint64) {
	(*bits)[pos/8] |= 0x80 >> (pos % 8)
}
func (bits *bitvec) set8(pos uint64) {
	(*bits)[pos/8] |= 0xFF >> (pos % 8)
	(*bits)[pos/8+1] |= ^(0xFF >> (pos % 8))
}

//   pos
func (bits *bitvec) codeSegment(pos uint64) bool {
	return ((*bits)[pos/8] & (0x80 >> (pos % 8))) == 0
}

//              ，      1
func codeBitmap(code []byte) bitvec {
	// bitmap      4，     push32  （            4 byte，   32  ）
	//
	bits := make(bitvec, len(code)/8+1+4)
	for pc := uint64(0); pc < uint64(len(code)); {
		op := OpCode(code[pc])

		//  PUSH1 PUSH32         ，
		if op >= PUSH1 && op <= PUSH32 {
			numbits := op - PUSH1 + 1
			pc++

			//      8
			for ; numbits >= 8; numbits -= 8 {
				bits.set8(pc) // 8
				pc += 8
			}

			//           8  bit
			for ; numbits > 0; numbits-- {
				bits.set(pc)
				pc++
			}
		} else {
			pc++
		}
	}
	return bits
}
