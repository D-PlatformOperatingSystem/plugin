// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

/*
trade     trade      ，

           ：
1）    ；
2）       ；
3）    ；
4）    ；
5）       ；
6）    ；
*/

import (
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"

	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/trade/types"
)

var (
	tradelog         = log.New("module", "execs.trade")
	defaultAssetExec = "token"
	driverName       = "trade"
	defaultPriceExec = "coins"
)

// Init :     trade
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	drivers.Register(cfg, GetName(), newTrade, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&trade{}))
}

// GetName :   trade
func GetName() string {
	return newTrade().GetName()
}

type trade struct {
	drivers.DriverBase
}

func newTrade() drivers.Driver {
	t := &trade{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

func (t *trade) GetDriverName() string {
	return driverName
}

func (t *trade) getSellOrderFromDb(sellID []byte) *pty.SellOrder {
	value, err := t.GetStateDB().Get(sellID)
	if err != nil {
		panic(err)
	}
	var sellorder pty.SellOrder
	types.Decode(value, &sellorder)
	return &sellorder
}

// sell limit
func (t *trade) saveSell(base *pty.ReceiptSellBase, ty int32, tx *types.Transaction, txIndex string, ldb *table.Table) {
	sellorder := t.getSellOrderFromDb([]byte(base.SellID))

	if ty == pty.TyLogTradeSellLimit && sellorder.SoldBoardlot == 0 {
		newOrder := t.genSellLimit(tx, base, sellorder, txIndex)
		tradelog.Info("Table", "sell-add", newOrder)
		ldb.Add(newOrder)
	} else {
		t.updateSellLimit(tx, base, sellorder, txIndex, ldb)
	}
}

func (t *trade) deleteSell(base *pty.ReceiptSellBase, ty int32, tx *types.Transaction, txIndex string, ldb *table.Table, tradedBoardlot int64) {
	sellorder := t.getSellOrderFromDb([]byte(base.SellID))
	if ty == pty.TyLogTradeSellLimit && sellorder.SoldBoardlot == 0 {
		ldb.Del([]byte(txIndex))
	} else {
		t.rollBackSellLimit(tx, base, sellorder, txIndex, ldb, tradedBoardlot)
	}
}

func (t *trade) saveBuy(receiptTradeBuy *pty.ReceiptBuyBase, tx *types.Transaction, txIndex string, ldb *table.Table) {
	order := t.genBuyMarket(tx, receiptTradeBuy, txIndex)
	tradelog.Debug("trade BuyMarket save local", "order", order)
	ldb.Add(order)
}

func (t *trade) deleteBuy(receiptTradeBuy *pty.ReceiptBuyBase, txIndex string, ldb *table.Table) {
	ldb.Del([]byte(txIndex))
}

// BuyLimit Local
func (t *trade) getBuyOrderFromDb(buyID []byte) *pty.BuyLimitOrder {
	value, err := t.GetStateDB().Get(buyID)
	if err != nil {
		panic(err)
	}
	var buyOrder pty.BuyLimitOrder
	types.Decode(value, &buyOrder)
	return &buyOrder
}

func (t *trade) saveBuyLimit(buy *pty.ReceiptBuyBase, ty int32, tx *types.Transaction, txIndex string, ldb *table.Table) {
	buyOrder := t.getBuyOrderFromDb([]byte(buy.BuyID))
	tradelog.Debug("Table", "buy-add", buyOrder)
	if buyOrder.Status == pty.TradeOrderStatusOnBuy && buy.BoughtBoardlot == 0 {
		order := t.genBuyLimit(tx, buy, txIndex)
		tradelog.Info("Table", "buy-add", order)
		ldb.Add(order)
	} else {
		t.updateBuyLimit(tx, buy, buyOrder, txIndex, ldb)
	}
}

func (t *trade) deleteBuyLimit(buy *pty.ReceiptBuyBase, ty int32, tx *types.Transaction, txIndex string, ldb *table.Table, traded int64) {
	buyOrder := t.getBuyOrderFromDb([]byte(buy.BuyID))
	if ty == pty.TyLogTradeBuyLimit && buy.BoughtBoardlot == 0 {
		ldb.Del([]byte(txIndex))
	} else {
		t.rollbackBuyLimit(tx, buy, buyOrder, txIndex, ldb, traded)
	}
}

func (t *trade) saveSellMarket(receiptTradeBuy *pty.ReceiptSellBase, tx *types.Transaction, txIndex string, ldb *table.Table) {
	order := t.genSellMarket(tx, receiptTradeBuy, txIndex)
	ldb.Add(order)

}

func (t *trade) deleteSellMarket(receiptTradeBuy *pty.ReceiptSellBase, txIndex string, ldb *table.Table) {
	ldb.Del([]byte(txIndex))
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (t *trade) CheckReceiptExecOk() bool {
	return true
}
