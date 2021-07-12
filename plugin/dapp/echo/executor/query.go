package executor

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	echotypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/echo/types/echo"
)

// Query_GetPing    ping
func (h *Echo) Query_GetPing(in *echotypes.Query) (types.Message, error) {
	var pingLog echotypes.PingLog
	localKey := []byte(fmt.Sprintf(KeyPrefixPingLocal, in.Msg))
	value, err := h.GetLocalDB().Get(localKey)
	if err != nil {
		return nil, err
	}
	types.Decode(value, &pingLog)
	res := echotypes.QueryResult{Msg: in.Msg, Count: pingLog.Count}
	return &res, nil
}

// Query_GetPang    pang
func (h *Echo) Query_GetPang(in *echotypes.Query) (types.Message, error) {
	var pangLog echotypes.PangLog
	localKey := []byte(fmt.Sprintf(KeyPrefixPangLocal, in.Msg))
	value, err := h.GetLocalDB().Get(localKey)
	if err != nil {
		return nil, err
	}
	types.Decode(value, &pangLog)
	res := echotypes.QueryResult{Msg: in.Msg, Count: pangLog.Count}
	return &res, nil
}
