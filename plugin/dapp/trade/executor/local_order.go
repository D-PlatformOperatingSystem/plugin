// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/trade/types"
)

/*

 1.            （   ）
   1.1         -> owner
   1.2        token  -> owner_asset
   1.3               -> owner
   1.4                  token     -> owner_asset
 2.           ：       （   ） -> owner_status
 3.     token         GetTokenBuyOrderByStatus  -> asset_inBuy_status
 4.     token         token      token     （   ） -> owner_asset/owner_asset_isSell
 5.               （   ）  -> owner_isSell_status
 6.     token            -> asset_isSell
 7.               （      ） owner_status
*/
var opt_order_table = &table.Option{
	Prefix:  "LODB-trade",
	Name:    "order",
	Primary: "txIndex",
	// asset = asset_exec+asset_symbol
	//
	// status:                ,      ，            ，
	//       ，    ，       ，     ，          .
	//      00     10     11          12         1*
	//       ：     ，status&isSell
	Index: []string{
		"key",                 //
		"asset",               //
		"asset_isSell_status", //    3
		// "asset_status",     ，
		// "asset_isSell",
		"owner",              //    1.1， 1.3
		"owner_asset",        //    1.2, 1.4, 4, 7
		"owner_asset_isSell", //    4
		"owner_asset_status", //    ，
		"owner_isSell",       //    6
		// "owner_isSell_status",      ，
		// "owner_isSell_statusPrefix", //         ,
		"owner_status",             //    2
		"assset_isSell_isFinished", //   isFinish,
		"owner_asset_isFinished",
		"owner_isFinished",
		// "owner_statusPrefix", //          ,
	},
}

// OrderRow order row
type OrderRow struct {
	*pty.LocalOrder
}

// NewOrderRow create row
func NewOrderRow() *OrderRow {
	return &OrderRow{LocalOrder: nil}
}

// CreateRow create row
func (r *OrderRow) CreateRow() *table.Row {
	return &table.Row{Data: &pty.LocalOrder{}}
}

// SetPayload set payload
func (r *OrderRow) SetPayload(data types.Message) error {
	if d, ok := data.(*pty.LocalOrder); ok {
		r.LocalOrder = d
		return nil
	}
	return types.ErrTypeAsset
}

// Get get index key
func (r *OrderRow) Get(key string) ([]byte, error) {
	switch key {
	case "txIndex":
		return []byte(r.TxIndex), nil
	case "key":
		return []byte(r.Key), nil
	case "asset":
		return []byte(r.asset()), nil
	case "asset_isSell_status":
		return []byte(fmt.Sprintf("%s_%d_%s", r.asset(), r.isSell(), r.status())), nil
	case "owner":
		return []byte(r.Owner), nil
	case "owner_asset":
		return []byte(fmt.Sprintf("%s_%s", r.Owner, r.asset())), nil
	case "owner_asset_isSell":
		return []byte(fmt.Sprintf("%s_%s_%d", r.Owner, r.asset(), r.isSell())), nil
	case "owner_asset_status":
		return []byte(fmt.Sprintf("%s_%s_%s", r.Owner, r.asset(), r.status())), nil
	case "owner_isSell":
		return []byte(fmt.Sprintf("%s_%d", r.Owner, r.isSell())), nil
	//case "owner_isSell_statusPrefix":
	//	return []byte(fmt.Sprintf("%s_%d_%s", r.Owner, r.asset(), r.isSell())), nil
	case "owner_status":
		return []byte(fmt.Sprintf("%s_%s", r.Owner, r.status())), nil
	//case "owner_statusPrefix":
	//	return []byte(fmt.Sprintf("%s_%d", r.Owner, r.isSell())), nil
	case "assset_isSell_isFinished":
		return []byte(fmt.Sprintf("%s_%d_%d", r.Owner, r.isSell(), r.isFinished())), nil
	case "owner_asset_isFinished":
		return []byte(fmt.Sprintf("%s_%s_%d", r.Owner, r.asset(), r.isFinished())), nil
	case "owner_isFinished":
		return []byte(fmt.Sprintf("%s_%d", r.Owner, r.isFinished())), nil
	default:
		return nil, types.ErrNotFound
	}
}

func (r *OrderRow) asset() string {
	return r.LocalOrder.AssetExec + "." + r.LocalOrder.AssetSymbol
}

func (r *OrderRow) isSell() int {
	if r.IsSellOrder {
		return 1
	}
	return 0
}

func (r *OrderRow) isFinished() int {
	if r.IsFinished {
		return 1
	}
	return 0
}

// status:                ,      ，            ，
//       ，    ，       ，     ，          .
//      01     10     11          12        19 -> 1*
func (r *OrderRow) status() string {
	if r.Status == pty.TradeOrderStatusOnBuy || r.Status == pty.TradeOrderStatusOnSale {
		return "01" //    1
	} else if r.Status == pty.TradeOrderStatusSoldOut || r.Status == pty.TradeOrderStatusBoughtOut {
		return "12"
	} else if r.Status == pty.TradeOrderStatusRevoked || r.Status == pty.TradeOrderStatusBuyRevoked {
		return "10"
	} else if r.Status == pty.TradeOrderStatusSellHalfRevoked || r.Status == pty.TradeOrderStatusBuyHalfRevoked {
		return "11"
	} else if r.Status == pty.TradeOrderStatusGroupComplete {
		return "1" // 1* match complete
	}

	return "XX"
}

// NewOrderTable create order table
func NewOrderTable(kvdb dbm.KV) *table.Table {
	rowMeta := NewOrderRow()
	rowMeta.SetPayload(&pty.LocalOrder{})
	t, err := table.NewTable(rowMeta, kvdb, opt_order_table)
	if err != nil {
		panic(err)
	}
	return t
}

// gen order from tx and receipt
func (t *trade) genSellLimit(tx *types.Transaction, sell *pty.ReceiptSellBase,
	sellorder *pty.SellOrder, txIndex string) *pty.LocalOrder {

	order := &pty.LocalOrder{
		AssetSymbol:       sellorder.TokenSymbol,
		TxIndex:           txIndex,
		Owner:             sellorder.Address,
		AmountPerBoardlot: sellorder.AmountPerBoardlot,
		MinBoardlot:       sellorder.MinBoardlot,
		PricePerBoardlot:  sellorder.PricePerBoardlot,
		TotalBoardlot:     sellorder.TotalBoardlot,
		TradedBoardlot:    sellorder.SoldBoardlot,
		BuyID:             "",
		Status:            sellorder.Status,
		SellID:            sell.SellID,
		TxHash:            []string{common.ToHex(tx.Hash())},
		Height:            sell.Height,
		Key:               sell.SellID,
		BlockTime:         t.GetBlockTime(),
		IsSellOrder:       true,
		AssetExec:         sellorder.AssetExec,
		IsFinished:        false,
		PriceExec:         sellorder.PriceExec,
		PriceSymbol:       sellorder.PriceSymbol,
	}
	return order
}

func (t *trade) updateSellLimit(tx *types.Transaction, sell *pty.ReceiptSellBase,
	sellorder *pty.SellOrder, txIndex string, ldb *table.Table) *pty.LocalOrder {

	xs, err := ldb.ListIndex("key", []byte(sell.SellID), nil, 1, 0)
	if err != nil || len(xs) != 1 {
		return nil
	}
	order, ok := xs[0].Data.(*pty.LocalOrder)
	tradelog.Debug("Table dbg", "sell-update", order, "data", xs[0].Data)
	if !ok {
		tradelog.Error("Table failed", "sell-update", order)
		return nil

	}
	status := sellorder.Status
	if status == pty.TradeOrderStatusRevoked && sell.SoldBoardlot > 0 {
		status = pty.TradeOrderStatusSellHalfRevoked
	}
	order.Status = status
	order.TxHash = append(order.TxHash, common.ToHex(tx.Hash()))
	order.TradedBoardlot = sellorder.SoldBoardlot
	order.IsFinished = (status != pty.TradeOrderStatusOnSale)

	tradelog.Debug("Table", "sell-update", order)

	ldb.Replace(order)

	return order
}

func (t *trade) rollBackSellLimit(tx *types.Transaction, sell *pty.ReceiptSellBase,
	sellorder *pty.SellOrder, txIndex string, ldb *table.Table, tradedBoardlot int64) *pty.LocalOrder {

	xs, err := ldb.ListIndex("key", []byte(sell.SellID), nil, 1, 0)
	if err != nil || len(xs) != 1 {
		return nil
	}
	order, ok := xs[0].Data.(*pty.LocalOrder)
	if !ok {
		return nil

	}
	//       ,
	//
	order.Status = pty.TradeOrderStatusOnSale
	order.TxHash = order.TxHash[:len(order.TxHash)-1]
	order.TradedBoardlot = order.TradedBoardlot - tradedBoardlot
	order.IsFinished = (order.Status != pty.TradeOrderStatusOnSale)

	ldb.Replace(order)

	return order
}

func parseOrderAmountFloat(s string) int64 {
	x, err := strconv.ParseFloat(s, 64)
	if err != nil {
		tradelog.Error("parseOrderAmountFloat", "decode receipt", err)
		return 0
	}
	return int64(x * float64(types.TokenPrecision))
}

func parseOrderPriceFloat(s string) int64 {
	x, err := strconv.ParseFloat(s, 64)
	if err != nil {
		tradelog.Error("parseOrderPriceFloat", "decode receipt", err)
		return 0
	}
	return int64(x * float64(types.Coin))
}

func (t *trade) genSellMarket(tx *types.Transaction, sell *pty.ReceiptSellBase, txIndex string) *pty.LocalOrder {
	order := &pty.LocalOrder{
		AssetSymbol:       sell.TokenSymbol,
		TxIndex:           txIndex,
		Owner:             sell.Owner,
		AmountPerBoardlot: parseOrderAmountFloat(sell.AmountPerBoardlot),
		MinBoardlot:       sell.MinBoardlot,
		PricePerBoardlot:  parseOrderPriceFloat(sell.PricePerBoardlot),
		TotalBoardlot:     sell.TotalBoardlot,
		TradedBoardlot:    sell.SoldBoardlot,
		BuyID:             sell.BuyID,
		Status:            pty.TradeOrderStatusSoldOut,
		SellID:            calcTokenSellID(hex.EncodeToString(tx.Hash())),
		TxHash:            []string{common.ToHex(tx.Hash())},
		Height:            sell.Height,
		Key:               calcTokenSellID(hex.EncodeToString(tx.Hash())),
		BlockTime:         t.GetBlockTime(),
		IsSellOrder:       true,
		AssetExec:         sell.AssetExec,

		IsFinished:  true,
		PriceExec:   sell.PriceExec,
		PriceSymbol: sell.PriceSymbol,
	}
	return order
}

func (t *trade) genBuyLimit(tx *types.Transaction, buy *pty.ReceiptBuyBase, txIndex string) *pty.LocalOrder {
	order := &pty.LocalOrder{
		AssetSymbol:       buy.TokenSymbol,
		TxIndex:           txIndex,
		Owner:             buy.Owner,
		AmountPerBoardlot: parseOrderAmountFloat(buy.AmountPerBoardlot),
		MinBoardlot:       buy.MinBoardlot,
		PricePerBoardlot:  parseOrderPriceFloat(buy.PricePerBoardlot),
		TotalBoardlot:     buy.TotalBoardlot,
		TradedBoardlot:    buy.BoughtBoardlot,
		BuyID:             buy.BuyID,
		Status:            pty.TradeOrderStatusOnBuy,
		SellID:            "",
		TxHash:            []string{common.ToHex(tx.Hash())},
		Height:            buy.Height,
		Key:               buy.BuyID,
		BlockTime:         t.GetBlockTime(),
		IsSellOrder:       false,
		AssetExec:         buy.AssetExec,
		IsFinished:        false,
		PriceExec:         buy.PriceExec,
		PriceSymbol:       buy.PriceSymbol,
	}
	return order
}

func (t *trade) updateBuyLimit(tx *types.Transaction, buy *pty.ReceiptBuyBase,
	buyorder *pty.BuyLimitOrder, txIndex string, ldb *table.Table) *pty.LocalOrder {

	xs, err := ldb.ListIndex("key", []byte(buy.BuyID), nil, 1, 0)
	if err != nil || len(xs) != 1 {
		return nil
	}
	order, ok := xs[0].Data.(*pty.LocalOrder)
	if !ok {
		return nil

	}
	status := buyorder.Status
	if status == pty.TradeOrderStatusBuyRevoked && buy.BoughtBoardlot > 0 {
		status = pty.TradeOrderStatusBuyHalfRevoked
	}
	order.Status = status
	order.TxHash = append(order.TxHash, common.ToHex(tx.Hash()))
	order.TradedBoardlot = buyorder.BoughtBoardlot
	order.IsFinished = (status != pty.TradeOrderStatusOnBuy)

	ldb.Replace(order)

	return order
}

func (t *trade) rollbackBuyLimit(tx *types.Transaction, buy *pty.ReceiptBuyBase,
	buyorder *pty.BuyLimitOrder, txIndex string, ldb *table.Table, traded int64) *pty.LocalOrder {

	xs, err := ldb.ListIndex("key", []byte(buy.BuyID), nil, 1, 0)
	if err != nil || len(xs) != 1 {
		return nil
	}
	order, ok := xs[0].Data.(*pty.LocalOrder)
	if !ok {
		return nil
	}

	order.Status = pty.TradeOrderStatusOnBuy
	order.TxHash = order.TxHash[:len(order.TxHash)-1]
	order.TradedBoardlot = order.TradedBoardlot - traded
	order.IsFinished = false

	ldb.Replace(order)

	return order
}

func (t *trade) genBuyMarket(tx *types.Transaction, buy *pty.ReceiptBuyBase, txIndex string) *pty.LocalOrder {
	order := &pty.LocalOrder{
		AssetSymbol:       buy.TokenSymbol,
		TxIndex:           txIndex,
		Owner:             buy.Owner,
		AmountPerBoardlot: parseOrderAmountFloat(buy.AmountPerBoardlot),
		MinBoardlot:       buy.MinBoardlot,
		PricePerBoardlot:  parseOrderPriceFloat(buy.PricePerBoardlot),
		TotalBoardlot:     buy.TotalBoardlot,
		TradedBoardlot:    buy.BoughtBoardlot,
		BuyID:             calcTokenBuyID(hex.EncodeToString(tx.Hash())),
		Status:            pty.TradeOrderStatusBoughtOut,
		SellID:            buy.SellID,
		TxHash:            []string{common.ToHex(tx.Hash())},
		Height:            buy.Height,
		Key:               calcTokenBuyID(hex.EncodeToString(tx.Hash())),
		BlockTime:         t.GetBlockTime(),
		IsSellOrder:       false,
		AssetExec:         buy.AssetExec,
		IsFinished:        true,
		PriceExec:         buy.PriceExec,
		PriceSymbol:       buy.PriceSymbol,
	}
	return order
}
