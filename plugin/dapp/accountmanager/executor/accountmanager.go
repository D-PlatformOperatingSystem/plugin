package executor

import (
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	et "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/accountmanager/types"
)

/*
 *
 *
 */

var (
	//
	elog = log.New("module", "accountmanager.executor")
)

var driverName = et.AccountmanagerX

// Init register dapp
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	drivers.Register(cfg, GetName(), newAccountmanager, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

// InitExecType Init Exec Type
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Accountmanager{}))
}

//Accountmanager ...
type Accountmanager struct {
	drivers.DriverBase
}

func newAccountmanager() drivers.Driver {
	t := &Accountmanager{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetName get driver name
func GetName() string {
	return newAccountmanager().GetName()
}

//GetDriverName ...
func (a *Accountmanager) GetDriverName() string {
	return driverName
}

//ExecutorOrder Exec          ExecLocal
func (a *Accountmanager) ExecutorOrder() int64 {
	return drivers.ExecLocalSameTime
}

// CheckTx            ，
func (a *Accountmanager) CheckTx(tx *types.Transaction, index int) error {
	//          payload,
	var ama et.AccountmanagerAction
	err := types.Decode(tx.GetPayload(), &ama)
	if err != nil {
		return err
	}
	switch ama.Ty {
	case et.TyRegisterAction:
		register := ama.GetRegister()
		if a.CheckAccountIDIsExist(register.GetAccountID()) {
			return et.ErrAccountIDExist
		}
	case et.TySuperviseAction:

	case et.TyApplyAction:

	case et.TyTransferAction:

	case et.TyResetAction:

	}
	return nil
}

//CheckAccountIDIsExist ...
func (a *Accountmanager) CheckAccountIDIsExist(accountID string) bool {
	_, err := findAccountByID(a.GetLocalDB(), accountID)
	return err != types.ErrNotFound
}
