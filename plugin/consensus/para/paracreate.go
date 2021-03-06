// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package para

import (
	"time"

	"encoding/hex"

	"bytes"

	"sync/atomic"

	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	paraexec "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/executor"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/types"
)

type emptyBlockInterval struct {
	startHeight int64
	interval    int64
}

type downloadClient struct {
	emptyInterval []*emptyBlockInterval
}

func (client *client) createLocalGenesisBlock(genesis *types.Block) error {
	return client.alignLocalBlock2ChainBlock(genesis)
}

func getNewBlock(lastBlock *pt.ParaLocalDbBlock, txs []*types.Transaction, mainBlock *types.ParaTxDetail) *pt.ParaLocalDbBlock {
	var newblock pt.ParaLocalDbBlock

	newblock.Height = lastBlock.Height + 1
	newblock.MainHash = mainBlock.Header.Hash
	newblock.MainHeight = mainBlock.Header.Height
	newblock.ParentMainHash = lastBlock.MainHash
	newblock.BlockTime = mainBlock.Header.BlockTime
	newblock.Txs = txs

	return &newblock
}

func (client *client) createLocalBlock(lastBlock *pt.ParaLocalDbBlock, txs []*types.Transaction, mainBlock *types.ParaTxDetail) error {
	err := client.addLocalBlock(getNewBlock(lastBlock, txs, mainBlock))
	if err != nil {
		return err
	}
	return err
}

func (client *client) getLocalBlockSeq(height int64) (int64, []byte, error) {
	lastBlock, err := client.getLocalBlockByHeight(height)
	if err != nil {
		return -2, nil, err
	}

	//    mainHash  seq    ，  0 seq，   hash， switchLocalHashMatchedBlock
	mainSeq, err := client.GetSeqByHashOnMainChain(lastBlock.MainHash)
	if err != nil {
		return 0, lastBlock.MainHash, nil
	}
	return mainSeq, lastBlock.MainHash, nil

}

//      chainblock，    localdb block
func (client *client) alignLocalBlock2ChainBlock(chainBlock *types.Block) error {
	localBlock := &pt.ParaLocalDbBlock{
		Height:     chainBlock.Height,
		MainHeight: chainBlock.MainHeight,
		MainHash:   chainBlock.MainHash,
		BlockTime:  chainBlock.BlockTime,
	}

	return client.addLocalBlock(localBlock)

}

//  localBlock    ，  chain block  ，  block    ，   seq    ，        hash    seq=0,
//
func (client *client) getLastLocalBlockSeq() (int64, []byte, error) {
	height, err := client.getLastLocalHeight()
	if err == nil {
		mainSeq, mainHash, err := client.getLocalBlockSeq(height)
		if err == nil {
			return mainSeq, mainHash, nil
		}
	}

	plog.Info("Parachain getLastLocalBlockSeq from block")
	//  localDb      ， chain
	mainSeq, chainBlock, err := client.getLastBlockMainInfo()
	if err != nil {
		return -2, nil, err
	}

	//chain block     ，  last local block    chainBlock main   mainhash
	err = client.alignLocalBlock2ChainBlock(chainBlock)
	if err != nil {
		return -2, nil, err
	}
	return mainSeq, chainBlock.MainHash, nil

}

func (client *client) getLastLocalBlock() (*pt.ParaLocalDbBlock, error) {
	height, err := client.getLastLocalHeight()
	if err != nil {
		return nil, err
	}

	return client.getLocalBlockByHeight(height)
}

//genesis block scenario
func (client *client) syncFromGenesisBlock() (int64, *types.Block, error) {
	lastSeq, lastBlock, err := client.getLastBlockMainInfo()
	if err != nil {
		plog.Error("Parachain getLastBlockInfo fail", "err", err)
		return -2, nil, err
	}
	plog.Info("syncFromGenesisBlock sync from height 0")
	return lastSeq, lastBlock, nil
}

func (client *client) getMatchedBlockOnChain(startHeight int64) (int64, *types.Block, error) {
	lastBlock, err := client.RequestLastBlock()
	if err != nil {
		plog.Error("Parachain RequestLastBlock fail", "err", err)
		return -2, nil, err
	}

	if lastBlock.Height == 0 {
		return client.syncFromGenesisBlock()
	}

	if startHeight == 0 || startHeight > lastBlock.Height {
		startHeight = lastBlock.Height
	}

	depth := defaultSearchMatchedBlockDepth
	for height := startHeight; height > 0 && depth > 0; height-- {
		block, err := client.GetBlockByHeight(height)
		if err != nil {
			return -2, nil, err
		}
		//  block     mainHash MainHeight   blockchain   block     ，       ，     minerTx
		plog.Info("switchHashMatchedBlock", "lastParaBlockHeight", height, "mainHeight",
			block.MainHeight, "mainHash", hex.EncodeToString(block.MainHash))
		mainSeq, err := client.GetSeqByHashOnMainChain(block.MainHash)
		if err != nil {
			depth--
			if depth == 0 {
				plog.Error("switchHashMatchedBlock depth overflow", "last info:mainHeight", block.MainHeight,
					"mainHash", hex.EncodeToString(block.MainHash), "search startHeight", lastBlock.Height, "curHeight", height,
					"search depth", defaultSearchMatchedBlockDepth)
				panic("search HashMatchedBlock overflow, re-setting search depth and restart to try")
			}
			if height == 1 {
				plog.Error("switchHashMatchedBlock search to height=1 not found", "lastBlockHeight", lastBlock.Height,
					"height1 mainHash", hex.EncodeToString(block.MainHash))
				return client.syncFromGenesisBlock()

			}
			continue
		}

		plog.Info("getMatchedBlockOnChain succ", "currHeight", height, "initHeight", lastBlock.Height,
			"new currSeq", mainSeq, "new preMainBlockHash", hex.EncodeToString(block.MainHash))
		return mainSeq, block, nil
	}
	return -2, nil, pt.ErrParaCurHashNotMatch
}

func (client *client) switchMatchedBlockOnChain(startHeight int64) (int64, []byte, error) {
	mainSeq, chainBlock, err := client.getMatchedBlockOnChain(startHeight)
	if err != nil {
		return -2, nil, err
	}
	//chain block     ，  last local block    chainBlock main   mainhash
	err = client.alignLocalBlock2ChainBlock(chainBlock)
	if err != nil {
		return -2, nil, err
	}
	return mainSeq, chainBlock.MainHash, nil
}

func (client *client) switchHashMatchedBlock() (int64, []byte, error) {
	mainSeq, localBlock, err := client.switchLocalHashMatchedBlock()
	if err != nil {
		return client.switchMatchedBlockOnChain(0)
	}
	return mainSeq, localBlock.MainHash, nil
}

//
func (client *client) switchLocalHashMatchedBlock() (int64, *pt.ParaLocalDbBlock, error) {
	lastBlock, err := client.getLastLocalBlock()
	if err != nil {
		plog.Error("Parachain RequestLastBlock fail", "err", err)
		return -2, nil, err
	}

	for height := lastBlock.Height; height >= 0; height-- {
		block, err := client.getLocalBlockByHeight(height)
		if err != nil {
			return -2, nil, err
		}
		//  block     mainHash MainHeight   blockchain   block     ，       ，     minerTx
		plog.Info("switchLocalHashMatchedBlock", "height", height, "mainHeight", block.MainHeight, "mainHash", hex.EncodeToString(block.MainHash))
		mainHash, err := client.GetHashByHeightOnMainChain(block.MainHeight)
		if err != nil || !bytes.Equal(mainHash, block.MainHash) {
			continue
		}

		mainSeq, err := client.GetSeqByHashOnMainChain(block.MainHash)
		if err != nil {
			continue
		}

		//remove fail, the para chain may be remove part, set the preMainBlockHash to nil, to match nothing, force to search from last
		err = client.removeLocalBlocks(height)
		if err != nil {
			return -2, nil, err
		}

		plog.Info("switchLocalHashMatchedBlock succ", "currHeight", height, "initHeight", lastBlock.Height,
			"currSeq", mainSeq, "mainHeight", block.MainHeight, "currMainBlockHash", hex.EncodeToString(block.MainHash))
		return mainSeq, block, nil
	}
	return -2, nil, pt.ErrParaCurHashNotMatch
}

func (client *client) getBatchSeqCount(currSeq int64) (int64, error) {
	lastSeq, err := client.GetLastSeqOnMainChain()
	if err != nil {
		return 0, err
	}

	if lastSeq > currSeq {
		if lastSeq-currSeq > client.dldCfg.emptyInterval[0].interval {
			atomic.StoreInt32(&client.caughtUp, 0)
		} else {
			atomic.StoreInt32(&client.caughtUp, 1)
		}
		if lastSeq-currSeq > client.subCfg.BatchFetchBlockCount {
			return client.subCfg.BatchFetchBlockCount, nil
		}
		return 1, nil
	}

	if lastSeq == currSeq {
		return 1, nil
	}

	// lastSeq = currSeq -1
	if lastSeq+1 == currSeq {
		plog.Debug("Waiting new sequence from main chain")
		return 0, pt.ErrParaWaitingNewSeq
	}

	// lastSeq < currSeq-1
	return 0, pt.ErrParaCurHashNotMatch

}

func getParentHash(block *types.ParaTxDetail) []byte {
	if block.Type == types.AddBlock {
		return block.Header.ParentHash
	}

	return block.Header.Hash
}

func getVerifyHash(block *types.ParaTxDetail) []byte {
	if block.Type == types.AddBlock {
		return block.Header.Hash
	}

	return block.Header.ParentHash
}
func verifyMainBlockHash(preMainBlockHash []byte, mainBlock *types.ParaTxDetail) error {
	if bytes.Equal(preMainBlockHash, getParentHash(mainBlock)) {
		return nil
	}
	plog.Error("verifyMainBlockHash", "preMainBlockHash", hex.EncodeToString(preMainBlockHash),
		"mainParentHash", hex.EncodeToString(mainBlock.Header.ParentHash), "mainHash", hex.EncodeToString(mainBlock.Header.Hash),
		"type", mainBlock.Type, "height", mainBlock.Header.Height)
	return pt.ErrParaCurHashNotMatch
}

func verifyMainBlocks(preMainBlockHash []byte, mainBlocks *types.ParaTxDetails) error {
	pre := preMainBlockHash
	for _, block := range mainBlocks.Items {
		err := verifyMainBlockHash(pre, block)
		if err != nil {
			return err
		}
		pre = getVerifyHash(block)
	}
	return nil
}

func verifyMainBlocksInternal(mainBlocks *types.ParaTxDetails) error {
	return verifyMainBlocks(getParentHash(mainBlocks.Items[0]), mainBlocks)
}

func isValidSeqType(ty int64) bool {
	return ty == types.AddBlock || ty == types.DelBlock
}

func validMainBlocks(txs *types.ParaTxDetails) *types.ParaTxDetails {
	for i, item := range txs.Items {
		if item == nil || !isValidSeqType(item.Type) {
			txs.Items = txs.Items[:i]
			return txs
		}
	}
	return txs
}

func (client *client) requestTxsFromBlock(currSeq int64, preMainBlockHash []byte) (*types.ParaTxDetails, error) {
	cfg := client.GetAPI().GetConfig()
	blockSeq, err := client.GetBlockOnMainBySeq(currSeq)
	if err != nil {
		return nil, err
	}

	txDetail := blockSeq.Detail.FilterParaTxsByTitle(cfg, cfg.GetTitle())
	txDetail.Type = blockSeq.Seq.Type

	if !isValidSeqType(txDetail.Type) {
		return nil, types.ErrInvalidParam
	}

	err = verifyMainBlockHash(preMainBlockHash, txDetail)
	if err != nil {
		plog.Error("requestTxsFromBlock", "curr seq", currSeq, "preMainBlockHash", hex.EncodeToString(preMainBlockHash))
		return nil, err
	}
	return &types.ParaTxDetails{Items: []*types.ParaTxDetail{txDetail}}, nil
}

func (client *client) requestFilterParaTxs(currSeq int64, count int64, preMainBlockHash []byte) (*types.ParaTxDetails, error) {
	cfg := client.GetAPI().GetConfig()
	req := &types.ReqParaTxByTitle{IsSeq: true, Start: currSeq, End: currSeq + count - 1, Title: cfg.GetTitle()}
	details, err := client.GetParaTxByTitle(req)
	if err != nil {
		return nil, err
	}

	details = validMainBlocks(details)
	err = verifyMainBlocks(preMainBlockHash, details)
	if err != nil {
		plog.Error("requestFilterParaTxs", "curSeq", currSeq, "count", count, "preMainBlockHash", hex.EncodeToString(preMainBlockHash))
		return nil, err
	}
	//      １
	if len(details.Items) == 0 {
		plog.Error("requestFilterParaTxs ret nil", "curSeq", currSeq, "count", count, "preMainBlockHash", hex.EncodeToString(preMainBlockHash))
		return nil, types.ErrNotFound
	}

	return details, nil
}

func (client *client) RequestTx(currSeq int64, count int64, preMainBlockHash []byte) (*types.ParaTxDetails, error) {
	return client.requestFilterParaTxs(currSeq, count, preMainBlockHash)
}

func (client *client) processHashNotMatchError(currSeq int64, lastSeqMainHash []byte, err error) (int64, []byte, error) {
	if err == pt.ErrParaCurHashNotMatch {
		preSeq, preSeqMainHash, err := client.switchHashMatchedBlock()
		if err == nil {
			return preSeq + 1, preSeqMainHash, nil
		}
	}
	return currSeq, lastSeqMainHash, err
}

func (client *client) getEmptyInterval(lastBlock *pt.ParaLocalDbBlock) int64 {
	for i := len(client.dldCfg.emptyInterval) - 1; i >= 0; i-- {
		if lastBlock.Height >= client.dldCfg.emptyInterval[i].startHeight {
			return client.dldCfg.emptyInterval[i].interval
		}
	}
	panic(fmt.Sprintf("emptyBlockInterval not set for height=%d", lastBlock.Height))
}

func (client *client) procLocalBlock(mainBlock *types.ParaTxDetail) (bool, error) {
	cfg := client.GetAPI().GetConfig()
	lastSeqMainHeight := mainBlock.Header.Height

	lastBlock, err := client.getLastLocalBlock()
	if err != nil {
		plog.Error("Parachain getLastLocalBlock", "err", err)
		return false, err
	}
	emptyInterval := client.getEmptyInterval(lastBlock)

	txs := paraexec.FilterTxsForPara(cfg, mainBlock)

	plog.Info("Parachain process block", "lastBlockHeight", lastBlock.Height, "lastBlockMainHeight", lastBlock.MainHeight,
		"lastBlockMainHash", common.ToHex(lastBlock.MainHash), "currMainHeight", lastSeqMainHeight,
		"curMainHash", common.ToHex(mainBlock.Header.Hash), "emptyIntval", emptyInterval, "seqTy", mainBlock.Type)

	if mainBlock.Type == types.DelBlock {
		if len(txs) == 0 {
			if lastSeqMainHeight > lastBlock.MainHeight {
				return false, nil
			}
			plog.Info("Delete empty block", "height", lastBlock.Height)
		}
		return true, client.delLocalBlock(lastBlock.Height)

	}
	//AddAct
	if len(txs) == 0 {
		if lastSeqMainHeight-lastBlock.MainHeight < emptyInterval {
			return false, nil
		}
		plog.Info("Create empty block", "newHeight", lastBlock.Height+1)
	}
	return true, client.createLocalBlock(lastBlock, txs, mainBlock)
}

func (client *client) procLocalBlocks(mainBlocks *types.ParaTxDetails) error {
	var notify bool
	for _, main := range mainBlocks.Items {
		changed, err := client.procLocalBlock(main)
		if err != nil {
			return err
		}
		if changed {
			notify = true
		}
	}
	if notify {
		client.blockSyncClient.handleLocalChangedMsg()
	}

	return nil
}

func (client *client) procLocalAddBlock(mainBlock *types.ParaTxDetail, lastBlock *pt.ParaLocalDbBlock) *pt.ParaLocalDbBlock {
	cfg := client.GetAPI().GetConfig()
	curMainHeight := mainBlock.Header.Height

	emptyInterval := client.getEmptyInterval(lastBlock)

	txs := paraexec.FilterTxsForPara(cfg, mainBlock)

	plog.Debug("Parachain process block", "lastBlockHeight", lastBlock.Height, "lastBlockMainHeight", lastBlock.MainHeight,
		"lastBlockMainHash", common.ToHex(lastBlock.MainHash), "currMainHeight", curMainHeight,
		"curMainHash", common.ToHex(mainBlock.Header.Hash), "emptyIntval", emptyInterval, "seqTy", mainBlock.Type)

	if mainBlock.Type != types.AddBlock {
		panic("para chain quick sync,not addBlock type")

	}
	//AddAct
	if len(txs) == 0 {
		if curMainHeight-lastBlock.MainHeight < emptyInterval {
			return nil
		}
		plog.Debug("Create empty block", "newHeight", lastBlock.Height+1)
	}
	return getNewBlock(lastBlock, txs, mainBlock)

}

//     AddType block，       １w      blocks，      ，   addType   ，     ，
func (client *client) procLocalAddBlocks(mainBlocks *types.ParaTxDetails) error {
	var blocks []*pt.ParaLocalDbBlock
	lastBlock, err := client.getLastLocalBlock()
	if err != nil {
		plog.Error("procLocalAddBlocks getLastLocalBlock", "err", err)
		return err
	}

	for _, main := range mainBlocks.Items {
		b := client.procLocalAddBlock(main, lastBlock)
		if b == nil {
			continue
		}
		lastBlock = b
		blocks = append(blocks, b)
	}
	if len(blocks) <= 0 {
		return nil
	}
	err = client.saveBatchLocalBlocks(blocks)
	if err != nil {
		plog.Error("procLocalAddBlocks saveBatchLocalBlocks", "err", err)
		panic(err)
	}
	plog.Info("procLocalAddBlocks.saveLocalBlocks", "start", blocks[0].Height, "end", blocks[len(blocks)-1].Height)
	client.blockSyncClient.handleLocalChangedMsg()
	return nil
}

func (client *client) CreateBlock() {
	defer client.wg.Done()

	if !client.subCfg.JumpDownloadClose {
		client.jumpDldCli.tryJumpDownload()
	}

	if client.subCfg.MultiDownloadOpen {
		client.multiDldCli.tryMultiServerDownload()
	}

	lastSeq, lastSeqMainHash, err := client.getLastLocalBlockSeq()
	if err != nil {
		plog.Error("Parachain CreateBlock getLastLocalBlockSeq fail", "err", err.Error())
		return
	}
	currSeq := lastSeq + 1

out:
	for {
		select {
		case <-client.quit:
			break out
		default:
			count, err := client.getBatchSeqCount(currSeq)
			if err != nil {
				currSeq, lastSeqMainHash, err = client.processHashNotMatchError(currSeq, lastSeqMainHash, err)
				if err == nil {
					continue
				}
				time.Sleep(time.Second * time.Duration(client.subCfg.WriteBlockSeconds))
				continue
			}

			plog.Debug("Parachain CreateBlock", "curSeq", currSeq, "count", count, "lastSeqMainHash", common.ToHex(lastSeqMainHash))
			paraTxs, err := client.RequestTx(currSeq, count, lastSeqMainHash)
			if err != nil {
				currSeq, lastSeqMainHash, err = client.processHashNotMatchError(currSeq, lastSeqMainHash, err)
				continue
			}

			if count != int64(len(paraTxs.Items)) {
				plog.Debug("para CreateBlock count not match", "count", count, "items", len(paraTxs.Items))
				count = int64(len(paraTxs.Items))
			}
			//        ，
			if client.commitMsgClient.authAccount != "" && client.isCaughtUp() && len(paraTxs.Items) > 0 {
				//      ，  seq     ，
				client.commitMsgClient.commitTxCheckNotify(paraTxs.Items[0])
			}

			err = client.procLocalBlocks(paraTxs)
			if err != nil {
				//  localblock，
				lastSeqMainHash = nil
				plog.Error("para CreateBlock.procLocalBlocks", "err", err.Error())
				continue
			}

			//    seq lastSeqMainHash
			lastSeqMainHash = paraTxs.Items[count-1].Header.Hash
			if paraTxs.Items[count-1].Type == types.DelBlock {
				lastSeqMainHash = paraTxs.Items[count-1].Header.ParentHash
			}
			currSeq = currSeq + count

		}
	}

	plog.Info("para CreateBlock quit")
}
