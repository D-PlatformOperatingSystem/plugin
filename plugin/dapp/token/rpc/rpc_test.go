// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/client/mocks"
	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	tokenty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	context "golang.org/x/net/context"
)

func newTestChannelClient() *channelClient {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	api := &mocks.QueueProtocolAPI{}
	api.On("GetConfig", mock.Anything).Return(cfg)
	return &channelClient{
		ChannelClient: rpctypes.ChannelClient{QueueProtocolAPI: api},
	}
}

func newTestJrpcClient() *Jrpc {
	return &Jrpc{cli: newTestChannelClient()}
}

func testChannelClientGetTokenBalanceToken(t *testing.T) {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	api := new(mocks.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(cfg)

	client := &channelClient{
		ChannelClient: rpctypes.ChannelClient{QueueProtocolAPI: api},
	}

	head := &types.Header{StateHash: []byte("sdfadasds")}
	api.On("GetLastHeader").Return(head, nil)

	var acc = &types.Account{Addr: "1Jn2qu84Z1SUUosWjySggBS9pKWdAP3tZt", Balance: 100}
	accv := types.Encode(acc)
	storevalue := &types.StoreReplyValue{}
	storevalue.Values = append(storevalue.Values, accv)
	api.On("StoreGet", mock.Anything).Return(storevalue, nil)

	var addrs = make([]string, 1)
	addrs = append(addrs, "1Jn2qu84Z1SUUosWjySggBS9pKWdAP3tZt")
	var in = &tokenty.ReqTokenBalance{
		Execer:      cfg.ExecName(tokenty.TokenX),
		Addresses:   addrs,
		TokenSymbol: "xxx",
	}
	data, err := client.GetTokenBalance(context.Background(), in)
	assert.Nil(t, err)
	accounts := data.Acc
	assert.Equal(t, acc.Addr, accounts[0].Addr)

}

func testChannelClientGetTokenBalanceOther(t *testing.T) {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	api := new(mocks.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(cfg)

	client := &channelClient{
		ChannelClient: rpctypes.ChannelClient{QueueProtocolAPI: api},
	}

	head := &types.Header{StateHash: []byte("sdfadasds")}
	api.On("GetLastHeader").Return(head, nil)

	var acc = &types.Account{Addr: "1Jn2qu84Z1SUUosWjySggBS9pKWdAP3tZt", Balance: 100}
	accv := types.Encode(acc)
	storevalue := &types.StoreReplyValue{}
	storevalue.Values = append(storevalue.Values, accv)
	api.On("StoreGet", mock.Anything).Return(storevalue, nil)

	var addrs = make([]string, 1)
	addrs = append(addrs, "1Jn2qu84Z1SUUosWjySggBS9pKWdAP3tZt")
	var in = &tokenty.ReqTokenBalance{
		Execer:      cfg.ExecName("trade"),
		Addresses:   addrs,
		TokenSymbol: "xxx",
	}
	data, err := client.GetTokenBalance(context.Background(), in)
	assert.Nil(t, err)
	accounts := data.Acc
	assert.Equal(t, acc.Addr, accounts[0].Addr)

}

func TestChannelClientGetTokenBalance(t *testing.T) {
	testChannelClientGetTokenBalanceToken(t)
	testChannelClientGetTokenBalanceOther(t)

}

func TestChannelClientCreateRawTokenPreCreateTx(t *testing.T) {
	client := newTestJrpcClient()
	var data interface{}
	err := client.CreateRawTokenPreCreateTx(nil, &data)
	assert.NotNil(t, err)
	assert.Nil(t, data)

	token := &tokenty.TokenPreCreate{
		Owner:  "asdf134",
		Symbol: "CNY",
	}
	err = client.CreateRawTokenPreCreateTx(token, &data)
	assert.NotNil(t, data)
	assert.Nil(t, err)
}

func TestChannelClientCreateRawTokenRevokeTx(t *testing.T) {
	client := newTestJrpcClient()
	var data interface{}
	err := client.CreateRawTokenRevokeTx(nil, &data)
	assert.NotNil(t, err)
	assert.Nil(t, data)

	token := &tokenty.TokenRevokeCreate{
		Owner:  "asdf134",
		Symbol: "CNY",
	}
	err = client.CreateRawTokenRevokeTx(token, &data)
	assert.NotNil(t, data)
	assert.Nil(t, err)
}

func TestChannelClientCreateRawTokenFinishTx(t *testing.T) {
	client := newTestJrpcClient()
	var data interface{}
	err := client.CreateRawTokenFinishTx(nil, &data)
	assert.NotNil(t, err)
	assert.Nil(t, data)

	token := &tokenty.TokenFinishCreate{
		Owner:  "asdf134",
		Symbol: "CNY",
	}
	err = client.CreateRawTokenFinishTx(token, &data)
	assert.NotNil(t, data)
	assert.Nil(t, err)
}
