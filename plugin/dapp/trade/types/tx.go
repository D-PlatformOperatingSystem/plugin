// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//TradeSellTx : info for sell order
type TradeSellTx struct {
	TokenSymbol       string `json:"tokenSymbol"`
	AmountPerBoardlot int64  `json:"amountPerBoardlot"`
	MinBoardlot       int64  `json:"minBoardlot"`
	PricePerBoardlot  int64  `json:"pricePerBoardlot"`
	TotalBoardlot     int64  `json:"totalBoardlot"`
	Fee               int64  `json:"fee"`
	AssetExec         string `json:"assetExec"`
	PriceExec         string `json:"priceExec"`
	PriceSymbol       string `json:"priceSymbol"`
}

//TradeBuyTx :info for buy order to speficied order
type TradeBuyTx struct {
	SellID      string `json:"sellID"`
	BoardlotCnt int64  `json:"boardlotCnt"`
	Fee         int64  `json:"fee"`
}

//TradeRevokeTx :
type TradeRevokeTx struct {
	SellID string `json:"sellID,"`
	Fee    int64  `json:"fee"`
}

//TradeBuyLimitTx :
type TradeBuyLimitTx struct {
	TokenSymbol       string `json:"tokenSymbol"`
	AmountPerBoardlot int64  `json:"amountPerBoardlot"`
	MinBoardlot       int64  `json:"minBoardlot"`
	PricePerBoardlot  int64  `json:"pricePerBoardlot"`
	TotalBoardlot     int64  `json:"totalBoardlot"`
	Fee               int64  `json:"fee"`
	AssetExec         string `json:"assetExec"`
	PriceExec         string `json:"priceExec"`
	PriceSymbol       string `json:"priceSymbol"`
}

//TradeSellMarketTx :         token
type TradeSellMarketTx struct {
	BuyID       string `json:"buyID"`
	BoardlotCnt int64  `json:"boardlotCnt"`
	Fee         int64  `json:"fee"`
}

//TradeRevokeBuyTx :
type TradeRevokeBuyTx struct {
	BuyID string `json:"buyID,"`
	Fee   int64  `json:"fee"`
}
