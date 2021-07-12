// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet_test

import (
	"fmt"
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/blockchain"
	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/log"
	"github.com/D-PlatformOperatingSystem/dpos/mempool"
	"github.com/D-PlatformOperatingSystem/dpos/queue"
	"github.com/D-PlatformOperatingSystem/dpos/store"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	"github.com/D-PlatformOperatingSystem/dpos/wallet"
	"github.com/D-PlatformOperatingSystem/dpos/wallet/bipwallet"
	wcom "github.com/D-PlatformOperatingSystem/dpos/wallet/common"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/types"

	privacy "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/privacy/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//      
	testPrivateKeys = []string{
		"0x8dea7332c7bb3e3b0ce542db41161fd021e3cfda9d7dabacf24f98f2dfd69558",
		"0x920976ffe83b5a98f603b999681a0bc790d97e22ffc4e578a707c2234d55cc8a",
		"0xb59f2b02781678356c231ad565f73699753a28fd3226f1082b513ebf6756c15c",
	}
	//      
	testAddrs = []string{
		"1EDDghAtgBsamrNEtNmYdQzC1QEhLkr87t",
		"13cS5G1BDN2YfGudsxRxr7X25yu6ZdgxMU",
		"1JSRSwp16NvXiTjYBYK9iUQ9wqp3sCxz2p",
	}
	//         
	testPubkeyPairs = []string{
		"92fe6cfec2e19cd15f203f83b5d440ddb63d0cb71559f96dc81208d819fea85886b08f6e874fca15108d244b40f9086d8c03260d4b954a40dfb3cbe41ebc7389",
		"6326126c968a93a546d8f67d623ad9729da0e3e4b47c328a273dfea6930ffdc87bcc365822b80b90c72d30e955e7870a7a9725e9a946b9e89aec6db9455557eb",
		"44bf54abcbae297baf3dec4dd998b313eafb01166760f0c3a4b36509b33d3b50239de0a5f2f47c2fc98a98a382dcd95a2c5bf1f4910467418a3c2595b853338e",
	}
)

func setLogLevel(level string) {
	log.SetLogLevel(level)
}

func init() {
	queue.DisableLog()
	//setLogLevel("err")
	setLogLevel("crit")
}

type testDataMock struct {
	policy wcom.WalletBizPolicy

	wallet  *wallet.Wallet
	modules []queue.Module

	accdb            *account.DB
	mockMempool      bool
	mockBlockChain   bool
	blockChainHeight int64
	password         string
}

func (mock *testDataMock) init() {
	mock.initMember()
	mock.initAccounts()
}

func (mock *testDataMock) initMember() {
	cfg := testnode.GetDefaultConfig()
	mcfg := cfg.GetModuleConfig()

	var q = queue.New("channel")
	q.SetConfig(cfg)

	util.ResetDatadir(mcfg, "$TEMP/")
	wallet := wallet.New(cfg)
	wallet.SetQueueClient(q.Client())
	mock.modules = append(mock.modules, wallet)
	mock.wallet = wallet

	store := store.New(cfg)
	store.SetQueueClient(q.Client())
	mock.modules = append(mock.modules, store)

	if mock.mockBlockChain {
		mock.mockBlockChainProc(q)
	} else {
		chain := blockchain.New(cfg)
		chain.SetQueueClient(q.Client())
		mock.modules = append(mock.modules, chain)
	}

	if mock.mockMempool {
		mock.mockMempoolProc(q)
	} else {
		mempool := mempool.New(cfg)
		mempool.SetQueueClient(q.Client())
		mock.modules = append(mock.modules, mempool)
	}

	mock.accdb = account.NewCoinsAccount(cfg)
	mock.policy = privacy.New()
	sub := cfg.GetSubConfig()
	mock.policy.Init(wallet, sub.Wallet["privacy"])
	mock.password = "ab123456"
}

func (mock *testDataMock) importPrivateKey(PrivKey *types.ReqWalletImportPrivkey) {
	wallet := mock.wallet
	ok, err := wallet.CheckWalletStatus()
	if !ok || err != nil {
		return
	}

	if PrivKey == nil || len(PrivKey.GetLabel()) == 0 || len(PrivKey.GetPrivkey()) == 0 {
		return
	}

	//  label       
	Account, err := wallet.GetAccountByLabel(PrivKey.GetLabel())
	if Account != nil || err != nil {
		return
	}

	var cointype uint32
	signType := wallet.GetSignType()
	if signType == 1 {
		cointype = bipwallet.TypeDpos
	} else if signType == 2 {
		cointype = bipwallet.TypeYcc
	} else {
		cointype = bipwallet.TypeDpos
	}

	privkeybyte, err := common.FromHex(PrivKey.Privkey)
	if err != nil || len(privkeybyte) == 0 {
		return
	}

	pub, err := bipwallet.PrivkeyToPub(cointype, uint32(signType), privkeybyte)
	if err != nil {
		return
	}
	addr, err := bipwallet.PubToAddress(pub)
	if err != nil {
		return
	}

	//     
	Encryptered := wcom.CBCEncrypterPrivkey([]byte(wallet.Password), privkeybyte)
	Encrypteredstr := common.ToHex(Encryptered)
	//  PrivKey   addr         
	Account, err = wallet.GetAccountByAddr(addr)
	if Account != nil || err != nil {
		if Account.Privkey == Encrypteredstr {
			return
		}
	}

	var walletaccount types.WalletAccount
	var WalletAccStore types.WalletAccountStore
	WalletAccStore.Privkey = Encrypteredstr //        
	WalletAccStore.Label = PrivKey.GetLabel()
	WalletAccStore.Addr = addr
	//  Addr:label+privkey+addr    
	err = wallet.SetWalletAccount(false, addr, &WalletAccStore)
	if err != nil {
		return
	}

	//            account  
	addrs := make([]string, 1)
	addrs[0] = addr
	accounts, err := mock.accdb.LoadAccounts(wallet.GetAPI(), addrs)
	if err != nil {
		return
	}
	//         
	if len(accounts[0].Addr) == 0 {
		accounts[0].Addr = addr
	}
	walletaccount.Acc = accounts[0]
	walletaccount.Label = PrivKey.Label
}

func (mock *testDataMock) initAccounts() {
	wallet := mock.wallet
	replySeed, _ := wallet.GenSeed(1)
	wallet.SaveSeed(mock.password, replySeed.Seed)
	wallet.ProcWalletUnLock(&types.WalletUnLock{
		Passwd: mock.password,
	})

	for index, key := range testPrivateKeys {
		privKey := &types.ReqWalletImportPrivkey{
			Label:   fmt.Sprintf("Label%d", index+1),
			Privkey: key,
		}
		mock.importPrivateKey(privKey)
	}
	cfg := mock.wallet.GetAPI().GetConfig()
	accCoin := account.NewCoinsAccount(cfg)
	accCoin.SetDB(wallet.GetDBStore())
	accounts, _ := mock.accdb.LoadAccounts(wallet.GetAPI(), testAddrs)
	for _, account := range accounts {
		account.Balance = 1000 * types.Coin
		accCoin.SaveAccount(account)
	}
}

func (mock *testDataMock) enablePrivacy() {
	mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "EnablePrivacy", &ty.ReqEnablePrivacy{Addrs: testAddrs})
}

func (mock *testDataMock) setBlockChainHeight(height int64) {
	mock.blockChainHeight = height
}

func (mock *testDataMock) mockBlockChainProc(q queue.Queue) {
	// blockchain
	go func() {
		topic := "blockchain"
		client := q.Client()
		client.Sub(topic)
		for msg := range client.Recv() {
			switch msg.Ty {
			case types.EventGetBlockHeight:
				msg.Reply(client.NewMessage(topic, types.EventReplyBlockHeight, &types.ReplyBlockHeight{Height: mock.blockChainHeight}))
			default:
				msg.ReplyErr("Do not support", types.ErrNotSupport)
			}
		}
	}()
}

func (mock *testDataMock) mockMempoolProc(q queue.Queue) {
	// mempool
	go func() {
		topic := "mempool"
		client := q.Client()
		client.Sub(topic)
		for msg := range client.Recv() {
			switch msg.Ty {
			case types.EventTx:
				msg.Reply(client.NewMessage(topic, types.EventReply, &types.Reply{IsOk: true, Msg: []byte("word")}))
			default:
				msg.ReplyErr("Do not support", types.ErrNotSupport)
			}
		}
	}()
}

func Test_EnablePrivacy(t *testing.T) {
	mock := &testDataMock{}
	mock.init()

	testCases := []struct {
		req       *ty.ReqEnablePrivacy
		needReply *ty.RepEnablePrivacy
		needError error
	}{
		{
			needError: types.ErrInvalidParam,
		},
		{
			req: &ty.ReqEnablePrivacy{Addrs: []string{testAddrs[0]}},
			needReply: &ty.RepEnablePrivacy{
				Results: []*ty.PriAddrResult{
					{Addr: testAddrs[0], Msg: "ErrAddrNotExist"}},
			},
		},
	}
	for index, testCase := range testCases {
		getReply, getErr := mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "EnablePrivacy", testCase.req)
		require.Equalf(t, getErr, testCase.needError, "EnablePrivacy test case index %d", index)
		if testCase.needReply == nil {
			assert.Nil(t, getReply)
		} else {
			require.Equal(t, getReply, testCase.needReply)
		}
	}
}

func Test_ShowPrivacyKey(t *testing.T) {
	mock := &testDataMock{}
	mock.init()
	//    0         
	mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "EnablePrivacy", &ty.ReqEnablePrivacy{Addrs: []string{testAddrs[0]}})

	testCases := []struct {
		req       *types.ReqString
		needReply *ty.ReplyPrivacyPkPair
		needError error
	}{
		{
			req:       &types.ReqString{Data: testAddrs[1]},
			needError: types.ErrAddrNotExist,
		},
		{
			req: &types.ReqString{Data: testAddrs[0]},
			/*needReply: &ty.ReplyPrivacyPkPair{
				ShowSuccessful: true,
				Pubkeypair:     "92fe6cfec2e19cd15f203f83b5d440ddb63d0cb71559f96dc81208d819fea85886b08f6e874fca15108d244b40f9086d8c03260d4b954a40dfb3cbe41ebc7389",
			},*/
			needError: types.ErrAddrNotExist,
		},
	}

	for index, testCase := range testCases {
		getReply, getErr := mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "ShowPrivacyKey", testCase.req)
		require.Equalf(t, getErr, testCase.needError, "ShowPrivacyKey test case index %d", index)
		if testCase.needReply == nil {
			assert.Nil(t, getReply)
			continue
		}
		require.Equal(t, getReply, testCase.needReply)
	}
}

func Test_CreateTransaction(t *testing.T) {
	mock := &testDataMock{
		mockMempool:    true,
		mockBlockChain: true,
	}
	mock.init()
	mock.enablePrivacy()
	//       
	privacyMock := privacy.PrivacyMock{}
	privacyMock.Init(mock.wallet, mock.password)
	//       UTXO
	privacyMock.CreateUTXOs(testAddrs[0], testPubkeyPairs[0], 17*types.Coin, 10000, 5)
	mock.setBlockChainHeight(10020)

	testCases := []struct {
		req       *ty.ReqCreatePrivacyTx
		needReply *types.Transaction
		needError error
	}{
		{
			needError: types.ErrInvalidParam,
		},
		{ //      
			req: &ty.ReqCreatePrivacyTx{
				AssetExec:  "coins",
				Tokenname:  types.DOM,
				ActionType: ty.ActionPublic2Privacy,
				Amount:     100 * types.Coin,
				From:       testAddrs[0],
				Pubkeypair: testPubkeyPairs[0],
			},
			//needError:types.ErrAddrNotExist,
		},
		{ //      
			req: &ty.ReqCreatePrivacyTx{
				AssetExec:  "coins",
				Tokenname:  types.DOM,
				ActionType: ty.ActionPrivacy2Privacy,
				Amount:     10 * types.Coin,
				From:       testAddrs[0],
				Pubkeypair: testPubkeyPairs[1],
			},
			needError: types.ErrAddrNotExist,
		},
		{ //      
			req: &ty.ReqCreatePrivacyTx{
				AssetExec:  "coins",
				Tokenname:  types.DOM,
				ActionType: ty.ActionPrivacy2Public,
				Amount:     10 * types.Coin,
				From:       testAddrs[0],
				Pubkeypair: testPubkeyPairs[0],
			},
			needError: types.ErrAddrNotExist,
		},
	}
	for index, testCase := range testCases {
		_, getErr := mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "CreateTransaction", testCase.req)
		require.Equalf(t, getErr, testCase.needError, "CreateTransaction test case index %d", index)
	}
}

func Test_PrivacyAccountInfo(t *testing.T) {
	mock := &testDataMock{}
	mock.init()

	testCases := []struct {
		req       *ty.ReqPrivacyAccount
		needReply *ty.ReplyPrivacyAccount
		needError error
	}{
		{
			needError: types.ErrInvalidParam,
		},
		{
			req: &ty.ReqPrivacyAccount{
				Addr:        testAddrs[0],
				Token:       types.DOM,
				Displaymode: 0,
			},
		},
	}
	for index, testCase := range testCases {
		_, getErr := mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "ShowPrivacyAccountInfo", testCase.req)
		require.Equalf(t, getErr, testCase.needError, "ShowPrivacyAccoPrivacyAccountInfountInfo test case index %d", index)
	}
}

func Test_ShowPrivacyAccountSpend(t *testing.T) {
	mock := &testDataMock{}
	mock.init()

	testCases := []struct {
		req       *ty.ReqPrivBal4AddrToken
		needReply *ty.UTXOHaveTxHashs
		needError error
	}{
		{
			needError: types.ErrInvalidParam,
		},
		{
			req: &ty.ReqPrivBal4AddrToken{
				Addr:  testAddrs[0],
				Token: types.DOM,
			},
			//needError: types.ErrNotFound,
		},
	}
	for index, testCase := range testCases {
		_, getErr := mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "ShowPrivacyAccountSpend", testCase.req)
		require.Equalf(t, getErr, testCase.needError, "ShowPrivacyAccountSpend test case index %d", index)
	}
}

func Test_PrivacyTransactionList(t *testing.T) {
	mock := &testDataMock{}
	mock.init()

	testCases := []struct {
		req       *ty.ReqPrivacyTransactionList
		needReply *types.WalletTxDetails
		needError error
	}{
		{
			needError: types.ErrInvalidParam,
		},
		{
			req: &ty.ReqPrivacyTransactionList{
				Tokenname:    types.DOM,
				SendRecvFlag: 1,
				Direction:    0,
				Count:        10,
				Address:      testAddrs[0],
			},
			//needError: types.ErrTxNotExist,
		},
	}
	for index, testCase := range testCases {
		_, getErr := mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "PrivacyTransactionList", testCase.req)
		require.Equalf(t, getErr, testCase.needError, "PrivacyTransactionList test case index %d", index)
	}
}

func Test_RescanUTXOs(t *testing.T) {
	mock := &testDataMock{}
	mock.init()

	testCases := []struct {
		enable    bool
		req       *ty.ReqRescanUtxos
		needReply *ty.RepRescanUtxos
		needError error
	}{
		{
			needError: types.ErrInvalidParam,
		},
		{
			req: &ty.ReqRescanUtxos{
				Addrs: testAddrs,
				Flag:  0,
			},
			needError: types.ErrAccountNotExist,
		},
		{
			enable: true,
			req: &ty.ReqRescanUtxos{
				Addrs: testAddrs,
				Flag:  0,
			},
			needError: types.ErrAccountNotExist,
		},
	}
	for index, testCase := range testCases {
		if testCase.enable {
			mock.enablePrivacy()
		}
		_, getErr := mock.wallet.GetAPI().ExecWalletFunc(ty.PrivacyX, "RescanUtxos", testCase.req)
		require.Equalf(t, getErr, testCase.needError, "RescanUtxos test case index %d", index)
	}
}
