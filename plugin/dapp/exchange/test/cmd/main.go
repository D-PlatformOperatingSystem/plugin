package main

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/test"
	et "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
)

//setting ...
var (
	cli      *test.GRPCCli
	PrivKeyA = "0x6da92a632ab7deb67d38c0f6560bcfed28167998f6496db64c258d5e8393a81b" // 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4
	PrivKeyB = "0x19c069234f9d3e61135fefbeb7791b149cdf6af536f26bebb310d4cd22c3fee4" // 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR
)

//      ，
func main() {
	cli = test.NewGRPCCli("localhost:28804")
	go buy()
	go sell()
	select {}
}

func sell() {
	req := &et.LimitOrder{
		LeftAsset:  &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"},
		Price:      1,
		Amount:     types.Coin,
		Op:         et.OpSell,
	}
	ety := types.LoadExecutorType(et.ExchangeX)
	//   2000  ，  2000*1=2000   dpos
	for i := 0; i < 2000; i++ {
		fmt.Println("sell ", i)
		tx, err := ety.Create("LimitOrder", req)
		if err != nil {
			panic(err)
		}
		go cli.SendTx(tx, PrivKeyA)
	}
}

func buy() {
	req := &et.LimitOrder{
		LeftAsset:  &et.Asset{Symbol: "dpos", Execer: "coins"},
		RightAsset: &et.Asset{Execer: "token", Symbol: "CCNY"},
		Price:      1,
		Amount:     types.Coin,
		Op:         et.OpBuy,
	}
	ety := types.LoadExecutorType(et.ExchangeX)
	//   2000  ，  2000*1=2000   ccny
	for i := 0; i < 2000; i++ {
		fmt.Println("buy ", i)
		tx, err := ety.Create("LimitOrder", req)
		if err != nil {
			panic(err)
		}
		go cli.SendTx(tx, PrivKeyB)
	}
}
