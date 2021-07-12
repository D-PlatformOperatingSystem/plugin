// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package para

import (
	"errors"
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/types"
)

//IsCaughtUp         ,
func (client *client) Query_IsCaughtUp(req *types.ReqNil) (types.Message, error) {
	if client == nil {
		return nil, fmt.Errorf("%s", "client not bind message queue.")
	}

	return &types.IsCaughtUp{Iscaughtup: client.isCaughtUp()}, nil
}

func (client *client) Query_LocalBlockInfo(req *types.ReqInt) (types.Message, error) {
	if client == nil {
		return nil, fmt.Errorf("%s", "client not bind message queue.")
	}

	var block *pt.ParaLocalDbBlock
	var err error
	if req.Height <= -1 {
		block, err = client.getLastLocalBlock()
		if err != nil {
			return nil, err
		}
	} else {
		block, err = client.getLocalBlockByHeight(req.Height)
		if err != nil {
			return nil, err
		}
	}

	blockInfo := &pt.ParaLocalDbBlockInfo{
		Height:         block.Height,
		MainHash:       common.ToHex(block.MainHash),
		MainHeight:     block.MainHeight,
		ParentMainHash: common.ToHex(block.ParentMainHash),
		BlockTime:      block.BlockTime,
	}

	for _, tx := range block.Txs {
		blockInfo.Txs = append(blockInfo.Txs, common.ToHex(tx.Hash()))
	}

	return blockInfo, nil
}

func (client *client) Query_LeaderInfo(req *types.ReqNil) (types.Message, error) {
	if client == nil {
		return nil, fmt.Errorf("%s", "client not bind message queue.")
	}
	nodes, leader, base, off, isLeader := client.blsSignCli.getLeaderInfo()
	return &pt.ElectionStatus{IsLeader: isLeader, Leader: &pt.LeaderSyncInfo{ID: nodes[leader], BaseIdx: base, Offset: off}}, nil
}

func (client *client) Query_CommitTxInfo(req *types.ReqNil) (types.Message, error) {
	if client == nil {
		return nil, fmt.Errorf("%s", "client not bind message queue.")
	}

	rt := client.blsSignCli.showTxBuffInfo()
	return rt, nil
}

func (client *client) Query_BlsPubKey(req *types.ReqString) (types.Message, error) {
	if client == nil || req == nil {
		return nil, fmt.Errorf("%s", "client not bind message queue.")
	}

	var pub pt.BlsPubKey
	if len(req.Data) > 0 {
		p, err := client.blsSignCli.secp256Prikey2BlsPub(req.Data)
		if err != nil {
			return nil, err
		}
		pub.Key = p
		return &pub, nil
	}
	//
	if nil != client.blsSignCli.blsPubKey {
		t := client.blsSignCli.blsPubKey.Bytes()
		pub.Key = common.ToHex(t[:])
		return &pub, nil
	}

	return nil, errors.New("no bls prikey init")
}

// Query_CreateNewAccount   para
func (client *client) Query_CreateNewAccount(acc *types.Account) (types.Message, error) {
	if acc == nil {
		return nil, types.ErrInvalidParam
	}
	plog.Info("Query_CreateNewAccount", "acc", acc.Addr)
	//   para                     commit
	client.commitMsgClient.onWalletAccount(acc)
	return &types.Reply{IsOk: true, Msg: []byte("OK")}, nil
}

// Query_WalletStatus   para
func (client *client) Query_WalletStatus(walletStatus *types.WalletStatus) (types.Message, error) {
	if walletStatus == nil {
		return nil, types.ErrInvalidParam
	}
	plog.Info("Query_WalletStatus", "walletStatus", walletStatus.IsWalletLock)
	//   para      walletStatus.IsWalletLock      /
	client.commitMsgClient.onWalletStatus(walletStatus)
	return &types.Reply{IsOk: true, Msg: []byte("OK")}, nil
}
