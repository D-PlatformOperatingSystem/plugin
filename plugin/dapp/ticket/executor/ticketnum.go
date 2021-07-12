// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	tickettypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket/types"
)

const (
	minBlockNum = 3
	maxBlockNum = 10
)

// GetRandNum for ticket executor
func (ticket *Ticket) GetRandNum(blockHash []byte, blockNum int64) (types.Message, error) {
	tlog.Debug("GetRandNum", "blockHash", common.ToHex(blockHash), "blockNum", blockNum)

	if blockNum < minBlockNum {
		blockNum = minBlockNum
	} else if blockNum > maxBlockNum {
		blockNum = maxBlockNum
	}

	if len(blockHash) == 0 {
		return nil, types.ErrBlockNotFound
	}

	txActions, err := ticket.getTxActions(blockHash, blockNum)
	if err != nil {
		return nil, err
	}
	//   genesis block            ，         
	if txActions == nil && err == nil {
		modify := common.Sha256([]byte("hello"))
		return &types.ReplyHash{Hash: modify}, nil
	}
	var modifies []byte
	var bits uint32
	var ticketIds string
	var privHashs []byte
	var vrfHashs []byte

	for _, ticketAction := range txActions {
		//tlog.Debug("GetRandNum", "modify", ticketAction.GetMiner().GetModify(), "bits", ticketAction.GetMiner().GetBits(), "ticketId", ticketAction.GetMiner().GetTicketId(), "PrivHash", ticketAction.GetMiner().GetPrivHash())
		modifies = append(modifies, ticketAction.GetMiner().GetModify()...)
		bits += ticketAction.GetMiner().GetBits()
		ticketIds += ticketAction.GetMiner().GetTicketId()
		privHashs = append(privHashs, ticketAction.GetMiner().GetPrivHash()...)
		vrfHashs = append(vrfHashs, ticketAction.GetMiner().GetVrfHash()...)
	}

	newmodify := fmt.Sprintf("%s:%s:%d:%s", string(modifies), ticketIds, bits, string(privHashs))
	if len(vrfHashs) != 0 {
		newmodify = fmt.Sprintf("%s:%x", newmodify, vrfHashs)
	}

	modify := common.Sha256([]byte(newmodify))

	return &types.ReplyHash{Hash: modify}, nil
}

func (ticket *Ticket) getTxActions(blockHash []byte, blockNum int64) ([]*tickettypes.TicketAction, error) {
	var txActions []*tickettypes.TicketAction
	var reqHashes types.ReqHashes
	currHash := blockHash
	tlog.Debug("getTxActions", "blockHash", common.ToHex(blockHash), "blockNum", blockNum)

	//  blockHash，  block，  blockNum
	for blockNum > 0 {
		req := types.ReqHash{Hash: currHash}

		tempBlock, err := ticket.GetAPI().GetBlockOverview(&req)
		if err != nil {
			return txActions, err
		}
		if tempBlock.Head.Height <= 0 {
			return nil, nil
		}
		reqHashes.Hashes = append(reqHashes.Hashes, currHash)
		currHash = tempBlock.Head.ParentHash
		if tempBlock.Head.Height < 0 && blockNum > 1 {
			return txActions, types.ErrBlockNotFound
		}
		if tempBlock.Head.Height <= 1 {
			break
		}
		blockNum--
	}
	blockDetails, err := ticket.GetAPI().GetBlockByHashes(&reqHashes)
	if err != nil {
		tlog.Error("getTxActions", "blockHash", blockHash, "blockNum", blockNum, "err", err)
		return txActions, err
	}
	cfg := ticket.GetAPI().GetConfig()
	for _, block := range blockDetails.Items {
		tlog.Debug("getTxActions", "blockHeight", block.Block.Height, "blockhash", common.ToHex(block.Block.Hash(cfg)))
		ticketAction, err := ticket.getMinerTx(block.Block)
		if err != nil {
			return txActions, err
		}
		txActions = append(txActions, ticketAction)
	}
	return txActions, nil
}

func (ticket *Ticket) getMinerTx(current *types.Block) (*tickettypes.TicketAction, error) {
	//         execs,       
	if len(current.Txs) == 0 {
		return nil, types.ErrEmptyTx
	}
	baseTx := current.Txs[0]
	//           
	var ticketAction tickettypes.TicketAction
	err := types.Decode(baseTx.GetPayload(), &ticketAction)
	if err != nil {
		return nil, err
	}
	if ticketAction.GetTy() != tickettypes.TicketActionMiner {
		return nil, types.ErrCoinBaseTxType
	}
	//        OK
	if ticketAction.GetMiner() == nil {
		return nil, tickettypes.ErrEmptyMinerTx
	}
	return &ticketAction, nil
}
