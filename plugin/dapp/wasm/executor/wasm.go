package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	types2 "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/wasm/types"
)

var driverName = types2.WasmX
var log = log15.New("module", "execs."+types2.WasmX)

func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	if name != driverName {
		panic("system dapp can not be rename")
	}

	drivers.Register(cfg, name, newWasm, cfg.GetDappFork(name, "Enable"))
	initExecType()
}

func initExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Wasm{}))
}

type Wasm struct {
	drivers.DriverBase

	tx           *types.Transaction
	stateKVC     *dapp.KVCreator
	localCache   []*types2.LocalDataLog
	kvs          []*types.KeyValue
	receiptLogs  []*types.ReceiptLog
	customLogs   []string
	execAddr     string
	contractName string
}

func newWasm() drivers.Driver {
	d := &Wasm{}
	d.SetChild(d)
	d.SetExecutorType(types.LoadExecutorType(driverName))
	return d
}

// GetName
func GetName() string {
	return newWasm().GetName()
}

func (w *Wasm) GetDriverName() string {
	return driverName
}
