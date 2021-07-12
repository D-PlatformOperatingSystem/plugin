// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet_test

import (
	"fmt"
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/rpc/jsonclient"
	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
	"github.com/stretchr/testify/assert"
)

var (
	Symbol     = "DOM"
	Asset      = "coins"
	PrivKeyA   = "0x06c0fa653c719275d1baa365c7bc0b9306447287499a715b541b930482eaa504" // 1C5xK2ytuoFqxmVGMcyz9XFKFWcDA8T3rK
	PrivKeyB   = "0x4c8663cded61093af20339ae038b3c6bfa58a33e65874a655022f82eaf3f2fa0" // 1LDGrokrZjo1HtSmSnw8ef3oy5Vm1nctbj
	PrivKeyC   = "0x9abcf378b397682109c174b37a45bfc8a459c9514dd2ef719e22a9815373047d" // 1DkrXbz2bK6XMpY4v9z2YUnhwWTXT6V5jd
	PrivKeyD   = "0xbf8f865a03fec64f30d2243847807e88d2dbc8104e77925e4fc11c4d4380f3da" // 166po3ghRbRu53hu8jBBQzddp7kUJ9Ynyf
	PrivKeyE   = "0x5b8ca316cf073aa94f1056a9e3f6e0b9a9ec11ae45862d58c7a09640b4d55302" // 1KHwX7ZadNeQDjBGpnweb4k2dqj2CWtAYo
	PrivKeyGen = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
	AddrA      = "1C5xK2ytuoFqxmVGMcyz9XFKFWcDA8T3rK"
	AddrB      = "1LDGrokrZjo1HtSmSnw8ef3oy5Vm1nctbj"
	AddrC      = "1DkrXbz2bK6XMpY4v9z2YUnhwWTXT6V5jd"
	AddrD      = "166po3ghRbRu53hu8jBBQzddp7kUJ9Ynyf"
	AddrE      = "1KHwX7ZadNeQDjBGpnweb4k2dqj2CWtAYo"

	GenAddr = "16ERTbYtKKQ64wMthAY9J4La4nAiidG45A"
)

//TestPrivkeyHex ：
var TestPrivkeyHex = []string{
	"0x06c0fa653c719275d1baa365c7bc0b9306447287499a715b541b930482eaa504",
	"0x4c8663cded61093af20339ae038b3c6bfa58a33e65874a655022f82eaf3f2fa0",
	"0x9abcf378b397682109c174b37a45bfc8a459c9514dd2ef719e22a9815373047d",
	"0xbf8f865a03fec64f30d2243847807e88d2dbc8104e77925e4fc11c4d4380f3da",
	"0x5b8ca316cf073aa94f1056a9e3f6e0b9a9ec11ae45862d58c7a09640b4d55302",
}

func getRPCClient(t *testing.T, mocker *testnode.DplatformOSMock) *jsonclient.JSONClient {
	jrpcClient := mocker.GetJSONC()
	assert.NotNil(t, jrpcClient)
	return jrpcClient
}

func getTx(t *testing.T, hex string) *types.Transaction {
	data, err := common.FromHex(hex)
	assert.Nil(t, err)
	var tx types.Transaction
	err = types.Decode(data, &tx)
	assert.Nil(t, err)
	return &tx
}

func TestMultiSigAccount(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)

	//
	for i, priv := range TestPrivkeyHex {
		privkey := &types.ReqWalletImportPrivkey{Privkey: priv, Label: fmt.Sprintf("heyubin%d", i)}
		_, err := mocker.GetAPI().ExecWalletFunc("wallet", "WalletImportPrivkey", privkey)
		if err != nil {
			panic(err)
		}
		//t.Log("import", "index", i, "addr", acc.Acc.Addr)
	}
	//        ,owner:AddrA,AddrB,GenAddr,weight:20,10,30;coins:DOM 1000000000 RequestWeight:15
	multiSigAccAddr := testAccCreateTx(t, mocker, jrpcClient)

	//owner add AddrE
	testAddOwner(t, mocker, jrpcClient, multiSigAccAddr)
	//owner del AddrE
	testDelOwner(t, mocker, jrpcClient, multiSigAccAddr)
	//owner AddrA modify weight to 30
	testModifyOwnerWeight(t, mocker, jrpcClient, multiSigAccAddr)
	//owner AddrA replace by  AddrE
	testReplaceOwner(t, mocker, jrpcClient, multiSigAccAddr)
}

//
func testAccCreateTx(t *testing.T, mocker *testnode.DplatformOSMock, jrpcClient *jsonclient.JSONClient) string {
	gen := mocker.GetGenesisKey()
	var params rpctypes.Query4Jrpc
	//1. MultiSigAccCreateTx
	var owners []*mty.Owner
	owmer1 := &mty.Owner{OwnerAddr: AddrA, Weight: 20}
	owmer2 := &mty.Owner{OwnerAddr: AddrB, Weight: 10}
	owmer3 := &mty.Owner{OwnerAddr: GenAddr, Weight: 30}
	owners = append(owners, owmer1)
	owners = append(owners, owmer2)
	owners = append(owners, owmer3)

	symboldailylimit := &mty.SymbolDailyLimit{
		Symbol:     Symbol,
		Execer:     Asset,
		DailyLimit: 1000000000,
	}

	req := &mty.MultiSigAccCreate{
		Owners:         owners,
		RequiredWeight: 15,
		DailyLimit:     symboldailylimit,
	}

	var res string
	err := jrpcClient.Call("multisig.MultiSigAccCreateTx", req, &res)
	assert.Nil(t, err)
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	_, err = mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)
	//  account
	params.Execer = mty.MultiSigX
	params.FuncName = "MultiSigAccCount"
	params.Payload = types.MustPBToJSON(&types.ReqNil{})
	rep := &types.Int64{}
	err = jrpcClient.Call("DplatformOS.Query", &params, rep)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rep.Data)

	//  account addr
	//t.Log("MultiSigAccounts ")
	req1 := mty.ReqMultiSigAccs{
		Start: 0,
		End:   0,
	}
	params.Execer = mty.MultiSigX
	params.FuncName = "MultiSigAccounts"
	params.Payload = types.MustPBToJSON(&req1)
	rep1 := &mty.ReplyMultiSigAccs{}
	err = jrpcClient.Call("DplatformOS.Query", params, rep1)
	assert.Nil(t, err)
	//t.Log(rep1)

	multiSigAccAddr := rep1.Address[0]
	//  owner
	req4 := &types.ReqString{
		Data: GenAddr,
	}
	var res4 mty.OwnerAttrs
	err = jrpcClient.Call("multisig.MultiSigAddresList", req4, &res4)
	assert.Nil(t, err)
	assert.Equal(t, res4.Items[0].OwnerAddr, GenAddr)
	return multiSigAccAddr
}

//owner add AddrE
func testAddOwner(t *testing.T, mocker *testnode.DplatformOSMock, jrpcClient *jsonclient.JSONClient, multiSigAccAddr string) {
	gen := mocker.GetGenesisKey()

	params9 := &mty.MultiSigOwnerOperate{
		MultiSigAccAddr: multiSigAccAddr,
		NewOwner:        AddrE,
		NewWeight:       8,
		OperateFlag:     mty.OwnerAdd,
	}
	var res9 string
	err := jrpcClient.Call("multisig.MultiSigOwnerOperateTx", params9, &res9)
	assert.Nil(t, err)
	tx := getTx(t, res9)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	_, err = mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)

	//  owner
	req4 := &types.ReqString{
		Data: AddrE,
	}
	var res4 mty.OwnerAttrs
	err = jrpcClient.Call("multisig.MultiSigAddresList", req4, &res4)
	assert.Nil(t, err)
	//t.Log(res4)
	assert.Equal(t, res4.Items[0].OwnerAddr, AddrE)
	assert.Equal(t, res4.Items[0].MultiSigAddr, multiSigAccAddr)
	assert.Equal(t, res4.Items[0].Weight, uint64(8))
}

//owner del AddrE
func testDelOwner(t *testing.T, mocker *testnode.DplatformOSMock, jrpcClient *jsonclient.JSONClient, multiSigAccAddr string) {
	gen := mocker.GetGenesisKey()

	param := &mty.MultiSigOwnerOperate{
		MultiSigAccAddr: multiSigAccAddr,
		OldOwner:        AddrE,
		OperateFlag:     mty.OwnerDel,
	}
	var res string
	err := jrpcClient.Call("multisig.MultiSigOwnerOperateTx", param, &res)

	assert.Nil(t, err)
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	_, err = mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)

	//  owner
	req4 := &types.ReqString{
		Data: AddrE,
	}
	var res4 mty.OwnerAttrs
	err = jrpcClient.Call("multisig.MultiSigAddresList", req4, &res4)
	//t.Log(res4)
	assert.Equal(t, err, types.ErrNotFound)
}

//ModifyOwnerWeight
func testModifyOwnerWeight(t *testing.T, mocker *testnode.DplatformOSMock, jrpcClient *jsonclient.JSONClient, multiSigAccAddr string) {
	gen := mocker.GetGenesisKey()

	param := &mty.MultiSigOwnerOperate{
		MultiSigAccAddr: multiSigAccAddr,
		OldOwner:        AddrA,
		NewWeight:       30,
		OperateFlag:     mty.OwnerModify,
	}
	var res string
	err := jrpcClient.Call("multisig.MultiSigOwnerOperateTx", param, &res)

	assert.Nil(t, err)
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	_, err = mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)

	//  owner
	req4 := &types.ReqString{
		Data: AddrA,
	}
	var res4 mty.OwnerAttrs
	err = jrpcClient.Call("multisig.MultiSigAddresList", req4, &res4)
	assert.Nil(t, err)
	assert.Equal(t, res4.Items[0].OwnerAddr, AddrA)
	assert.Equal(t, res4.Items[0].MultiSigAddr, multiSigAccAddr)
	assert.Equal(t, res4.Items[0].Weight, uint64(30))
}

//testReplaceOwner owner AddrA replace by  AddrE
func testReplaceOwner(t *testing.T, mocker *testnode.DplatformOSMock, jrpcClient *jsonclient.JSONClient, multiSigAccAddr string) {
	gen := mocker.GetGenesisKey()

	param := &mty.MultiSigOwnerOperate{
		MultiSigAccAddr: multiSigAccAddr,
		OldOwner:        AddrA,
		NewOwner:        AddrE,
		OperateFlag:     mty.OwnerReplace,
	}
	var res string
	err := jrpcClient.Call("multisig.MultiSigOwnerOperateTx", param, &res)

	assert.Nil(t, err)
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	_, err = mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)

	//  owner AddrE
	req4 := &types.ReqString{
		Data: AddrE,
	}
	var res4 mty.OwnerAttrs
	err = jrpcClient.Call("multisig.MultiSigAddresList", req4, &res4)
	assert.Nil(t, err)
	assert.Equal(t, res4.Items[0].OwnerAddr, AddrE)
	assert.Equal(t, res4.Items[0].MultiSigAddr, multiSigAccAddr)
	assert.Equal(t, res4.Items[0].Weight, uint64(30))

	//  owner AddrA
	req5 := &types.ReqString{
		Data: AddrA,
	}
	var res5 mty.OwnerAttrs
	err = jrpcClient.Call("multisig.MultiSigAddresList", req5, &res5)
	assert.Equal(t, err, types.ErrNotFound)
}
