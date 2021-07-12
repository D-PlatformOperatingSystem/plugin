package executor

import (
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/common/db"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util"

	"github.com/D-PlatformOperatingSystem/dpos/client"
	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	"github.com/D-PlatformOperatingSystem/dpos/queue"
	et "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
	"github.com/stretchr/testify/assert"
)

type execEnv struct {
	blockTime   int64
	blockHeight int64
	difficulty  uint64
}

var (
	PrivKeyA = "0x6da92a632ab7deb67d38c0f6560bcfed28167998f6496db64c258d5e8393a81b" // 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4
	PrivKeyB = "0x19c069234f9d3e61135fefbeb7791b149cdf6af536f26bebb310d4cd22c3fee4" // 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR
	PrivKeyC = "0x7a80a1f75d7360c6123c32a78ecf978c1ac55636f87892df38d8b85a9aeff115" // 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k
	PrivKeyD = "0xcacb1f5d51700aea07fca2246ab43b0917d70405c65edea9b5063d72eb5c6b71" // 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs
	Nodes    = []string{
		"1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4",
		"1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR",
		"1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k",
		"1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs",
	}
)

func TestExchange(t *testing.T) {
	//
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	Init(et.ExchangeX, cfg, nil)
	total := 100 * types.Coin
	accountA := types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[0],
	}
	accountB := types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[1],
	}

	accountC := types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[2],
	}
	accountD := types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[3],
	}
	dir, stateDB, kvdb := util.CreateTestDB()
	//defer util.CloseTestDB(dir, stateDB)
	execAddr := address.ExecAddress(et.ExchangeX)

	accA, _ := account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accA.SaveExecAccount(execAddr, &accountA)

	accB, _ := account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accB.SaveExecAccount(execAddr, &accountB)

	accC, _ := account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accC.SaveExecAccount(execAddr, &accountC)

	accD, _ := account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accD.SaveExecAccount(execAddr, &accountD)

	accA1, _ := account.NewAccountDB(cfg, "token", "CCNY", stateDB)
	accA1.SaveExecAccount(execAddr, &accountA)

	accB1, _ := account.NewAccountDB(cfg, "paracross", "coins.dpos", stateDB)
	accB1.SaveExecAccount(execAddr, &accountB)

	accC1, _ := account.NewAccountDB(cfg, "paracross", "token.CCNY", stateDB)
	accC1.SaveExecAccount(execAddr, &accountC)

	accD1, _ := account.NewAccountDB(cfg, "token", "CCNY", stateDB)
	accD1.SaveExecAccount(execAddr, &accountD)

	env := &execEnv{
		10,
		1,
		1539918074,
	}

	/*
	         ，        ，
	     ：
	   1.     10   。
	   2.       5
	   3.
	*/

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 4, Amount: 10 * types.Coin, Op: et.OpBuy}, PrivKeyA, stateDB, kvdb, env)
	//          ,          list[0],
	orderList, err := Exec_QueryOrderList(et.Ordered, Nodes[0], "", stateDB, kvdb)
	assert.Equal(t, nil, err)

	orderID1 := orderList.List[0].OrderID
	//     ，
	order, err := Exec_QueryOrder(orderID1, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Ordered), order.Status)
	assert.Equal(t, 10*types.Coin, order.GetBalance())

	//  op
	marketDepthList, err := Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpBuy}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, 10*types.Coin, marketDepthList.List[0].GetAmount())

	//
	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 4, Amount: 5 * types.Coin, Op: et.OpSell}, PrivKeyB, stateDB, kvdb, env)
	//          ,          list[0],
	orderList, err = Exec_QueryOrderList(et.Completed, Nodes[1], "", stateDB, kvdb)
	assert.Equal(t, nil, err)
	orderID2 := orderList.List[0].OrderID
	//    1
	order, err = Exec_QueryOrder(orderID1, stateDB, kvdb)
	assert.Equal(t, nil, err)
	//  1       ordered
	assert.Equal(t, int32(et.Ordered), order.Status)
	assert.Equal(t, 5*types.Coin, order.Balance)

	//order2   completed
	order, err = Exec_QueryOrder(orderID2, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Completed), order.Status)
	//  op
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpBuy}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	//
	assert.Equal(t, 5*types.Coin, marketDepthList.List[0].GetAmount())

	//QueryHistoryOrderList
	orderList, err = Exec_QueryHistoryOrder(&et.QueryHistoryOrderList{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, orderID2, orderList.List[0].OrderID)
	//
	Exec_RevokeOrder(t, orderID1, PrivKeyA, stateDB, kvdb, env)
	//     ，
	order, err = Exec_QueryOrder(orderID1, stateDB, kvdb)
	assert.Equal(t, nil, err)
	//  1     Revoked
	assert.Equal(t, int32(et.Revoked), order.Status)
	assert.Equal(t, 5*types.Coin, order.Balance)

	//  op
	//  dpos,CCNY     ,
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpBuy}, stateDB, kvdb)
	assert.NotEqual(t, nil, err)
	//
	//  ordered
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[0], "", stateDB, kvdb)
	assert.Equal(t, types.ErrNotFound, err)

	/*
			       ，        ，

			     ：
			   1.     10   。
		       2.       10
		       3.     5
		       4.     15
	*/
	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"}, RightAsset: &et.Asset{Execer: "paracross", Symbol: "coins.dpos"}, Price: 50000000, Amount: 10 * types.Coin, Op: et.OpSell}, PrivKeyA, stateDB, kvdb, env)
	//
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[0], "", stateDB, kvdb)
	assert.Nil(t, err)
	orderID3 := orderList.List[0].OrderID

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "paracross", Symbol: "coins.dpos"}, Price: 50000000, Amount: 10 * types.Coin, Op: et.OpSell}, PrivKeyA, stateDB, kvdb, env)
	//
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[0], "", stateDB, kvdb)
	assert.Nil(t, err)
	orderID4 := orderList.List[0].OrderID

	//  op
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "paracross", Symbol: "coins.dpos"}, Op: et.OpSell}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	//
	assert.Equal(t, 20*types.Coin, marketDepthList.List[0].GetAmount())

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "paracross", Symbol: "coins.dpos"}, Price: 50000000, Amount: 5 * types.Coin, Op: et.OpBuy}, PrivKeyB, stateDB, kvdb, env)
	//
	orderList, err = Exec_QueryOrderList(et.Completed, Nodes[1], "", stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Completed), orderList.List[1].Status)
	//
	//    3
	order, err = Exec_QueryOrder(orderID3, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Ordered), order.Status)
	//
	assert.Equal(t, 5*types.Coin, order.Balance)
	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "paracross", Symbol: "coins.dpos"}, Price: 50000000, Amount: 15 * types.Coin, Op: et.OpBuy}, PrivKeyB, stateDB, kvdb, env)
	//order3,order4
	order, err = Exec_QueryOrder(orderID3, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Completed), order.Status)
	order, err = Exec_QueryOrder(orderID4, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Completed), order.Status)
	//  op      ,
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "paracross", Symbol: "coins.dpos"}, Op: et.OpSell}, stateDB, kvdb)
	assert.Equal(t, types.ErrNotFound, err)

	/*
	            /
	     ：
	   1.     5,   4
	   2.     10,   3
	   3.     5,   4
	   4.     5,   5
	   5.    15,   4.5
	   6.     2,   1
	   7.     8,   1
	   8.
	*/

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 400000000, Amount: 5 * types.Coin, Op: et.OpBuy}, PrivKeyD, stateDB, kvdb, env)
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[3], "", stateDB, kvdb)
	assert.Nil(t, err)
	orderID6 := orderList.List[0].OrderID

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 300000000, Amount: 10 * types.Coin, Op: et.OpSell}, PrivKeyC, stateDB, kvdb, env)
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[2], "", stateDB, kvdb)
	assert.Nil(t, err)
	orderID7 := orderList.List[0].OrderID

	//    6
	order, err = Exec_QueryOrder(orderID6, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Completed), order.Status)

	order, err = Exec_QueryOrder(orderID7, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Ordered), order.Status)
	//      ,
	acc := accD1.LoadExecAccount(Nodes[3], execAddr)
	assert.Equal(t, 85*types.Coin, acc.Balance)

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 400000000, Amount: 5 * types.Coin, Op: et.OpSell}, PrivKeyC, stateDB, kvdb, env)
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[2], "", stateDB, kvdb)
	assert.Nil(t, err)
	orderID8 := orderList.List[0].OrderID
	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 500000000, Amount: 5 * types.Coin, Op: et.OpSell}, PrivKeyC, stateDB, kvdb, env)
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[2], "", stateDB, kvdb)
	assert.Nil(t, err)
	orderID9 := orderList.List[0].OrderID

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 450000000, Amount: 15 * types.Coin, Op: et.OpBuy}, PrivKeyD, stateDB, kvdb, env)
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[3], "", stateDB, kvdb)
	//orderID10 := orderList.List[0].OrderID
	assert.Equal(t, 5*types.Coin, orderList.List[0].Balance)
	assert.Nil(t, err)
	//order7 order8
	order, err = Exec_QueryOrder(orderID7, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Completed), order.Status)
	order, err = Exec_QueryOrder(orderID8, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Completed), order.Status)

	order, err = Exec_QueryOrder(orderID9, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(et.Ordered), order.Status)
	assert.Equal(t, 5*types.Coin, order.Balance)
	//
	acc = accD1.LoadExecAccount(Nodes[3], execAddr)
	// 100-3*10-5*4-4.5*5   = 27.5
	assert.Equal(t, int64(2750000000), acc.Balance)
	acc = accC.LoadExecAccount(Nodes[2], execAddr)
	assert.Equal(t, 80*types.Coin, acc.Balance)

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 100000000, Amount: 2 * types.Coin, Op: et.OpSell}, PrivKeyC, stateDB, kvdb, env)
	orderList, err = Exec_QueryOrderList(et.Completed, Nodes[2], "", stateDB, kvdb)
	assert.Nil(t, err)
	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 100000000, Amount: 8 * types.Coin, Op: et.OpSell}, PrivKeyC, stateDB, kvdb, env)
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[2], "", stateDB, kvdb)
	assert.Nil(t, err)
	orderID10 := orderList.List[0].OrderID
	assert.Equal(t, int32(et.Ordered), orderList.List[0].Status)
	assert.Equal(t, 5*types.Coin, orderList.List[0].Balance)
	//
	acc = accD1.LoadExecAccount(Nodes[3], execAddr)
	// 100-3*10-5*4-1*5   = 45
	assert.Equal(t, 45*types.Coin, acc.Balance)
	acc = accC.LoadExecAccount(Nodes[2], execAddr)
	assert.Equal(t, 70*types.Coin, acc.Balance)
	//orderID9 order10
	Exec_RevokeOrder(t, orderID9, PrivKeyC, stateDB, kvdb, env)
	Exec_RevokeOrder(t, orderID10, PrivKeyC, stateDB, kvdb, env)
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpSell}, stateDB, kvdb)
	assert.NotEqual(t, nil, err)
	acc = accC.LoadExecAccount(Nodes[2], execAddr)
	assert.Equal(t, 80*types.Coin, acc.Balance)

	//    ,
	util.CloseTestDB(dir, stateDB)
	total = 1000 * types.Coin
	accountA = types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[0],
	}
	accountB = types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[1],
	}

	dir, stateDB, kvdb = util.CreateTestDB()
	defer util.CloseTestDB(dir, stateDB)
	//execAddr := address.ExecAddress(et.ExchangeX)

	accA, _ = account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accA.SaveExecAccount(execAddr, &accountA)

	accB, _ = account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accB.SaveExecAccount(execAddr, &accountB)

	accA1, _ = account.NewAccountDB(cfg, "token", "CCNY", stateDB)
	accA1.SaveExecAccount(execAddr, &accountA)

	accB1, _ = account.NewAccountDB(cfg, "token", "CCNY", stateDB)
	accB1.SaveExecAccount(execAddr, &accountB)

	env = &execEnv{
		10,
		1,
		1539918074,
	}
	/*
	        ：
	      :
	    1.  200 ，   1   1
	    2.       1,   200
	    3.         ,
	    4.
	*/

	for i := 0; i < 200; i++ {
		Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 100000000, Amount: 1 * types.Coin, Op: et.OpBuy}, PrivKeyA, stateDB, kvdb, env)
	}

	Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 100000000, Amount: 200 * types.Coin, Op: et.OpSell}, PrivKeyB, stateDB, kvdb, env)
	if et.MaxMatchCount > 200 {
		return
	}
	orderList, err = Exec_QueryOrderList(et.Ordered, Nodes[1], "", stateDB, kvdb)
	orderID := orderList.List[0].OrderID
	assert.Equal(t, nil, err)
	assert.Equal(t, (200-et.MaxMatchCount)*types.Coin, orderList.List[0].Balance)
	//  op
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpBuy}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, (200-et.MaxMatchCount)*types.Coin, marketDepthList.List[0].GetAmount())
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpSell}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, (200-et.MaxMatchCount)*types.Coin, marketDepthList.List[0].GetAmount())

	//
	//
	var count int
	var primaryKey string
	for {
		orderList, err := Exec_QueryOrderList(et.Completed, Nodes[0], primaryKey, stateDB, kvdb)
		if err != nil {
			break
		}
		count = count + len(orderList.List)
		if orderList.PrimaryKey == "" {
			break
		}
		primaryKey = orderList.PrimaryKey
	}
	assert.Equal(t, et.MaxMatchCount, count)

	//
	//
	count = 0
	primaryKey = ""
	for {
		query := &et.QueryHistoryOrderList{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, PrimaryKey: primaryKey}
		orderList, err := Exec_QueryHistoryOrder(query, stateDB, kvdb)
		if err != nil {
			break
		}
		count = count + len(orderList.List)
		if orderList.PrimaryKey == "" {
			break
		}
		primaryKey = orderList.PrimaryKey
	}
	assert.Equal(t, et.MaxMatchCount, count)
	//         ,
	err = Exec_LimitOrder(t, &et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: 100000000, Amount: 100 * types.Coin, Op: et.OpSell}, PrivKeyA, stateDB, kvdb, env)
	assert.Equal(t, nil, err)
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpSell}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, (200-et.MaxMatchCount+100)*types.Coin, marketDepthList.List[0].GetAmount())
	//
	err = Exec_RevokeOrder(t, orderID, PrivKeyA, stateDB, kvdb, env)
	assert.NotEqual(t, nil, err)
	err = Exec_RevokeOrder(t, orderID, PrivKeyB, stateDB, kvdb, env)
	assert.Equal(t, nil, err)

	//    ,
	util.CloseTestDB(dir, stateDB)
	total = 1000 * types.Coin
	accountA = types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[0],
	}
	accountB = types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[1],
	}

	dir, stateDB, kvdb = util.CreateTestDB()
	defer util.CloseTestDB(dir, stateDB)
	//execAddr := address.ExecAddress(et.ExchangeX)

	accA, _ = account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accA.SaveExecAccount(execAddr, &accountA)

	accB, _ = account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accB.SaveExecAccount(execAddr, &accountB)

	accA1, _ = account.NewAccountDB(cfg, "token", "CCNY", stateDB)
	accA1.SaveExecAccount(execAddr, &accountA)

	accB1, _ = account.NewAccountDB(cfg, "token", "CCNY", stateDB)
	accB1.SaveExecAccount(execAddr, &accountB)

	env = &execEnv{
		10,
		1,
		1539918074,
	}
	/*
	  //    ,
	      :
	    1.       ,    ：
	        10
	        20
	        50
	        20
	        50
	        100

	*/
	acc = accB1.LoadExecAccount(Nodes[1], execAddr)
	assert.Equal(t, 1000*types.Coin, acc.Balance)
	acc = accA.LoadExecAccount(Nodes[0], execAddr)
	assert.Equal(t, 1000*types.Coin, acc.Balance)
	var txs []*types.Transaction
	for i := 0; i < 10; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpBuy}, PrivKeyB)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}

	for i := 0; i < 20; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpSell}, PrivKeyA)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}

	for i := 0; i < 50; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpBuy}, PrivKeyB)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}

	for i := 0; i < 20; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpSell}, PrivKeyA)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}
	for i := 0; i < 50; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpBuy}, PrivKeyB)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}

	for i := 0; i < 100; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpSell}, PrivKeyA)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}
	err = Exec_Block(t, stateDB, kvdb, env, txs...)
	assert.Equal(t, nil, err)
	acc = accB1.LoadExecAccount(Nodes[1], execAddr)
	assert.Equal(t, 890*types.Coin, acc.Balance)
	acc = accA.LoadExecAccount(Nodes[0], execAddr)
	assert.Equal(t, 860*types.Coin, acc.Balance)
	assert.Equal(t, 30*types.Coin, acc.Frozen)

	//  op
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpBuy}, stateDB, kvdb)
	assert.NotEqual(t, nil, err)
	//assert.Equal(t, (200-et.MaxMatchCount)*types.Coin, marketDepthList.List[0].GetAmount())
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpSell}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, 30*types.Coin, marketDepthList.List[0].GetAmount())

	//    ,
	util.CloseTestDB(dir, stateDB)
	total = 1000 * types.Coin
	accountA = types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[0],
	}
	accountB = types.Account{
		Balance: total,
		Frozen:  0,
		Addr:    Nodes[1],
	}

	dir, stateDB, kvdb = util.CreateTestDB()
	defer util.CloseTestDB(dir, stateDB)
	//execAddr := address.ExecAddress(et.ExchangeX)

	accA, _ = account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accA.SaveExecAccount(execAddr, &accountA)

	accB, _ = account.NewAccountDB(cfg, "coins", "dpos", stateDB)
	accB.SaveExecAccount(execAddr, &accountB)

	accA1, _ = account.NewAccountDB(cfg, "token", "CCNY", stateDB)
	accA1.SaveExecAccount(execAddr, &accountA)

	accB1, _ = account.NewAccountDB(cfg, "token", "CCNY", stateDB)
	accB1.SaveExecAccount(execAddr, &accountB)

	env = &execEnv{
		10,
		1,
		1539918074,
	}
	/*
			  //    ,
			      :
			    1.      ,    ：
		            100
			        50
			        20
			        100
	*/
	acc = accB1.LoadExecAccount(Nodes[1], execAddr)
	assert.Equal(t, 1000*types.Coin, acc.Balance)
	acc = accA.LoadExecAccount(Nodes[0], execAddr)
	assert.Equal(t, 1000*types.Coin, acc.Balance)
	txs = nil
	for i := 0; i < 100; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpSell}, PrivKeyA)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}

	for i := 0; i < 50; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpBuy}, PrivKeyB)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}

	for i := 0; i < 20; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpSell}, PrivKeyA)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}
	for i := 0; i < 100; i++ {
		tx, err := CreateLimitOrder(&et.LimitOrder{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
			RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Price: types.Coin, Amount: types.Coin, Op: et.OpBuy}, PrivKeyB)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)
	}
	err = Exec_Block(t, stateDB, kvdb, env, txs...)
	assert.Equal(t, nil, err)
	acc = accB1.LoadExecAccount(Nodes[1], execAddr)
	assert.Equal(t, 850*types.Coin, acc.Balance)
	assert.Equal(t, 30*types.Coin, acc.Frozen)
	acc = accA.LoadExecAccount(Nodes[0], execAddr)
	assert.Equal(t, 880*types.Coin, acc.Balance)

	//  op
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpSell}, stateDB, kvdb)
	assert.NotEqual(t, nil, err)
	marketDepthList, err = Exec_QueryMarketDepth(&et.QueryMarketDepth{LeftAsset: &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"}, Op: et.OpBuy}, stateDB, kvdb)
	assert.Equal(t, nil, err)
	assert.Equal(t, 30*types.Coin, marketDepthList.List[0].GetAmount())

	//
	//
	count = 0
	primaryKey = ""
	for {
		orderList, err := Exec_QueryOrderList(et.Completed, Nodes[1], primaryKey, stateDB, kvdb)
		if err != nil {
			break
		}
		count = count + len(orderList.List)
		if orderList.PrimaryKey == "" {
			break
		}
		primaryKey = orderList.PrimaryKey
	}
	assert.Equal(t, 120, count)

	count = 0
	primaryKey = ""
	for {
		orderList, err := Exec_QueryOrderList(et.Ordered, Nodes[1], primaryKey, stateDB, kvdb)
		if err != nil {
			break
		}
		count = count + len(orderList.List)
		if orderList.PrimaryKey == "" {
			break
		}
		primaryKey = orderList.PrimaryKey
	}
	assert.Equal(t, 30, count)

}

func CreateLimitOrder(limitOrder *et.LimitOrder, privKey string) (tx *types.Transaction, err error) {
	ety := types.LoadExecutorType(et.ExchangeX)
	tx, err = ety.Create("LimitOrder", limitOrder)
	if err != nil {
		return nil, err
	}
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	tx, err = types.FormatTx(cfg, et.ExchangeX, tx)
	if err != nil {
		return nil, err
	}
	tx, err = signTx(tx, privKey)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
func CreateRevokeOrder(orderID int64, privKey string) (tx *types.Transaction, err error) {
	ety := types.LoadExecutorType(et.ExchangeX)
	tx, err = ety.Create("RevokeOrder", &et.RevokeOrder{OrderID: orderID})
	if err != nil {
		return nil, err
	}
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	tx, err = types.FormatTx(cfg, et.ExchangeX, tx)
	if err != nil {
		return nil, err
	}
	tx, err = signTx(tx, privKey)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

//
func Exec_Block(t *testing.T, stateDB db.DB, kvdb db.KVDB, env *execEnv, txs ...*types.Transaction) error {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	exec := NewExchange()
	e := exec.(*exchange)
	for index, tx := range txs {
		err := e.CheckTx(tx, index)
		if err != nil {
			t.Log(err.Error())
			return err
		}

	}
	q := queue.New("channel")
	q.SetConfig(cfg)
	api, _ := client.New(q.Client(), nil)
	exec.SetAPI(api)
	exec.SetStateDB(stateDB)
	exec.SetLocalDB(kvdb)
	env.blockHeight = env.blockHeight + 1
	env.blockTime = env.blockTime + 20
	env.difficulty = env.difficulty + 1
	exec.SetEnv(env.blockHeight, env.blockTime, env.difficulty)
	for index, tx := range txs {
		receipt, err := exec.Exec(tx, index)
		if err != nil {
			t.Log(err.Error())
			return err
		}
		for _, kv := range receipt.KV {
			stateDB.Set(kv.Key, kv.Value)
		}
		receiptData := &types.ReceiptData{Ty: receipt.Ty, Logs: receipt.Logs}
		set, err := exec.ExecLocal(tx, receiptData, index)
		if err != nil {
			t.Log(err.Error())
			return err
		}
		for _, kv := range set.KV {
			kvdb.Set(kv.Key, kv.Value)
		}
		//save to database
		util.SaveKVList(stateDB, set.KV)
		assert.Equal(t, types.ExecOk, int(receipt.Ty))
	}
	return nil
}
func Exec_LimitOrder(t *testing.T, limitOrder *et.LimitOrder, privKey string, stateDB db.DB, kvdb db.KVDB, env *execEnv) error {
	tx, err := CreateLimitOrder(limitOrder, privKey)
	if err != nil {
		return err
	}
	return Exec_Block(t, stateDB, kvdb, env, tx)
}

func Exec_RevokeOrder(t *testing.T, orderID int64, privKey string, stateDB db.DB, kvdb db.KVDB, env *execEnv) error {
	tx, err := CreateRevokeOrder(orderID, privKey)
	if err != nil {
		return err
	}
	return Exec_Block(t, stateDB, kvdb, env, tx)
}

func Exec_QueryOrderList(status int32, addr string, primaryKey string, stateDB db.KV, kvdb db.KVDB) (*et.OrderList, error) {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	exec := NewExchange()
	q := queue.New("channel")
	q.SetConfig(cfg)
	api, _ := client.New(q.Client(), nil)
	exec.SetAPI(api)
	exec.SetStateDB(stateDB)
	exec.SetLocalDB(kvdb)
	//
	msg, err := exec.Query(et.FuncNameQueryOrderList, types.Encode(&et.QueryOrderList{Status: status, Address: addr, PrimaryKey: primaryKey}))
	if err != nil {
		return nil, err
	}
	return msg.(*et.OrderList), nil
}
func Exec_QueryOrder(orderID int64, stateDB db.KV, kvdb db.KVDB) (*et.Order, error) {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	exec := NewExchange()
	q := queue.New("channel")
	q.SetConfig(cfg)
	api, _ := client.New(q.Client(), nil)
	exec.SetAPI(api)
	exec.SetStateDB(stateDB)
	exec.SetLocalDB(kvdb)
	//  orderID
	msg, err := exec.Query(et.FuncNameQueryOrder, types.Encode(&et.QueryOrder{OrderID: orderID}))
	if err != nil {
		return nil, err
	}
	return msg.(*et.Order), err
}

func Exec_QueryMarketDepth(query *et.QueryMarketDepth, stateDB db.KV, kvdb db.KVDB) (*et.MarketDepthList, error) {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	exec := NewExchange()
	q := queue.New("channel")
	q.SetConfig(cfg)
	api, _ := client.New(q.Client(), nil)
	exec.SetAPI(api)
	exec.SetStateDB(stateDB)
	exec.SetLocalDB(kvdb)
	//  QueryMarketDepth
	msg, err := exec.Query(et.FuncNameQueryMarketDepth, types.Encode(query))
	if err != nil {
		return nil, err
	}
	return msg.(*et.MarketDepthList), err
}

func Exec_QueryHistoryOrder(query *et.QueryHistoryOrderList, stateDB db.KV, kvdb db.KVDB) (*et.OrderList, error) {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	exec := NewExchange()
	q := queue.New("channel")
	q.SetConfig(cfg)
	api, _ := client.New(q.Client(), nil)
	exec.SetAPI(api)
	exec.SetStateDB(stateDB)
	exec.SetLocalDB(kvdb)
	//  QueryMarketDepth
	msg, err := exec.Query(et.FuncNameQueryHistoryOrderList, types.Encode(query))
	return msg.(*et.OrderList), err
}
func signTx(tx *types.Transaction, hexPrivKey string) (*types.Transaction, error) {
	signType := types.SECP256K1
	c, err := crypto.New(types.GetSignName("", signType))
	if err != nil {
		return tx, err
	}

	bytes, err := common.FromHex(hexPrivKey[:])
	if err != nil {
		return tx, err
	}

	privKey, err := c.PrivKeyFromBytes(bytes)
	if err != nil {
		return tx, err
	}

	tx.Sign(int32(signType), privKey)
	return tx, nil
}

func TestCheckPrice(t *testing.T) {
	t.Log(CheckPrice(1e8))
	t.Log(CheckPrice(-1))
	t.Log(CheckPrice(1e17))
	t.Log(CheckPrice(0))
}

func TestRawMeta(t *testing.T) {
	HistoryOrderRow := NewHistoryOrderRow()
	t.Log(HistoryOrderRow.Get("index"))
	MarketDepthRow := NewMarketDepthRow()
	t.Log(MarketDepthRow.Get("price"))
	marketOrderRow := NewOrderRow()
	t.Log(marketOrderRow.Get("orderID"))
}

func TestKV(t *testing.T) {
	a := &types.KeyValue{Key: []byte("1111111"), Value: nil}
	t.Log(a.Key, a.Value)
}

func TestSafeMul(t *testing.T) {
	t.Log(SafeMul(1e8, 1e7))
	t.Log(SafeMul(1e10, 1e16))
	t.Log(SafeMul(1e7, 1e6))
}
