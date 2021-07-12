// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package gossip ...
package gossip

import (
	"bytes"
	"encoding/hex"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/p2p/utils"

	"github.com/D-PlatformOperatingSystem/dpos/common/merkle"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var (

	//           ,        mempool blockchain
	txHashFilter    = utils.NewFilter(TxRecvFilterCacheNum)
	blockHashFilter = utils.NewFilter(BlockFilterCacheNum)

	//            ,
	txSendFilter    = utils.NewFilter(TxSendFilterCacheNum)
	blockSendFilter = utils.NewFilter(BlockFilterCacheNum)

	//         ,             ,
	totalBlockCache = utils.NewSpaceLimitCache(BlockCacheNum, MaxBlockCacheByteSize)
	//          ,          ,    ,
	ltBlockCache = utils.NewSpaceLimitCache(BlockCacheNum/2, MaxBlockCacheByteSize/2)
)

type sendFilterInfo struct {
	//                 ,               ,             ,
	ignoreSendPeers map[string]bool
}

type pubFuncType func(interface{}, string)

func (n *Node) pubToPeer(data interface{}, pid string) {
	n.pubsub.FIFOPub(data, pid)
}

func (n *Node) processSendP2P(rawData interface{}, peerVersion int32, pid, peerAddr string) (sendData *types.BroadCastData, doSend bool) {
	//
	defer func() {
		if r := recover(); r != nil {
			log.Error("processSendP2P_Panic", "sendData", rawData, "peerAddr", peerAddr, "recoverErr", r)
			doSend = false
		}
	}()
	log.Debug("ProcessSendP2PBegin", "peerID", pid, "peerAddr", peerAddr)
	sendData = &types.BroadCastData{}
	doSend = false
	if tx, ok := rawData.(*types.P2PTx); ok {
		doSend = n.sendTx(tx, sendData, peerVersion, pid, peerAddr)
	} else if blc, ok := rawData.(*types.P2PBlock); ok {
		doSend = n.sendBlock(blc, sendData, peerVersion, pid, peerAddr)
	} else if query, ok := rawData.(*types.P2PQueryData); ok {
		doSend = n.sendQueryData(query, sendData, peerAddr)
	} else if rep, ok := rawData.(*types.P2PBlockTxReply); ok {
		doSend = n.sendQueryReply(rep, sendData, peerAddr)
	} else if ping, ok := rawData.(*types.P2PPing); ok {
		doSend = true
		sendData.Value = &types.BroadCastData_Ping{Ping: ping}
	}
	log.Debug("ProcessSendP2PEnd", "peerAddr", peerAddr, "doSend", doSend)
	return
}

func (n *Node) processRecvP2P(data *types.BroadCastData, pid string, pubPeerFunc pubFuncType, peerAddr string) (handled bool) {

	//
	defer func() {
		if r := recover(); r != nil {
			log.Error("ProcessRecvP2P_Panic", "recvData", data, "peerAddr", peerAddr, "recoverErr", r)
		}
	}()
	log.Debug("ProcessRecvP2P", "peerID", pid, "peerAddr", peerAddr)
	if pid == "" {
		return false
	}
	handled = true
	if tx := data.GetTx(); tx != nil {
		n.recvTx(tx, pid, peerAddr)
	} else if ltTx := data.GetLtTx(); ltTx != nil {
		n.recvLtTx(ltTx, pid, peerAddr, pubPeerFunc)
	} else if ltBlc := data.GetLtBlock(); ltBlc != nil {
		n.recvLtBlock(ltBlc, pid, peerAddr, pubPeerFunc)
	} else if blc := data.GetBlock(); blc != nil {
		n.recvBlock(blc, pid, peerAddr)
	} else if query := data.GetQuery(); query != nil {
		n.recvQueryData(query, pid, peerAddr, pubPeerFunc)
	} else if rep := data.GetBlockRep(); rep != nil {
		n.recvQueryReply(rep, pid, peerAddr, pubPeerFunc)
	} else {
		handled = false
	}
	log.Debug("ProcessRecvP2P", "peerAddr", peerAddr, "handled", handled)
	return
}

func (n *Node) sendBlock(block *types.P2PBlock, p2pData *types.BroadCastData, peerVersion int32, pid, peerAddr string) (doSend bool) {

	byteHash := block.Block.Hash(n.chainCfg)
	blockHash := hex.EncodeToString(byteHash)
	//
	ignoreSend := n.addIgnoreSendPeerAtomic(blockSendFilter, blockHash, pid)
	log.Debug("P2PSendBlock", "blockHash", blockHash, "peerIsLtVersion", peerVersion >= lightBroadCastVersion,
		"peerAddr", peerAddr, "ignoreSend", ignoreSend)
	if ignoreSend {
		return false
	}

	if peerVersion >= lightBroadCastVersion && types.Size(block.GetBlock()) >= int(n.nodeInfo.cfg.MinLtBlockSize*1024) {

		ltBlock := &types.LightBlock{}
		ltBlock.Size = int64(types.Size(block.Block))
		ltBlock.Header = block.Block.GetHeader(n.chainCfg)
		ltBlock.Header.Hash = byteHash[:]
		ltBlock.Header.Signature = block.Block.Signature
		ltBlock.MinerTx = block.Block.Txs[0]
		for _, tx := range block.Block.Txs[1:] {
			//tx short hash
			ltBlock.STxHashes = append(ltBlock.STxHashes, types.CalcTxShortHash(tx.Hash()))
		}

		// cache block
		if !totalBlockCache.Contains(blockHash) {
			totalBlockCache.Add(blockHash, block.Block, int(ltBlock.Size))
		}

		p2pData.Value = &types.BroadCastData_LtBlock{LtBlock: ltBlock}
	} else {
		p2pData.Value = &types.BroadCastData_Block{Block: block}
	}

	return true
}

func (n *Node) sendQueryData(query *types.P2PQueryData, p2pData *types.BroadCastData, peerAddr string) bool {
	log.Debug("P2PSendQueryData", "peerAddr", peerAddr)
	p2pData.Value = &types.BroadCastData_Query{Query: query}
	return true
}

func (n *Node) sendQueryReply(rep *types.P2PBlockTxReply, p2pData *types.BroadCastData, peerAddr string) bool {
	log.Debug("P2PSendQueryReply", "peerAddr", peerAddr)
	p2pData.Value = &types.BroadCastData_BlockRep{BlockRep: rep}
	return true
}

func (n *Node) sendTx(tx *types.P2PTx, p2pData *types.BroadCastData, peerVersion int32, pid, peerAddr string) (doSend bool) {

	txHash := hex.EncodeToString(tx.Tx.Hash())
	ttl := tx.GetRoute().GetTTL()

	//     ttl,
	if ttl > n.nodeInfo.cfg.MaxTTL {
		return false
	}

	isLightSend := peerVersion >= lightBroadCastVersion && ttl >= n.nodeInfo.cfg.LightTxTTL
	//
	if n.addIgnoreSendPeerAtomic(txSendFilter, txHash, pid) {
		return false
	}

	//log.Debug("P2PSendTx", "txHash", txHash, "ttl", ttl, "isLightSend", isLightSend, "peerAddr", peerAddr, "ignoreSend", ignoreSend)

	//    ttl
	if isLightSend {
		p2pData.Value = &types.BroadCastData_LtTx{ //     ttl,
			LtTx: &types.LightTx{
				TxHash: tx.Tx.Hash(),
				Route:  tx.GetRoute(),
			},
		}
	} else {
		p2pData.Value = &types.BroadCastData_Tx{Tx: tx}
	}
	return true
}

func (n *Node) recvTx(tx *types.P2PTx, pid, peerAddr string) {
	if tx.GetTx() == nil {
		return
	}
	txHash := hex.EncodeToString(tx.GetTx().Hash())
	//   id       ,
	n.addIgnoreSendPeerAtomic(txSendFilter, txHash, pid)
	//
	isDuplicate := txHashFilter.AddWithCheckAtomic(txHash, true)
	//log.Debug("recvTx", "tx", txHash, "ttl", tx.GetRoute().GetTTL(), "peerAddr", peerAddr, "duplicateTx", isDuplicate)
	if isDuplicate {
		return
	}
	//             ,  route
	if tx.GetRoute() == nil {
		tx.Route = &types.P2PRoute{TTL: 1}
	}
	txHashFilter.Add(txHash, tx.GetRoute())

	errs := n.postMempool(txHash, tx.GetTx())
	if errs != nil {
		log.Error("recvTx", "process post mempool EventTx msg Error", errs.Error())
	}

}

func (n *Node) recvLtTx(tx *types.LightTx, pid, peerAddr string, pubPeerFunc pubFuncType) {

	txHash := hex.EncodeToString(tx.TxHash)
	//   id       ,
	n.addIgnoreSendPeerAtomic(txSendFilter, txHash, pid)
	exist := txHashFilter.Contains(txHash)
	//log.Debug("recvLtTx", "txHash", txHash, "ttl", tx.GetRoute().GetTTL(), "peerAddr", peerAddr, "exist", exist)
	//     ,                .                  ,
	if !exist {

		query := &types.P2PQueryData{}
		query.Value = &types.P2PQueryData_TxReq{
			TxReq: &types.P2PTxReq{
				TxHash: tx.TxHash,
			},
		}
		//
		pubPeerFunc(query, pid)
	}
}

func (n *Node) recvBlock(block *types.P2PBlock, pid, peerAddr string) {

	if block.GetBlock() == nil {
		return
	}
	blockHash := hex.EncodeToString(block.GetBlock().Hash(n.chainCfg))
	//   id       ,
	n.addIgnoreSendPeerAtomic(blockSendFilter, blockHash, pid)
	//      ,      blockchain
	isDuplicate := blockHashFilter.AddWithCheckAtomic(blockHash, true)
	log.Debug("recvBlock", "blockHeight", block.GetBlock().GetHeight(), "peerAddr", peerAddr,
		"block size(KB)", float32(block.Block.Size())/1024, "blockHash", blockHash, "duplicateBlock", isDuplicate)
	if isDuplicate {
		return
	}
	//   blockchain
	if err := n.postBlockChain(blockHash, pid, block.GetBlock()); err != nil {
		log.Error("recvBlock", "send block to blockchain Error", err.Error())
	}

}

func (n *Node) recvLtBlock(ltBlock *types.LightBlock, pid, peerAddr string, pubPeerFunc pubFuncType) {

	blockHash := hex.EncodeToString(ltBlock.Header.Hash)
	//   id       ,
	n.addIgnoreSendPeerAtomic(blockSendFilter, blockHash, pid)
	//         block
	isDuplicate := blockHashFilter.AddWithCheckAtomic(blockHash, true)
	log.Debug("recvLtBlock", "blockHash", blockHash, "blockHeight", ltBlock.GetHeader().GetHeight(),
		"peerAddr", peerAddr, "duplicateBlock", isDuplicate)
	if isDuplicate {
		return
	}
	//  block
	block := &types.Block{}
	block.TxHash = ltBlock.Header.TxHash
	block.Signature = ltBlock.Header.Signature
	block.ParentHash = ltBlock.Header.ParentHash
	block.Height = ltBlock.Header.Height
	block.BlockTime = ltBlock.Header.BlockTime
	block.Difficulty = ltBlock.Header.Difficulty
	block.Version = ltBlock.Header.Version
	block.StateHash = ltBlock.Header.StateHash
	//add miner tx
	block.Txs = append(block.Txs, ltBlock.MinerTx)

	txList := &types.ReplyTxList{}
	ok := false
	//get tx list from mempool
	if len(ltBlock.STxHashes) > 0 {
		resp, err := n.queryMempool(types.EventTxListByHash, &types.ReqTxHashList{Hashes: ltBlock.STxHashes, IsShortHash: true})
		if err != nil {
			log.Error("recvLtBlock", "queryTxListByHashErr", err)
			return
		}

		txList, ok = resp.(*types.ReplyTxList)
		if !ok {
			log.Error("recvLtBlock", "queryMemPool", "nilReplyTxList")
		}
	}
	nilTxIndices := make([]int32, 0)
	for i := 0; ok && i < len(txList.Txs); i++ {
		tx := txList.Txs[i]
		if tx == nil {
			//tx not exist in mempool
			nilTxIndices = append(nilTxIndices, int32(i+1))
			tx = &types.Transaction{}
		} else if count := tx.GetGroupCount(); count > 0 {

			group, err := tx.GetTxGroup()
			if err != nil {
				log.Error("recvLtBlock", "getTxGroupErr", err)
				//
				nilTxIndices = nilTxIndices[:0]
				break
			}
			block.Txs = append(block.Txs, group.Txs...)
			//
			i += len(group.Txs) - 1
			continue
		}

		block.Txs = append(block.Txs, tx)
	}
	nilTxLen := len(nilTxIndices)
	//             ,
	if nilTxLen == 0 && len(block.Txs) == int(ltBlock.Header.TxCount) {
		if bytes.Equal(block.TxHash, merkle.CalcMerkleRoot(n.chainCfg, block.Height, block.Txs)) {
			log.Debug("recvLtBlock", "height", block.GetHeight(), "peerAddr", peerAddr,
				"blockHash", blockHash, "block size(KB)", float32(ltBlock.Size)/1024)
			//   blockchain
			if err := n.postBlockChain(blockHash, pid, block); err != nil {
				log.Error("recvLtBlock", "send block to blockchain Error", err.Error())
			}
			return
		}
		log.Debug("recvLtBlock:TxHashCheckFail", "height", block.GetHeight(), "peerAddr", peerAddr,
			"blockHash", blockHash, "block.Txs", block.Txs)
	}
	//            1/3           2/3,
	if nilTxLen > 0 && (float32(nilTxLen) > float32(ltBlock.Header.TxCount)/3 ||
		float32(block.Size()) < float32(ltBlock.Size)/3) {
		nilTxIndices = nilTxIndices[:0]
	}
	log.Debug("recvLtBlock", "queryBlockHash", blockHash, "queryHeight", ltBlock.GetHeader().GetHeight(), "queryTxNum", len(nilTxIndices))

	// query not exist txs
	query := &types.P2PQueryData{
		Value: &types.P2PQueryData_BlockTxReq{
			BlockTxReq: &types.P2PBlockTxReq{
				BlockHash: blockHash,
				TxIndices: nilTxIndices,
			},
		},
	}
	//       block
	ltBlockCache.Add(blockHash, block, block.Size())
	//pub to specified peer
	pubPeerFunc(query, pid)
}

func (n *Node) recvQueryData(query *types.P2PQueryData, pid, peerAddr string, pubPeerFunc pubFuncType) {

	if txReq := query.GetTxReq(); txReq != nil {

		txHash := hex.EncodeToString(txReq.TxHash)
		log.Debug("recvQueryTx", "txHash", txHash, "peerAddr", peerAddr)
		// mempool
		resp, err := n.queryMempool(types.EventTxListByHash, &types.ReqTxHashList{Hashes: []string{string(txReq.TxHash)}})
		if err != nil {
			log.Error("recvQuery", "queryMempoolErr", err)
			return
		}

		txList, _ := resp.(*types.ReplyTxList)
		//
		if len(txList.GetTxs()) != 1 || txList.GetTxs()[0] == nil {
			log.Error("recvQueryTx", "txHash", txHash, "err", "recvNilTxFromMempool")
			return
		}
		p2pTx := &types.P2PTx{Tx: txList.Txs[0]}
		//           , ttl   1
		p2pTx.Route = &types.P2PRoute{TTL: 1}
		n.removeIgnoreSendPeerAtomic(txSendFilter, txHash, pid)
		pubPeerFunc(p2pTx, pid)

	} else if blcReq := query.GetBlockTxReq(); blcReq != nil {

		log.Debug("recvQueryBlockTx", "blockHash", blcReq.BlockHash, "queryTxCount", len(blcReq.TxIndices), "peerAddr", peerAddr)
		if block, ok := totalBlockCache.Get(blcReq.BlockHash).(*types.Block); ok {

			blockRep := &types.P2PBlockTxReply{BlockHash: blcReq.BlockHash}

			blockRep.TxIndices = blcReq.TxIndices
			for _, idx := range blcReq.TxIndices {
				blockRep.Txs = append(blockRep.Txs, block.Txs[idx])
			}
			//
			if len(blockRep.TxIndices) == 0 {
				blockRep.Txs = block.Txs
			}
			pubPeerFunc(blockRep, pid)
		}
	}
}

func (n *Node) recvQueryReply(rep *types.P2PBlockTxReply, pid, peerAddr string, pubPeerFunc pubFuncType) {

	log.Debug("recvQueryReplyBlock", "blockHash", rep.GetBlockHash(), "queryTxsCount", len(rep.GetTxIndices()), "peerAddr", peerAddr)
	val, exist := ltBlockCache.Remove(rep.BlockHash)
	block, _ := val.(*types.Block)
	//not exist in cache or nil block
	if !exist || block == nil {
		return
	}
	for i, idx := range rep.TxIndices {
		block.Txs[idx] = rep.Txs[i]
	}

	//
	if len(rep.TxIndices) == 0 {
		block.Txs = rep.Txs
	}

	//   root hash
	if bytes.Equal(block.TxHash, merkle.CalcMerkleRoot(n.chainCfg, block.Height, block.Txs)) {

		log.Debug("recvQueryReplyBlock", "blockHeight", block.GetHeight(), "peerAddr", peerAddr,
			"block size(KB)", float32(block.Size())/1024, "blockHash", rep.BlockHash)
		//   blockchain
		if err := n.postBlockChain(rep.BlockHash, pid, block); err != nil {
			log.Error("recvQueryReplyBlock", "send block to blockchain Error", err.Error())
		}
	} else if len(rep.TxIndices) != 0 {
		log.Debug("recvQueryReplyBlock", "GetTotalBlock", block.GetHeight())
		//              ,
		query := &types.P2PQueryData{
			Value: &types.P2PQueryData_BlockTxReq{
				BlockTxReq: &types.P2PBlockTxReq{
					BlockHash: rep.BlockHash,
					TxIndices: nil,
				},
			},
		}
		block.Txs = nil
		ltBlockCache.Add(rep.BlockHash, block, block.Size())
		//pub to specified peer
		pubPeerFunc(query, pid)
	}
}

func (n *Node) queryMempool(ty int64, data interface{}) (interface{}, error) {

	client := n.nodeInfo.client

	msg := client.NewMessage("mempool", ty, data)
	err := client.Send(msg, true)
	if err != nil {
		return nil, err
	}
	resp, err := client.WaitTimeout(msg, time.Second*10)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (n *Node) postBlockChain(blockHash, pid string, block *types.Block) error {
	return n.p2pMgr.PubBroadCast(blockHash, &types.BlockPid{Pid: pid, Block: block}, types.EventBroadcastAddBlock)
}

func (n *Node) postMempool(txHash string, tx *types.Transaction) error {
	return n.p2pMgr.PubBroadCast(txHash, tx, types.EventTx)
}

//        ,          (               ,  filter lru          )
func (n *Node) addIgnoreSendPeerAtomic(filter *utils.Filterdata, key, pid string) (exist bool) {

	filter.GetAtomicLock()
	defer filter.ReleaseAtomicLock()
	var info *sendFilterInfo
	if !filter.Contains(key) {
		info = &sendFilterInfo{ignoreSendPeers: make(map[string]bool)}
		filter.Add(key, info)
	} else {
		data, _ := filter.Get(key)
		info = data.(*sendFilterInfo)
	}
	_, exist = info.ignoreSendPeers[pid]
	info.ignoreSendPeers[pid] = true
	return exist
}

//
func (n *Node) removeIgnoreSendPeerAtomic(filter *utils.Filterdata, key, pid string) {

	filter.GetAtomicLock()
	defer filter.ReleaseAtomicLock()
	if filter.Contains(key) {
		data, _ := filter.Get(key)
		info := data.(*sendFilterInfo)
		delete(info.ignoreSendPeers, pid)
	}
}
