// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"bytes"
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/types"
)

//1,     paracross ，     ExecOk，        ，     PACK，    TyLogErr，OK       
//2,   paracross+other， other  PACK，      OK，      OK，     PACK，  TyLogErr
//3,     other，   PACK
func checkReceiptExecOk(receipt *types.ReceiptData) bool {
	if receipt.Ty == types.ExecOk {
		return true
	}
	//    allow    tx            paracross
	for _, log := range receipt.Logs {
		if log.Ty == types.TyLogErr {
			return false
		}
	}
	return true
}

//1.         ，                   ，         ,  PACK。（      ，            ,       block  ）
//2.        ，   paracross+user.p.xx.paracross  ，    user.p.xx.paracross  ，       
//3.         ExecOk,        ok ，      
//4.           ,      tx    group  ，         ，        ，      
//5.      ExecPack，     ，                ，                 ，    LogErr，    ，     
// para filter  ，           tx：
// 1,   	paracross	+  	     user.p.xx.paracross        
// 2,      paracross	+  	     user.p.xx.other 		      
// 3,      other  		+ 	     user.p.xx.paracross 	       
// 4,    	other 		+ 	     user.p.xx.other 		      
// 5,   +    user.p.xx.paracross    				        
// 6,    	    user.p.xx.paracross + user.p.xx.other          
// 7,         all user.p.xx.other  					       
/////                    tx，    tx
// para filter  ，           tx：
// 1,   +    user.p.xx.paracross    				          paracross      
// 2,    	    user.p.xx.paracross + user.p.xx.other              paracross      
// 3,         user.p.xx.other     					           other  pack
func filterParaTxGroup(cfg *types.DplatformOSConfig, tx *types.Transaction, allTxs []*types.TxDetail, index int, mainBlockHeight, forkHeight int64) ([]*types.Transaction, int) {
	var headIdx int

	for i := index; i >= 0; i-- {
		if bytes.Equal(tx.Header, allTxs[i].Tx.Hash()) {
			headIdx = i
			break
		}
	}

	endIdx := headIdx + int(tx.GroupCount)
	for i := headIdx; i < endIdx; i++ {
		//    forkHeight         ，        ,        6.2.0              blockhash  ，   6.2.0    ，   
		if cfg.IsPara() && mainBlockHeight < forkHeight && !types.Conf(cfg, pt.ParaPrefixConsSubConf).IsEnable(pt.ParaFilterIgnoreTxGroup) {
			if types.IsParaExecName(string(allTxs[i].Tx.Execer)) {
				continue
			}
		}

		if !checkReceiptExecOk(allTxs[i].Receipt) {
			clog.Error("filterParaTxGroup rmv tx group", "txhash", hex.EncodeToString(allTxs[i].Tx.Hash()))
			return nil, endIdx
		}
	}
	//                     tx
	var retTxs []*types.Transaction
	for _, retTx := range allTxs[headIdx:endIdx] {
		retTxs = append(retTxs, retTx.Tx)
	}
	return retTxs, endIdx
}

//FilterTxsForPara include some main tx in tx group before ForkParacrossCommitTx
func FilterTxsForPara(cfg *types.DplatformOSConfig, main *types.ParaTxDetail) []*types.Transaction {
	var txs []*types.Transaction
	forkHeight := pt.GetDappForkHeight(cfg, pt.ForkCommitTx)
	for i := 0; i < len(main.TxDetails); i++ {
		tx := main.TxDetails[i].Tx
		if tx.GroupCount >= 2 {
			mainTxs, endIdx := filterParaTxGroup(cfg, tx, main.TxDetails, i, main.Header.Height, forkHeight)
			txs = append(txs, mainTxs...)
			i = endIdx - 1
			continue
		}
		//   paracross tx             , 6.2fork         user.p.xx.paracross      
		if main.Header.Height >= forkHeight && bytes.HasSuffix(tx.Execer, []byte(pt.ParaX)) && !checkReceiptExecOk(main.TxDetails[i].Receipt) {
			clog.Error("FilterTxsForPara rmv tx", "txhash", hex.EncodeToString(tx.Hash()))
			continue
		}

		txs = append(txs, tx)
	}
	return txs
}

// FilterParaCrossTxHashes only all para chain cross txs like xx.paracross exec
func FilterParaCrossTxHashes(txs []*types.Transaction) [][]byte {
	var txHashs [][]byte
	for _, tx := range txs {
		if types.IsParaExecName(string(tx.Execer)) && bytes.HasSuffix(tx.Execer, []byte(pt.ParaX)) {
			txHashs = append(txHashs, tx.Hash())
		}
	}
	return txHashs
}

// para filter  ，           tx：
// 1,   	paracross	+  	     user.p.xx.paracross        
// 2,      paracross	+  	     user.p.xx.other 		      
// 3,      other  		+ 	     user.p.xx.paracross 	      
// 4,    	other 		+ 	     user.p.xx.other 		      
// 5,   +    user.p.xx.paracross    				        
// 6,    	    user.p.xx.paracross + user.p.xx.other          
// 7,         all user.p.xx.other  					       
//             user.p.xx.paracross       ，                         paracross       
func crossTxGroupProc(title string, txs []*types.Transaction, index int) ([]*types.Transaction, int32) {
	var headIdx, endIdx int32

	for i := index; i >= 0; i-- {
		if bytes.Equal(txs[index].Header, txs[i].Hash()) {
			headIdx = int32(i)
			break
		}
	}
	//cross mix tx, contain main and para tx, main prefix with pt.paraX
	//              ，  paracross    ，              unfreeze  ，              
	//           ，        trade  
	endIdx = headIdx + txs[index].GroupCount
	for i := headIdx; i < endIdx; i++ {
		if bytes.HasPrefix(txs[i].Execer, []byte(pt.ParaX)) {
			return txs[headIdx:endIdx], endIdx
		}
	}
	//cross asset transfer in tx group
	var transfers []*types.Transaction
	for i := headIdx; i < endIdx; i++ {
		if types.IsSpecificParaExecName(title, string(txs[i].Execer)) && bytes.HasSuffix(txs[i].Execer, []byte(pt.ParaX)) {
			transfers = append(transfers, txs[i])

		}
	}
	return transfers, endIdx

}

//FilterParaMainCrossTxHashes ForkParacrossCommitTx    txgroup   main chain tx   
func FilterParaMainCrossTxHashes(title string, txs []*types.Transaction) [][]byte {
	var crossTxHashs [][]byte
	//  tx    paracross   user.p.  ， user.p.xx.  paracross      
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		if tx.GroupCount > 1 {
			groupTxs, end := crossTxGroupProc(title, txs, i)
			for _, tx := range groupTxs {
				crossTxHashs = append(crossTxHashs, tx.Hash())

			}
			i = int(end) - 1
			continue
		}
		if types.IsSpecificParaExecName(title, string(tx.Execer)) && bytes.HasSuffix(tx.Execer, []byte(pt.ParaX)) {
			crossTxHashs = append(crossTxHashs, tx.Hash())
		}
	}
	return crossTxHashs

}

//CalcTxHashsHash     txhash hash       
func CalcTxHashsHash(txHashs [][]byte) []byte {
	if len(txHashs) == 0 {
		return nil
	}
	totalTxHash := &types.ReqHashes{}
	totalTxHash.Hashes = append(totalTxHash.Hashes, txHashs...)
	data := types.Encode(totalTxHash)
	return common.Sha256(data)
}
