// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csp

import "crypto"

const (
	// ECDSAP256KeyGen ECDSA
	ECDSAP256KeyGen = 1
	// SM2P256KygGen SM2
	SM2P256KygGen = 2
)

// Key   key
type Key interface {
	Bytes() ([]byte, error)
	SKI() []byte
	Symmetric() bool
	Private() bool
	PublicKey() (Key, error)
}

// SignerOpts
type SignerOpts interface {
	crypto.SignerOpts
}

// CSP
type CSP interface {
	KeyGen(opts int) (k Key, err error)
	Sign(k Key, digest []byte, opts SignerOpts) (signature []byte, err error)
}

// KeyStore key
type KeyStore interface {
	ReadOnly() bool

	StoreKey(k Key) (err error)
}

// Signer
type Signer interface {
	Sign(k Key, digest []byte, opts SignerOpts) (signature []byte, err error)
}

// KeyGenerator key
type KeyGenerator interface {
	KeyGen(opts int) (k Key, err error)
}
