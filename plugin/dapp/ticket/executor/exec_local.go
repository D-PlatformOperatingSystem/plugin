// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket/types"
)

func (t *Ticket) execLocal(receiptData *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	for _, item := range receiptData.Logs {
		//    ticket  log
		if item.Ty == ty.TyLogNewTicket || item.Ty == ty.TyLogMinerTicket || item.Ty == ty.TyLogCloseTicket {
			var ticketlog ty.ReceiptTicket
			err := types.Decode(item.Log, &ticketlog)
			if err != nil {
				panic(err) //     ，
			}
			kv := t.saveTicket(&ticketlog)
			dbSet.KV = append(dbSet.KV, kv...)
		} else if item.Ty == ty.TyLogTicketBind {
			var ticketlog ty.ReceiptTicketBind
			err := types.Decode(item.Log, &ticketlog)
			if err != nil {
				panic(err) //     ，
			}
			kv := t.saveTicketBind(&ticketlog)
			dbSet.KV = append(dbSet.KV, kv...)
		}
	}
	return dbSet, nil
}

// ExecLocal_Genesis exec local genesis
func (t *Ticket) ExecLocal_Genesis(payload *ty.TicketGenesis, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return t.execLocal(receiptData)
}

// ExecLocal_Topen exec local open
func (t *Ticket) ExecLocal_Topen(payload *ty.TicketOpen, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return t.execLocal(receiptData)
}

// ExecLocal_Tbind exec local bind
func (t *Ticket) ExecLocal_Tbind(payload *ty.TicketBind, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return t.execLocal(receiptData)
}

// ExecLocal_Tclose exec local close
func (t *Ticket) ExecLocal_Tclose(payload *ty.TicketClose, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return t.execLocal(receiptData)
}

// ExecLocal_Miner exec local miner
func (t *Ticket) ExecLocal_Miner(payload *ty.TicketMiner, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return t.execLocal(receiptData)
}
