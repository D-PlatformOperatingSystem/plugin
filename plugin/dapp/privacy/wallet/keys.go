// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import "fmt"

const (
	// PrivacyDBVersion          ，                     KEY
	PrivacyDBVersion = "Privacy-DBVersion"
	// Privacy4Addr                    KEY
	// KEY      	Privacy4Addr-
	// VALUE     types.WalletAccountPrivacy，
	Privacy4Addr = "Privacy-Addr"
	// AvailUTXOs             UTXO     KEY
	// KEY      	AvailUTXOs-tokenname-address-outtxhash-outindex   outtxhash    UTXO     ，  common.Byte2Hex()
	// VALUE     types.PrivacyDBStore，           UTXO
	AvailUTXOs = "Privacy-UTXO"
	// UTXOsSpentInTx          ，        UTXO  ，         ，     UTXO
	// KEY     	UTXOsSpentInTx：costtxhash 	  costtxhash          ，  common.Byte2Hex()
	// VALUE    	types.FTXOsSTXOsInOneTx
	UTXOsSpentInTx = "Privacy-UTXOsSpentInTx"
	// FrozenUTXOs          ，         UTXO      ，       KEY
	// KEY      	FrozenUTXOs:tokenname-address-costtxhash   costtxhash   UTXO     ，  common.Byte2Hex()
	// VALUE    	   UTXOsSpentInTx KEY
	FrozenUTXOs = "Privacy-FUTXO4Tx"
	// PrivacySTXO                   ，  UTXO     UTXO，       UTXO   KEY
	// KEY    	PrivacySTXO-tokenname-address-costtxhash	  costtxhash   UTXO     ，  common.Byte2Hex()
	// VALUE    	   UTXOsSpentInTx KEY
	PrivacySTXO = "Privacy-SUTXO"
	// STXOs4Tx      UTXO      KEY
	// KEY    	STXOs4Tx：costtxhash	  costtxhash   UTXO     ，  common.Byte2Hex()
	// VALUE    	   UTXOsSpentInTx KEY
	STXOs4Tx = "Privacy-SUTXO4Tx"
	// RevertSendtx                   UTXO
	// KEY    	RevertSendtx:tokenname-address-costtxhash	  costtxhash   UTXO     ，  common.Byte2Hex()
	// VALUE    	   UTXOsSpentInTx KEY
	RevertSendtx = "Privacy-RevertSendtx"
	// RecvPrivacyTx                      KEY
	// KEY    	RecvPrivacyTx:tokenname-address-heighstr	  heighstr       types.MaxTxsPerBlock              index
	// VALUE    	  PrivacyTX   KEY
	RecvPrivacyTx = "Privacy-RecvTX"
	// SendPrivacyTx                 KEY
	// KEY    	SendPrivacyTx:tokenname-address-heighstr	  heighstr       types.MaxTxsPerBlock              index
	// VALUE    	  PrivacyTX   KEY
	SendPrivacyTx = "Privacy-SendTX"
	// PrivacyTX                       KEY
	// KEY    	PrivacyTX:heighstr	  heighstr       types.MaxTxsPerBlock              index
	// VALUE    	types.WalletTxDetail
	PrivacyTX = "Privacy-TX"
	// ScanPrivacyInput        ，                      UTXO
	// KEY    	ScanPrivacyInput-outtxhash-outindex	  outtxhash    UTXO     ，  common.Byte2Hex()
	// VALUE    	types.UTXOGlobalIndex
	ScanPrivacyInput = "Privacy-ScaneInput"
	// ReScanUtxosFlag           UTXO
	// KEY    	ReScanUtxosFlag
	// VALUE    	types.Int64，
	//		UtxoFlagNoScan  int32 = 0
	//		UtxoFlagScaning int32 = 1
	//		UtxoFlagScanEnd int32 = 2
	ReScanUtxosFlag = "Privacy-RescanFlag"
)

func calcPrivacyDBVersion() []byte {
	return []byte(PrivacyDBVersion)
}

// calcUTXOKey     UTXO   ,       +
//key and prefix for privacy
//types.PrivacyDBStore      calcUTXOKey  key，
//1.  utxo                    ， key  value，   calcUTXOKey4TokenAddr    key   kv ；
//2.      ，calcUTXOKey4TokenAddr   kv   ，   calcPrivacyFUTXOKey   key  kv ，       key，
//                ，  utxo   futxo；
//3.             ，         futxo ，        ，  key   stxo ，
//4.           ，   del block    ，
// 4.a            stxo ，    stxo    ftxo ，
// 4.b            utxo ftxo  ，    utxo ftxo       ，    types.PrivacyDBStore
// 4.c            stxo  ，      ，     ，    utxo         ，        ，
func calcUTXOKey(txhash string, index int) []byte {
	return []byte(fmt.Sprintf("%s-%s-%d", AvailUTXOs, txhash, index))
}

func calcKey4UTXOsSpentInTx(key string) []byte {
	return []byte(fmt.Sprintf("%s:%s", UTXOsSpentInTx, key))
}

// calcPrivacyAddrKey
func calcPrivacyAddrKey(addr string) []byte {
	return []byte(fmt.Sprintf("%s-%s", Privacy4Addr, addr))
}

//calcAddrKey   addr    Account
func calcAddrKey(addr string) []byte {
	return []byte(fmt.Sprintf("Addr:%s", addr))
}

// calcPrivacyUTXOPrefix4Addr          UTXO     KEY
func calcPrivacyUTXOPrefix4Addr(assetExec, token, addr string) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s-%s-", AvailUTXOs, assetExec, token, addr))
}

// calcFTXOsKeyPrefix                        UTXO         KEY
func calcFTXOsKeyPrefix(assetExec, token, addr string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-", FrozenUTXOs, assetExec, token, addr))
}

// calcSendPrivacyTxKey
// addr
// key   calcTxKey(heightstr)
func calcSendPrivacyTxKey(assetExec, tokenname, addr, key string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-%s", SendPrivacyTx, assetExec, tokenname, addr, key))
}

// calcRecvPrivacyTxKey
// addr
// key   calcTxKey(heightstr)
func calcRecvPrivacyTxKey(assetExec, tokenname, addr, key string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-%s", RecvPrivacyTx, assetExec, tokenname, addr, key))
}

// calcUTXOKey4TokenAddr         UTXO Key
func calcUTXOKey4TokenAddr(assetExec, token, addr, txhash string, index int) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s-%s-%s-%d", AvailUTXOs, assetExec, token, addr, txhash, index))
}

// calcKey4FTXOsInTx       ,   UTXO
func calcKey4FTXOsInTx(assetExec, token, addr, txhash string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-%s", FrozenUTXOs, assetExec, token, addr, txhash))
}

// calcRescanUtxosFlagKey                  UTXO
func calcRescanUtxosFlagKey(addr string) []byte {
	return []byte(fmt.Sprintf("%s-%s", ReScanUtxosFlag, addr))
}

func calcScanPrivacyInputUTXOKey(txhash string, index int) []byte {
	return []byte(fmt.Sprintf("%s-%s-%d", ScanPrivacyInput, txhash, index))
}

func calcKey4STXOsInTx(txhash string) []byte {
	return []byte(fmt.Sprintf("%s:%s", STXOs4Tx, txhash))
}

// calcSTXOTokenAddrTxKey           UTXO
func calcSTXOTokenAddrTxKey(assetExec, token, addr, txhash string) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s-%s-%s", PrivacySTXO, assetExec, token, addr, txhash))
}

func calcSTXOPrefix4Addr(assetExec, token, addr string) []byte {
	return []byte(fmt.Sprintf("%s-%s-%s-%s-", PrivacySTXO, assetExec, token, addr))
}

// calcRevertSendTxKey                UTXO     UTXO
func calcRevertSendTxKey(assetExec, tokenname, addr, txhash string) []byte {
	return []byte(fmt.Sprintf("%s:%s-%s-%s-%s", RevertSendtx, assetExec, tokenname, addr, txhash))
}

//  height*100000+index   Tx
//key:Tx:height*100000+index
func calcTxKey(key string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrivacyTX, key))
}
