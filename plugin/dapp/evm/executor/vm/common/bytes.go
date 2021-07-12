// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"encoding/hex"
	"math/big"
	"sort"
	"strings"
)

// RightPadBytes
func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded, slice)

	return padded
}

// LeftPadBytes
func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}

// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

// FromHex
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// HexToBytes
func HexToBytes(str string) ([]byte, error) {
	if len(str) > 1 && (strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X")) {
		str = str[2:]
	}
	return hex.DecodeString(str)
}

// Bytes2Hex         16
func Bytes2Hex(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

// Bytes2HexTrim         16
//         0
func Bytes2HexTrim(b []byte) string {
	//
	idx := sort.Search(len(b), func(i int) bool {
		return b[i] != 0
	})

	//    0，      ，     0x
	if idx == len(b) {
		return "0x00"
	}
	data := b[idx:]
	enc := make([]byte, len(data)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], data)
	return string(enc)
}

// CopyBytes returns an exact copy of the provided bytes.
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}
