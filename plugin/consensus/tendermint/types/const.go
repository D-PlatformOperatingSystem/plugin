// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"errors"

	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

//authbls register
const (
	AuthBLS = 259
)

var (
	// ErrHeightLessThanOne error type
	ErrHeightLessThanOne = errors.New("ErrHeightLessThanOne")
	// ErrBaseTxType error type
	ErrBaseTxType = errors.New("ErrBaseTxType")
	// ErrBlockInfoTx error type
	ErrBlockInfoTx = errors.New("ErrBlockInfoTx")
	// ErrBaseExecErr error type
	ErrBaseExecErr = errors.New("ErrBaseExecErr")
	// ErrLastBlockID error type
	ErrLastBlockID = errors.New("ErrLastBlockID")
)

var (
	ttlog = log15.New("module", "tendermint-types")
	//ConsensusCrypto define
	ConsensusCrypto crypto.Crypto
	//CryptoName ...
	CryptoName string
	// SignMap define sign type
	SignMap = map[string]int{
		"secp256k1": types.SECP256K1,
		"ed25519":   types.ED25519,
		"sm2":       types.SM2,
		"bls":       AuthBLS,
	}
)
