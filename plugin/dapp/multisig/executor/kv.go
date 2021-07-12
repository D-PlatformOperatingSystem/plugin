// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"
)

//       key
const (
	//MultiSigPrefix statedb
	MultiSigPrefix   = "mavl-multisig-"
	MultiSigTxPrefix = "mavl-multisig-tx-"

	//MultiSigLocalPrefix localdb           multisig account count
	MultiSigLocalPrefix = "LODB-multisig-"
	MultiSigAccCount    = "acccount"
	MultiSigAcc         = "account"
	MultiSigAllAcc      = "allacc"
	MultiSigTx          = "tx"
	MultiSigRecvAssets  = "assets"
	MultiSigAccCreate   = "create"
)

//statedb
func calcMultiSigAccountKey(multiSigAccAddr string) (key []byte) {
	return []byte(fmt.Sprintf(MultiSigPrefix+"%s", multiSigAccAddr))
}

//    ："mavl-multisig-tx-accaddr-000000000000"
func calcMultiSigAccTxKey(multiSigAccAddr string, txid uint64) (key []byte) {
	txstr := fmt.Sprintf("%018d", txid)
	return []byte(fmt.Sprintf(MultiSigTxPrefix+"%s-%s", multiSigAccAddr, txstr))
}

//localdb

//         ，key:Msac value：count。
func calcMultiSigAccCountKey() []byte {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s", MultiSigAccCount))
}

//     MultiSig    ：             : key:Ms:allacc:index,value:accaddr
func calcMultiSigAllAcc(accindex int64) (key []byte) {
	accstr := fmt.Sprintf("%018d", accindex)
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s", MultiSigAllAcc, accstr))
}

//           value：MultiSig。key:Ms:acc
func calcMultiSigAcc(addr string) (key []byte) {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s", MultiSigAcc, addr))
}

//                 。key:Ms:create:createAddr，value：[]string。
func calcMultiSigAccCreateAddr(createAddr string) (key []byte) {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s:-%s", MultiSigAccCreate, createAddr))
}

//localdb

//           key:Ms:tx:addr:txid  value：MultiSigTx。
func calcMultiSigAccTx(addr string, txid uint64) (key []byte) {
	accstr := fmt.Sprintf("%018d", txid)

	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s-%s", MultiSigTx, addr, accstr))
}

//
//MultiSig              key:Ms:assets:addr:execname:symbol  value：AccountAssets。
//message AccountAssets {
//	string multiSigAddr = 1;
//	string execer 		= 2;
//	string symbol 		= 3;
//	int64  amount 		= 4;
func calcAddrRecvAmountKey(addr, execname, symbol string) []byte {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s-%s-%s", MultiSigRecvAssets, addr, execname, symbol))
}

//
func calcAddrRecvAmountPrefix(addr string) []byte {
	return []byte(fmt.Sprintf(MultiSigLocalPrefix+"%s-%s-", MultiSigRecvAssets, addr))
}
