// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"encoding/hex"
	"unsafe"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	privacy "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/crypto"
	privacytypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/types"
)

func checkAmountValid(amount int64) bool {
	if amount <= 0 {
		return false
	}
	//      ，       types.Coin    
	//                      
	if (int64(float64(amount)/float64(types.Coin)) * types.Coin) != amount {
		return false
	}
	return true
}

func makeViewSpendPubKeyPairToString(viewPubKey, spendPubKey []byte) string {
	pair := viewPubKey
	pair = append(pair, spendPubKey...)
	return hex.EncodeToString(pair)
}

// amount   1,2,5   ，     amount                 utxo
func decomAmount2Nature(amount int64, order int64) []int64 {
	res := make([]int64, 0)
	if order == 0 {
		return res
	}
	mul := amount / order
	switch mul {
	case 3:
		res = append(res, order)
		res = append(res, 2*order)
	case 4:
		res = append(res, 2*order)
		res = append(res, 2*order)
	case 6:
		res = append(res, 5*order)
		res = append(res, order)
	case 7:
		res = append(res, 5*order)
		res = append(res, 2*order)
	case 8:
		res = append(res, 5*order)
		res = append(res, 2*order)
		res = append(res, 1*order)
	case 9:
		res = append(res, 5*order)
		res = append(res, 2*order)
		res = append(res, 2*order)
	default:
		res = append(res, mul*order)
		return res
	}
	return res
}

// 62387455827 -> 455827 + 7000000 + 80000000 + 300000000 + 2000000000 + 60000000000, where 455827 <= dustThreshold
//res:[455827, 7000000, 80000000, 300000000, 2000000000, 60000000000]
func decomposeAmount2digits(amount, dustThreshold int64) []int64 {
	res := make([]int64, 0)
	if 0 >= amount {
		return res
	}

	isDustHandled := false
	var dust int64
	var order int64 = 1
	var chunk int64

	for 0 != amount {
		chunk = (amount % 10) * order
		amount /= 10
		order *= 10
		if dust+chunk < dustThreshold {
			dust += chunk //    ，    dust_threshold  
		} else {
			if !isDustHandled && 0 != dust {
				//1st      ，  dust    
				res = append(res, dust)
				isDustHandled = true
			}
			if 0 != chunk {
				//2nd                
				goodAmount := decomAmount2Nature(chunk, order/10)
				res = append(res, goodAmount...)
			}
		}
	}

	//           < dustThreshold，         
	if !isDustHandled && 0 != dust {
		res = append(res, dust)
	}

	return res
}

func parseViewSpendPubKeyPair(in string) (viewPubKey, spendPubKey []byte, err error) {
	src, err := common.FromHex(in)
	if err != nil {
		return nil, nil, err
	}
	if 64 != len(src) {
		bizlog.Error("parseViewSpendPubKeyPair", "pair with len", len(src))
		return nil, nil, types.ErrPubKeyLen
	}
	viewPubKey = src[:32]
	spendPubKey = src[32:]
	return
}

// genCustomOuts          
//      ，P=Hs(rA)G+B
//func genCustomOuts(viewpubTo, spendpubto *[32]byte, transAmount int64, count int32) (*privacytypes.PrivacyOutput, error) {
//	decomDigit := make([]int64, count)
//	for i := range decomDigit {
//		decomDigit[i] = transAmount
//	}
//
//	pk := &privacy.PubKeyPrivacy{}
//	sk := &privacy.PrivKeyPrivacy{}
//	privacy.GenerateKeyPair(sk, pk)
//	RtxPublicKey := pk.Bytes()
//
//	sktx := (*[32]byte)(unsafe.Pointer(&sk[0]))
//	var privacyOutput privacytypes.PrivacyOutput
//	privacyOutput.RpubKeytx = RtxPublicKey
//	privacyOutput.Keyoutput = make([]*privacytypes.KeyOutput, len(decomDigit))
//
//	//             （UTXO），          
//	for index, digit := range decomDigit {
//		pubkeyOnetime, err := privacy.GenerateOneTimeAddr(viewpubTo, spendpubto, sktx, int64(index))
//		if err != nil {
//			bizlog.Error("genCustomOuts", "Fail to GenerateOneTimeAddr due to cause", err)
//			return nil, err
//		}
//		keyOutput := &privacytypes.KeyOutput{
//			Amount:        digit,
//			Onetimepubkey: pubkeyOnetime[:],
//		}
//		privacyOutput.Keyoutput[index] = keyOutput
//	}
//
//	return &privacyOutput, nil
//}

//       utxo   2   ，      utxo，        
//1.      utxo
//2.      utxo
func generateOuts(viewpubTo, spendpubto, viewpubChangeto, spendpubChangeto *[32]byte, transAmount, selectedAmount, fee int64) (*privacytypes.PrivacyOutput, error) {
	decomDigit := decomposeAmount2digits(transAmount, privacytypes.DOMDustThreshold)
	//    
	changeAmount := selectedAmount - transAmount - fee
	var decomChange []int64
	if 0 < changeAmount {
		decomChange = decomposeAmount2digits(changeAmount, privacytypes.DOMDustThreshold)
	}
	bizlog.Info("generateOuts", "decompose digit for amount", selectedAmount-fee, "decomDigit", decomDigit)

	pk := &privacy.PubKeyPrivacy{}
	sk := &privacy.PrivKeyPrivacy{}
	privacy.GenerateKeyPair(sk, pk)
	RtxPublicKey := pk.Bytes()

	sktx := (*[32]byte)(unsafe.Pointer(&sk[0]))
	var privacyOutput privacytypes.PrivacyOutput
	privacyOutput.RpubKeytx = RtxPublicKey
	privacyOutput.Keyoutput = make([]*privacytypes.KeyOutput, len(decomDigit)+len(decomChange))

	//             （UTXO），          
	for index, digit := range decomDigit {
		pubkeyOnetime, err := privacy.GenerateOneTimeAddr(viewpubTo, spendpubto, sktx, int64(index))
		if err != nil {
			bizlog.Error("generateOuts", "Fail to GenerateOneTimeAddr due to cause", err)
			return nil, err
		}
		keyOutput := &privacytypes.KeyOutput{
			Amount:        digit,
			Onetimepubkey: pubkeyOnetime[:],
		}
		privacyOutput.Keyoutput[index] = keyOutput
	}
	//         UTXO      UTXO
	for index, digit := range decomChange {
		pubkeyOnetime, err := privacy.GenerateOneTimeAddr(viewpubChangeto, spendpubChangeto, sktx, int64(index+len(decomDigit)))
		if err != nil {
			bizlog.Error("generateOuts", "Fail to GenerateOneTimeAddr for change due to cause", err)
			return nil, err
		}
		keyOutput := &privacytypes.KeyOutput{
			Amount:        digit,
			Onetimepubkey: pubkeyOnetime[:],
		}
		privacyOutput.Keyoutput[index+len(decomDigit)] = keyOutput
	}
	//         utxo，                
	if 0 != fee {
		//viewPub, _ := common.Hex2Bytes(types.ViewPubFee)
		//spendPub, _ := common.Hex2Bytes(types.SpendPubFee)
		//viewPublic := (*[32]byte)(unsafe.Pointer(&viewPub[0]))
		//spendPublic := (*[32]byte)(unsafe.Pointer(&spendPub[0]))
		//
		//pubkeyOnetime, err := privacy.GenerateOneTimeAddr(viewPublic, spendPublic, sktx, int64(len(privacyOutput.Keyoutput)))
		//if err != nil {
		//	bizlog.Error("transPub2PriV2", "Fail to GenerateOneTimeAddr for fee due to cause", err)
		//	return nil, nil, err
		//}
		//keyOutput := &types.KeyOutput{
		//	Amount:        fee,
		//	Ometimepubkey: pubkeyOnetime[:],
		//}
		//privacyOutput.Keyoutput = append(privacyOutput.Keyoutput, keyOutput)
	}

	return &privacyOutput, nil
}
