package executor

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	echotypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/echo/types/echo"
)

// Exec_Ping    ping
func (h *Echo) Exec_Ping(ping *echotypes.Ping, tx *types.Transaction, index int) (*types.Receipt, error) {
	msg := ping.Msg
	res := fmt.Sprintf("%s, ping ping ping!", msg)
	xx := &echotypes.PingLog{Msg: msg, Echo: res}
	logs := []*types.ReceiptLog{{Ty: echotypes.TyLogPing, Log: types.Encode(xx)}}
	kv := []*types.KeyValue{{Key: []byte(fmt.Sprintf(KeyPrefixPing, msg)), Value: []byte(res)}}
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// Exec_Pang    pang
func (h *Echo) Exec_Pang(ping *echotypes.Pang, tx *types.Transaction, index int) (*types.Receipt, error) {
	msg := ping.Msg
	res := fmt.Sprintf("%s, pang pang pang!", msg)
	xx := &echotypes.PangLog{Msg: msg, Echo: res}
	logs := []*types.ReceiptLog{{Ty: echotypes.TyLogPang, Log: types.Encode(xx)}}
	kv := []*types.KeyValue{{Key: []byte(fmt.Sprintf(KeyPrefixPang, msg)), Value: []byte(res)}}
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}
