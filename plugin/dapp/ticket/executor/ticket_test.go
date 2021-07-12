package executor_test

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket/executor"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket/types"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/consensus/init"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket"
)

var mockDOM *testnode.DplatformOSMock

func TestMain(m *testing.M) {
	mockDOM = testnode.New("testdata/dplatformos.cfg.toml", nil)
	mockDOM.Listen()
	code := m.Run()
	mockDOM.Close()
	os.Exit(code)
}

func TestTicketPrice(t *testing.T) {
	cfg := mockDOM.GetAPI().GetConfig()
	//test price
	ti := &executor.DB{}
	assert.Equal(t, ti.GetRealPrice(cfg), 10000*types.Coin)

	ti = &executor.DB{}
	ti.Price = 10
	assert.Equal(t, ti.GetRealPrice(cfg), int64(10))
}

func TestCheckFork(t *testing.T) {
	cfg := mockDOM.GetAPI().GetConfig()
	assert.Equal(t, int64(1), cfg.GetFork("ForkChainParamV2"))
	p1 := ty.GetTicketMinerParam(cfg, 0)
	assert.Equal(t, 10000*types.Coin, p1.TicketPrice)
	p1 = ty.GetTicketMinerParam(cfg, 1)
	assert.Equal(t, 3000*types.Coin, p1.TicketPrice)
	p1 = ty.GetTicketMinerParam(cfg, 2)
	assert.Equal(t, 3000*types.Coin, p1.TicketPrice)
	p1 = ty.GetTicketMinerParam(cfg, 3)
	assert.Equal(t, 3000*types.Coin, p1.TicketPrice)
}

func TestTicket(t *testing.T) {
	cfg := mockDOM.GetAPI().GetConfig()
	reply, err := mockDOM.GetAPI().ExecWalletFunc("ticket", "WalletAutoMiner", &ty.MinerFlag{Flag: 1})
	assert.Nil(t, err)
	assert.Equal(t, true, reply.(*types.Reply).IsOk)
	acc := account.NewCoinsAccount(cfg)
	addr := mockDOM.GetGenesisAddress()
	accounts, err := acc.GetBalance(mockDOM.GetAPI(), &types.ReqBalance{Execer: "ticket", Addresses: []string{addr}})
	assert.Nil(t, err)
	assert.Equal(t, accounts[0].Balance, int64(0))
	hotaddr := mockDOM.GetHotAddress()
	_, err = acc.GetBalance(mockDOM.GetAPI(), &types.ReqBalance{Execer: "coins", Addresses: []string{hotaddr}})
	assert.Nil(t, err)
	//assert.Equal(t, accounts[0].Balance, int64(1000000000000))
	//send to address
	tx := util.CreateCoinsTx(cfg, mockDOM.GetHotKey(), mockDOM.GetGenesisAddress(), types.Coin/100)
	mockDOM.SendTx(tx)
	mockDOM.Wait()
	//bind miner
	tx = createBindMiner(t, cfg, hotaddr, addr, mockDOM.GetGenesisKey())
	hash := mockDOM.SendTx(tx)
	detail, err := mockDOM.WaitTx(hash)
	assert.Nil(t, err)
	//debug:
	//js, _ := json.MarshalIndent(detail, "", " ")
	//fmt.Println(string(js))
	_, err = mockDOM.GetAPI().ExecWalletFunc("ticket", "WalletAutoMiner", &ty.MinerFlag{Flag: 0})
	assert.Nil(t, err)
	status, err := mockDOM.GetAPI().ExecWalletFunc("wallet", "GetWalletStatus", &types.ReqNil{})
	assert.Nil(t, err)
	assert.Equal(t, false, status.(*types.WalletStatus).IsAutoMining)
	assert.Equal(t, int32(2), detail.Receipt.Ty)
	_, err = mockDOM.GetAPI().ExecWalletFunc("ticket", "WalletAutoMiner", &ty.MinerFlag{Flag: 1})
	assert.Nil(t, err)
	status, err = mockDOM.GetAPI().ExecWalletFunc("wallet", "GetWalletStatus", &types.ReqNil{})
	assert.Nil(t, err)
	assert.Equal(t, true, status.(*types.WalletStatus).IsAutoMining)
	for i := mockDOM.GetLastBlock().Height; i < 100; i++ {
		err = mockDOM.WaitHeight(i)
		assert.Nil(t, err)
		//       close，
		req := &types.ReqWalletTransactionList{Count: 1000}
		resp, err := mockDOM.GetAPI().ExecWalletFunc("wallet", "WalletTransactionList", req)
		assert.Nil(t, err)
		list := resp.(*types.WalletTxDetails)
		hastclose := false
		hastopen := false
		for _, tx := range list.TxDetails {
			if tx.Height < 1 {
				continue
			}
			if tx.ActionName == "tclose" && tx.Receipt.Ty == 2 {
				hastclose = true
			}
			if tx.ActionName == "topen" && tx.Receipt.Ty == 2 {
				hastopen = true
				fmt.Println(tx)
				list := ticketList(t, mockDOM, &ty.TicketList{Addr: tx.Fromaddr, Status: 1})
				for _, ti := range list.GetTickets() {
					if strings.Contains(ti.TicketId, hex.EncodeToString(tx.Txhash)) {
						assert.Equal(t, 3000*types.Coin, ti.Price)
					}
				}
			}
		}
		if hastclose && hastopen {
			return
		}
	}
	t.Error("wait 100 , open and close not happened")
}

func createBindMiner(t *testing.T, cfg *types.DplatformOSConfig, m, r string, priv crypto.PrivKey) *types.Transaction {
	ety := types.LoadExecutorType("ticket")
	tx, err := ety.Create("Tbind", &ty.TicketBind{MinerAddress: m, ReturnAddress: r})
	assert.Nil(t, err)
	tx, err = types.FormatTx(cfg, "ticket", tx)
	assert.Nil(t, err)
	tx.Sign(types.SECP256K1, priv)
	return tx
}

func ticketList(t *testing.T, mockDOM *testnode.DplatformOSMock, req proto.Message) *ty.ReplyTicketList {
	data, err := mockDOM.GetAPI().Query("ticket", "TicketList", req)
	assert.Nil(t, err)
	return data.(*ty.ReplyTicketList)
}
