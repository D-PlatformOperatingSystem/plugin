// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package types ...
package types

import "github.com/D-PlatformOperatingSystem/dpos/types"

const (
	// InvalidAction invalid action type
	InvalidAction = 0
	//action type for privacy

	// ActionPublic2Privacy public to privacy action type
	ActionPublic2Privacy = iota + 100
	// ActionPrivacy2Privacy privacy to privacy action type
	ActionPrivacy2Privacy
	// ActionPrivacy2Public privacy to public action type
	ActionPrivacy2Public

	// log for privacy

	// TyLogPrivacyFee privacy fee log type
	TyLogPrivacyFee = 500
	// TyLogPrivacyInput privacy input type
	TyLogPrivacyInput = 501
	// TyLogPrivacyOutput privacy output type
	TyLogPrivacyOutput = 502
)

const (

	//SignNameOnetimeED25519 privacy name of crypto
	SignNameOnetimeED25519 = "privacy.onetimeed25519"
	// SignNameRing signature name ring
	SignNameRing = "privacy.RingSignatue"
	// OnetimeED25519 one time ED25519
	OnetimeED25519 = 4
	// RingBaseonED25519 ring raseon ED25519
	RingBaseonED25519 = 5
	// PrivacyMaxCount max mix utxo cout
	PrivacyMaxCount = 16
	// PrivacyTxFee privacy tx fee
	PrivacyTxFee = 1e7
)

//const ...
const (
	// utxo
	UTXOCacheCount = 256
	// UtxoMaturityDegree utxo
	UtxoMaturityDegree = 12
	DOMDustThreshold  = types.Coin
	ConfirmedHeight    = 12
	SignatureSize      = (4 + 33 + 65)
	// Size1Kshiftlen tx    1k
	Size1Kshiftlen = 10
)
