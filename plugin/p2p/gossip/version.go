// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gossip

//    ：
// 1.p2p     nat   ，   peer stream，ping,version

//2018-3-26
// 1. p2p       ，  blockhash   block height
// 2.   p2p

//p2p     10020, 11000

//
const (
	//p2p
	lightBroadCastVersion = 10030
)

// VERSION number
const VERSION = lightBroadCastVersion

// MainNet Channel = 0x0000

const (
	defaultTestNetChannel = 256
)
