package executor

import (
	"fmt"
	"math/big"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/client"
	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	tab "github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	et "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
)

// Action action struct
type Action struct {
	statedb   dbm.KV
	txhash    []byte
	fromaddr  string
	blocktime int64
	height    int64
	execaddr  string
	localDB   dbm.KVDB
	index     int
	api       client.QueueProtocolAPI
}

//NewAction ...
func NewAction(e *exchange, tx *types.Transaction, index int) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &Action{e.GetStateDB(), hash, fromaddr,
		e.GetBlockTime(), e.GetHeight(), dapp.ExecAddress(string(tx.Execer)), e.GetLocalDB(), index, e.GetAPI()}
}

//GetIndex get index
func (a *Action) GetIndex() int64 {
	//  4 0,      matchorder
	return (a.height*types.MaxTxsPerBlock + int64(a.index)) * 1e4
}

//GetKVSet get kv set
func (a *Action) GetKVSet(order *et.Order) (kvset []*types.KeyValue) {
	kvset = append(kvset, &types.KeyValue{Key: calcOrderKey(order.OrderID), Value: types.Encode(order)})
	return kvset
}

//OpSwap
func (a *Action) OpSwap(op int32) int32 {
	if op == et.OpBuy {
		return et.OpSell
	}
	return et.OpBuy
}

//CalcActualCost
func CalcActualCost(op int32, amount int64, price int64) int64 {
	if op == et.OpBuy {
		return SafeMul(amount, price)
	}
	return amount
}

//CheckPrice price        1<=price<=1e16
func CheckPrice(price int64) bool {
	if price > 1e16 || price < 1 {
		return false
	}
	return true
}

//CheckOp ...
func CheckOp(op int32) bool {
	if op == et.OpBuy || op == et.OpSell {
		return true
	}
	return false
}

//CheckCount ...
func CheckCount(count int32) bool {
	return count <= 20 && count >= 0
}

//CheckAmount     1e8
func CheckAmount(amount int64) bool {
	if amount < types.Coin || amount >= types.MaxCoin {
		return false
	}
	return true
}

//CheckDirection ...
func CheckDirection(direction int32) bool {
	if direction == et.ListASC || direction == et.ListDESC {
		return true
	}
	return false
}

//CheckStatus ...
func CheckStatus(status int32) bool {
	if status == et.Ordered || status == et.Completed || status == et.Revoked {
		return true
	}
	return false
}

//CheckExchangeAsset
func CheckExchangeAsset(left, right *et.Asset) bool {
	if left.Execer == "" || left.Symbol == "" || right.Execer == "" || right.Symbol == "" {
		return false
	}
	if (left.Execer == "coins" && right.Execer == "coins") || (left.Symbol == right.Symbol) {
		return false
	}
	return true
}

//LimitOrder ...
func (a *Action) LimitOrder(payload *et.LimitOrder) (*types.Receipt, error) {
	leftAsset := payload.GetLeftAsset()
	rightAsset := payload.GetRightAsset()
	//TODO      ，        ，       checkTx
	//coins       dpos
	if !CheckExchangeAsset(leftAsset, rightAsset) {
		return nil, et.ErrAsset
	}
	if !CheckAmount(payload.GetAmount()) {
		return nil, et.ErrAssetAmount
	}
	if !CheckPrice(payload.GetPrice()) {
		return nil, et.ErrAssetPrice
	}
	if !CheckOp(payload.GetOp()) {
		return nil, et.ErrAssetOp
	}
	//TODO   symbol
	cfg := a.api.GetConfig()
	leftAssetDB, err := account.NewAccountDB(cfg, leftAsset.GetExecer(), leftAsset.GetSymbol(), a.statedb)
	if err != nil {
		return nil, err
	}
	rightAssetDB, err := account.NewAccountDB(cfg, rightAsset.GetExecer(), rightAsset.GetSymbol(), a.statedb)
	if err != nil {
		return nil, err
	}
	//
	if payload.GetOp() == et.OpBuy {
		amount := SafeMul(payload.GetAmount(), payload.GetPrice())
		rightAccount := rightAssetDB.LoadExecAccount(a.fromaddr, a.execaddr)
		if rightAccount.Balance < amount {
			elog.Error("limit check right balance", "addr", a.fromaddr, "avail", rightAccount.Balance, "need", amount)
			return nil, et.ErrAssetBalance
		}
		return a.matchLimitOrder(payload, leftAssetDB, rightAssetDB)

	}
	if payload.GetOp() == et.OpSell {
		amount := payload.GetAmount()
		leftAccount := leftAssetDB.LoadExecAccount(a.fromaddr, a.execaddr)
		if leftAccount.Balance < amount {
			elog.Error("limit check left balance", "addr", a.fromaddr, "avail", leftAccount.Balance, "need", amount)
			return nil, et.ErrAssetBalance
		}
		return a.matchLimitOrder(payload, leftAssetDB, rightAssetDB)
	}
	return nil, fmt.Errorf("unknow op")
}

//RevokeOrder ...
func (a *Action) RevokeOrder(payload *et.RevokeOrder) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	order, err := findOrderByOrderID(a.statedb, a.localDB, payload.GetOrderID())
	if err != nil {
		return nil, err
	}
	if order.Addr != a.fromaddr {
		elog.Error("RevokeOrder.OrderCheck", "addr", a.fromaddr, "order.addr", order.Addr, "order.status", order.Status)
		return nil, et.ErrAddr
	}
	if order.Status == et.Completed || order.Status == et.Revoked {
		elog.Error("RevokeOrder.OrderCheck", "addr", a.fromaddr, "order.addr", order.Addr, "order.status", order.Status)
		return nil, et.ErrOrderSatus
	}
	leftAsset := order.GetLimitOrder().GetLeftAsset()
	rightAsset := order.GetLimitOrder().GetRightAsset()
	price := order.GetLimitOrder().GetPrice()
	balance := order.GetBalance()

	cfg := a.api.GetConfig()

	if order.GetLimitOrder().GetOp() == et.OpBuy {
		rightAssetDB, err := account.NewAccountDB(cfg, rightAsset.GetExecer(), rightAsset.GetSymbol(), a.statedb)
		if err != nil {
			return nil, err
		}
		amount := CalcActualCost(et.OpBuy, balance, price)
		rightAccount := rightAssetDB.LoadExecAccount(a.fromaddr, a.execaddr)
		if rightAccount.Frozen < amount {
			elog.Error("revoke check right frozen", "addr", a.fromaddr, "avail", rightAccount.Frozen, "amount", amount)
			return nil, et.ErrAssetBalance
		}
		receipt, err := rightAssetDB.ExecActive(a.fromaddr, a.execaddr, amount)
		if err != nil {
			elog.Error("RevokeOrder.ExecActive", "addr", a.fromaddr, "amount", amount, "err", err.Error())
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kvs = append(kvs, receipt.KV...)
	}
	if order.GetLimitOrder().GetOp() == et.OpSell {
		leftAssetDB, err := account.NewAccountDB(cfg, leftAsset.GetExecer(), leftAsset.GetSymbol(), a.statedb)
		if err != nil {
			return nil, err
		}
		amount := CalcActualCost(et.OpSell, balance, price)
		leftAccount := leftAssetDB.LoadExecAccount(a.fromaddr, a.execaddr)
		if leftAccount.Frozen < amount {
			elog.Error("revoke check left frozen", "addr", a.fromaddr, "avail", leftAccount.Frozen, "amount", amount)
			return nil, et.ErrAssetBalance
		}
		receipt, err := leftAssetDB.ExecActive(a.fromaddr, a.execaddr, amount)
		if err != nil {
			elog.Error("RevokeOrder.ExecActive", "addr", a.fromaddr, "amount", amount, "err", err.Error())
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kvs = append(kvs, receipt.KV...)
	}

	//  order
	order.Status = et.Revoked
	order.UpdateTime = a.blocktime
	kvs = append(kvs, a.GetKVSet(order)...)
	re := &et.ReceiptExchange{
		Order: order,
		Index: a.GetIndex(),
	}
	receiptlog := &types.ReceiptLog{Ty: et.TyRevokeOrderLog, Log: types.Encode(re)}
	logs = append(logs, receiptlog)
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil

}

//
//   ：
//1.       ，         。
//2.       ，           。
//3.
//4.
func (a *Action) matchLimitOrder(payload *et.LimitOrder, leftAccountDB, rightAccountDB *account.DB) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var orderKey string
	var priceKey string
	var count int

	or := &et.Order{
		OrderID:    a.GetIndex(),
		Value:      &et.Order_LimitOrder{LimitOrder: payload},
		Ty:         et.TyLimitOrderAction,
		Executed:   0,
		AVGPrice:   0,
		Balance:    payload.GetAmount(),
		Status:     et.Ordered,
		Addr:       a.fromaddr,
		UpdateTime: a.blocktime,
		Index:      a.GetIndex(),
	}
	re := &et.ReceiptExchange{
		Order: or,
		Index: a.GetIndex(),
	}

	//        100     ,        ，
	//
	for {
		//
		if count >= et.MaxMatchCount {
			break
		}
		//
		marketDepthList, err := QueryMarketDepth(a.localDB, payload.GetLeftAsset(), payload.GetRightAsset(), a.OpSwap(payload.Op), priceKey, et.Count)
		if err == types.ErrNotFound {
			break
		}
		for _, marketDepth := range marketDepthList.List {
			if count >= et.MaxMatchCount {
				break
			}
			//
			if payload.Op == et.OpBuy && marketDepth.Price > payload.GetPrice() {
				continue
			}
			//
			if payload.Op == et.OpSell && marketDepth.Price < payload.GetPrice() {
				continue
			}
			//
			for {
				//
				if count >= et.MaxMatchCount {
					break
				}
				orderList, err := findOrderIDListByPrice(a.localDB, payload.GetLeftAsset(), payload.GetRightAsset(), marketDepth.Price, a.OpSwap(payload.Op), et.ListASC, orderKey)
				if err == types.ErrNotFound {
					break
				}

				for _, matchorder := range orderList.List {
					//
					if count >= et.MaxMatchCount {
						break
					}
					//
					if matchorder.Addr == a.fromaddr {
						continue
					}
					//  ,
					log, kv, err := a.matchModel(leftAccountDB, rightAccountDB, payload, matchorder, or, re) // payload, or redundant
					if err != nil {
						return nil, err
					}
					logs = append(logs, log...)
					kvs = append(kvs, kv...)
					//    ,    ，      ，     ，  count
					if or.Status == et.Completed {
						receiptlog := &types.ReceiptLog{Ty: et.TyLimitOrderLog, Log: types.Encode(re)}
						logs = append(logs, receiptlog)
						receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
						return receipts, nil
					}
					//TODO            ?        ，   receipt      ，           ，          ，
					//
					count = count + 1
				}
				//       10      ,
				if orderList.PrimaryKey == "" {
					break
				}
				orderKey = orderList.PrimaryKey
			}
		}

		//         primaryKey         ,
		if marketDepthList.PrimaryKey == "" {
			break
		}
		priceKey = marketDepthList.PrimaryKey
	}

	//
	if payload.Op == et.OpBuy {
		amount := CalcActualCost(et.OpBuy, or.Balance, payload.Price)
		receipt, err := rightAccountDB.ExecFrozen(a.fromaddr, a.execaddr, amount)
		if err != nil {
			elog.Error("LimitOrder.ExecFrozen", "addr", a.fromaddr, "amount", amount, "err", err.Error())
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kvs = append(kvs, receipt.KV...)
	}
	if payload.Op == et.OpSell {
		amount := CalcActualCost(et.OpSell, or.Balance, payload.Price)
		receipt, err := leftAccountDB.ExecFrozen(a.fromaddr, a.execaddr, amount)
		if err != nil {
			elog.Error("LimitOrder.ExecFrozen", "addr", a.fromaddr, "amount", amount, "err", err.Error())
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kvs = append(kvs, receipt.KV...)
	}
	//  order
	kvs = append(kvs, a.GetKVSet(or)...)
	re.Order = or
	receiptlog := &types.ReceiptLog{Ty: et.TyLimitOrderLog, Log: types.Encode(re)}
	logs = append(logs, receiptlog)
	receipts := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipts, nil
}

//
func (a *Action) matchModel(leftAccountDB, rightAccountDB *account.DB, payload *et.LimitOrder, matchorder *et.Order, or *et.Order, re *et.ReceiptExchange) ([]*types.ReceiptLog, []*types.KeyValue, error) {
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	var matched int64

	if matchorder.GetBalance() >= or.GetBalance() {
		matched = or.GetBalance()
	} else {
		matched = matchorder.GetBalance()
	}

	elog.Info("try match", "activeId", or.OrderID, "passiveId", matchorder.OrderID, "activeAddr", or.Addr, "passiveAddr",
		matchorder.Addr, "amount", matched, "price", payload.Price)

	if payload.Op == et.OpSell {
		//
		amount := CalcActualCost(matchorder.GetLimitOrder().Op, matched, payload.Price)
		receipt, err := rightAccountDB.ExecTransferFrozen(matchorder.Addr, a.fromaddr, a.execaddr, amount)
		if err != nil {
			elog.Error("matchModel.ExecTransferFrozen", "from", matchorder.Addr, "to", a.fromaddr, "amount", amount, "err", err)
			return nil, nil, err
		}
		logs = append(logs, receipt.Logs...)
		kvs = append(kvs, receipt.KV...)
		//
		if payload.Price < matchorder.GetLimitOrder().Price {
			amount := CalcActualCost(matchorder.GetLimitOrder().Op, matched, matchorder.GetLimitOrder().Price-payload.Price)
			receipt, err := rightAccountDB.ExecActive(matchorder.Addr, a.execaddr, amount)
			if err != nil {
				elog.Error("matchModel.ExecActive", "addr", matchorder.Addr, "amount", amount, "err", err.Error())
				return nil, nil, err
			}
			logs = append(logs, receipt.Logs...)
			kvs = append(kvs, receipt.KV...)
		}
		//
		amount = CalcActualCost(payload.Op, matched, payload.Price)
		receipt, err = leftAccountDB.ExecTransfer(a.fromaddr, matchorder.Addr, a.execaddr, amount)
		if err != nil {
			elog.Error("matchModel.ExecTransfer", "from", a.fromaddr, "to", matchorder.Addr, "amount", amount, "err", err.Error())
			return nil, nil, err
		}
		logs = append(logs, receipt.Logs...)
		kvs = append(kvs, receipt.KV...)

		//
		or.AVGPrice = payload.Price
		//  matchOrder
		matchorder.AVGPrice = caclAVGPrice(matchorder, payload.Price, matched) //TODO
	}
	if payload.Op == et.OpBuy {
		//
		amount := CalcActualCost(matchorder.GetLimitOrder().Op, matched, matchorder.GetLimitOrder().Price)
		receipt, err := leftAccountDB.ExecTransferFrozen(matchorder.Addr, a.fromaddr, a.execaddr, amount)
		if err != nil {
			elog.Error("matchModel.ExecTransferFrozen2", "from", matchorder.Addr, "to", a.fromaddr, "amount", amount, "err", err.Error())
			return nil, nil, err
		}
		logs = append(logs, receipt.Logs...)
		kvs = append(kvs, receipt.KV...)
		//
		amount = CalcActualCost(payload.Op, matched, matchorder.GetLimitOrder().Price)
		receipt, err = rightAccountDB.ExecTransfer(a.fromaddr, matchorder.Addr, a.execaddr, amount)
		if err != nil {
			elog.Error("matchModel.ExecTransfer2", "from", a.fromaddr, "to", matchorder.Addr, "amount", amount, "err", err.Error())
			return nil, nil, err
		}
		logs = append(logs, receipt.Logs...)
		kvs = append(kvs, receipt.KV...)

		//    ，
		or.AVGPrice = matchorder.GetLimitOrder().Price
		//  matchOrder
		matchorder.AVGPrice = caclAVGPrice(matchorder, matchorder.GetLimitOrder().Price, matched) //TODO
	}

	if matched == matchorder.GetBalance() {
		matchorder.Status = et.Completed
	} else {
		matchorder.Status = et.Ordered
	}

	if matched == or.GetBalance() {
		or.Status = et.Completed
	} else {
		or.Status = et.Ordered
	}

	if matched == or.GetBalance() {
		matchorder.Balance -= matched
		matchorder.Executed = matched
		kvs = append(kvs, a.GetKVSet(matchorder)...)

		or.Executed += matched
		or.Balance = 0
		kvs = append(kvs, a.GetKVSet(or)...) //or complete
	} else {
		or.Balance -= matched
		or.Executed += matched

		matchorder.Executed = matched
		matchorder.Balance = 0
		kvs = append(kvs, a.GetKVSet(matchorder)...) //matchorder complete
	}

	re.Order = or
	re.MatchOrders = append(re.MatchOrders, matchorder)
	return logs, kvs, nil
}

//       ，    ，   localdb   ，
// 1.          orderID localdb
// 2.    ，     ，  orderID localdb          ，
func findOrderByOrderID(statedb dbm.KV, localdb dbm.KV, orderID int64) (*et.Order, error) {
	table := NewMarketOrderTable(localdb)
	primaryKey := []byte(fmt.Sprintf("%022d", orderID))
	row, err := table.GetData(primaryKey)
	if err != nil {
		data, err := statedb.Get(calcOrderKey(orderID))
		if err != nil {
			elog.Error("findOrderByOrderID.Get", "orderID", orderID, "err", err.Error())
			return nil, err
		}
		var order et.Order
		err = types.Decode(data, &order)
		if err != nil {
			elog.Error("findOrderByOrderID.Decode", "orderID", orderID, "err", err.Error())
			return nil, err
		}
		return &order, nil
	}
	return row.Data.(*et.Order), nil

}

func findOrderIDListByPrice(localdb dbm.KV, left, right *et.Asset, price int64, op, direction int32, primaryKey string) (*et.OrderList, error) {
	table := NewMarketOrderTable(localdb)
	prefix := []byte(fmt.Sprintf("%s:%s:%d:%016d", left.GetSymbol(), right.GetSymbol(), op, price))

	var rows []*tab.Row
	var err error
	if primaryKey == "" { //     ,
		rows, err = table.ListIndex("market_order", prefix, nil, et.Count, direction)
	} else {
		rows, err = table.ListIndex("market_order", prefix, []byte(primaryKey), et.Count, direction)
	}
	if err != nil {
		elog.Error("findOrderIDListByPrice.", "left", left, "right", right, "price", price, "err", err.Error())
		return nil, err
	}
	var orderList et.OrderList
	for _, row := range rows {
		order := row.Data.(*et.Order)
		//
		order.Executed = order.GetLimitOrder().Amount - order.Balance
		orderList.List = append(orderList.List, order)
	}
	//
	if len(rows) == int(et.Count) {
		orderList.PrimaryKey = string(rows[len(rows)-1].Primary)
	}
	return &orderList, nil
}

//Direction           ，
func Direction(op int32) int32 {
	if op == et.OpBuy {
		return et.ListDESC
	}
	return et.ListASC
}

//QueryMarketDepth   primaryKey        ，         ,         ，
func QueryMarketDepth(localdb dbm.KV, left, right *et.Asset, op int32, primaryKey string, count int32) (*et.MarketDepthList, error) {
	table := NewMarketDepthTable(localdb)
	prefix := []byte(fmt.Sprintf("%s:%s:%d", left.GetSymbol(), right.GetSymbol(), op))
	if count == 0 {
		count = et.Count
	}
	var rows []*tab.Row
	var err error
	if primaryKey == "" { //     ,
		rows, err = table.ListIndex("price", prefix, nil, count, Direction(op))
	} else {
		rows, err = table.ListIndex("price", prefix, []byte(primaryKey), count, Direction(op))
	}
	if err != nil {
		//elog.Error("QueryMarketDepth.", "left", left, "right", right, "err", err.Error())
		return nil, err
	}

	var list et.MarketDepthList
	for _, row := range rows {
		list.List = append(list.List, row.Data.(*et.MarketDepth))
	}
	//
	if len(rows) == int(count) {
		list.PrimaryKey = string(rows[len(rows)-1].Primary)
	}
	return &list, nil
}

//QueryHistoryOrderList
func QueryHistoryOrderList(localdb dbm.KV, left, right *et.Asset, primaryKey string, count, direction int32) (types.Message, error) {
	table := NewHistoryOrderTable(localdb)
	prefix := []byte(fmt.Sprintf("%s:%s", left.Symbol, right.Symbol))
	indexName := "name"
	if count == 0 {
		count = et.Count
	}
	var rows []*tab.Row
	var err error
	var orderList et.OrderList
HERE:
	if primaryKey == "" { //     ,
		rows, err = table.ListIndex(indexName, prefix, nil, count, direction)
	} else {
		rows, err = table.ListIndex(indexName, prefix, []byte(primaryKey), count, direction)
	}
	if err != nil && err != types.ErrNotFound {
		elog.Error("QueryCompletedOrderList.", "left", left, "right", right, "err", err.Error())
		return nil, err
	}
	if err == types.ErrNotFound {
		return &orderList, nil
	}
	for _, row := range rows {
		order := row.Data.(*et.Order)
		//           completed,revoked        ，
		if order.Status == et.Revoked {
			continue
		}
		//
		order.Executed = order.GetLimitOrder().Amount - order.Balance
		orderList.List = append(orderList.List, order)
		if len(orderList.List) == int(count) {
			//
			orderList.PrimaryKey = string(row.Primary)
			return &orderList, nil
		}
	}
	if len(orderList.List) != int(count) && len(rows) == int(count) {
		primaryKey = string(rows[len(rows)-1].Primary)
		goto HERE
	}
	return &orderList, nil
}

//QueryOrderList
func QueryOrderList(localdb dbm.KV, addr string, status, count, direction int32, primaryKey string) (types.Message, error) {
	var table *tab.Table
	if status == et.Completed || status == et.Revoked {
		table = NewHistoryOrderTable(localdb)
	} else {
		table = NewMarketOrderTable(localdb)
	}
	prefix := []byte(fmt.Sprintf("%s:%d", addr, status))
	indexName := "addr_status"
	if count == 0 {
		count = et.Count
	}
	var rows []*tab.Row
	var err error
	if primaryKey == "" { //     ,
		rows, err = table.ListIndex(indexName, prefix, nil, count, direction)
	} else {
		rows, err = table.ListIndex(indexName, prefix, []byte(primaryKey), count, direction)
	}
	if err != nil {
		elog.Error("QueryOrderList.", "addr", addr, "err", err.Error())
		return nil, err
	}
	var orderList et.OrderList
	for _, row := range rows {
		order := row.Data.(*et.Order)
		//
		order.Executed = order.GetLimitOrder().Amount - order.Balance
		orderList.List = append(orderList.List, order)
	}
	//
	if len(rows) == int(count) {
		orderList.PrimaryKey = string(rows[len(rows)-1].Primary)
	}
	return &orderList, nil
}

func queryMarketDepth(localdb dbm.KV, left, right *et.Asset, op int32, price int64) (*et.MarketDepth, error) {
	table := NewMarketDepthTable(localdb)
	primaryKey := []byte(fmt.Sprintf("%s:%s:%d:%016d", left.GetSymbol(), right.GetSymbol(), op, price))
	row, err := table.GetData(primaryKey)
	if err != nil {
		return nil, err
	}
	return row.Data.(*et.MarketDepth), nil
}

//SafeMul math         ，
func SafeMul(x, y int64) int64 {
	res := big.NewInt(0).Mul(big.NewInt(x), big.NewInt(y))
	res = big.NewInt(0).Div(res, big.NewInt(types.Coin))
	return res.Int64()
}

//
func caclAVGPrice(order *et.Order, price int64, amount int64) int64 {
	x := big.NewInt(0).Mul(big.NewInt(order.AVGPrice), big.NewInt(order.GetLimitOrder().Amount-order.GetBalance()))
	y := big.NewInt(0).Mul(big.NewInt(price), big.NewInt(amount))
	total := big.NewInt(0).Add(x, y)
	div := big.NewInt(0).Add(big.NewInt(order.GetLimitOrder().Amount-order.GetBalance()), big.NewInt(amount))
	avg := big.NewInt(0).Div(total, div)
	return avg.Int64()
}
