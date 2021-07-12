package types

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/accountmanager/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/accountmanager/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/accountmanager/rpc"
	accountmanagertypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/accountmanager/types"
)

/*
 *    dapp
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     accountmanagertypes.AccountmanagerX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
