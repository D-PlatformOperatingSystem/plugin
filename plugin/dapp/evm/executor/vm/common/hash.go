// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"math/big"

	"github.com/holiman/uint256"

	"github.com/D-PlatformOperatingSystem/dpos/common"
)

const (
	// HashLength
	HashLength = 32

	// Hash160Length Hash160
	Hash160Length = 20
	// AddressLength
	AddressLength = 20
)

// Hash
type Hash common.Hash

// Str
func (h Hash) Str() string { return string(h[:]) }

// Bytes
func (h Hash) Bytes() []byte { return h[:] }

// Big
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex
func (h Hash) Hex() string { return Bytes2Hex(h[:]) }

// SetBytes          ，              ，    ，
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// BigToHash
func BigToHash(b *big.Int) Hash {
	return Hash(common.BytesToHash(b.Bytes()))
}

// Uint256ToHash
func Uint256ToHash(u *uint256.Int) Hash {
	return Hash(common.BytesToHash(u.Bytes()))
}

// BytesToHash  []byte
func BytesToHash(b []byte) Hash {
	return Hash(common.BytesToHash(b))
}

// ToHash  []byte
func ToHash(data []byte) Hash {
	return BytesToHash(common.Sha256(data))
}
