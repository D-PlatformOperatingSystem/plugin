// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "github.com/D-PlatformOperatingSystem/dpos/types"

// IsExpire   FTXO    ，    true
func (ftxos *FTXOsSTXOsInOneTx) IsExpire(blockheight, blocktime int64) bool {
	valid := ftxos.GetExpire()
	if valid == 0 {
		// Expire 0，  false
		return false
	}
	//   expireBound
	if valid <= types.ExpireBound {
		return valid <= blockheight
	}
	return valid <= blocktime
}

// SetExpire
func (ftxos *FTXOsSTXOsInOneTx) SetExpire(expire int64) {
	if expire > types.ExpireBound {
		// FTXO       ，  Tx       12
		ftxos.Expire = expire + 12
	} else {
		// FTXO           ，  Tx   +1
		ftxos.Expire = expire + 1
	}
}
