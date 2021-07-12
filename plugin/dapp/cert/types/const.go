// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

var (
	// CertX cert
	CertX = "cert"
	// ExecerCert cert
	ExecerCert = []byte(CertX)
	actionName = map[string]int32{
		"New":    CertActionNew,
		"Update": CertActionUpdate,
		"Normal": CertActionNormal,
	}

	AdminKey = "Auth-cert-admin"
)
