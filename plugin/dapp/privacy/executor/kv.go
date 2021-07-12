// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common"
)

const (
	privacyOutputKeyPrefix  = "mavl-privacy-UTXO-tahi"
	privacyKeyImagePrefix   = "mavl-privacy-UTXO-keyimage"
	privacyUTXOKEYPrefix    = "LODB-privacy-UTXO-tahhi"
	privacyAmountTypePrefix = "LODB-privacy-UTXO-atype"
	privacyTokenTypesPrefix = "LODB-privacy-UTXO-token"
	keyImageSpentAlready    = 0x01
	invalidIndex            = -1
)

//      utxo   ,  exec,token
func calcUtxoAssetPrefix(exec, token string) string {
	//  coins   key  exec  ,
	if exec == "" || exec == "coins" {
		return token
	}
	return exec + ":" + token
}

//CalcPrivacyOutputKey  key    types.KeyOutput
// kv  store
func CalcPrivacyOutputKey(exec, token string, amount int64, txhash string, outindex int) (key []byte) {
	return []byte(fmt.Sprintf(privacyOutputKeyPrefix+"-%s-%d-%s-%d", calcUtxoAssetPrefix(exec, token), amount, txhash, outindex))
}

func calcPrivacyKeyImageKey(exec, token string, keyimage []byte) []byte {
	return []byte(fmt.Sprintf(privacyKeyImagePrefix+"-%s-%s", calcUtxoAssetPrefix(exec, token), common.ToHex(keyimage)))
}

//CalcPrivacyUTXOkeyHeight                  amount    utxo global index
func CalcPrivacyUTXOkeyHeight(exec, token string, amount, height int64, txhash string, txindex, outindex int) (key []byte) {
	return []byte(fmt.Sprintf(privacyUTXOKEYPrefix+"-%s-%s-%d-%d-%s-%d-%d", exec, token, amount, height, txhash, txindex, outindex))
}

// CalcPrivacyUTXOkeyHeightPrefix get privacy utxo key by height and prefix
func CalcPrivacyUTXOkeyHeightPrefix(exec, token string, amount int64) (key []byte) {
	return []byte(fmt.Sprintf(privacyUTXOKEYPrefix+"-%s-%s-%d-", exec, token, amount))
}

//CalcprivacyKeyTokenAmountType          token amount   ，   1,3,5,100...     ,
func CalcprivacyKeyTokenAmountType(exec, token string) (key []byte) {
	return []byte(fmt.Sprintf(privacyAmountTypePrefix+"-%s-%s-", exec, token))
}

// CalcprivacyKeyTokenTypes get privacy token types key
func CalcprivacyKeyTokenTypes() (key []byte) {
	return []byte(privacyTokenTypesPrefix)
}

func calcExecLocalAssetKey(exec, symbol string) string {
	return exec + "-" + symbol
}
