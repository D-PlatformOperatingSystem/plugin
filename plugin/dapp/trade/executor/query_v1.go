package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/trade/types"
)

//    1.8   token
func (t *trade) Query_GetTokenSellOrderByStatus(req *pty.ReqTokenSellOrder) (types.Message, error) {
	return t.GetTokenSellOrderByStatus(req, req.Status)
}

// GetTokenSellOrderByStatus by status
// sell & TokenSymbol & status  sort by price
func (t *trade) GetTokenSellOrderByStatus(req *pty.ReqTokenSellOrder, status int32) (types.Message, error) {
	return t.GetTokenOrderByStatus(true, req, status)
}

func (t *trade) GetTokenOrderByStatus(isSell bool, req *pty.ReqTokenSellOrder, status int32) (types.Message, error) {
	if req.Count <= 0 || (req.Direction != 1 && req.Direction != 0) {
		return nil, types.ErrInvalidParam
	}

	var order pty.LocalOrder
	if len(req.FromKey) > 0 {
		order.TxIndex = req.FromKey
	}

	t.setQueryAsset(&order, req.TokenSymbol)
	order.IsSellOrder = isSell
	order.Status = req.Status

	rows, err := listV2(t.GetLocalDB(), "asset_isSell_status_price", &order, req.Count, req.Direction)
	if err != nil {
		tradelog.Error("GetOnesOrderWithStatus", "err", err)
		return nil, err
	}

	return t.toTradeOrders(rows)
}

// 1.3   token
func (t *trade) Query_GetTokenBuyOrderByStatus(req *pty.ReqTokenBuyOrder) (types.Message, error) {
	if req.Status == 0 {
		req.Status = pty.TradeOrderStatusOnBuy
	}
	return t.GetTokenBuyOrderByStatus(req, req.Status)
}

// GetTokenBuyOrderByStatus by status
// buy & TokenSymbol & status buy sort by price
func (t *trade) GetTokenBuyOrderByStatus(req *pty.ReqTokenBuyOrder, status int32) (types.Message, error) {
	// List Direction    ，       ，         ，        ，       。
	direction := 1 - req.Direction
	req2 := pty.ReqTokenSellOrder{
		TokenSymbol: req.TokenSymbol,
		FromKey:     req.FromKey,
		Count:       req.Count,
		Direction:   direction,
		Status:      req.Status,
	}
	return t.GetTokenOrderByStatus(false, &req2, status)
}

// addr part
// 1.4 addr(-token)      ，
func (t *trade) Query_GetOnesSellOrder(req *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesSellOrder(req)
}

// 1.1 addr(-token)      ，
func (t *trade) Query_GetOnesBuyOrder(req *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesBuyOrder(req)
}

// GetOnesSellOrder by address or address-token
func (t *trade) GetOnesSellOrder(addrTokens *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesOrder(true, addrTokens)
}

// GetOnesBuyOrder by address or address-token
func (t *trade) GetOnesBuyOrder(addrTokens *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesOrder(false, addrTokens)
}

// GetOnesSellOrder by address or address-token
func (t *trade) GetOnesOrder(isSell bool, addrTokens *pty.ReqAddrAssets) (types.Message, error) {
	var order pty.LocalOrder
	order.Owner = addrTokens.Addr
	order.IsSellOrder = isSell

	if 0 == len(addrTokens.Token) {
		rows, err := listV2(t.GetLocalDB(), "owner_isSell", &order, 0, 0)
		if err != nil {
			tradelog.Error("GetOnesSellOrder", "err", err)
			return nil, err
		}
		return t.toTradeOrders(rows)
	}

	var replys pty.ReplyTradeOrders
	for _, token := range addrTokens.Token {
		t.setQueryAsset(&order, token)
		rows, err := listV2(t.GetLocalDB(), "owner_asset_isSell", &order, 0, 0)
		if err != nil && err != types.ErrNotFound {
			return nil, err
		}
		if len(rows) == 0 {
			continue
		}
		rs, err := t.toTradeOrders(rows)
		if err != nil {
			return nil, err
		}
		replys.Orders = append(replys.Orders, rs.Orders...)

	}
	return &replys, nil
}

// Query_GetOnesBuyOrderWithStatus 1.5
//         addr-status
func (t *trade) Query_GetOnesSellOrderWithStatus(req *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesSellOrdersWithStatus(req)
}

// Query_GetOnesBuyOrderWithStatus 1.2         addr-status
func (t *trade) Query_GetOnesBuyOrderWithStatus(req *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesBuyOrdersWithStatus(req)
}

// GetOnesSellOrdersWithStatus by address-status
func (t *trade) GetOnesSellOrdersWithStatus(req *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesStatusOrder(true, req)
}

// GetOnesBuyOrdersWithStatus by address-status
func (t *trade) GetOnesBuyOrdersWithStatus(req *pty.ReqAddrAssets) (types.Message, error) {
	return t.GetOnesStatusOrder(false, req)
}

// GetOnesStatusOrder Get Ones Status Order
func (t *trade) GetOnesStatusOrder(isSell bool, req *pty.ReqAddrAssets) (types.Message, error) {
	var order pty.LocalOrder
	order.Owner = req.Addr
	order.Status = req.Status
	order.IsSellOrder = isSell

	rows, err := listV2(t.GetLocalDB(), "owner_isSell_status", &order, 0, 0)
	if err != nil {
		tradelog.Error("GetOnesStatusOrder", "err", err)
		return nil, err
	}
	return t.toTradeOrders(rows)
}

// util
//      ，    token
func (t *trade) setQueryAsset(order *pty.LocalOrder, tokenSymbol string) {
	order.AssetSymbol = tokenSymbol
	order.AssetExec = defaultAssetExec
	order.PriceSymbol = t.GetAPI().GetConfig().GetCoinSymbol()
	order.PriceExec = defaultPriceExec
}

//       ，
func (t *trade) toTradeOrders(rows []*table.Row) (*pty.ReplyTradeOrders, error) {
	var replys pty.ReplyTradeOrders
	cfg := t.GetAPI().GetConfig()
	for _, row := range rows {
		o, ok := row.Data.(*pty.LocalOrder)
		if !ok {
			tradelog.Error("toTradeOrders", "err", "bad row type")
			return nil, types.ErrTypeAsset
		}
		reply := fmtReply(cfg, o)
		replys.Orders = append(replys.Orders, reply)
	}
	return &replys, nil
}
