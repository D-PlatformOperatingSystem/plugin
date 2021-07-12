// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

/*
privacy            ，

           ：
1）             ， ：public address -> one-time addrss
2）    ，                  one-time address -> one-time address；
3）            ，  ：one-time address -> public address

    ：
1）      coin token     ，       balance   privacy     ；
2）            ，              (A,B),      ，balance           ；
3）    ，

*/

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"sort"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/types"
)

var privacylog = log.New("module", "execs.privacy")

var driverName = "privacy"

// Init initialize executor driver
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	drivers.Register(cfg, GetName(), newPrivacy, cfg.GetDappFork(driverName, "Enable"))
	//                 ，           ，
	//drivers.Register(newPrivacy().GetName(), newPrivacy, 0)
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&privacy{}))
}

// GetName get privacy name
func GetName() string {
	return newPrivacy().GetName()
}

type privacy struct {
	drivers.DriverBase
}

func newPrivacy() drivers.Driver {
	t := &privacy{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetDriverName get driver name
func (p *privacy) GetDriverName() string {
	return driverName
}

func (p *privacy) getUtxosByTokenAndAmount(exec, tokenName string, amount int64, count int32) ([]*pty.LocalUTXOItem, error) {
	localDB := p.GetLocalDB()
	var utxos []*pty.LocalUTXOItem
	prefix := CalcPrivacyUTXOkeyHeightPrefix(exec, tokenName, amount)
	values, err := localDB.List(prefix, nil, count, 0)
	if err != nil {
		return utxos, err
	}

	for _, value := range values {
		var utxo pty.LocalUTXOItem
		err := types.Decode(value, &utxo)
		if err != nil {
			privacylog.Info("getUtxosByTokenAndAmount", "decode to LocalUTXOItem failed because of", err)
			return utxos, err
		}
		utxos = append(utxos, &utxo)
	}

	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Height <= utxos[j].Height
	})
	return utxos, nil
}

func (p *privacy) getGlobalUtxoIndex(req *pty.ReqUTXOGlobalIndex) (types.Message, error) {
	debugBeginTime := time.Now()
	utxoGlobalIndexResp := &pty.ResUTXOGlobalIndex{}
	currentHeight := p.GetHeight()
	for _, amount := range req.GetAmount() {
		utxos, err := p.getUtxosByTokenAndAmount(req.GetAssetExec(), req.GetAssetSymbol(), amount, pty.UTXOCacheCount)
		if err != nil {
			return utxoGlobalIndexResp, err
		}

		index := len(utxos) - 1
		for ; index >= 0; index-- {
			if utxos[index].GetHeight()+pty.ConfirmedHeight <= currentHeight {
				break
			}
		}

		mixCount := req.GetMixCount()
		totalCnt := int32(index + 1)
		if mixCount > totalCnt {
			mixCount = totalCnt
		}

		utxoIndex4Amount := &pty.UTXOIndex4Amount{
			Amount: amount,
		}

		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		positions := random.Perm(int(totalCnt))
		for i := int(mixCount - 1); i >= 0; i-- {
			position := positions[i]
			item := utxos[position]
			utxoGlobalIndex := &pty.UTXOGlobalIndex{
				Outindex: item.GetOutindex(),
				Txhash:   item.GetTxhash(),
			}
			utxo := &pty.UTXOBasic{
				UtxoGlobalIndex: utxoGlobalIndex,
				OnetimePubkey:   item.GetOnetimepubkey(),
			}
			utxoIndex4Amount.Utxos = append(utxoIndex4Amount.Utxos, utxo)
		}
		utxoGlobalIndexResp.UtxoIndex4Amount = append(utxoGlobalIndexResp.UtxoIndex4Amount, utxoIndex4Amount)
	}

	duration := time.Since(debugBeginTime)
	privacylog.Debug("getGlobalUtxoIndex cost", duration)
	return utxoGlobalIndexResp, nil
}

//ShowAmountsOfUTXO     amount    utxo，             amout    UTXO,
//
//             UTXO, 1,3,5,10,20,30,100...
func (p *privacy) ShowAmountsOfUTXO(reqtoken *pty.ReqPrivacyToken) (types.Message, error) {
	querydb := p.GetLocalDB()

	key := CalcprivacyKeyTokenAmountType(reqtoken.GetAssetExec(), reqtoken.GetAssetSymbol())
	replyAmounts := &pty.ReplyPrivacyAmounts{}
	value, err := querydb.Get(key)
	if err != nil {
		return replyAmounts, err
	}
	if value != nil {
		var amountTypes pty.AmountsOfUTXO
		err := types.Decode(value, &amountTypes)
		if err == nil {
			for amount, count := range amountTypes.AmountMap {
				amountDetail := &pty.AmountDetail{
					Amount: amount,
					Count:  count,
				}
				replyAmounts.AmountDetail = append(replyAmounts.AmountDetail, amountDetail)
			}
		}

	}
	return replyAmounts, nil
}

//ShowUTXOs4SpecifiedAmount          UTXO     ，     ，  hash，
func (p *privacy) ShowUTXOs4SpecifiedAmount(reqtoken *pty.ReqPrivacyToken) (types.Message, error) {
	querydb := p.GetLocalDB()

	var replyUTXOsOfAmount pty.ReplyUTXOsOfAmount
	values, err := querydb.List(CalcPrivacyUTXOkeyHeightPrefix(reqtoken.GetAssetExec(), reqtoken.GetAssetSymbol(), reqtoken.Amount), nil, 0, 0)
	if err != nil {
		return &replyUTXOsOfAmount, err
	}
	if len(values) != 0 {
		for _, value := range values {
			var localUTXOItem pty.LocalUTXOItem
			err := types.Decode(value, &localUTXOItem)
			if err == nil {
				replyUTXOsOfAmount.LocalUTXOItems = append(replyUTXOsOfAmount.LocalUTXOItems, &localUTXOItem)
			}
		}
	}

	return &replyUTXOsOfAmount, nil
}

// CheckTx check transaction
func (p *privacy) CheckTx(tx *types.Transaction, index int) error {
	txhashstr := hex.EncodeToString(tx.Hash())
	var action pty.PrivacyAction
	err := types.Decode(tx.Payload, &action)
	if err != nil {
		privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr, "Decode tx.Payload error", err)
		return types.ErrActionNotSupport
	}
	privacylog.Debug("PrivacyTrading CheckTx", "txhash", txhashstr, "action type ", action.Ty)
	assertExec, token := action.GetAssetExecSymbol()
	if token == "" {
		return types.ErrInvalidParam
	}
	if pty.ActionPublic2Privacy == action.Ty && action.GetPublic2Privacy() != nil {
		return nil
	}
	input := action.GetInput()
	//           , input
	if len(input.GetKeyinput()) == 0 {
		privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr)
		return pty.ErrNilUtxoInput
	}

	output := action.GetOutput()
	//      utxo
	if action.GetPrivacy2Privacy() != nil && len(output.GetKeyoutput()) == 0 {
		privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr)
		return pty.ErrNilUtxoOutput
	}
	// check sign
	var ringSignature types.RingSignature
	if err := types.Decode(tx.Signature.Signature, &ringSignature); err != nil {
		privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr, "Decode tx.Signature.Signature error ", err)
		return pty.ErrRingSign
	}

	totalInput := int64(0)
	keyinput := input.GetKeyinput()
	keyImages := make([][]byte, len(keyinput))
	keys := make([][]byte, 0)
	pubkeys := make([][]byte, 0)
	for i, input := range keyinput {
		totalInput += input.Amount
		keyImages[i] = calcPrivacyKeyImageKey(assertExec, token, input.KeyImage)
		for j, globalIndex := range input.UtxoGlobalIndex {
			keys = append(keys, CalcPrivacyOutputKey(assertExec, token, input.Amount, common.ToHex(globalIndex.Txhash), int(globalIndex.Outindex)))
			pubkeys = append(pubkeys, ringSignature.Items[i].Pubkey[j])
		}
	}
	res, errIndex := p.checkUTXOValid(keyImages)
	if !res {
		if errIndex >= 0 && errIndex < int32(len(keyinput)) {
			input := keyinput[errIndex]
			privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr, "UTXO spent already errindex", errIndex, "utxo amout", input.Amount/types.Coin, "utxo keyimage", common.ToHex(input.KeyImage))
		}
		privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr, "err", "checkUTXOValid failed ")
		return pty.ErrDoubleSpendOccur
	}

	res, errIndex = p.checkPubKeyValid(keys, pubkeys)
	if !res {
		if errIndex >= 0 && errIndex < int32(len(pubkeys)) {
			privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr, "Wrong pubkey errIndex ", errIndex)
		}
		privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr, "checkPubKeyValid ", false)
		return pty.ErrPubkeysOfUTXO
	}

	//    coins            , assertExec
	cfg := p.GetAPI().GetConfig()
	if !cfg.IsPara() && (assertExec == "" || assertExec == "coins") {

		totalOutput := int64(0)
		for _, output := range output.GetKeyoutput() {
			totalOutput += output.GetAmount()
		}
		if tx.Fee < pty.PrivacyTxFee {
			privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr, "fee set:", tx.Fee, "required:", pty.PrivacyTxFee, " error ErrPrivacyTxFeeNotEnough")
			return pty.ErrPrivacyTxFeeNotEnough
		}
		//            ，        utxo  ,          UTXO,
		var feeAmount int64
		if action.Ty == pty.ActionPrivacy2Privacy {
			feeAmount = totalInput - totalOutput
		} else if action.Ty == pty.ActionPrivacy2Public && action.GetPrivacy2Public() != nil {
			feeAmount = totalInput - totalOutput - action.GetPrivacy2Public().Amount
		}

		if feeAmount < pty.PrivacyTxFee {
			privacylog.Error("PrivacyTrading CheckTx", "txhash", txhashstr, "fee available:", feeAmount, "required:", pty.PrivacyTxFee)
			return pty.ErrPrivacyTxFeeNotEnough
		}
	}
	return nil
}

func batchGet(stateDB db.KV, keyImages [][]byte) (values [][]byte, err error) {
	for i := 0; i < len(keyImages); i++ {
		v, err := stateDB.Get(keyImages[i])
		if err != nil && err != types.ErrNotFound {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

//  keyImage        ，        ，  true，     false
func (p *privacy) checkUTXOValid(keyImages [][]byte) (bool, int32) {
	stateDB := p.GetStateDB()
	values, err := batchGet(stateDB, keyImages)
	if err != nil {
		privacylog.Error("exec module", "checkUTXOValid failed to get value from statDB", err)
		return false, invalidIndex
	}
	if len(values) != len(keyImages) {
		privacylog.Error("exec module", "err", "checkUTXOValid return different count value with keys")
		return false, invalidIndex
	}
	for i, value := range values {
		if value != nil {
			privacylog.Error("exec module", "checkUTXOValid i=", i, " value=", value)
			return false, int32(i)
		}
	}

	return true, invalidIndex
}

func (p *privacy) checkPubKeyValid(keys [][]byte, pubkeys [][]byte) (bool, int32) {
	values, err := batchGet(p.GetStateDB(), keys)
	if err != nil {
		privacylog.Error("exec module", "checkPubKeyValid failed to get value from statDB with err", err)
		return false, invalidIndex
	}

	if len(values) != len(pubkeys) {
		privacylog.Error("exec module", "err", "checkPubKeyValid return different count value with keys")
		return false, invalidIndex
	}

	for i, value := range values {
		var keyoutput pty.KeyOutput
		types.Decode(value, &keyoutput)
		if !bytes.Equal(keyoutput.Onetimepubkey, pubkeys[i]) {
			privacylog.Error("exec module", "Invalid pubkey for tx with hash", string(keys[i]))
			return false, int32(i)
		}
	}

	return true, invalidIndex
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (p *privacy) CheckReceiptExecOk() bool {
	return true
}
