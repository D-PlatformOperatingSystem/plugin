// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"math/big"

	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto/sha3"
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/holiman/uint256"
)

// Address        ，
//               Address<->big.Int， Address<->[]byte
//                 Hash160  ，   EVM      [20]byte，
type Address struct {
	Addr *address.Address
}

// Hash160Address EVM
type Hash160Address [Hash160Length]byte

// String
func (a Address) String() string { return a.Addr.String() }

// Bytes
func (a Address) Bytes() []byte {
	return a.Addr.Hash160[:]
}

// Big
func (a Address) Big() *big.Int {
	ret := new(big.Int).SetBytes(a.Bytes())
	return ret
}

// Hash
func (a Address) Hash() Hash { return ToHash(a.Bytes()) }

// ToHash160   EVM
func (a Address) ToHash160() Hash160Address {
	var h Hash160Address
	h.SetBytes(a.Bytes())
	return h
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (h *Hash160Address) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-Hash160Length:]
	}
	copy(h[Hash160Length-len(b):], b)
}

// String implements fmt.Stringer.
func (h Hash160Address) String() string {
	return h.Hex()
}

// Hex returns an EIP55-compliant hex string representation of the address.
func (h Hash160Address) Hex() string {
	unchecksummed := hex.EncodeToString(h[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

// ToAddress   DplatformOS
func (h Hash160Address) ToAddress() Address {
	return BytesToAddress(h[:])
}

// NewAddress xHash  EVM
func NewAddress(cfg *types.DplatformOSConfig, txHash []byte) Address {
	execAddr := address.GetExecAddress(cfg.ExecName("user.evm.") + BytesToHash(txHash).Hex())
	return Address{Addr: execAddr}
}

// ExecAddress
func ExecAddress(execName string) Address {
	execAddr := address.GetExecAddress(execName)
	return Address{Addr: execAddr}
}

// BytesToAddress
func BytesToAddress(b []byte) Address {
	a := new(address.Address)
	a.Version = 0
	a.SetBytes(copyBytes(LeftPadBytes(b, 20)))
	return Address{Addr: a}
}

// BytesToHash160Address
func BytesToHash160Address(b []byte) Hash160Address {
	var h Hash160Address
	h.SetBytes(b)
	return h
}

// StringToAddress
func StringToAddress(s string) *Address {
	addr, err := address.NewAddrFromString(s)
	if err != nil {
		log15.Error("create address form string error", "string:", s)
		return nil
	}
	return &Address{Addr: addr}
}

func copyBytes(data []byte) (out []byte) {
	out = make([]byte, 20)
	copy(out[:], data)
	return
}

func bigBytes(b *big.Int) (out []byte) {
	out = make([]byte, 20)
	copy(out[:], b.Bytes())
	return
}

// BigToAddress
func BigToAddress(b *big.Int) Address {
	a := new(address.Address)
	a.Version = 0
	a.SetBytes(bigBytes(b))
	return Address{Addr: a}
}

// EmptyAddress
func EmptyAddress() Address { return BytesToAddress([]byte{0}) }

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Hash160Address { return BytesToHash160Address(FromHex(s)) }

// Uint256ToAddress
func Uint256ToAddress(b *uint256.Int) Address {
	a := new(address.Address)
	a.Version = 0
	out := make([]byte, 20)
	copy(out[:], b.Bytes())
	a.SetBytes(out)
	return Address{Addr: a}
}

// HexToAddr
func HexToAddr(s string) Address {
	a := new(address.Address)
	a.Version = 0
	out := make([]byte, 20)
	copy(out[:], FromHex(s))
	a.SetBytes(out)
	return Address{Addr: a}
}
