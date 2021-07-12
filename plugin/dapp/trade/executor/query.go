// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/trade/types"
)

//     trade  query，
// 1.  token
//         token      (     )： OnBuy/OnSale
//    token       （     ）: SoldOut/BoughtOut--> TODO        （         ）
// 2.   addr  。
//             （      ）
//            （addr      ）
//
//      /    orderID ，  txhash     key
// key     orderID， txhash (0xAAAAAAAAAAAAAAA)

// 1.15 both buy/sell order
func (t *trade) Query_GetOnesOrderWithStatus(req *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesOrderWithStatus(req)
}

// get order by id
func (t *trade) Query_GetOneOrder(req *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOneOrder(req)
}

// query reply utils

const (
	orderStatusInvalid = iota
	orderStatusOn
	orderStatusDone
	orderStatusRevoke
)

const (
	orderTypeInvalid = iota
	orderTypeSell
	orderTypeBuy
)

func fromStatus(status int32) (st, ty int32) {
	switch status {
	case pty.TradeOrderStatusOnSale:
		return orderStatusOn, orderTypeSell
	case pty.TradeOrderStatusSoldOut:
		return orderStatusDone, orderTypeSell
	case pty.TradeOrderStatusRevoked:
		return orderStatusRevoke, orderTypeSell
	case pty.TradeOrderStatusOnBuy:
		return orderStatusOn, orderTypeBuy
	case pty.TradeOrderStatusBoughtOut:
		return orderStatusDone, orderTypeBuy
	case pty.TradeOrderStatusBuyRevoked:
		return orderStatusRevoke, orderTypeBuy
	}
	return orderStatusInvalid, orderTypeInvalid
}

// GetOnesOrderWithStatus by address-status
func (t *trade) GetOnesOrderWithStatus(req *pty.ReqAddrAssets) (types.Message, error) {
	orderStatus, orderType := fromStatus(req.Status)
	if orderStatus == orderStatusInvalid || orderType == orderTypeInvalid {
		return nil, types.ErrInvalidParam
	}

	//    owner isFinished
	var order pty.LocalOrder
	if orderStatus == orderStatusOn {
		order.IsFinished = false
	} else {
		order.IsFinished = true
	}
	order.Owner = req.Addr
	if len(req.FromKey) > 0 {
		order.TxIndex = req.FromKey
	}
	rows, err := listV2(t.GetLocalDB(), "owner_isFinished", &order, req.Count, req.Direction)
	if err != nil {
		tradelog.Error("GetOnesOrderWithStatus", "err", err)
		return nil, err
	}
	return t.toTradeOrders(rows)
}

func fmtReply(cfg *types.DplatformOSConfig, order *pty.LocalOrder) *pty.ReplyTradeOrder {
	priceExec := order.PriceExec
	priceSymbol := order.PriceSymbol
	if priceExec == "" {
		priceExec = defaultPriceExec
		priceSymbol = cfg.GetCoinSymbol()
	}

	return &pty.ReplyTradeOrder{
		TokenSymbol:       order.AssetSymbol,
		Owner:             order.Owner,
		AmountPerBoardlot: order.AmountPerBoardlot,
		MinBoardlot:       order.MinBoardlot,
		PricePerBoardlot:  order.PricePerBoardlot,
		TotalBoardlot:     order.TotalBoardlot,
		TradedBoardlot:    order.TradedBoardlot,
		BuyID:             order.BuyID,
		Status:            order.Status,
		SellID:            order.SellID,
		TxHash:            order.TxHash[0],
		Height:            order.Height,
		Key:               order.TxIndex,
		BlockTime:         order.BlockTime,
		IsSellOrder:       order.IsSellOrder,
		AssetExec:         order.AssetExec,
		PriceExec:         priceExec,
		PriceSymbol:       priceSymbol,
	}
}

func (t *trade) GetOneOrder(req *pty.ReqAddrAssets) (types.Message, error) {
	query := NewOrderTableV2(t.GetLocalDB())
	tradelog.Debug("query GetData dbg", "primary", req.FromKey)
	row, err := query.GetData([]byte(req.FromKey))
	if err != nil {
		tradelog.Error("query GetData failed", "key", req.FromKey, "err", err)
		return nil, err
	}

	o, ok := row.Data.(*pty.LocalOrder)
	if !ok {
		tradelog.Error("query GetData failed", "err", "bad row type")
		return nil, types.ErrTypeAsset
	}
	cfg := t.GetAPI().GetConfig()
	reply := fmtReply(cfg, o)

	return reply, nil
}
