package types

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"

	"reflect"
)

/*
 *
 *   action      log  ，
 *    action log   id   name
 */

// action  id name，
const (
	TyUnknowAction = iota
	TyContentStorageAction
	TyHashStorageAction
	TyLinkStorageAction
	TyEncryptStorageAction
	TyEncryptShareStorageAction
	TyEncryptAddAction

	NameContentStorageAction      = "ContentStorage"
	NameHashStorageAction         = "HashStorage"
	NameLinkStorageAction         = "LinkStorage"
	NameEncryptStorageAction      = "EncryptStorage"
	NameEncryptShareStorageAction = "EncryptShareStorage"
	NameEncryptAddAction          = "EncryptAdd"

	FuncNameQueryStorage      = "QueryStorage"
	FuncNameBatchQueryStorage = "BatchQueryStorage"
)

// log  id
const (
	TyUnknownLog = iota
	TyContentStorageLog
	TyHashStorageLog
	TyLinkStorageLog
	TyEncryptStorageLog
	TyEncryptShareStorageLog
	TyEncryptAddLog
)

//storage op
const (
	OpCreate = int32(iota)
	OpAdd
)

//fork
var (
	ForkStorageLocalDB = "ForkStorageLocalDB"
)
var (
	//StorageX
	StorageX = "storage"
	//  actionMap
	actionMap = map[string]int32{
		NameContentStorageAction:      TyContentStorageAction,
		NameHashStorageAction:         TyHashStorageAction,
		NameLinkStorageAction:         TyLinkStorageAction,
		NameEncryptStorageAction:      TyEncryptStorageAction,
		NameEncryptShareStorageAction: TyEncryptShareStorageAction,
		NameEncryptAddAction:          TyEncryptAddAction,
	}
	//  log id   log     ，       log
	logMap = map[int64]*types.LogInfo{
		TyContentStorageLog:      {Ty: reflect.TypeOf(Storage{}), Name: "LogContentStorage"},
		TyHashStorageLog:         {Ty: reflect.TypeOf(Storage{}), Name: "LogHashStorage"},
		TyLinkStorageLog:         {Ty: reflect.TypeOf(Storage{}), Name: "LogLinkStorage"},
		TyEncryptStorageLog:      {Ty: reflect.TypeOf(Storage{}), Name: "LogEncryptStorage"},
		TyEncryptShareStorageLog: {Ty: reflect.TypeOf(Storage{}), Name: "LogEncryptShareStorage"},
		TyEncryptAddLog:          {Ty: reflect.TypeOf(Storage{}), Name: "LogEncryptAdd"},
	}
)

// init defines a register function
func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(StorageX))
	//
	types.RegFork(StorageX, InitFork)
	types.RegExec(StorageX, InitExecutor)
}

// InitFork defines register fork
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(StorageX, "Enable", 0)
	cfg.RegisterDappFork(StorageX, ForkStorageLocalDB, 0)
}

// InitExecutor defines register executor
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(StorageX, NewType(cfg))
}

//StorageType ...
type StorageType struct {
	types.ExecTypeBase
}

//NewType ...
func NewType(cfg *types.DplatformOSConfig) *StorageType {
	c := &StorageType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetPayload     action
func (s *StorageType) GetPayload() types.Message {
	return &StorageAction{}
}

// GetTypeMap     action id name
func (s *StorageType) GetTypeMap() map[string]int32 {
	return actionMap
}

// GetLogMap     log
func (s *StorageType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
