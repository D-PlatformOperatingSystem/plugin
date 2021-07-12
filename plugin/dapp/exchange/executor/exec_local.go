package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	ety "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
)

/*
 *             ，     
 *      ，    (localDB),       ，   
 */

func (e *exchange) ExecLocal_LimitOrder(payload *ety.LimitOrder, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.Ty == types.ExecOk {
		for _, log := range receiptData.Logs {
			switch log.Ty {
			case ety.TyLimitOrderLog:
				receipt := &ety.ReceiptExchange{}
				if err := types.Decode(log.Log, receipt); err != nil {
					return nil, err
				}
				kv := e.updateIndex(receipt)
				dbSet.KV = append(dbSet.KV, kv...)
			}
		}
	}
	return e.addAutoRollBack(tx, dbSet.KV), nil
}

func (e *exchange) ExecLocal_MarketOrder(payload *ety.MarketOrder, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.Ty == types.ExecOk {
		for _, log := range receiptData.Logs {
			switch log.Ty {
			case ety.TyMarketOrderLog:
				receipt := &ety.ReceiptExchange{}
				if err := types.Decode(log.Log, receipt); err != nil {
					return nil, err
				}
				kv := e.updateIndex(receipt)
				dbSet.KV = append(dbSet.KV, kv...)
			}
		}
	}
	return e.addAutoRollBack(tx, dbSet.KV), nil
}

func (e *exchange) ExecLocal_RevokeOrder(payload *ety.RevokeOrder, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.Ty == types.ExecOk {
		for _, log := range receiptData.Logs {
			switch log.Ty {
			case ety.TyRevokeOrderLog:
				receipt := &ety.ReceiptExchange{}
				if err := types.Decode(log.Log, receipt); err != nil {
					return nil, err
				}
				kv := e.updateIndex(receipt)
				dbSet.KV = append(dbSet.KV, kv...)
			}
		}
	}
	return e.addAutoRollBack(tx, dbSet.KV), nil
}

//      
func (e *exchange) addAutoRollBack(tx *types.Transaction, kv []*types.KeyValue) *types.LocalDBSet {
	dbSet := &types.LocalDBSet{}
	dbSet.KV = e.AddRollbackKV(tx, tx.Execer, kv)
	return dbSet
}

func (e *exchange) updateIndex(receipt *ety.ReceiptExchange) (kvs []*types.KeyValue) {
	historyTable := NewHistoryOrderTable(e.GetLocalDB())
	marketTable := NewMarketDepthTable(e.GetLocalDB())
	orderTable := NewMarketOrderTable(e.GetLocalDB())
	switch receipt.Order.Status {
	case ety.Ordered:
		err := e.updateOrder(marketTable, orderTable, historyTable, receipt.GetOrder(), receipt.GetIndex())
		if err != nil {
			return nil
		}
		err = e.updateMatchOrders(marketTable, orderTable, historyTable, receipt.GetOrder(), receipt.GetMatchOrders(), receipt.GetIndex())
		if err != nil {
			return nil
		}
	case ety.Completed:
		err := e.updateOrder(marketTable, orderTable, historyTable, receipt.GetOrder(), receipt.GetIndex())
		if err != nil {
			return nil
		}
		err = e.updateMatchOrders(marketTable, orderTable, historyTable, receipt.GetOrder(), receipt.GetMatchOrders(), receipt.GetIndex())
		if err != nil {
			return nil
		}
	case ety.Revoked:
		err := e.updateOrder(marketTable, orderTable, historyTable, receipt.GetOrder(), receipt.GetIndex())
		if err != nil {
			return nil
		}
	}

	//  KV
	kv, err := marketTable.Save()
	if err != nil {
		elog.Error("updateIndex", "marketTable.Save", err.Error())
		return nil
	}
	kvs = append(kvs, kv...)
	kv, err = orderTable.Save()
	if err != nil {
		elog.Error("updateIndex", "orderTable.Save", err.Error())
		return nil
	}
	kvs = append(kvs, kv...)
	kv, err = historyTable.Save()
	if err != nil {
		elog.Error("updateIndex", "historyTable.Save", err.Error())
		return nil
	}
	kvs = append(kvs, kv...)

	return
}

func (e *exchange) updateOrder(marketTable, orderTable, historyTable *table.Table, order *ety.Order, index int64) error {
	left := order.GetLimitOrder().GetLeftAsset()
	right := order.GetLimitOrder().GetRightAsset()
	op := order.GetLimitOrder().GetOp()
	price := order.GetLimitOrder().GetPrice()
	switch order.Status {
	case ety.Ordered:
		var markDepth ety.MarketDepth
		depth, err := queryMarketDepth(e.GetLocalDB(), left, right, op, price)
		if err == types.ErrNotFound {
			markDepth.Price = price
			markDepth.LeftAsset = left
			markDepth.RightAsset = right
			markDepth.Op = op
			markDepth.Amount = order.Balance
		} else {
			markDepth.Price = price
			markDepth.LeftAsset = left
			markDepth.RightAsset = right
			markDepth.Op = op
			markDepth.Amount = depth.Amount + order.Balance
		}
		//marketDepth
		err = marketTable.Replace(&markDepth)
		if err != nil {
			elog.Error("updateIndex", "marketTable.Replace", err.Error())
			return err
		}
		err = orderTable.Replace(order)
		if err != nil {
			elog.Error("updateIndex", "orderTable.Replace", err.Error())
			return err
		}

	case ety.Completed:
		err := historyTable.Replace(order)
		if err != nil {
			elog.Error("updateIndex", "historyTable.Replace", err.Error())
			return err
		}
	case ety.Revoked:
		//     ordered          
		var marketDepth ety.MarketDepth
		depth, err := queryMarketDepth(e.GetLocalDB(), left, right, op, price)
		if err == nil {
			//marketDepth
			marketDepth.Price = price
			marketDepth.LeftAsset = left
			marketDepth.RightAsset = right
			marketDepth.Op = op
			marketDepth.Amount = depth.Amount - order.Balance
			err = marketTable.Replace(&marketDepth)
			if err != nil {
				elog.Error("updateIndex", "marketTable.Replace", err.Error())
				return err
			}
		}
		if marketDepth.Amount <= 0 {
			//  
			err = marketTable.DelRow(&marketDepth)
			if err != nil {
				elog.Error("updateIndex", "marketTable.DelRow", err.Error())
				return err
			}
		}
		//      orderID
		order.Status = ety.Ordered
		err = orderTable.DelRow(order)
		if err != nil {
			elog.Error("updateIndex", "orderTable.DelRow", err.Error())
			return err
		}
		order.Status = ety.Revoked
		order.Index = index
		//       
		err = historyTable.Replace(order)
		if err != nil {
			elog.Error("updateIndex", "historyTable.Replace", err.Error())
			return err
		}
	}
	return nil
}
func (e *exchange) updateMatchOrders(marketTable, orderTable, historyTable *table.Table, order *ety.Order, matchOrders []*ety.Order, index int64) error {
	left := order.GetLimitOrder().GetLeftAsset()
	right := order.GetLimitOrder().GetRightAsset()
	op := order.GetLimitOrder().GetOp()
	if len(matchOrders) > 0 {
		//      
		cache := make(map[int64]int64)
		for i, matchOrder := range matchOrders {
			if matchOrder.Status == ety.Completed {
				//       orderID
				matchOrder.Status = ety.Ordered
				err := orderTable.DelRow(matchOrder)
				if err != nil {
					elog.Error("updateIndex", "orderTable.DelRow", err.Error())
					return err
				}
				//  index,     index
				matchOrder.Status = ety.Completed
				matchOrder.Index = index + int64(i+1)
				err = historyTable.Replace(matchOrder)
				if err != nil {
					elog.Error("updateIndex", "historyTable.Replace", err.Error())
					return err
				}
			}
			if matchOrder.Status == ety.Ordered {
				//    
				err := orderTable.Replace(matchOrder)
				if err != nil {
					elog.Error("updateIndex", "orderTable.Replace", err.Error())
					return err
				}
			}
			executed := cache[matchOrder.GetLimitOrder().Price]
			executed = executed + matchOrder.Executed
			cache[matchOrder.GetLimitOrder().Price] = executed
		}

		//        
		for pr, executed := range cache {
			var matchDepth ety.MarketDepth
			depth, err := queryMarketDepth(e.GetLocalDB(), left, right, OpSwap(op), pr)
			if err == types.ErrNotFound {
				continue
			} else {
				matchDepth.Price = pr
				matchDepth.LeftAsset = left
				matchDepth.RightAsset = right
				matchDepth.Op = OpSwap(op)
				matchDepth.Amount = depth.Amount - executed
			}
			//marketDepth
			err = marketTable.Replace(&matchDepth)
			if err != nil {
				elog.Error("updateIndex", "marketTable.Replace", err.Error())
				return err
			}
			if matchDepth.Amount <= 0 {
				//  
				err = marketTable.DelRow(&matchDepth)
				if err != nil {
					elog.Error("updateIndex", "marketTable.DelRow", err.Error())
					return err
				}
			}
		}
	}
	return nil
}

//OpSwap ...
func OpSwap(op int32) int32 {
	if op == ety.OpBuy {
		return ety.OpSell
	}
	return ety.OpBuy
}
