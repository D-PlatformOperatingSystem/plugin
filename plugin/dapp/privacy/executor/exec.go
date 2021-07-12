// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/account"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/types"
)

// Exec_Public2Privacy execute public to privacy
func (p *privacy) Exec_Public2Privacy(payload *ty.Public2Privacy, tx *types.Transaction, index int) (*types.Receipt, error) {

	accDB, err := p.createAccountDB(payload.GetAssetExec(), payload.GetTokenname())
	if err != nil {
		privacylog.Error("Exec_pub2priv_newAccountDB", "exec", payload.GetAssetExec(),
			"symbol", payload.GetTokenname(), "err", err)
		return nil, err
	}
	txhashstr := hex.EncodeToString(tx.Hash())
	from := tx.From()
	receipt, err := accDB.ExecWithdraw(address.ExecAddress(string(tx.Execer)), from, payload.Amount)
	if err != nil {
		privacylog.Error("PrivacyTrading Exec", "txhash", txhashstr, "ExecWithdraw error ", err)
		return nil, err
	}

	txhash := common.ToHex(tx.Hash())
	output := payload.GetOutput().GetKeyoutput()
	//           block       ，          ，
	//      utxo  input，    ，          KV
	//executor        ，    kv   blockchain
	// ：       UTXO
	for index, keyOutput := range output {
		key := CalcPrivacyOutputKey(payload.AssetExec, payload.Tokenname, keyOutput.Amount, txhash, index)
		value := types.Encode(keyOutput)
		receipt.KV = append(receipt.KV, &types.KeyValue{Key: key, Value: value})
	}
	privacylog.Debug("testkey", "output", payload.GetOutput().Keyoutput)
	receiptLogs := p.buildPrivacyReceiptLog(payload.GetAssetExec(), payload.GetTokenname(), payload.GetOutput())
	execlog := &types.ReceiptLog{Ty: ty.TyLogPrivacyOutput, Log: types.Encode(receiptLogs)}
	receipt.Logs = append(receipt.Logs, execlog)

	//////////////////debug code begin///////////////
	privacylog.Debug("PrivacyTrading Exec", "ActionPublic2Privacy txhash", txhashstr, "receipt is", receipt)
	//////////////////debug code end///////////////

	return receipt, nil
}

// Exec_Privacy2Privacy execute privacy to privacy transaction
func (p *privacy) Exec_Privacy2Privacy(payload *ty.Privacy2Privacy, tx *types.Transaction, index int) (*types.Receipt, error) {

	txhashstr := hex.EncodeToString(tx.Hash())
	receipt := &types.Receipt{KV: make([]*types.KeyValue, 0)}
	privacyInput := payload.Input
	for _, keyInput := range privacyInput.Keyinput {
		value := []byte{keyImageSpentAlready}
		key := calcPrivacyKeyImageKey(payload.AssetExec, payload.Tokenname, keyInput.KeyImage)
		stateDB := p.GetStateDB()
		stateDB.Set(key, value)
		receipt.KV = append(receipt.KV, &types.KeyValue{Key: key, Value: value})
	}

	execlog := &types.ReceiptLog{Ty: ty.TyLogPrivacyInput, Log: types.Encode(payload.GetInput())}
	receipt.Logs = append(receipt.Logs, execlog)

	txhash := common.ToHex(tx.Hash())
	output := payload.GetOutput().GetKeyoutput()
	for index, keyOutput := range output {
		key := CalcPrivacyOutputKey(payload.AssetExec, payload.Tokenname, keyOutput.Amount, txhash, index)
		value := types.Encode(keyOutput)
		receipt.KV = append(receipt.KV, &types.KeyValue{Key: key, Value: value})
	}
	receiptLogs := p.buildPrivacyReceiptLog(payload.GetAssetExec(), payload.GetTokenname(), payload.GetOutput())
	execlog = &types.ReceiptLog{Ty: ty.TyLogPrivacyOutput, Log: types.Encode(receiptLogs)}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = types.ExecOk

	//////////////////debug code begin///////////////
	privacylog.Debug("PrivacyTrading Exec", "ActionPrivacy2Privacy txhash", txhashstr, "receipt is", receipt)
	//////////////////debug code end///////////////
	return receipt, nil
}

// Exec_Privacy2Public execute privacy to public transaction
func (p *privacy) Exec_Privacy2Public(payload *ty.Privacy2Public, tx *types.Transaction, index int) (*types.Receipt, error) {
	accDB, err := p.createAccountDB(payload.GetAssetExec(), payload.GetTokenname())
	if err != nil {
		privacylog.Error("Exec_pub2priv_newAccountDB", "exec", payload.GetAssetExec(),
			"symbol", payload.GetTokenname(), "err", err)
		return nil, err
	}
	txhashstr := hex.EncodeToString(tx.Hash())
	receipt, err := accDB.ExecDeposit(payload.To, address.ExecAddress(string(tx.Execer)), payload.Amount)
	if err != nil {
		privacylog.Error("PrivacyTrading Exec", "ActionPrivacy2Public txhash", txhashstr, "ExecDeposit error ", err)
		return nil, err
	}
	privacyInput := payload.Input
	for _, keyInput := range privacyInput.Keyinput {
		value := []byte{keyImageSpentAlready}
		key := calcPrivacyKeyImageKey(payload.AssetExec, payload.Tokenname, keyInput.KeyImage)
		stateDB := p.GetStateDB()
		stateDB.Set(key, value)
		receipt.KV = append(receipt.KV, &types.KeyValue{Key: key, Value: value})
	}

	execlog := &types.ReceiptLog{Ty: ty.TyLogPrivacyInput, Log: types.Encode(payload.GetInput())}
	receipt.Logs = append(receipt.Logs, execlog)

	txhash := common.ToHex(tx.Hash())
	output := payload.GetOutput().GetKeyoutput()
	for index, keyOutput := range output {
		key := CalcPrivacyOutputKey(payload.AssetExec, payload.Tokenname, keyOutput.Amount, txhash, index)
		value := types.Encode(keyOutput)
		receipt.KV = append(receipt.KV, &types.KeyValue{Key: key, Value: value})
	}

	receiptLog := p.buildPrivacyReceiptLog(payload.GetAssetExec(), payload.GetTokenname(), payload.GetOutput())
	execlog = &types.ReceiptLog{Ty: ty.TyLogPrivacyOutput, Log: types.Encode(receiptLog)}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = types.ExecOk

	//////////////////debug code begin///////////////
	privacylog.Debug("PrivacyTrading Exec", "ActionPrivacy2Privacy txhash", txhashstr, "receipt is", receipt)
	//////////////////debug code end///////////////
	return receipt, nil
}

func (p *privacy) createAccountDB(exec, symbol string) (*account.DB, error) {

	if exec == "" || exec == "coins" {
		return p.GetCoinsAccount(), nil
	}
	cfg := p.GetAPI().GetConfig()
	return account.NewAccountDB(cfg, exec, symbol, p.GetStateDB())
}

func (p *privacy) buildPrivacyReceiptLog(assetExec, assetSymbol string, output *ty.PrivacyOutput) *ty.ReceiptPrivacyOutput {
	if assetExec == "" {
		assetExec = "coins"
	}
	receipt := &ty.ReceiptPrivacyOutput{
		AssetExec:   assetExec,
		AssetSymbol: assetSymbol,
		Keyoutput:   output.Keyoutput,
	}

	return receipt
}
