// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc_test

import (
	"testing"

	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin"
	tokenty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/types"
	"github.com/stretchr/testify/assert"
)

func TestRPCTokenPreCreate(t *testing.T) {
	//   RPCmocker
	mockDOM := testnode.New("", nil)
	cfg := mockDOM.GetAPI().GetConfig()
	defer mockDOM.Close()
	mockDOM.Listen()
	//precreate
	err := mockDOM.SendHot()
	assert.Nil(t, err)
	block := mockDOM.GetLastBlock()
	acc := mockDOM.GetAccount(block.StateHash, mockDOM.GetGenesisAddress())
	assert.Equal(t, acc.Balance, int64(9998999999900000))
	acc = mockDOM.GetAccount(block.StateHash, mockDOM.GetHotAddress())
	assert.Equal(t, acc.Balance, 10000*types.Coin)

	tx := util.CreateManageTx(cfg, mockDOM.GetHotKey(), "token-blacklist", "add", "DOM")
	reply, err := mockDOM.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	detail, err := mockDOM.WaitTx(reply.GetMsg())
	assert.Nil(t, err)
	assert.Equal(t, detail.Receipt.Ty, int32(types.ExecOk))
	//    percreate
	param := tokenty.TokenPreCreate{
		Name:   "Test",
		Symbol: "TEST",
		Total:  10000 * types.Coin,
		Owner:  mockDOM.GetHotAddress(),
	}
	var txhex string
	err = mockDOM.GetJSONC().Call("token.CreateRawTokenPreCreateTx", param, &txhex)
	assert.Nil(t, err)
	hash, err := mockDOM.SendAndSign(mockDOM.GetHotKey(), txhex)
	assert.Nil(t, err)
	assert.NotNil(t, hash)
	detail, err = mockDOM.WaitTx(hash)
	assert.Nil(t, err)
	assert.Equal(t, detail.Receipt.Ty, int32(types.ExecOk))
}
