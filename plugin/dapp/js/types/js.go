package types

import (
	"errors"
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/js/types/jsproto"
)

// action for executor
const (
	jsActionCreate = 0
	jsActionCall   = 1
)

//
const (
	TyLogJs = 10000
)

// JsCreator       js
const JsCreator = "js-creator"

var (
	typeMap = map[string]int32{
		"Create": jsActionCreate,
		"Call":   jsActionCall,
	}
	logMap = map[int64]*types.LogInfo{
		TyLogJs: {Ty: reflect.TypeOf(jsproto.JsLog{}), Name: "TyLogJs"},
	}
)

//JsX
var JsX = "jsvm"

//
var (
	ErrDupName            = errors.New("ErrDupName")
	ErrJsReturnNotObject  = errors.New("ErrJsReturnNotObject")
	ErrJsReturnKVSFormat  = errors.New("ErrJsReturnKVSFormat")
	ErrJsReturnLogsFormat = errors.New("ErrJsReturnLogsFormat")
	//ErrInvalidFuncFormat          (  _)
	ErrInvalidFuncFormat = errors.New("dplatformos.js: invalid function name format")
	//ErrInvalidFuncPrefix not exec_ execloal_ query_
	ErrInvalidFuncPrefix = errors.New("dplatformos.js: invalid function prefix format")
	//ErrFuncNotFound
	ErrFuncNotFound = errors.New("dplatformos.js: invalid function name not found")
	ErrSymbolName   = errors.New("dplatformos.js: ErrSymbolName")
	ErrExecerName   = errors.New("dplatformos.js: ErrExecerName")
	ErrDBType       = errors.New("dplatformos.js: ErrDBType")
	// ErrJsCreator
	ErrJsCreator = errors.New("ErrJsCreator")
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(JsX))
	types.RegFork(JsX, InitFork)
	types.RegExec(JsX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(JsX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(JsX, NewType(cfg))
}

//JsType
type JsType struct {
	types.ExecTypeBase
}

//NewType     plugin
func NewType(cfg *types.DplatformOSConfig) *JsType {
	c := &JsType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

//GetPayload
func (t *JsType) GetPayload() types.Message {
	return &jsproto.JsAction{}
}

//GetTypeMap
func (t *JsType) GetTypeMap() map[string]int32 {
	return typeMap
}

//GetLogMap
func (t *JsType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
