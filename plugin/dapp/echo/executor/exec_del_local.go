package executor

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	echotypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/echo/types/echo"
)

// ExecDelLocal_Ping       ，          1
func (h *Echo) ExecDelLocal_Ping(ping *echotypes.Ping, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	//       ，             
	var pingLog echotypes.PingLog
	types.Decode(receipt.Logs[0].Log, &pingLog)
	localKey := []byte(fmt.Sprintf(KeyPrefixPingLocal, pingLog.Msg))
	oldValue, err := h.GetLocalDB().Get(localKey)
	if err != nil {
		return nil, err
	}
	types.Decode(oldValue, &pingLog)
	if pingLog.Count > 0 {
		pingLog.Count--
	}
	val := types.Encode(&pingLog)
	if pingLog.Count == 0 {
		val = nil
	}
	kv := []*types.KeyValue{{Key: localKey, Value: val}}
	return &types.LocalDBSet{KV: kv}, nil
}

// ExecDelLocal_Pang       ，          1
func (h *Echo) ExecDelLocal_Pang(ping *echotypes.Pang, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	//       ，             
	var pangLog echotypes.PangLog
	types.Decode(receipt.Logs[0].Log, &pangLog)
	localKey := []byte(fmt.Sprintf(KeyPrefixPangLocal, pangLog.Msg))
	oldValue, err := h.GetLocalDB().Get(localKey)
	if err != nil {
		return nil, err
	}
	types.Decode(oldValue, &pangLog)
	if pangLog.Count > 0 {
		pangLog.Count--
	}
	val := types.Encode(&pangLog)
	if pangLog.Count == 0 {
		val = nil
	}
	kv := []*types.KeyValue{{Key: localKey, Value: val}}
	return &types.LocalDBSet{KV: kv}, nil
}
