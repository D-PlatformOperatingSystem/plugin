package executor

import (
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	storagetypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/storage/types"
)

/*
 *        
 *         
 */

var (
	//  
	elog = log.New("module", "storage.executor")
)

var driverName = storagetypes.StorageX

// Init register dapp
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	drivers.Register(cfg, GetName(), newStorage, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

// InitExecType Init Exec Type
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&storage{}))
}

type storage struct {
	drivers.DriverBase
}

func newStorage() drivers.Driver {
	t := &storage{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetName get driver name
func GetName() string {
	return newStorage().GetName()
}

func (s *storage) GetDriverName() string {
	return driverName
}

//ExecutorOrder Exec          ExecLocal
func (s *storage) ExecutorOrder() int64 {
	cfg := s.GetAPI().GetConfig()
	if cfg.IsDappFork(s.GetHeight(), storagetypes.StorageX, storagetypes.ForkStorageLocalDB) {
		return drivers.ExecLocalSameTime
	}
	return s.DriverBase.ExecutorOrder()
}

// CheckTx            ，     
func (s *storage) CheckTx(tx *types.Transaction, index int) error {
	// implement code
	return nil
}
