package types

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/storage/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/storage/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/storage/rpc"
	storagetypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/storage/types"
)

/*
 *    dapp
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     storagetypes.StorageX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
