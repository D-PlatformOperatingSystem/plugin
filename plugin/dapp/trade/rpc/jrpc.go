// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/types"

	ptypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/trade/types"
)

//CreateRawTradeSellTx :     token
func (jrpc *Jrpc) CreateRawTradeSellTx(in *ptypes.TradeSellTx, result *interface{}) error {
	if in == nil {
		return types.ErrInvalidParam
	}
	param := &ptypes.TradeForSell{
		TokenSymbol:       in.TokenSymbol,
		AmountPerBoardlot: in.AmountPerBoardlot,
		MinBoardlot:       in.MinBoardlot,
		PricePerBoardlot:  in.PricePerBoardlot,
		TotalBoardlot:     in.TotalBoardlot,
		Starttime:         0,
		Stoptime:          0,
		Crowdfund:         false,
		AssetExec:         in.AssetExec,
		PriceExec:         in.PriceExec,
		PriceSymbol:       in.PriceSymbol,
	}

	reply, err := jrpc.cli.CreateRawTradeSellTx(context.Background(), param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

//CreateRawTradeBuyTx :     token      ,
func (jrpc *Jrpc) CreateRawTradeBuyTx(in *ptypes.TradeBuyTx, result *interface{}) error {
	if in == nil {
		return types.ErrInvalidParam
	}
	param := &ptypes.TradeForBuy{
		SellID:      in.SellID,
		BoardlotCnt: in.BoardlotCnt,
	}

	reply, err := jrpc.cli.CreateRawTradeBuyTx(context.Background(), param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

//CreateRawTradeRevokeTx :
func (jrpc *Jrpc) CreateRawTradeRevokeTx(in *ptypes.TradeRevokeTx, result *interface{}) error {
	if in == nil {
		return types.ErrInvalidParam
	}
	param := &ptypes.TradeForRevokeSell{
		SellID: in.SellID,
	}

	reply, err := jrpc.cli.CreateRawTradeRevokeTx(context.Background(), param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

//CreateRawTradeBuyLimitTx :      token
func (jrpc *Jrpc) CreateRawTradeBuyLimitTx(in *ptypes.TradeBuyLimitTx, result *interface{}) error {
	if in == nil {
		return types.ErrInvalidParam
	}
	param := &ptypes.TradeForBuyLimit{
		TokenSymbol:       in.TokenSymbol,
		AmountPerBoardlot: in.AmountPerBoardlot,
		MinBoardlot:       in.MinBoardlot,
		PricePerBoardlot:  in.PricePerBoardlot,
		TotalBoardlot:     in.TotalBoardlot,
		AssetExec:         in.AssetExec,
		PriceExec:         in.PriceExec,
		PriceSymbol:       in.PriceSymbol,
	}

	reply, err := jrpc.cli.CreateRawTradeBuyLimitTx(context.Background(), param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

//CreateRawTradeSellMarketTx :        token
func (jrpc *Jrpc) CreateRawTradeSellMarketTx(in *ptypes.TradeSellMarketTx, result *interface{}) error {
	if in == nil {
		return types.ErrInvalidParam
	}
	param := &ptypes.TradeForSellMarket{
		BuyID:       in.BuyID,
		BoardlotCnt: in.BoardlotCnt,
	}

	reply, err := jrpc.cli.CreateRawTradeSellMarketTx(context.Background(), param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

//CreateRawTradeRevokeBuyTx :
func (jrpc *Jrpc) CreateRawTradeRevokeBuyTx(in *ptypes.TradeRevokeBuyTx, result *interface{}) error {
	if in == nil {
		return types.ErrInvalidParam
	}
	param := &ptypes.TradeForRevokeBuy{
		BuyID: in.BuyID,
	}

	reply, err := jrpc.cli.CreateRawTradeRevokeBuyTx(context.Background(), param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}
