package test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/executor"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	et "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
	"github.com/stretchr/testify/assert"
)

var (
	PrivKeyA   = "0x6da92a632ab7deb67d38c0f6560bcfed28167998f6496db64c258d5e8393a81b" // 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4
	PrivKeyB   = "0x19c069234f9d3e61135fefbeb7791b149cdf6af536f26bebb310d4cd22c3fee4" // 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR
	PrivKeyC   = "0x7a80a1f75d7360c6123c32a78ecf978c1ac55636f87892df38d8b85a9aeff115" // 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k
	PrivKeyD   = "0xcacb1f5d51700aea07fca2246ab43b0917d70405c65edea9b5063d72eb5c6b71" // 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs
	coin       = "dpos"
	token      = "CCNY"
	leftAsset  = &et.Asset{Symbol: coin, Execer: "coins"}
	rightAsset = &et.Asset{Symbol: token, Execer: "token"}

	cli     Cli
	orderID int64
)

func init() {
	cli = NewExecCli()
	//cli = NewGRPCCli(":28804")
}

func TestLimitOrder(t *testing.T) {
	//A    4x10
	req := &et.LimitOrder{LeftAsset: leftAsset, RightAsset: rightAsset, Price: 4, Amount: 10 * types.Coin, Op: et.OpBuy}
	testPlaceLimitOrder(t, req, Nodes[0], PrivKeyA)
}

func TestOrderList(t *testing.T) {
	orderList, err := getOrderList(et.Ordered, Nodes[0], "")
	assert.Nil(t, err)
	t.Log(orderList)
	orderID = orderList.List[0].OrderID
}

func TestGetOrder(t *testing.T) {
	order, err := getOrder(orderID)
	assert.Nil(t, err)
	t.Log(order)
}

func TestMarketDepth(t *testing.T) {
	depth, err := getMarketDepth(&et.QueryMarketDepth{LeftAsset: leftAsset, RightAsset: rightAsset, Op: et.OpBuy})
	t.Log(depth, err)
}

func TestMatch(t *testing.T) {
	//B    4x5
	req := &et.LimitOrder{LeftAsset: leftAsset, RightAsset: rightAsset, Price: 4, Amount: 5 * types.Coin, Op: et.OpSell}
	doLimitOrder(req, PrivKeyB)
}

func TestHistoryOrderList(t *testing.T) {
	historyq := &et.QueryHistoryOrderList{
		LeftAsset:  leftAsset,
		RightAsset: rightAsset,
	}
	historyOrderList, err := getHistoryOrderList(historyq)
	assert.Nil(t, err)
	t.Log(historyOrderList)
}

func TestRevokeOrder(t *testing.T) {
	//A
	testRevokeLimitOrder(t, orderID, Nodes[0], PrivKeyA)
}

func TestSample0(t *testing.T) {
	depth, _ := getMarketDepth(&et.QueryMarketDepth{LeftAsset: leftAsset, RightAsset: rightAsset, Op: et.OpBuy})
	assert.Nil(t, depth)

	depth, _ = getMarketDepth(&et.QueryMarketDepth{LeftAsset: leftAsset, RightAsset: rightAsset, Op: et.OpBuy})
	assert.Nil(t, depth)
}

//       ，        ，
//1.     10   。
//2.       5
//3.
func TestCase1(t *testing.T) {
	//     10
	req := &et.LimitOrder{LeftAsset: leftAsset, RightAsset: rightAsset, Price: 4, Amount: 10 * types.Coin, Op: et.OpBuy}
	_, err := doLimitOrder(req, PrivKeyA)
	assert.Nil(t, err)

	orderList, err := getOrderList(et.Ordered, Nodes[0], "")
	assert.Nil(t, err)

	//     ，
	orderID1 := orderList.List[0].OrderID
	order, err := getOrder(orderID1)
	assert.Nil(t, err)
	assert.Equal(t, int32(et.Ordered), order.Status)
	assert.Equal(t, 10*types.Coin, order.GetBalance())

	//  op
	q := &et.QueryMarketDepth{LeftAsset: leftAsset, RightAsset: rightAsset, Op: et.OpBuy}
	marketDepthList, err := getMarketDepth(q)
	assert.Nil(t, err)
	assert.Equal(t, 10*types.Coin, marketDepthList.List[0].GetAmount())

	//       5
	req = &et.LimitOrder{LeftAsset: leftAsset, RightAsset: rightAsset, Price: 4, Amount: 5 * types.Coin, Op: et.OpSell}
	_, err = doLimitOrder(req, PrivKeyB)
	assert.Nil(t, err)

	orderList, err = getOrderList(et.Completed, Nodes[1], "")
	assert.Nil(t, err)
	orderID2 := orderList.List[0].OrderID

	//    1
	order, err = getOrder(orderID1)
	assert.Nil(t, err)
	//  1       ordered
	assert.Equal(t, int32(et.Ordered), order.Status)
	assert.Equal(t, 5*types.Coin, order.Balance)

	//order2   completed
	order, err = getOrder(orderID2)
	assert.Nil(t, err)
	assert.Equal(t, int32(et.Completed), order.Status)

	//  op
	q = &et.QueryMarketDepth{LeftAsset: leftAsset, RightAsset: rightAsset, Op: et.OpBuy}
	marketDepthList, err = getMarketDepth(q)
	assert.Nil(t, err)
	//
	assert.Equal(t, 5*types.Coin, marketDepthList.List[0].GetAmount())

	//
	q2 := &et.QueryHistoryOrderList{LeftAsset: leftAsset, RightAsset: rightAsset}
	orderList, err = getHistoryOrderList(q2)
	assert.Nil(t, err)
	assert.Equal(t, orderID2, orderList.List[0].OrderID)

	//
	_, err = doRevokeOrder(orderID1, PrivKeyA)
	assert.Nil(t, err)

	//    1
	order, err = getOrder(orderID1)
	assert.Nil(t, err)
	//  1     Revoked
	assert.Equal(t, int32(et.Revoked), order.Status)
	assert.Equal(t, 5*types.Coin, order.Balance)

	//  op
	q = &et.QueryMarketDepth{LeftAsset: leftAsset, RightAsset: rightAsset, Op: et.OpBuy}
	_, err = getMarketDepth(q)
	assert.NotNil(t, err)

	//
	//  ordered
	_, err = getOrderList(et.Ordered, Nodes[0], "")
	assert.NotNil(t, err)
}

func BenchmarkOrder(b *testing.B) {
	req := &et.LimitOrder{LeftAsset: leftAsset, RightAsset: rightAsset, Price: 1, Amount: 10 * types.Coin, Op: et.OpSell}
	ety := types.LoadExecutorType(et.ExchangeX)
	tx, _ := ety.Create("LimitOrder", req)
	for i := 0; i < b.N; i++ {
		cli.Send(tx, PrivKeyA)
	}
}

func doLimitOrder(req *et.LimitOrder, privKey string) ([]*types.ReceiptLog, error) {
	ety := types.LoadExecutorType(et.ExchangeX)
	tx, err := ety.Create("LimitOrder", req)
	if err != nil {
		return nil, err
	}
	return cli.Send(tx, privKey)
}

func doRevokeOrder(orderID int64, privKey string) ([]*types.ReceiptLog, error) {
	ety := types.LoadExecutorType(et.ExchangeX)
	tx, err := ety.Create("RevokeOrder", &et.RevokeOrder{OrderID: orderID})
	if err != nil {
		return nil, err
	}
	return cli.Send(tx, privKey)
}

func getOrderList(status int32, addr string, primaryKey string) (*et.OrderList, error) {
	msg, err := cli.Query(et.FuncNameQueryOrderList, &et.QueryOrderList{Status: status, Address: addr, PrimaryKey: primaryKey})
	if err != nil {
		return nil, err
	}

	var resp et.OrderList
	err = types.Decode(msg, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func getOrder(orderID int64) (*et.Order, error) {
	msg, err := cli.Query(et.FuncNameQueryOrder, &et.QueryOrder{OrderID: orderID})
	if err != nil {
		return nil, err
	}

	var resp et.Order
	err = types.Decode(msg, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func getMarketDepth(q proto.Message) (*et.MarketDepthList, error) {
	msg, err := cli.Query(et.FuncNameQueryMarketDepth, q)
	if err != nil {
		return nil, err
	}

	var resp et.MarketDepthList
	err = types.Decode(msg, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func getHistoryOrderList(q proto.Message) (*et.OrderList, error) {
	msg, err := cli.Query(et.FuncNameQueryHistoryOrderList, q)
	if err != nil {
		return nil, err
	}

	var resp et.OrderList
	err = types.Decode(msg, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func testPlaceLimitOrder(t *testing.T, req *et.LimitOrder, addr string, privkey string) {
	accPrev, err := cli.GetExecAccount(addr, "coins", coin)
	assert.Nil(t, err)
	t.Log(accPrev)

	tokenPrev, err := cli.GetExecAccount(addr, "token", token)
	assert.Nil(t, err)
	t.Log(tokenPrev)

	_, err = doLimitOrder(req, privkey)
	assert.Nil(t, err)

	accAfter, err := cli.GetExecAccount(addr, "coins", coin)
	assert.Nil(t, err)
	t.Log(accAfter)

	tokenAfter, err := cli.GetExecAccount(addr, "token", token)
	assert.Nil(t, err)
	t.Log(tokenAfter)

	cost := executor.CalcActualCost(req.Op, req.Amount, req.Price)
	t.Log(req.Amount, req.Price, cost)
	// dpos/ccny
	if req.Op == et.OpBuy {
		// dpos
		assert.Equal(t, accAfter.Balance, accPrev.Balance)
		assert.Equal(t, accAfter.Frozen, accPrev.Frozen)
		// ccny
		assert.Equal(t, tokenAfter.Balance, tokenPrev.Balance-cost)
		assert.Equal(t, tokenAfter.Frozen, tokenPrev.Frozen+cost)
	} else {
		// dpos
		assert.Equal(t, accAfter.Balance, accPrev.Balance-cost)
		assert.Equal(t, accAfter.Frozen, accPrev.Frozen+cost)
		// ccny
		assert.Equal(t, tokenAfter.Balance, tokenPrev.Balance)
		assert.Equal(t, tokenAfter.Frozen, tokenPrev.Frozen)
	}
}

func testRevokeLimitOrder(t *testing.T, orderID int64, addr string, privkey string) {
	order, err := getOrder(orderID)
	assert.Nil(t, err)
	assert.NotNil(t, order)
	lo := order.Value.(*et.Order_LimitOrder).LimitOrder
	assert.NotNil(t, lo)

	accPrev, err := cli.GetExecAccount(addr, "coins", coin)
	assert.Nil(t, err)
	t.Log(accPrev)

	tokenPrev, err := cli.GetExecAccount(addr, "token", token)
	assert.Nil(t, err)
	t.Log(tokenPrev)

	_, err = doRevokeOrder(orderID, privkey)
	assert.Nil(t, err)

	accAfter, err := cli.GetExecAccount(addr, "coins", coin)
	assert.Nil(t, err)
	t.Log(accAfter)

	tokenAfter, err := cli.GetExecAccount(addr, "token", token)
	assert.Nil(t, err)
	t.Log(tokenAfter)

	cost := executor.CalcActualCost(lo.Op, order.Balance, lo.Price)
	// dpos/ccny
	if lo.Op == et.OpBuy {
		// dpos
		assert.Equal(t, accAfter.Balance, accPrev.Balance)
		assert.Equal(t, accAfter.Frozen, accPrev.Frozen)
		// ccny
		assert.Equal(t, tokenAfter.Balance, tokenPrev.Balance+cost)
		assert.Equal(t, tokenAfter.Frozen, tokenPrev.Frozen-cost)
	} else {
		// dpos
		assert.Equal(t, accAfter.Balance, accPrev.Balance+cost)
		assert.Equal(t, accAfter.Frozen, accPrev.Frozen-cost)
		// ccny
		assert.Equal(t, tokenAfter.Balance, tokenPrev.Balance)
		assert.Equal(t, tokenAfter.Frozen, tokenPrev.Frozen)
	}
}
