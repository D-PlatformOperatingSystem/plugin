// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package para

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/merkle"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/types"
)

const (
	//defaultMaxCacheCount   local
	defaultMaxCacheCount = int64(1000)
	//defaultMaxSyncErrCount
	defaultMaxSyncErrCount = int32(100)
)

//blockSyncClient
type blockSyncClient struct {
	paraClient *client
	//notifyChan
	notifyChan chan bool
	//quitChan
	quitChan chan struct{}
	//syncState
	syncState int32
	//syncErrMaxCount
	maxSyncErrCount int32
	//maxCacheCount
	maxCacheCount int64
	//isSyncCaughtUp
	isSyncCaughtUpAtom int32
	//isDownloadCaughtUpAtom
	isDownloadCaughtUpAtom int32
	//isSyncFirstCaughtUp      download   sync           ，                ，
	isSyncFirstCaughtUp bool
}

//nextActionType
type nextActionType int8

const (
	//nextActionKeep
	nextActionKeep nextActionType = iota
	//nextActionRollback
	nextActionRollback
	//nextActionAdd
	nextActionAdd
)

//blockSyncState
type blockSyncState int32

const (
	//blockSyncStateNone
	blockSyncStateNone blockSyncState = iota
	//blockSyncStateSyncing
	blockSyncStateSyncing
	//blockSyncStateFinished
	blockSyncStateFinished
)

func newBlockSyncCli(para *client, cfg *subConfig) *blockSyncClient {
	cli := &blockSyncClient{
		paraClient:      para,
		notifyChan:      make(chan bool, 1),
		quitChan:        make(chan struct{}),
		maxCacheCount:   defaultMaxCacheCount,
		maxSyncErrCount: defaultMaxSyncErrCount,
	}
	if cfg.MaxCacheCount > 0 {
		cli.maxCacheCount = cfg.MaxCacheCount
	}
	if cfg.MaxSyncErrCount > 0 {
		cli.maxSyncErrCount = cfg.MaxSyncErrCount
	}
	return cli
}

//syncHasCaughtUp           ，
func (client *blockSyncClient) syncHasCaughtUp() bool {
	return atomic.LoadInt32(&client.isSyncCaughtUpAtom) == 1
}

//handleLocalChangedMsg         ，
func (client *blockSyncClient) handleLocalChangedMsg() {
	client.printDebugInfo("Para sync - notify change")
	if client.getBlockSyncState() == blockSyncStateSyncing || client.paraClient.isCancel() {
		return
	}
	client.printDebugInfo("Para sync - notified change")
	client.notifyChan <- true
}

//handleLocalCaughtUpMsg           ，
func (client *blockSyncClient) handleLocalCaughtUpMsg() {
	client.printDebugInfo("Para sync -notify download has caughtUp")
	if !client.downloadHasCaughtUp() {
		client.setDownloadHasCaughtUp(true)
	}
}

//createGenesisBlock
func (client *blockSyncClient) createGenesisBlock(newblock *types.Block) error {
	return client.writeBlock(zeroHash[:], newblock)
}

//syncBlocks
//
func (client *blockSyncClient) syncBlocks() {

	client.syncInit()
	//    ，
	client.batchSyncBlocks()
	//      ,
out:
	for {
		select {
		case <-client.notifyChan:

			client.batchSyncBlocks()

		case <-client.quitChan:
			break out
		}
	}

	plog.Info("Para sync - quit block sync goroutine")
	client.paraClient.wg.Done()
}

//
func (client *blockSyncClient) batchSyncBlocks() {
	client.setBlockSyncState(blockSyncStateSyncing)
	client.printDebugInfo("Para sync - syncing")

	errCount := int32(0)
	for {
		if client.paraClient.isCancel() {
			return
		}
		//      ,
		curSyncCaughtState, err := client.syncBlocksIfNeed()

		//
		if err != nil {
			errCount++
			client.printError(err)
		} else {
			errCount = int32(0)
		}
		//          ,    ，
		if errCount > client.maxSyncErrCount {
			client.printError(errors.New(
				"para sync - sync has some errors,please check"))
			client.setBlockSyncState(blockSyncStateNone)
			return
		}
		//        ,        localCacheCount
		if err == nil && curSyncCaughtState {
			_, err := client.clearLocalOldBlocks()
			if err != nil {
				client.printError(err)
			}

			client.setBlockSyncState(blockSyncStateFinished)
			client.printDebugInfo("Para sync - finished")
			return
		}
	}

}

//
func (client *blockSyncClient) getNextAction() (nextActionType, *types.Block, *pt.ParaLocalDbBlock, int64, error) {
	lastBlock, err := client.paraClient.getLastBlockInfo()
	if err != nil {
		//            ，
		return nextActionKeep, nil, nil, -1, err
	}

	lastLocalHeight, err := client.paraClient.getLastLocalHeight()
	if err != nil {
		// db           ，
		return nextActionKeep, nil, nil, lastLocalHeight, err
	}

	if lastLocalHeight <= 0 {
		//db      0,      （    ）
		return nextActionKeep, nil, nil, lastLocalHeight, nil
	}

	switch {
	case lastLocalHeight < lastBlock.Height:
		//db                  ,
		return nextActionRollback, lastBlock, nil, lastLocalHeight, nil
	case lastLocalHeight == lastBlock.Height:
		localBlock, err := client.paraClient.getLocalBlockByHeight(lastBlock.Height)
		if err != nil {
			// db           ，
			return nextActionKeep, nil, nil, lastLocalHeight, err
		}
		if common.ToHex(localBlock.MainHash) == common.ToHex(lastBlock.MainHash) {
			//db                    hash  ,      (       )
			return nextActionKeep, nil, nil, lastLocalHeight, nil
		}
		//db                    hash  ,
		return nextActionRollback, lastBlock, nil, lastLocalHeight, nil
	default:
		// lastLocalHeight > lastBlock.Height
		localBlock, err := client.paraClient.getLocalBlockByHeight(lastBlock.Height + 1)
		if err != nil {
			// db           ，
			return nextActionKeep, nil, nil, lastLocalHeight, err
		}
		if common.ToHex(localBlock.ParentMainHash) != common.ToHex(lastBlock.MainHash) {
			//db         hash           hash,
			return nextActionRollback, lastBlock, nil, lastLocalHeight, nil
		}
		//db         hash          hash,
		return nextActionAdd, lastBlock, localBlock, lastLocalHeight, nil
	}
}

//
//
//bool
func (client *blockSyncClient) syncBlocksIfNeed() (bool, error) {
	nextAction, lastBlock, localBlock, lastLocalHeight, err := client.getNextAction()
	if err != nil {
		return false, err
	}

	switch nextAction {
	case nextActionAdd:
		//1 db         hash          hash
		plog.Info("Para sync -    add block",
			"lastBlock.Height", lastBlock.Height, "lastLocalHeight", lastLocalHeight)

		err := client.addBlock(lastBlock, localBlock)

		//
		if err == nil {
			isSyncCaughtUp := lastBlock.Height+1 == lastLocalHeight
			client.setSyncCaughtUp(isSyncCaughtUp)
			if client.paraClient.commitMsgClient.authAccount != "" {
				client.printDebugInfo("Para sync - add block commit", "isSyncCaughtUp", isSyncCaughtUp)
				client.paraClient.commitMsgClient.updateChainHeightNotify(lastBlock.Height+1, false)
			}
		}

		return false, err
	case nextActionRollback:
		//1 db
		//2 db                    hash
		//3 db         hash           hash
		plog.Info("Para sync -    rollback block",
			"lastBlock.Height", lastBlock.Height, "lastLocalHeight", lastLocalHeight)

		err := client.rollbackBlock(lastBlock)

		//
		if err == nil {
			client.setSyncCaughtUp(false)
			if client.paraClient.commitMsgClient.authAccount != "" {
				client.printDebugInfo("Para sync - rollback block commit", "isSyncCaughtUp", false)
				client.paraClient.commitMsgClient.updateChainHeightNotify(lastBlock.Height-1, true)
			}
		}

		return false, err
	default: //nextActionKeep
		//1      ，
		return true, nil
	}

}

//
func (client *blockSyncClient) delLocalBlocks(startHeight int64, endHeight int64) error {
	if startHeight > endHeight {
		return errors.New("para sync - startHeight > endHeight,can't clear local blocks")
	}

	index := startHeight
	set := &types.LocalDBSet{}
	cfg := client.paraClient.GetAPI().GetConfig()
	for {
		if index > endHeight {
			break
		}

		key := calcTitleHeightKey(cfg.GetTitle(), index)
		kv := &types.KeyValue{Key: key, Value: nil}
		set.KV = append(set.KV, kv)

		index++
	}

	key := calcTitleFirstHeightKey(cfg.GetTitle())
	kv := &types.KeyValue{Key: key, Value: types.Encode(&types.Int64{Data: endHeight + 1})}
	set.KV = append(set.KV, kv)

	client.printDebugInfo("Para sync - clear local blocks", "startHeight:", startHeight, "endHeight:", endHeight)

	return client.paraClient.setLocalDb(set)
}

//
func (client *blockSyncClient) initFirstLocalHeightIfNeed() error {
	height, err := client.getFirstLocalHeight()
	cfg := client.paraClient.GetAPI().GetConfig()
	if err != nil || height < 0 {
		set := &types.LocalDBSet{}
		key := calcTitleFirstHeightKey(cfg.GetTitle())
		kv := &types.KeyValue{Key: key, Value: types.Encode(&types.Int64{Data: 0})}
		set.KV = append(set.KV, kv)

		return client.paraClient.setLocalDb(set)
	}

	return err
}

//
func (client *blockSyncClient) getFirstLocalHeight() (int64, error) {
	cfg := client.paraClient.GetAPI().GetConfig()
	key := calcTitleFirstHeightKey(cfg.GetTitle())
	set := &types.LocalDBGet{Keys: [][]byte{key}}
	value, err := client.paraClient.getLocalDb(set, len(set.Keys))
	if err != nil {
		return -1, err
	}

	if len(value) == 0 {
		return -1, types.ErrNotFound
	}

	if value[0] == nil {
		return -1, types.ErrNotFound
	}

	height := &types.Int64{}
	err = types.Decode(value[0], height)
	if err != nil {
		return -1, err
	}
	return height.Data, nil
}

//      (localCacheCount)
func (client *blockSyncClient) clearLocalOldBlocks() (bool, error) {
	lastLocalHeight, err := client.paraClient.getLastLocalHeight()
	if err != nil {
		return false, err
	}

	firstLocalHeight, err := client.getFirstLocalHeight()
	if err != nil {
		return false, err
	}

	canDelCount := lastLocalHeight - firstLocalHeight - client.maxCacheCount + 1
	if canDelCount <= client.maxCacheCount {
		return false, nil
	}

	return true, client.delLocalBlocks(firstLocalHeight, firstLocalHeight+canDelCount-1)
}

// miner tx need all para node create, but not all node has auth account, here just not sign to keep align
func (client *blockSyncClient) addMinerTx(preStateHash []byte, block *types.Block, localBlock *pt.ParaLocalDbBlock) error {
	cfg := client.paraClient.GetAPI().GetConfig()
	status := &pt.ParacrossNodeStatus{
		Title:           cfg.GetTitle(),
		Height:          block.Height,
		MainBlockHash:   localBlock.MainHash,
		MainBlockHeight: localBlock.MainHeight,
	}

	maxHeight := pt.GetDappForkHeight(cfg, pt.ForkLoopCheckCommitTxDone)
	if maxHeight < client.paraClient.subCfg.RmCommitParamMainHeight {
		maxHeight = client.paraClient.subCfg.RmCommitParamMainHeight
	}
	if status.MainBlockHeight < maxHeight {
		status.PreBlockHash = block.ParentHash
		status.PreStateHash = preStateHash
	}

	//selfConsensEnablePreContract  ForkParaSelfConsStages            ，fork        ，
	tx, err := pt.CreateRawMinerTx(cfg, &pt.ParacrossMinerAction{
		Status:          status,
		IsSelfConsensus: client.paraClient.commitMsgClient.isSelfConsEnable(status.Height),
	})
	if err != nil {
		return err
	}

	tx.Sign(types.SECP256K1, client.paraClient.minerPrivateKey)
	block.Txs = append([]*types.Transaction{tx}, block.Txs...)

	return nil
}

//
func (client *blockSyncClient) addBlock(lastBlock *types.Block, localBlock *pt.ParaLocalDbBlock) error {
	cfg := client.paraClient.GetAPI().GetConfig()
	var newBlock types.Block
	newBlock.ParentHash = lastBlock.Hash(cfg)
	newBlock.Height = lastBlock.Height + 1
	newBlock.Txs = localBlock.Txs
	err := client.addMinerTx(lastBlock.StateHash, &newBlock, localBlock)
	if err != nil {
		return err
	}
	//
	newBlock.Difficulty = cfg.GetP(0).PowLimitBits

	//                TxHash
	if cfg.IsFork(newBlock.GetMainHeight(), "ForkRootHash") {
		newBlock.Txs = types.TransactionSort(newBlock.Txs)
	}
	newBlock.TxHash = merkle.CalcMerkleRoot(cfg, newBlock.GetMainHeight(), newBlock.Txs)
	newBlock.BlockTime = localBlock.BlockTime
	newBlock.MainHash = localBlock.MainHash
	newBlock.MainHeight = localBlock.MainHeight
	if newBlock.Height == 1 && newBlock.BlockTime < client.paraClient.cfg.GenesisBlockTime {
		panic("genesisBlockTime　bigger than the 1st block time, need rmv db and reset genesisBlockTime")
	}
	err = client.writeBlock(lastBlock.StateHash, &newBlock)

	client.printDebugInfo("Para sync - create new Block",
		"newblock.ParentHash", common.ToHex(newBlock.ParentHash),
		"newblock.Height", newBlock.Height, "newblock.TxHash", common.ToHex(newBlock.TxHash),
		"newblock.BlockTime", newBlock.BlockTime)

	return err
}

//  blockchain
func (client *blockSyncClient) rollbackBlock(block *types.Block) error {
	start := block.Height
	if start <= 0 {
		return errors.New("para sync - attempt to rollbackBlock genesisBlock")
	}

	reqBlocks := &types.ReqBlocks{Start: start, End: start, IsDetail: true, Pid: []string{""}}
	msg := client.paraClient.GetQueueClient().NewMessage("blockchain", types.EventGetBlocks, reqBlocks)
	err := client.paraClient.GetQueueClient().Send(msg, true)
	if err != nil {
		return err
	}

	resp, err := client.paraClient.GetQueueClient().Wait(msg)
	if err != nil {
		return err
	}

	blocks := resp.GetData().(*types.BlockDetails)
	if len(blocks.Items) == 0 {
		return errors.New("para sync -blocks Items len = 0 ")
	}

	paraBlockDetail := &types.ParaChainBlockDetail{Blockdetail: blocks.Items[0]}
	msg = client.paraClient.GetQueueClient().NewMessage("blockchain", types.EventDelParaChainBlockDetail, paraBlockDetail)
	err = client.paraClient.GetQueueClient().Send(msg, true)
	if err != nil {
		return err
	}

	resp, err = client.paraClient.GetQueueClient().Wait(msg)
	if err != nil {
		return err
	}

	if !resp.GetData().(*types.Reply).IsOk {
		reply := resp.GetData().(*types.Reply)
		return errors.New(string(reply.GetMsg()))
	}

	return nil
}

//  blockchain
func (client *blockSyncClient) writeBlock(prev []byte, paraBlock *types.Block) error {
	//       block，   blockchain    block       ，      blockdetail
	blockDetail := &types.BlockDetail{Block: paraBlock}
	//database    ，     ，      ，     download   sync             ，               ，     ，
	if !client.isSyncFirstCaughtUp && client.downloadHasCaughtUp() && client.syncHasCaughtUp() {
		client.isSyncFirstCaughtUp = true
		plog.Info("Para sync - SyncFirstCaughtUp", "Height", paraBlock.Height)
	}

	paraBlockDetail := &types.ParaChainBlockDetail{Blockdetail: blockDetail, IsSync: client.isSyncFirstCaughtUp}
	msg := client.paraClient.GetQueueClient().NewMessage("blockchain", types.EventAddParaChainBlockDetail, paraBlockDetail)
	err := client.paraClient.GetQueueClient().Send(msg, true)
	if err != nil {
		return err
	}

	resp, err := client.paraClient.GetQueueClient().Wait(msg)
	if err != nil {
		return err
	}

	respBlockDetail := resp.GetData().(*types.BlockDetail)
	if respBlockDetail == nil {
		return errors.New("para sync - block detail is nil")
	}

	client.paraClient.SetCurrentBlock(respBlockDetail.Block)

	return nil
}

//
func (client *blockSyncClient) getBlockSyncState() blockSyncState {
	return blockSyncState(atomic.LoadInt32(&client.syncState))
}

//
func (client *blockSyncClient) setBlockSyncState(state blockSyncState) {
	atomic.StoreInt32(&client.syncState, int32(state))
}

//
func (client *blockSyncClient) setSyncCaughtUp(isCaughtUp bool) {
	if isCaughtUp {
		atomic.StoreInt32(&client.isSyncCaughtUpAtom, 1)
	} else {
		atomic.StoreInt32(&client.isSyncCaughtUpAtom, 0)
	}
}

//
func (client *blockSyncClient) downloadHasCaughtUp() bool {
	return atomic.LoadInt32(&client.isDownloadCaughtUpAtom) == 1
}

//
func (client *blockSyncClient) setDownloadHasCaughtUp(isCaughtUp bool) {
	if isCaughtUp {
		atomic.CompareAndSwapInt32(&client.isDownloadCaughtUpAtom, 0, 1)
	} else {
		atomic.CompareAndSwapInt32(&client.isDownloadCaughtUpAtom, 1, 0)
	}
}

//
func (client *blockSyncClient) printError(err error) {
	plog.Error(fmt.Sprintf("Para sync - sync block error:%v", err.Error()))
}

//
func (client *blockSyncClient) printDebugInfo(msg string, ctx ...interface{}) {
	plog.Debug(msg, ctx...)
}

//
func (client *blockSyncClient) syncInit() {
	client.printDebugInfo("Para sync - init")
	client.setBlockSyncState(blockSyncStateNone)
	client.setSyncCaughtUp(false)
	client.setDownloadHasCaughtUp(false)
	err := client.initFirstLocalHeightIfNeed()
	if err != nil {
		client.printError(err)
	}

	//    chainHeight,      ，
	lastBlock, err := client.paraClient.getLastBlockInfo()
	if err != nil {
		//            ，
		plog.Info("Para sync init", "err", err)
	} else {
		client.paraClient.commitMsgClient.setInitChainHeight(lastBlock.Height)
	}

}
