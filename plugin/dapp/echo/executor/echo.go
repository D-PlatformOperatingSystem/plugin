package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	echotypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/echo/types/echo"
)

var (
	// KeyPrefixPing ping
	KeyPrefixPing = "mavl-echo-ping:%s"
	// KeyPrefixPang pang
	KeyPrefixPang = "mavl-echo-pang:%s"

	// KeyPrefixPingLocal local ping
	KeyPrefixPingLocal = "LODB-echo-ping:%s"
	// KeyPrefixPangLocal local pang
	KeyPrefixPangLocal = "LODB-echo-pang:%s"
)

// Init           ，         ，         0
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	dapp.Register(cfg, echotypes.EchoX, newEcho, 0)
	InitExecType()
}

// InitExecType
func InitExecType() {
	ety := types.LoadExecutorType(echotypes.EchoX)
	ety.InitFuncList(types.ListMethod(&Echo{}))
}

// Echo
type Echo struct {
	dapp.DriverBase
}

//             ，
func newEcho() dapp.Driver {
	c := &Echo{}
	c.SetChild(c)
	c.SetExecutorType(types.LoadExecutorType(echotypes.EchoX))
	return c
}

// GetDriverName
func (h *Echo) GetDriverName() string {
	return echotypes.EchoX
}
