package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	et "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
)

//      
func (s *exchange) Query_QueryMarketDepth(in *et.QueryMarketDepth) (types.Message, error) {
	if !CheckCount(in.Count) {
		return nil, et.ErrCount
	}
	if !CheckExchangeAsset(in.LeftAsset, in.RightAsset) {
		return nil, et.ErrAsset
	}

	if !CheckOp(in.Op) {
		return nil, et.ErrAssetOp
	}
	return QueryMarketDepth(s.GetLocalDB(), in.LeftAsset, in.RightAsset, in.Op, in.PrimaryKey, in.Count)
}

//         
func (s *exchange) Query_QueryHistoryOrderList(in *et.QueryHistoryOrderList) (types.Message, error) {
	if !CheckExchangeAsset(in.LeftAsset, in.RightAsset) {
		return nil, et.ErrAsset
	}
	if !CheckCount(in.Count) {
		return nil, et.ErrCount
	}

	if !CheckDirection(in.Direction) {
		return nil, et.ErrDirection
	}
	return QueryHistoryOrderList(s.GetLocalDB(), in.LeftAsset, in.RightAsset, in.PrimaryKey, in.Count, in.Direction)
}

//  orderID      
func (s *exchange) Query_QueryOrder(in *et.QueryOrder) (types.Message, error) {
	if in.OrderID == 0 {
		return nil, et.ErrOrderID
	}
	return findOrderByOrderID(s.GetStateDB(), s.GetLocalDB(), in.OrderID)
}

//      ，      （          ）
func (s *exchange) Query_QueryOrderList(in *et.QueryOrderList) (types.Message, error) {
	if !CheckStatus(in.Status) {
		return nil, et.ErrStatus
	}
	if !CheckCount(in.Count) {
		return nil, et.ErrCount
	}

	if !CheckDirection(in.Direction) {
		return nil, et.ErrDirection
	}

	if in.Address == "" {
		return nil, et.ErrAddr
	}
	return QueryOrderList(s.GetLocalDB(), in.Address, in.Status, in.Count, in.Direction, in.PrimaryKey)
}
