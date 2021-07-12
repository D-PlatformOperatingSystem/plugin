package echo

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

// CreateTx
func (e *Type) CreateTx(action string, message json.RawMessage) (*types.Transaction, error) {
	elog.Debug("echo.CreateTx", "action", action)
	//    ping/pang
	cfg := e.GetConfig()
	if action == "ping" || action == "pang" {
		var param Tx
		err := json.Unmarshal(message, &param)
		if err != nil {
			elog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return createPingTx(cfg, action, &param)
	}
	return nil, types.ErrNotSupport
}

func createPingTx(cfg *types.DplatformOSConfig, op string, parm *Tx) (*types.Transaction, error) {
	var action *EchoAction
	var err error
	if strings.EqualFold(op, "ping") {
		action, err = getPingAction(parm)
	} else {
		action, err = getPangAction(parm)
	}
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(EchoX)),
		Payload: types.Encode(action),
		Nonce:   rand.New(rand.NewSource(time.Now().UnixNano())).Int63(),
		To:      address.ExecAddress(cfg.ExecName(EchoX)),
	}
	return tx, nil
}

func getPingAction(parm *Tx) (*EchoAction, error) {
	pingAction := &Ping{Msg: parm.Message}
	action := &EchoAction{
		Value: &EchoAction_Ping{Ping: pingAction},
		Ty:    ActionPing,
	}
	return action, nil
}

func getPangAction(parm *Tx) (*EchoAction, error) {
	pangAction := &Pang{Msg: parm.Message}
	action := &EchoAction{
		Value: &EchoAction_Pang{Pang: pangAction},
		Ty:    ActionPang,
	}
	return action, nil
}
