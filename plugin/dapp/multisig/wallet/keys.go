// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import "fmt"

const (
	// MultisigAddr      owner           ，key:"multisig-addr-owneraddr, value [](multisigaddr,owneraddr,weight)
	MultisigAddr = "multisig-addr-"
)

func calcMultisigAddr(ownerAddr string) []byte {
	return []byte(fmt.Sprintf("%s%s", MultisigAddr, ownerAddr))
}

func calcPrefixMultisigAddr() []byte {
	return []byte(MultisigAddr)
}
