package types

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/types"
)

/*
 *
 *   action      log  ，
 *    action log   id   name
 */

// action  id name，
const (
	TyUnknowAction = iota + 200
	TyLimitOrderAction
	TyMarketOrderAction
	TyRevokeOrderAction

	NameLimitOrderAction  = "LimitOrder"
	NameMarketOrderAction = "MarketOrder"
	NameRevokeOrderAction = "RevokeOrder"

	FuncNameQueryMarketDepth      = "QueryMarketDepth"
	FuncNameQueryHistoryOrderList = "QueryHistoryOrderList"
	FuncNameQueryOrder            = "QueryOrder"
	FuncNameQueryOrderList        = "QueryOrderList"
)

// log  id
const (
	TyUnknownLog = iota + 200
	TyLimitOrderLog
	TyMarketOrderLog
	TyRevokeOrderLog
)

// OP
const (
	OpBuy = iota + 1
	OpSell
)

//order status
const (
	Ordered = iota
	Completed
	Revoked
)

//const
const (
	ListDESC = int32(0)
	ListASC  = int32(1)
	ListSeek = int32(2)
)

const (
	//Count   list
	Count = int32(10)
	//MaxMatchCount
	MaxMatchCount = 100
)

var (
	//ExchangeX
	ExchangeX = "exchange"
	//  actionMap
	actionMap = map[string]int32{
		NameLimitOrderAction:  TyLimitOrderAction,
		NameMarketOrderAction: TyMarketOrderAction,
		NameRevokeOrderAction: TyRevokeOrderAction,
	}
	//  log id   log     ，       log
	logMap = map[int64]*types.LogInfo{
		TyLimitOrderLog:  {Ty: reflect.TypeOf(ReceiptExchange{}), Name: "TyLimitOrderLog"},
		TyMarketOrderLog: {Ty: reflect.TypeOf(ReceiptExchange{}), Name: "TyMarketOrderLog"},
		TyRevokeOrderLog: {Ty: reflect.TypeOf(ReceiptExchange{}), Name: "TyRevokeOrderLog"},
	}
	//tlog = log.New("module", "exchange.types")
)

// init defines a register function
func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(ExchangeX))
	//
	types.RegFork(ExchangeX, InitFork)
	types.RegExec(ExchangeX, InitExecutor)
}

// InitFork defines register fork
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(ExchangeX, "Enable", 0)
}

// InitExecutor defines register executor
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(ExchangeX, NewType(cfg))
}

//ExchangeType ...
type ExchangeType struct {
	types.ExecTypeBase
}

//NewType ...
func NewType(cfg *types.DplatformOSConfig) *ExchangeType {
	c := &ExchangeType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetPayload     action
func (e *ExchangeType) GetPayload() types.Message {
	return &ExchangeAction{}
}

// GetTypeMap     action id name
func (e *ExchangeType) GetTypeMap() map[string]int32 {
	return actionMap
}

// GetLogMap     log
func (e *ExchangeType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
