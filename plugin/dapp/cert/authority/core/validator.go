// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// Validator
type Validator interface {
	Setup(config *AuthConfig) error

	Validate(cert []byte, pubKey []byte) error

	GetCertFromSignature(signature []byte) ([]byte, error)

	GetCertSnFromSignature(signature []byte) ([]byte, error)
}

// AuthConfig
type AuthConfig struct {
	RootCerts         [][]byte
	IntermediateCerts [][]byte
	RevocationList    [][]byte
}
