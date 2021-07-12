// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/types"
)

func (p *privacy) execDelLocal(tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	txhashstr := hex.EncodeToString(tx.Hash())
	dbSet := &types.LocalDBSet{}
	localDB := p.GetLocalDB()
	for i, item := range receiptData.Logs {
		if item.Ty != ty.TyLogPrivacyOutput {
			continue
		}
		var receiptPrivacyOutput ty.ReceiptPrivacyOutput
		err := types.Decode(item.Log, &receiptPrivacyOutput)
		if err != nil {
			privacylog.Error("PrivacyTrading ExecDelLocal", "txhash", txhashstr, "Decode item.Log error ", err)
			panic(err)
		}
		assetExec := receiptPrivacyOutput.GetAssetExec()
		assetSymbol := receiptPrivacyOutput.GetAssetSymbol()
		txhashInByte := tx.Hash()
		txhash := common.ToHex(txhashInByte)
		for m, keyOutput := range receiptPrivacyOutput.Keyoutput {
			//kv1，       UTXO，          token      ，     ，  txhash UTXO
			key := CalcPrivacyUTXOkeyHeight(assetExec, assetSymbol, keyOutput.Amount, p.GetHeight(), txhash, i, m)
			kv := &types.KeyValue{Key: key, Value: nil}
			dbSet.KV = append(dbSet.KV, kv)

			//kv2，         kv  ，                        UTXO
			var amountTypes ty.AmountsOfUTXO
			key2 := CalcprivacyKeyTokenAmountType(assetExec, assetSymbol)
			value2, err := localDB.Get(key2)
			//    token
			if err == nil && value2 != nil {
				err := types.Decode(value2, &amountTypes)
				if err == nil {
					//              ，
					if count, ok := amountTypes.AmountMap[keyOutput.Amount]; ok {
						count--
						if 0 == count {
							delete(amountTypes.AmountMap, keyOutput.Amount)
						} else {
							amountTypes.AmountMap[keyOutput.Amount] = count
						}

						value2 := types.Encode(&amountTypes)
						kv := &types.KeyValue{Key: key2, Value: value2}
						dbSet.KV = append(dbSet.KV, kv)
						//    query       ，           amout       kv,
						localDB.Set(key2, nil)
					}
				}
			}

			//kv3,        token
			assetKey := calcExecLocalAssetKey(assetExec, assetSymbol)
			var tokenNames ty.TokenNamesOfUTXO
			key3 := CalcprivacyKeyTokenTypes()
			value3, err := localDB.Get(key3)
			if err == nil && value3 != nil {
				err := types.Decode(value3, &tokenNames)
				if err == nil {
					if settxhash, ok := tokenNames.TokensMap[assetKey]; ok {
						if settxhash == txhash {
							delete(tokenNames.TokensMap, assetKey)
							value3 := types.Encode(&tokenNames)
							kv := &types.KeyValue{Key: key3, Value: value3}
							dbSet.KV = append(dbSet.KV, kv)
							localDB.Set(key3, nil)
						}
					}
				}
			}
		}
	}
	return dbSet, nil
}

// ExecDelLocal_Public2Privacy local delete execute public to privacy transaction
func (p *privacy) ExecDelLocal_Public2Privacy(payload *ty.Public2Privacy, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return p.execDelLocal(tx, receiptData, index)
}

// ExecDelLocal_Privacy2Privacy local delete execute privacy to privacy transaction
func (p *privacy) ExecDelLocal_Privacy2Privacy(payload *ty.Privacy2Privacy, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return p.execDelLocal(tx, receiptData, index)
}

// ExecDelLocal_Privacy2Public local delete execute public to public transaction
func (p *privacy) ExecDelLocal_Privacy2Public(payload *ty.Privacy2Public, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return p.execDelLocal(tx, receiptData, index)
}
