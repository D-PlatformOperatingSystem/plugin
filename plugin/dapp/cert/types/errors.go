// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "errors"

var (
	// ErrValidateCertFailed cert
	ErrValidateCertFailed = errors.New("ErrValidateCertFailed")
	// ErrGetHistoryCertData
	ErrGetHistoryCertData = errors.New("ErrGetHistoryCertData")
	// ErrUnknowAuthSignType
	ErrUnknowAuthSignType = errors.New("ErrUnknowAuthSignType")
	// ErrInitializeAuthority
	ErrInitializeAuthority = errors.New("ErrInitializeAuthority")
	// ErrPermissionDeny
	ErrPermissionDeny = errors.New("ErrPermissionDeny")
)
