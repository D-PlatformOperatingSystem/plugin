// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ticket

import (
	"crypto/ecdsa"
	"fmt"
	"testing"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	vrf "github.com/D-PlatformOperatingSystem/dpos/common/vrf/secp256k1"
	"github.com/D-PlatformOperatingSystem/dpos/queue"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket/types"
	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"

	apimocks "github.com/D-PlatformOperatingSystem/dpos/client/mocks"
	"github.com/D-PlatformOperatingSystem/dpos/common/merkle"
	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/consensus"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/init"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/store/init"
	"github.com/stretchr/testify/mock"
)

func TestTicket(t *testing.T) {
	testTicket(t)
}

func testTicket(t *testing.T) {
	mockDOM := testnode.New("testdata/dplatformos.cfg.toml", nil)
	defer mockDOM.Close()
	cfg := mockDOM.GetClient().GetConfig()
	mockDOM.Listen()
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
	tx = createBindMiner(cfg, t, hotaddr, addr, mockDOM.GetGenesisKey())
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
	start := time.Now()
	height := int64(0)
	hastclose := false
	hastopen := false
	for {
		height += 20
		err = mockDOM.WaitHeight(height)
		assert.Nil(t, err)
		//       close???
		req := &types.ReqWalletTransactionList{Count: 1000}
		resp, err := mockDOM.GetAPI().ExecWalletFunc("wallet", "WalletTransactionList", req)
		list := resp.(*types.WalletTxDetails)
		assert.Nil(t, err)
		for _, tx := range list.TxDetails {
			if tx.ActionName == "tclose" && tx.Receipt.Ty == 2 {
				hastclose = true
			}
			if tx.ActionName == "topen" && tx.Receipt.Ty == 2 {
				hastopen = true
			}
		}
		if hastopen == true && hastclose == true || time.Since(start) > 100*time.Second {
			break
		}
	}
	assert.Equal(t, true, hastclose)
	assert.Equal(t, true, hastopen)
	//
	accounts, err = acc.GetBalance(mockDOM.GetAPI(), &types.ReqBalance{Execer: "ticket", Addresses: []string{addr}})
	assert.Nil(t, err)
	fmt.Println(accounts[0])

	//         ,
	lastBlock := mockDOM.GetLastBlock()
	temblock := types.Clone(lastBlock)
	newblock := temblock.(*types.Block)
	newblock.GetTxs()[0].Nonce = newblock.GetTxs()[0].Nonce + 1
	newblock.TxHash = merkle.CalcMerkleRoot(cfg, newblock.GetHeight(), newblock.GetTxs())

	isbestBlock := util.CmpBestBlock(mockDOM.GetClient(), newblock, lastBlock.Hash(cfg))
	assert.Equal(t, isbestBlock, false)
}

func createBindMiner(cfg *types.DplatformOSConfig, t *testing.T, m, r string, priv crypto.PrivKey) *types.Transaction {
	ety := types.LoadExecutorType("ticket")
	tx, err := ety.Create("Tbind", &ty.TicketBind{MinerAddress: m, ReturnAddress: r})
	assert.Nil(t, err)
	tx, err = types.FormatTx(cfg, "ticket", tx)
	assert.Nil(t, err)
	tx.Sign(types.SECP256K1, priv)
	return tx
}

func TestTicketMap(t *testing.T) {
	c := Client{}
	ticketList := &ty.ReplyTicketList{}
	ticketList.Tickets = []*ty.Ticket{
		{TicketId: "1111"},
		{TicketId: "2222"},
		{TicketId: "3333"},
		{TicketId: "4444"},
	}
	privmap := make(map[string]crypto.PrivKey)
	//  privkey    pubkey        addr
	cr, _ := crypto.New("secp256k1")
	priv, _ := cr.PrivKeyFromBytes([]byte("2116459C0EC8ED01AA0EEAE35CAC5C96F94473F7816F114873291217303F6989"))
	privmap["1111"] = priv
	privmap["2222"] = priv
	privmap["3333"] = priv
	privmap["4444"] = priv

	assert.Equal(t, c.getTicketCount(), int64(0))
	c.setTicket(ticketList, privmap)
	assert.Equal(t, c.getTicketCount(), int64(4))
	c.delTicket("3333")
	assert.Equal(t, c.getTicketCount(), int64(3))

	c.setTicket(ticketList, nil)
	assert.Equal(t, c.getTicketCount(), int64(0))

	c.setTicket(nil, privmap)
	assert.Equal(t, c.getTicketCount(), int64(0))

	c.setTicket(nil, nil)
	assert.Equal(t, c.getTicketCount(), int64(0))
	_, err := c.Query_GetTicketCount(&types.ReqNil{})
	assert.Nil(t, err)
}

func TestProcEvent(t *testing.T) {
	c := Client{}
	ret := c.ProcEvent(&queue.Message{})
	assert.Equal(t, ret, true)
}

func Test_genPrivHash(t *testing.T) {
	c, err := crypto.New(types.GetSignName("", types.SECP256K1))
	assert.NoError(t, err)
	priv, _ := c.GenKey()

	bt, err := genPrivHash(priv, "AA:BB:CC:DD")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(bt))

	bt, err = genPrivHash(priv, "111:222:333:444")
	assert.NoError(t, err)
	assert.Equal(t, 32, len(bt))
}

func Test_getNextRequiredDifficulty(t *testing.T) {
	cfg := types.NewDplatformOSConfig(types.ReadFile("testdata/dplatformos.cfg.toml"))

	api := new(apimocks.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(cfg, nil)
	c := &Client{BaseClient: &drivers.BaseClient{}}
	c.SetAPI(api)

	bits, bt, err := c.getNextRequiredDifficulty(nil, 1)
	assert.NoError(t, err)
	assert.Equal(t, bt, defaultModify)
	assert.Equal(t, bits, cfg.GetP(0).PowLimitBits)
}

func Test_vrfVerify(t *testing.T) {
	c, err := crypto.New(types.GetSignName("", types.SECP256K1))
	assert.NoError(t, err)
	priv, err := c.GenKey()
	assert.NoError(t, err)
	pub := priv.PubKey().Bytes()

	privKey, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), priv.Bytes())
	vpriv := &vrf.PrivateKey{PrivateKey: (*ecdsa.PrivateKey)(privKey)}

	m1 := []byte("data1")
	m2 := []byte("data2")
	m3 := []byte("data2")
	hash1, proof1 := vpriv.Evaluate(m1)
	hash2, proof2 := vpriv.Evaluate(m2)
	hash3, proof3 := vpriv.Evaluate(m3)
	for _, tc := range []struct {
		m     []byte
		hash  [32]byte
		proof []byte
		err   error
	}{
		{m1, hash1, proof1, nil},
		{m2, hash2, proof2, nil},
		{m3, hash3, proof3, nil},
		{m3, hash3, proof2, nil},
		{m3, hash3, proof1, ty.ErrVrfVerify},
		{m3, hash1, proof3, ty.ErrVrfVerify},
	} {
		err := vrfVerify(pub, tc.m, tc.proof, tc.hash[:])
		if got, want := err, tc.err; got != want {
			t.Errorf("vrfVerify(%s, %x): %v, want %v", tc.m, tc.proof, got, want)
		}
	}
}
