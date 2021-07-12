// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pkt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/pokerbull/types"
)

func (c *PokerBull) updateIndex(log *pkt.ReceiptPBGame) (kvs []*types.KeyValue) {
	//     Action
	kvs = append(kvs, addPBGameStatusAndPlayer(log.Status, log.PlayerNum, log.Value, log.Index, log.GameId))
	kvs = append(kvs, addPBGameStatusIndexKey(log.Status, log.GameId, log.Index))
	kvs = append(kvs, addPBGameAddrIndexKey(log.Status, log.Addr, log.GameId, log.Index))

	/*
		//
		if log.Status == pkt.PBGameActionStart {
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionStart, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionStart, log.PrevIndex))
		}

		if log.Status == pkt.PBGameActionContinue {
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionStart, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionContinue, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionStart, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionContinue, log.PrevIndex))
		}

		if log.Status == pkt.PBGameActionQuit {
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionStart, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusAndPlayer(pkt.PBGameActionContinue, log.PlayerNum, log.Value, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionStart, log.PrevIndex))
			kvs = append(kvs, delPBGameStatusIndexKey(pkt.PBGameActionContinue, log.PrevIndex))

		}*/

	//    ，
	if !log.IsWaiting {
		for _, v := range log.Players {
			if v != log.Addr {
				kvs = append(kvs, addPBGameAddrIndexKey(log.Status, v, log.GameId, log.Index))
			}
			kvs = append(kvs, delPBGameAddrIndexKey(v, log.PrevIndex))
		}

		kvs = append(kvs, delPBGameStatusAndPlayer(log.PreStatus, log.PlayerNum, log.Value, log.PrevIndex))
		kvs = append(kvs, delPBGameStatusIndexKey(log.PreStatus, log.PrevIndex))
	}

	return kvs
}

func (c *PokerBull) execLocal(receipt *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	for i := 0; i < len(receipt.Logs); i++ {
		item := receipt.Logs[i]
		if item.Ty == pkt.TyLogPBGameStart || item.Ty == pkt.TyLogPBGameContinue || item.Ty == pkt.TyLogPBGameQuit || item.Ty == pkt.TyLogPBGamePlay {
			var Gamelog pkt.ReceiptPBGame
			err := types.Decode(item.Log, &Gamelog)
			if err != nil {
				panic(err) //     ，
			}
			kv := c.updateIndex(&Gamelog)
			dbSet.KV = append(dbSet.KV, kv...)
		}
	}
	return dbSet, nil
}

// ExecLocal_Start       local
func (c *PokerBull) ExecLocal_Start(payload *pkt.PBGameStart, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return c.execLocal(receiptData)
}

// ExecLocal_Continue       local
func (c *PokerBull) ExecLocal_Continue(payload *pkt.PBGameContinue, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return c.execLocal(receiptData)
}

// ExecLocal_Quit       local
func (c *PokerBull) ExecLocal_Quit(payload *pkt.PBGameQuit, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return c.execLocal(receiptData)
}

// ExecLocal_Play        local
func (c *PokerBull) ExecLocal_Play(payload *pkt.PBGamePlay, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return c.execLocal(receiptData)
}
