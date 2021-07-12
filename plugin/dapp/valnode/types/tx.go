// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// NodeUpdateTx for construction
type NodeUpdateTx struct {
	PubKey string `json:"pubKey"`
	Power  int64  `json:"power"`
}
