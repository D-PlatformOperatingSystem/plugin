// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

// BlackwhiteCreateTx
type BlackwhiteCreateTx struct {
	PlayAmount  int64  `json:"amount"`
	PlayerCount int32  `json:"playerCount"`
	Timeout     int64  `json:"timeout"`
	GameName    string `json:"gameName"`
	Fee         int64  `json:"fee"`
}

// BlackwhitePlayTx
type BlackwhitePlayTx struct {
	GameID     string   `json:"gameID"`
	Amount     int64    `json:"amount"`
	HashValues [][]byte `json:"hashValues"`
	Fee        int64    `json:"fee"`
}

// BlackwhiteShowTx
type BlackwhiteShowTx struct {
	GameID string `json:"gameID"`
	Secret string `json:"secret"`
	Fee    int64  `json:"fee"`
}

// BlackwhiteTimeoutDoneTx
type BlackwhiteTimeoutDoneTx struct {
	GameID string `json:"GameID"`
	Fee    int64  `json:"fee"`
}
