// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package para

import (
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/consensus"
	"github.com/D-PlatformOperatingSystem/dpos/types"

	"encoding/hex"
	"sync/atomic"

	"github.com/D-PlatformOperatingSystem/dpos/queue"
	typesmocks "github.com/D-PlatformOperatingSystem/dpos/types/mocks"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/testnode"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	//TestPrivateKey
	TestPrivateKey = "6da92a632ab7deb67d38c0f6560bcfed28167998f6496db64c258d5e8393a81b"
	//TestBlockTime
	TestBlockTime = 1514533390
	//TestMaxCacheCount     DB
	TestMaxCacheCount = 100
	//TestLoopCount
	TestMaxLoopCount = 3
)

var (
	//testLoopCountAtom   queue   Message
	testLoopCountAtom int32
	//actionReturnIndexAtom   getNextAction
	actionReturnIndexAtom int32
)

//
func initTestSyncBlock() {
	//println("initSyncBlock")
}

//    para
func createParaTestInstance(t *testing.T, q queue.Queue) *client {
	para := new(client)
	para.subCfg = new(subConfig)

	baseCli := drivers.NewBaseClient(&types.Consensus{Name: "name"})
	para.BaseClient = baseCli

	para.InitClient(q.Client(), initTestSyncBlock)

	//  rpc Client
	grpcClient := &typesmocks.DplatformOSClient{}
	para.grpcClient = grpcClient

	//
	pk, err := hex.DecodeString(TestPrivateKey)
	assert.Nil(t, err)
	secp, err := crypto.New(types.GetSignName("", types.SECP256K1))
	assert.Nil(t, err)
	priKey, err := secp.PrivKeyFromBytes(pk)
	assert.Nil(t, err)
	para.minerPrivateKey = priKey

	//   BlockSyncClient
	para.blockSyncClient = &blockSyncClient{
		paraClient:      para,
		notifyChan:      make(chan bool),
		quitChan:        make(chan struct{}),
		maxCacheCount:   TestMaxCacheCount,
		maxSyncErrCount: 100,
	}

	para.commitMsgClient = &commitMsgClient{
		paraClient: para,
	}
	return para
}

//
func makeGenesisBlockInputTestData() *types.Block {
	newBlock := &types.Block{}
	newBlock.Height = 0
	newBlock.BlockTime = TestBlockTime
	newBlock.ParentHash = zeroHash[:]
	newBlock.MainHash = []byte("genesisHash")
	newBlock.MainHeight = 0

	return newBlock
}

//
func makeGenesisBlockReplyTestData(testLoopCount int32) interface{} {
	switch testLoopCount {
	case 0:
		return &types.BlockDetail{}
	default:
		return errors.New("error")
	}
}

//  getNextAction
//index   getNextAction        ， return
//testLoopCount
func makeSyncReplyTestData(index int32, testLoopCount int32) (
	interface{}, //*types.Block, //GetLastBlock reply
	interface{}, //*types.LocalReplyValue, //GetLastLocalHeight reply
	interface{}, //*types.LocalReplyValue, //GetLocalBlockByHeight reply
	interface{}, //*types.BlockDetail, //writeBlock reply
	interface{}, //*types.BlockDetails, //rollbackBlock  reply
	interface{}) { //*types.Reply) { //rollbackBlock  reply

	detail := &types.BlockDetail{Block: &types.Block{}}
	details := &types.BlockDetails{Items: []*types.BlockDetail{detail}}

	err := errors.New("error")

	switch index {
	case 1:
		return err, err, err, err, err, err
	case 2:
		return &types.Block{},
			err, err, err, err, err
	case 3:
		return &types.Block{},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 0})}},
			err, err, err, err
	case 4:
		return &types.Block{Height: 2},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 1})}},
			err, err,
			details,
			&types.Reply{IsOk: testLoopCount == 0}
	case 5:
		return &types.Block{Height: 2},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 2})}},
			err, err, err, err
	case 6:
		localBlock := &pt.ParaLocalDbBlock{MainHash: []byte("hash1"), Height: 2}
		return &types.Block{Height: 2, MainHash: []byte("hash1")},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 2})}},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(localBlock)}},
			err, err, err

	case 7:
		localBlock := &pt.ParaLocalDbBlock{MainHash: []byte("hash2"), Height: 2}
		return &types.Block{Height: 2, MainHash: []byte("hash1")},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 2})}},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(localBlock)}},
			err,
			details,
			&types.Reply{IsOk: testLoopCount == 0}

	case 8:
		return &types.Block{Height: 2, MainHash: []byte("hash1")},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 3})}},
			err, err, err, err
	case 9:
		localBlock := &pt.ParaLocalDbBlock{ParentMainHash: []byte("hash2"), Height: 3}
		return &types.Block{Height: 2, MainHash: []byte("hash1")},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 3})}},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(localBlock)}},
			err,
			details,
			&types.Reply{IsOk: testLoopCount == 0}
	case 10:
		localBlock := &pt.ParaLocalDbBlock{ParentMainHash: []byte("hash1"), Height: 3}
		return &types.Block{Height: 2, MainHash: []byte("hash1")},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 3})}},
			&types.LocalReplyValue{Values: [][]byte{types.Encode(localBlock)}},
			&types.BlockDetail{},
			err, err
	default:
		return err, err, err, err, err, err
	}
}

//      Get Reply
func makeCleanDataGetReplyTestData(clearLocalDBCallCount int32, testLoopCount int32) interface{} {
	switch clearLocalDBCallCount {
	case 1: //testinitFirstLocalHeightIfNeed
		switch testLoopCount {
		case 0:
			return &types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 1})}}
		default:
			return &types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: -1})}}
		}

	case 2: //testclearLocalOldBlocks
		switch testLoopCount {
		case 0:
			return &types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 1 + 2*TestMaxCacheCount})}}
		default:
			return &types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 1 + 2*TestMaxCacheCount - 50})}}
		}
	case 3: //testclearLocalOldBlocks
		switch testLoopCount {
		case 0:
			return &types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 1})}}
		case 1:
			return &types.LocalReplyValue{Values: [][]byte{types.Encode(&types.Int64{Data: 1})}}
		default: //2
			return errors.New("error")
		}
	default:
		return errors.New("error")
	}
}

//      Set Reply
func makeCleanDataSetReplyTestData(testLoopCount int32) interface{} {
	reply := &types.Reply{}
	reply.IsOk = testLoopCount == 0

	return reply
}

//mock queue Message
func mockMessageReply(q queue.Queue) {

	blockChainKey := "blockchain"
	cli := q.Client()
	cli.Sub(blockChainKey)
	//    Call  ,  loop  ；quitEndCount
	//quitCount := int32(0)
	//quitEndCount := int32(111) //TODO: Need a nice loop quit way
	//           EventGetValueByKey
	useLocalReply := false
	usrLocalReplyStart := true
	//           EventGetValueByKey
	clearLocalDBCallCount := int32(0)

	for msg := range cli.Recv() {

		testLoopCount := atomic.LoadInt32(&testLoopCountAtom)
		getActionReturnIndex := atomic.LoadInt32(&actionReturnIndexAtom)

		switch {
		case getActionReturnIndex > 0:
			//mock          ,testsyncBlocksIfNeed
			lastBlockReply,
				lastLocalReply,
				localReply,
				writeBlockReply,
				getBlocksReply,
				rollBlockReply := makeSyncReplyTestData(getActionReturnIndex, testLoopCount)

			switch msg.Ty {
			case types.EventGetLastBlock:
				//quitCount++

				msg.Reply(cli.NewMessage(blockChainKey, types.EventBlock, lastBlockReply))

			case types.EventAddParaChainBlockDetail:
				//quitCount++

				msg.Reply(cli.NewMessage(blockChainKey, types.EventReply, writeBlockReply))

			case types.EventDelParaChainBlockDetail:
				//quitCount++

				msg.Reply(cli.NewMessage(blockChainKey, types.EventReply, rollBlockReply))

			case types.EventGetValueByKey:
				//quitCount++

				switch {
				case getActionReturnIndex > 4:
					if usrLocalReplyStart {
						usrLocalReplyStart = false
						useLocalReply = false
					} else {
						useLocalReply = !useLocalReply
					}
				default:
					useLocalReply = false
					usrLocalReplyStart = true

				}

				if !useLocalReply {
					msg.Reply(cli.NewMessage(blockChainKey, types.EventLocalReplyValue, lastLocalReply))
				} else {
					msg.Reply(cli.NewMessage(blockChainKey, types.EventLocalReplyValue, localReply))
				}

			case types.EventGetBlocks:
				//quitCount++

				msg.Reply(cli.NewMessage(blockChainKey, types.EventBlocks, getBlocksReply))
			default:
				//nothing
			}
		default:
			switch msg.Ty {
			case types.EventAddParaChainBlockDetail: //mock          ,testCreateGenesisBlock
				//quitCount++

				reply := makeGenesisBlockReplyTestData(testLoopCount)
				msg.Reply(cli.NewMessage(blockChainKey, types.EventReply, reply))

			case types.EventGetValueByKey: //mock
				//quitCount++

				clearLocalDBCallCount++
				reply := makeCleanDataGetReplyTestData(clearLocalDBCallCount, testLoopCount)
				msg.Reply(cli.NewMessage(blockChainKey, types.EventLocalReplyValue, reply))
				if clearLocalDBCallCount == 3 {
					//      ，
					clearLocalDBCallCount = 0
				}

			case types.EventSetValueByKey: //mock        ,testclearLocalOldBlocks
				//quitCount++

				reply := makeCleanDataSetReplyTestData(testLoopCount)

				msg.Reply(cli.NewMessage(blockChainKey, types.EventReply, reply))
			default:
				//nothing
			}
		}

		//println(quitCount)
		//if quitCount == quitEndCount {
		//	break
		//}
	}
}

//
func testCreateGenesisBlock(t *testing.T, para *client, testLoopCount int32) {
	genesisBlock := makeGenesisBlockInputTestData()
	err := para.blockSyncClient.createGenesisBlock(genesisBlock)

	switch testLoopCount {
	case 0:
		assert.Nil(t, err)
	default:
		assert.Error(t, err)
	}

}

//    localdb
func testClearLocalOldBlocks(t *testing.T, para *client, testLoopCount int32) {
	isCleaned, err := para.blockSyncClient.clearLocalOldBlocks()

	switch testLoopCount {
	case 0:
		assert.Nil(t, err)
	case 1:
		assert.Equal(t, true, !isCleaned && err == nil)
	default: //2
		assert.Error(t, err)
	}
}

//
func testInitFirstLocalHeightIfNeed(t *testing.T, para *client, testLoopCount int32) {
	err := para.blockSyncClient.initFirstLocalHeightIfNeed()

	switch testLoopCount {
	case 0:
		assert.Nil(t, err)
	default:
		assert.Error(t, err)
	}
}

//
func testSyncBlocksIfNeed(t *testing.T, para *client, testLoopCount int32) {
	errorCount := int32(0)
	for i := int32(1); i <= 10; i++ {
		atomic.StoreInt32(&actionReturnIndexAtom, i)
		isSynced, err := para.blockSyncClient.syncBlocksIfNeed()
		if err != nil {
			errorCount++
		}
		assert.Equal(t, isSynced, i == 3 || i == 6)
	}

	switch testLoopCount {
	case 0:
		assert.Equal(t, true, errorCount == 4)
	default:
		assert.Equal(t, true, errorCount == 7)
	}

	atomic.StoreInt32(&actionReturnIndexAtom, 0)
}

//  SyncHasCaughtUp
func testSyncHasCaughtUp(t *testing.T, para *client, testLoopCount int32) {
	oldValue := para.blockSyncClient.syncHasCaughtUp()
	para.blockSyncClient.setSyncCaughtUp(true)
	isSyncHasCaughtUp := para.blockSyncClient.syncHasCaughtUp()
	para.blockSyncClient.setSyncCaughtUp(oldValue)

	assert.Equal(t, true, isSyncHasCaughtUp)
}

//  getBlockSyncState
func testGetBlockSyncState(t *testing.T, para *client, testLoopCount int32) {
	oldValue := para.blockSyncClient.getBlockSyncState()
	para.blockSyncClient.setBlockSyncState(blockSyncStateFinished)
	syncState := para.blockSyncClient.getBlockSyncState()
	para.blockSyncClient.setBlockSyncState(oldValue)

	assert.Equal(t, true, syncState == blockSyncStateFinished)
}

//
func execTest(t *testing.T, para *client, testLoopCount int32) {
	atomic.StoreInt32(&actionReturnIndexAtom, 0)
	atomic.StoreInt32(&testLoopCountAtom, testLoopCount)

	testCreateGenesisBlock(t, para, testLoopCount)
	testSyncBlocksIfNeed(t, para, testLoopCount)
	testInitFirstLocalHeightIfNeed(t, para, testLoopCount)
	testClearLocalOldBlocks(t, para, testLoopCount)

	testSyncHasCaughtUp(t, para, testLoopCount)
	testGetBlockSyncState(t, para, testLoopCount)
}

//
func TestSyncBlocks(t *testing.T) {
	cfg := types.NewDplatformOSConfig(testnode.DefaultConfig)
	q := queue.New("channel")
	q.SetConfig(cfg)
	defer q.Close()
	para := createParaTestInstance(t, q)
	go q.Start()
	go mockMessageReply(q)
	//       ，               ,
	for i := int32(0); i <= TestMaxLoopCount-1; i++ {
		execTest(t, para, i)
	}

}
