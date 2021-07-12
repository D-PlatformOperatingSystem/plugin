// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc_test

import (
	"strings"
	"testing"

	commonlog "github.com/D-PlatformOperatingSystem/dpos/common/log"
	"github.com/D-PlatformOperatingSystem/dpos/rpc/jsonclient"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
	"github.com/stretchr/testify/assert"

	//   system plugin  
	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin"
)

func init() {
	commonlog.SetLogLevel("error")
}

func TestJRPCChannel(t *testing.T) {
	//   RPCmocker
	mocker := testnode.New("--notset--", nil)
	defer func() {
		mocker.Close()
	}()
	mocker.Listen()

	jrpcClient := mocker.GetJSONC()

	testCases := []struct {
		fn func(*testing.T, *jsonclient.JSONClient) error
	}{
		{fn: testCreateMultiSigAccCreateCmd},
		{fn: testCreateMultiSigAccOwnerAddCmd},
		{fn: testCreateMultiSigAccOwnerDelCmd},
		{fn: testCreateMultiSigAccOwnerModifyCmd},
		{fn: testCreateMultiSigAccOwnerReplaceCmd},
		{fn: testCreateMultiSigAccWeightModifyCmd},
		{fn: testCreateMultiSigAccDailyLimitModifyCmd},
		{fn: testCreateMultiSigConfirmTxCmd},
		{fn: testCreateMultiSigAccTransferInCmd},
		{fn: testCreateMultiSigAccTransferOutCmd},

		{fn: testGetMultiSigAccCountCmd},
		{fn: testGetMultiSigAccountsCmd},
		{fn: testGetMultiSigAccountInfoCmd},
		{fn: testGetMultiSigAccTxCountCmd},
		{fn: testGetMultiSigTxidsCmd},
		{fn: testGetMultiSigTxInfoCmd},
		{fn: testGetGetMultiSigTxConfirmedWeightCmd},
		{fn: testGetGetMultiSigAccUnSpentTodayCmd},
		{fn: testGetMultiSigAccAssetsCmd},
		{fn: testGetMultiSigAccAllAddressCmd},
	}
	for index, testCase := range testCases {
		err := testCase.fn(t, jrpcClient)
		if err == nil {
			continue
		}
		assert.NotEqualf(t, err, types.ErrActionNotSupport, "test index %d", index)
		if strings.Contains(err.Error(), "rpc: can't find") {
			assert.FailNowf(t, err.Error(), "test index %d", index)
		}
	}
}

//    
func testCreateMultiSigAccCreateCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigAccCreate{}
	return jrpc.Call("multisig.MultiSigAccCreateTx", params, nil)
}
func testCreateMultiSigAccOwnerAddCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigOwnerOperate{}
	return jrpc.Call("multisig.MultiSigOwnerOperateTx", params, nil)
}
func testCreateMultiSigAccOwnerDelCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigOwnerOperate{}
	return jrpc.Call("multisig.MultiSigOwnerOperateTx", params, nil)
}
func testCreateMultiSigAccOwnerModifyCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigOwnerOperate{}
	return jrpc.Call("multisig.MultiSigOwnerOperateTx", params, nil)
}
func testCreateMultiSigAccOwnerReplaceCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigOwnerOperate{}
	return jrpc.Call("multisig.MultiSigOwnerOperateTx", params, nil)
}
func testCreateMultiSigAccWeightModifyCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigAccOperate{}
	return jrpc.Call("multisig.MultiSigAccOperateTx", params, nil)
}

func testCreateMultiSigAccDailyLimitModifyCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigAccOperate{}
	return jrpc.Call("multisig.MultiSigAccOperateTx", params, nil)
}

func testCreateMultiSigConfirmTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigConfirmTx{}
	return jrpc.Call("multisig.MultiSigConfirmTx", params, nil)
}
func testCreateMultiSigAccTransferInCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigExecTransferTo{}
	return jrpc.Call("multisig.MultiSigAccTransferInTx", params, nil)
}
func testCreateMultiSigAccTransferOutCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &mty.MultiSigExecTransferFrom{}
	return jrpc.Call("multisig.MultiSigAccTransferOutTx", params, nil)
}

//get         
func testGetMultiSigAccCountCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &rpctypes.Query4Jrpc{
		Execer:   mty.MultiSigX,
		FuncName: "MultiSigAccCount",
		Payload:  types.MustPBToJSON(&types.ReqNil{}),
	}
	var res types.Int64
	return jrpc.Call("DplatformOS.Query", params, &res)
}

func testGetMultiSigAccountsCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &rpctypes.Query4Jrpc{
		Execer:   mty.MultiSigX,
		FuncName: "MultiSigAccounts",
		Payload:  types.MustPBToJSON(&mty.ReqMultiSigAccs{}),
	}
	var res mty.ReplyMultiSigAccs
	return jrpc.Call("DplatformOS.Query", params, &res)
}

func testGetMultiSigAccountInfoCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &rpctypes.Query4Jrpc{
		Execer:   mty.MultiSigX,
		FuncName: "MultiSigAccountInfo",
		Payload:  types.MustPBToJSON(&mty.ReqMultiSigAccInfo{}),
	}
	var res mty.MultiSig
	return jrpc.Call("DplatformOS.Query", params, &res)
}

func testGetMultiSigAccTxCountCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &rpctypes.Query4Jrpc{
		Execer:   mty.MultiSigX,
		FuncName: "MultiSigAccTxCount",
		Payload:  types.MustPBToJSON(&mty.ReqMultiSigAccInfo{}),
	}
	var res mty.Uint64
	return jrpc.Call("DplatformOS.Query", params, &res)
}

func testGetMultiSigTxidsCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &rpctypes.Query4Jrpc{
		Execer:   mty.MultiSigX,
		FuncName: "MultiSigTxids",
		Payload:  types.MustPBToJSON(&mty.ReqMultiSigTxids{}),
	}
	var res mty.ReplyMultiSigTxids
	return jrpc.Call("DplatformOS.Query", params, &res)
}

func testGetMultiSigTxInfoCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &mty.ReqMultiSigTxInfo{}
	params.Execer = mty.MultiSigX
	params.FuncName = "MultiSigTxInfo"
	params.Payload = types.MustPBToJSON(req)
	rep = &mty.MultiSigTx{}
	return jrpc.Call("DplatformOS.Query", &params, rep)
}

func testGetGetMultiSigTxConfirmedWeightCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &mty.ReqMultiSigTxInfo{}
	params.Execer = mty.MultiSigX
	params.FuncName = "MultiSigTxConfirmedWeight"
	params.Payload = types.MustPBToJSON(req)
	rep = &mty.Uint64{}
	return jrpc.Call("DplatformOS.Query", &params, rep)
}

func testGetGetMultiSigAccUnSpentTodayCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &mty.ReqAccAssets{}
	req.IsAll = true
	params.Execer = mty.MultiSigX
	params.FuncName = "MultiSigAccUnSpentToday"
	params.Payload = types.MustPBToJSON(req)
	rep = &mty.ReplyUnSpentAssets{}
	return jrpc.Call("DplatformOS.Query", &params, rep)
}

func testGetMultiSigAccAssetsCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc

	req := &mty.ReqAccAssets{}
	req.IsAll = true
	params.Execer = mty.MultiSigX
	params.FuncName = "MultiSigAccAssets"
	params.Payload = types.MustPBToJSON(req)
	rep = &mty.ReplyAccAssets{}
	return jrpc.Call("DplatformOS.Query", &params, rep)
}

func testGetMultiSigAccAllAddressCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc

	req := mty.ReqMultiSigAccInfo{
		MultiSigAccAddr: "14jv8WB7CwNQSnh4qo9WDBgRPRBjM5LQo6",
	}

	params.Execer = mty.MultiSigX
	params.FuncName = "MultiSigAccAllAddress"
	params.Payload = types.MustPBToJSON(&req)
	rep = &mty.AccAddress{}
	return jrpc.Call("DplatformOS.Query", &params, rep)
}
