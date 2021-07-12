/*
 * Copyright D-Platform Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

/*
    ：
1.sendAddPublisher   manage
2.sendPublishEvent
3.queryEventByeventID     ID
4.sendAbortPublishEvent
5.sendPrePublishResult
6.sendAbortPublishResult
7.sendPublishResult
    ：
1.                   ，  ：
  [exec.sub.manage]
  superManager=["16ERTbYtKKQ64wMthAY9J4La4nAiidG45A"]
2.TestPublishNomal
3.TestAbortPublishEvent
4.TestPrePublishResult
5.TestAbortPublishResult
6.TestPublishResult
7.TestQueryEventIDByStatus
8.TestQueryEventIDByAddrAndStatus
9.TestQueryEventIDByTypeAndStatus
*/

package rpc_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/rpc/jsonclient"
	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	oty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/oracle/types"
	"github.com/stretchr/testify/assert"
)

var (
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(types.Now().UnixNano()))
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

func TestPublishNomal(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)
	// publish event
	eventID := sendPublishEvent(t, jrpcClient, mocker)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)
	//pre publish result
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPrePublished)
	//publish result
	sendPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
}
func TestPublishEvent(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)
	// publish event
	// abort event
	// publish event
	eventID := sendPublishEvent(t, jrpcClient, mocker)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)
	eventoldID := eventID
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)
	assert.NotEqual(t, eventID, eventoldID)

	// publish event
	// pre publish result
	// publish event
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPrePublished)
	eventoldID = eventID
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)
	assert.NotEqual(t, eventID, eventoldID)

	// publish event
	// pre publish result
	// publilsh result
	// publish event
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	sendPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
	eventoldID = eventID
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)
	assert.NotEqual(t, eventID, eventoldID)
}

func TestAbortPublishEvent(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)

	// publish event
	// abort event
	// abort event
	eventID := sendPublishEvent(t, jrpcClient, mocker)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, oty.ErrEventAbortNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)

	// publish event
	// pre publish result
	// abort event
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPrePublished)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, oty.ErrEventAbortNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPrePublished)

	// publish event
	// pre publish result
	// abort pre publilsh result
	// abort event
	sendAbortPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultAborted)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, oty.ErrEventAbortNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)

	// publish event
	// pre publish result
	// publilsh result
	// abort event
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	sendPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, oty.ErrEventAbortNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
}

func TestPrePublishResult(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)

	// publish event
	// pre publish result
	// pre publish result
	eventID := sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPrePublished)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, oty.ErrResultPrePublishNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPrePublished)

	// publish event
	// abort event
	// pre publish
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, oty.ErrResultPrePublishNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)

	// publish event
	// pre publish result
	// abort pre publish result
	// pre publish result
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	sendAbortPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultAborted)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPrePublished)

	//publish result
	//pre publish result
	sendPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, oty.ErrResultPrePublishNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
}

func TestAbortPublishResult(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)

	//publish event
	//abort prepublish result
	eventID := sendPublishEvent(t, jrpcClient, mocker)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)
	sendAbortPublishResult(eventID, t, jrpcClient, mocker, oty.ErrPrePublishAbortNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)

	// publish event
	// abort event
	// abort pre publish result
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)
	sendAbortPublishResult(eventID, t, jrpcClient, mocker, oty.ErrPrePublishAbortNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)

	// publish event
	// pre publish result
	// abort pre publish result
	// abort pre publish result
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPrePublished)
	sendAbortPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultAborted)
	sendAbortPublishResult(eventID, t, jrpcClient, mocker, oty.ErrPrePublishAbortNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultAborted)

	// publish event
	// pre publish result
	// publish result
	// abort pre publish result
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	sendPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
	sendAbortPublishResult(eventID, t, jrpcClient, mocker, oty.ErrPrePublishAbortNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
}

func TestPublishResult(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)

	//publish event
	//publish result
	eventID := sendPublishEvent(t, jrpcClient, mocker)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)
	sendPublishResult(eventID, t, jrpcClient, mocker, oty.ErrResultPublishNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventPublished)

	// publish event
	// abort event
	// publish result
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendAbortPublishEvent(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)
	sendPublishResult(eventID, t, jrpcClient, mocker, oty.ErrResultPublishNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.EventAborted)

	// publish event
	// pre publish result
	// abort pre publish result
	// publish result
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	sendAbortPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultAborted)
	sendPublishResult(eventID, t, jrpcClient, mocker, oty.ErrResultPublishNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultAborted)

	// publish event
	// pre publish result
	// publish result
	// publish result
	eventID = sendPublishEvent(t, jrpcClient, mocker)
	sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
	sendPublishResult(eventID, t, jrpcClient, mocker, nil)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
	sendPublishResult(eventID, t, jrpcClient, mocker, oty.ErrResultPublishNotAllowed)
	queryEventByeventID(eventID, t, jrpcClient, oty.ResultPublished)
}

func createAllStatusEvent(t *testing.T, jrpcClient *jsonclient.JSONClient, mocker *testnode.DplatformOSMock) {
	//total loop*5
	loop := int(oty.DefaultCount + 10)
	for i := 0; i < loop; i++ {
		//EventPublished
		eventID := sendPublishEvent(t, jrpcClient, mocker)
		assert.NotEqual(t, "", eventID)

		//EventAborted
		eventID = sendPublishEvent(t, jrpcClient, mocker)
		sendAbortPublishEvent(eventID, t, jrpcClient, mocker, nil)

		//ResultPrePublished
		eventID = sendPublishEvent(t, jrpcClient, mocker)
		sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)

		//ResultAborted
		eventID = sendPublishEvent(t, jrpcClient, mocker)
		sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
		sendAbortPublishResult(eventID, t, jrpcClient, mocker, nil)

		//ResultPublished
		eventID = sendPublishEvent(t, jrpcClient, mocker)
		sendPrePublishResult(eventID, t, jrpcClient, mocker, nil)
		sendPublishResult(eventID, t, jrpcClient, mocker, nil)
	}
}

func TestQueryEventIDByStatus(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)
	createAllStatusEvent(t, jrpcClient, mocker)
	queryEventByStatus(t, jrpcClient)
}

func TestQueryEventIDByAddrAndStatus(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)

	createAllStatusEvent(t, jrpcClient, mocker)

	queryEventByStatusAndAddr(t, jrpcClient)

}

func TestQueryEventIDByTypeAndStatus(t *testing.T) {
	mocker := testnode.New("--free--", nil)
	defer mocker.Close()
	mocker.Listen()
	jrpcClient := getRPCClient(t, mocker)
	sendAddPublisher(t, jrpcClient, mocker)

	createAllStatusEvent(t, jrpcClient, mocker)

	queryEventByStatusAndType(t, jrpcClient)
}

func sendAddPublisher(t *testing.T, jrpcClient *jsonclient.JSONClient, mocker *testnode.DplatformOSMock) {
	//1.   createrawtransaction
	req := &rpctypes.CreateTxIn{
		Execer:     "manage",
		ActionName: "Modify",
		Payload:    []byte("{\"key\":\"oracle-publish-event\",\"op\":\"add\", \"value\":\"12oupcayRT7LvaC4qW4avxsTE7U41cKSio\"}"),
	}
	var res string
	err := jrpcClient.Call("DplatformOS.CreateTransaction", req, &res)
	assert.Nil(t, err)
	gen := mocker.GetHotKey()
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	_, err = mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)
}

func sendPublishEvent(t *testing.T, jrpcClient *jsonclient.JSONClient, mocker *testnode.DplatformOSMock) (eventID string) {
	ti := time.Now().AddDate(0, 0, 1)
	//1.   createrawtransaction
	req := &rpctypes.CreateTxIn{
		Execer:     oty.OracleX,
		ActionName: "EventPublish",
		Payload:    []byte(fmt.Sprintf("{\"type\":\"football\",\"subType\":\"Premier League\",\"time\":%d, \"content\":\"{\\\"team%d\\\":\\\"ChelSea\\\", \\\"team%d\\\":\\\"Manchester\\\",\\\"resultType\\\":\\\"score\\\"}\", \"introduction\":\"guess the sore result of football game between ChelSea and Manchester in 2019-01-21 14:00:00\"}", ti.Unix(), r.Int()%10, r.Int()%10)),
	}
	var res string
	err := jrpcClient.Call("DplatformOS.CreateTransaction", req, &res)
	assert.Nil(t, err)
	gen := mocker.GetHotKey()
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	result, err := mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)

	for _, log := range result.Receipt.Logs {
		if log.Ty >= oty.TyLogEventPublish && log.Ty <= oty.TyLogResultPublish {
			fmt.Println(log.TyName)
			fmt.Println(string(log.Log))
			status := oty.ReceiptOracle{}
			logData, err := common.FromHex(log.RawLog)
			assert.Nil(t, err)
			err = types.Decode(logData, &status)
			assert.Nil(t, err)
			eventID = status.EventID
		}
	}
	return eventID
}

func sendAbortPublishEvent(eventID string, t *testing.T, jrpcClient *jsonclient.JSONClient, mocker *testnode.DplatformOSMock, expectErr error) {
	req := &rpctypes.CreateTxIn{
		Execer:     oty.OracleX,
		ActionName: "EventAbort",
		Payload:    []byte(fmt.Sprintf("{\"eventID\":\"%s\"}", eventID)),
	}
	var res string
	err := jrpcClient.Call("DplatformOS.CreateTransaction", req, &res)
	assert.Nil(t, err)
	gen := mocker.GetHotKey()
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	result, err := mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)
	fmt.Println(string(result.Tx.Payload))
	for _, log := range result.Receipt.Logs {
		if log.Ty >= oty.TyLogEventPublish && log.Ty <= oty.TyLogResultPublish {
			fmt.Println(log.TyName)
			fmt.Println(string(log.Log))
		} else if log.Ty == 1 {
			logData, err := common.FromHex(log.RawLog)
			assert.Nil(t, err)
			assert.Equal(t, expectErr.Error(), string(logData))
		}
	}
}

func sendPrePublishResult(eventID string, t *testing.T, jrpcClient *jsonclient.JSONClient, mocker *testnode.DplatformOSMock, expectErr error) {
	req := &rpctypes.CreateTxIn{
		Execer:     oty.OracleX,
		ActionName: "ResultPrePublish",
		Payload:    []byte(fmt.Sprintf("{\"eventID\":\"%s\", \"source\":\"sina sport\", \"result\":\"%d:%d\"}", eventID, r.Int()%10, r.Int()%10)),
	}
	var res string
	err := jrpcClient.Call("DplatformOS.CreateTransaction", req, &res)
	assert.Nil(t, err)
	gen := mocker.GetHotKey()
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	result, err := mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)
	for _, log := range result.Receipt.Logs {
		if log.Ty >= oty.TyLogEventPublish && log.Ty <= oty.TyLogResultPublish {
			fmt.Println(log.TyName)
			fmt.Println(string(log.Log))
		} else if log.Ty == 1 {
			logData, err := common.FromHex(log.RawLog)
			assert.Nil(t, err)
			assert.Equal(t, expectErr.Error(), string(logData))
		}
	}
}

func sendAbortPublishResult(eventID string, t *testing.T, jrpcClient *jsonclient.JSONClient, mocker *testnode.DplatformOSMock, expectErr error) {
	req := &rpctypes.CreateTxIn{
		Execer:     oty.OracleX,
		ActionName: "ResultAbort",
		Payload:    []byte(fmt.Sprintf("{\"eventID\":\"%s\"}", eventID)),
	}
	var res string
	err := jrpcClient.Call("DplatformOS.CreateTransaction", req, &res)
	assert.Nil(t, err)
	gen := mocker.GetHotKey()
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	result, err := mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)
	for _, log := range result.Receipt.Logs {
		if log.Ty >= oty.TyLogEventPublish && log.Ty <= oty.TyLogResultPublish {
			fmt.Println(log.TyName)
			fmt.Println(string(log.Log))
		} else if log.Ty == 1 {
			logData, err := common.FromHex(log.RawLog)
			assert.Nil(t, err)
			assert.Equal(t, expectErr.Error(), string(logData))
		}
	}
}

func sendPublishResult(eventID string, t *testing.T, jrpcClient *jsonclient.JSONClient, mocker *testnode.DplatformOSMock, expectErr error) {
	req := &rpctypes.CreateTxIn{
		Execer:     oty.OracleX,
		ActionName: "ResultPublish",
		Payload:    []byte(fmt.Sprintf("{\"eventID\":\"%s\", \"source\":\"sina sport\", \"result\":\"%d:%d\"}", eventID, r.Int()%10, r.Int()%10)),
	}
	var res string
	err := jrpcClient.Call("DplatformOS.CreateTransaction", req, &res)
	assert.Nil(t, err)
	gen := mocker.GetHotKey()
	tx := getTx(t, res)
	tx.Sign(types.SECP256K1, gen)
	reply, err := mocker.GetAPI().SendTx(tx)
	assert.Nil(t, err)
	result, err := mocker.WaitTx(reply.GetMsg())
	assert.Nil(t, err)
	for _, log := range result.Receipt.Logs {
		if log.Ty >= oty.TyLogEventPublish && log.Ty <= oty.TyLogResultPublish {
			fmt.Println(log.TyName)
			fmt.Println(string(log.Log))
		} else if log.Ty == 1 {
			logData, err := common.FromHex(log.RawLog)
			assert.Nil(t, err)
			assert.Equal(t, expectErr.Error(), string(logData))
		}
	}
}

func queryEventByeventID(eventID string, t *testing.T, jrpcClient *jsonclient.JSONClient, expectedStatus int32) {
	//   ID
	params := rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryOracleListByIDs,
		Payload:  []byte(fmt.Sprintf("{\"eventID\":[\"%s\"]}", eventID)),
	}
	var resStatus oty.ReplyOracleStatusList
	err := jrpcClient.Call("DplatformOS.Query", params, &resStatus)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resStatus.Status[0].Status.Status)
	fmt.Println(resStatus.Status[0])

}

func queryEventByStatus(t *testing.T, jrpcClient *jsonclient.JSONClient) {
	for i := 1; i <= 5; i++ {
		//
		params := rpctypes.Query4Jrpc{
			Execer:   oty.OracleX,
			FuncName: oty.FuncNameQueryEventIDByStatus,
			Payload:  []byte(fmt.Sprintf("{\"status\":%d,\"addr\":\"\",\"type\":\"\",\"eventID\":\"\"}", i)),
		}
		var res oty.ReplyEventIDs
		err := jrpcClient.Call("DplatformOS.Query", params, &res)
		assert.Nil(t, err)
		assert.EqualValues(t, oty.DefaultCount, len(res.EventID))
		lastEventID := res.EventID[oty.DefaultCount-1]
		//
		params = rpctypes.Query4Jrpc{
			Execer:   oty.OracleX,
			FuncName: oty.FuncNameQueryEventIDByStatus,
			Payload:  []byte(fmt.Sprintf("{\"status\":%d,\"addr\":\"\",\"type\":\"\",\"eventID\":\"%s\"}", i, lastEventID)),
		}
		err = jrpcClient.Call("DplatformOS.Query", params, &res)
		assert.Nil(t, err)
		assert.Equal(t, 10, len(res.EventID))
		lastEventID = res.EventID[9]
		//         ,
		params = rpctypes.Query4Jrpc{
			Execer:   oty.OracleX,
			FuncName: oty.FuncNameQueryEventIDByStatus,
			Payload:  []byte(fmt.Sprintf("{\"status\":%d,\"addr\":\"\",\"type\":\"\",\"eventID\":\"%s\"}", i, lastEventID)),
		}
		err = jrpcClient.Call("DplatformOS.Query", params, &res)
		assert.Equal(t, types.ErrNotFound, err)
	}
}

func queryEventByStatusAndAddr(t *testing.T, jrpcClient *jsonclient.JSONClient) {
	//
	params := rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryEventIDByAddrAndStatus,
		Payload:  []byte("{\"status\":1,\"addr\":\"12oupcayRT7LvaC4qW4avxsTE7U41cKSio\",\"type\":\"\",\"eventID\":\"\"}"),
	}
	var res oty.ReplyEventIDs
	err := jrpcClient.Call("DplatformOS.Query", params, &res)
	assert.Nil(t, err)
	assert.EqualValues(t, oty.DefaultCount, len(res.EventID))
	lastEventID := res.EventID[oty.DefaultCount-1]
	//
	params = rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryEventIDByAddrAndStatus,
		Payload:  []byte(fmt.Sprintf("{\"status\":1,\"addr\":\"12oupcayRT7LvaC4qW4avxsTE7U41cKSio\",\"type\":\"\",\"eventID\":\"%s\"}", lastEventID)),
	}
	err = jrpcClient.Call("DplatformOS.Query", params, &res)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(res.EventID))
	lastEventID = res.EventID[9]

	//
	params = rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryEventIDByAddrAndStatus,
		Payload:  []byte(fmt.Sprintf("{\"status\":1,\"addr\":\"16ABEbYtKKQm5wMuySK9J4La5nAiidGuyt\",\"type\":\"\",\"eventID\":\"%s\"}", lastEventID)),
	}

	err = jrpcClient.Call("DplatformOS.Query", params, &res)
	assert.Equal(t, types.ErrNotFound, err)

	//       +  ，
	params = rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryEventIDByAddrAndStatus,
		Payload:  []byte("{\"status\":1,\"addr\":\"16ABEbYtKKQm5wMuySK9J4La5nAiidGuyt\",\"type\":\"\",\"eventID\":\"\"}"),
	}
	err = jrpcClient.Call("DplatformOS.Query", params, &res)
	assert.Equal(t, types.ErrNotFound, err)
}

func queryEventByStatusAndType(t *testing.T, jrpcClient *jsonclient.JSONClient) {
	//
	params := rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryEventIDByTypeAndStatus,
		Payload:  []byte("{\"status\":1,\"addr\":\"\",\"type\":\"football\",\"eventID\":\"\"}"),
	}
	var res oty.ReplyEventIDs
	err := jrpcClient.Call("DplatformOS.Query", params, &res)
	assert.Nil(t, err)
	assert.EqualValues(t, oty.DefaultCount, len(res.EventID))
	lastEventID := res.EventID[oty.DefaultCount-1]
	//
	params = rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryEventIDByTypeAndStatus,
		Payload:  []byte(fmt.Sprintf("{\"status\":1,\"addr\":\"\",\"type\":\"football\",\"eventID\":\"%s\"}", lastEventID)),
	}
	err = jrpcClient.Call("DplatformOS.Query", params, &res)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(res.EventID))
	lastEventID = res.EventID[9]

	//
	params = rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryEventIDByTypeAndStatus,
		Payload:  []byte(fmt.Sprintf("{\"status\":1,\"addr\":\"\",\"type\":\"football\",\"eventID\":\"%s\"}", lastEventID)),
	}

	err = jrpcClient.Call("DplatformOS.Query", params, &res)
	assert.Equal(t, types.ErrNotFound, err)

	//       +
	params = rpctypes.Query4Jrpc{
		Execer:   oty.OracleX,
		FuncName: oty.FuncNameQueryEventIDByTypeAndStatus,
		Payload:  []byte("{\"status\":1,\"addr\":\"\",\"type\":\"gambling\",\"eventID\":\"\"}"),
	}
	err = jrpcClient.Call("DplatformOS.Query", params, &res)
	assert.Equal(t, types.ErrNotFound, err)

}
