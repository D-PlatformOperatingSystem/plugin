package types

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/rpc"
	exchangetypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
)

/*
 *    dapp     
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     exchangetypes.ExchangeX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
