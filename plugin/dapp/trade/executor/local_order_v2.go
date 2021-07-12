// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/trade/types"
)

//
// 1. v0
// 2. v1  table   LocalOrder
// 3. v2  table   LocalOrderV2  v1       ，  v0       ，  v2    v0 v1
//          6.4    v0

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

//
var optV2 = &table.Option{
	Prefix:  "LODB-trade",
	Name:    "order_v2",
	Primary: "txIndex",
	// asset       price_exec + price_symbol + asset_exec+asset_symbol
	// status:                ,      ，            ，
	//       ，    ，       ，     ，          .
	//      00     10     11          12         1*
	//     ：    key   ，            ，   n         ，
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
		// "owner_isSell_statusPrefix", //         ,
		"owner_status",             //    2
		"assset_isSell_isFinished", //   isFinish,
		"owner_asset_isFinished",
		"owner_isFinished",
		// "owner_statusPrefix", //          ,
		//      key，         key    ，
		// https://dplatform.io/document/105   1.8 sell & asset-price & status, order by price
		// https://dplatform.io/document/105   1.3 buy  & asset-price & status, order by price
		"asset_isSell_status_price",
		//   1.2   1.5         addr-status buy or sell
		"owner_isSell_status",
	},
}

// OrderV2Row order row
type OrderV2Row struct {
	*pty.LocalOrder
}

// NewOrderV2Row create row
func NewOrderV2Row() *OrderV2Row {
	return &OrderV2Row{LocalOrder: nil}
}

// CreateRow create row
func (r *OrderV2Row) CreateRow() *table.Row {
	return &table.Row{Data: &pty.LocalOrder{}}
}

// SetPayload set payload
func (r *OrderV2Row) SetPayload(data types.Message) error {
	if d, ok := data.(*pty.LocalOrder); ok {
		r.LocalOrder = d
		return nil
	}

	return types.ErrTypeAsset
}

// Get get index key
func (r *OrderV2Row) Get(key string) ([]byte, error) {
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
	case "owner_isSell_status":
		return []byte(fmt.Sprintf("%s_%d_%s", r.Owner, r.isSell(), r.status())), nil
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
	case "asset_isSell_status_price":
		return []byte(fmt.Sprintf("%s_%d_%s_%s", r.asset(), r.isSell(), r.status(), r.price())), nil
	default:
		return nil, types.ErrNotFound
	}
}

//        ： Price
func (r *OrderV2Row) asset() string {
	return r.LocalOrder.PriceExec + "." + r.LocalOrder.PriceSymbol + "_" + r.LocalOrder.AssetExec + "." + r.LocalOrder.AssetSymbol
}

func (r *OrderV2Row) isSell() int {
	if r.IsSellOrder {
		return 1
	}
	return 0
}

func (r *OrderV2Row) isFinished() int {
	if r.IsFinished {
		return 1
	}
	return 0
}

func (r *OrderV2Row) price() string {
	//       ，
	if r.AmountPerBoardlot == 0 {
		return ""
	}
	p := calcPriceOfToken(r.PricePerBoardlot, r.AmountPerBoardlot)
	return fmt.Sprintf("%018d", p)
}

// status:                ,      ，            ，
//       ，    ，       ，     ，          .
//      01     10     11          12        19 -> 1*
func (r *OrderV2Row) status() string {
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

// NewOrderTableV2 create order table
func NewOrderTableV2(kvdb dbm.KV) *table.Table {
	rowMeta := NewOrderV2Row()
	rowMeta.SetPayload(&pty.LocalOrder{})
	t, err := table.NewTable(rowMeta, kvdb, optV2)
	if err != nil {
		panic(err)
	}
	return t
}

func listV2(db dbm.KVDB, indexName string, data *pty.LocalOrder, count, direction int32) ([]*table.Row, error) {
	query := NewOrderTableV2(db).GetQuery(db)
	var primary []byte
	if len(data.TxIndex) > 0 {
		primary = []byte(data.TxIndex)
	}

	cur := &OrderV2Row{LocalOrder: data}
	index, err := cur.Get(indexName)
	if err != nil {
		tradelog.Error("query List failed", "key", string(primary), "param", data, "err", err)
		return nil, err
	}
	tradelog.Debug("query List dbg", "indexName", indexName, "index", string(index), "primary", primary, "count", count, "direction", direction)
	rows, err := query.ListIndex(indexName, index, primary, count, direction)
	if err != nil {
		tradelog.Error("query List failed", "key", string(primary), "param", data, "err", err)
		return nil, err
	}
	if len(rows) == 0 {
		return nil, types.ErrNotFound
	}
	return rows, nil
}
