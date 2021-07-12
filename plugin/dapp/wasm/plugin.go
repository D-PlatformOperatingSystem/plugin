package wasm

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/wasm/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/wasm/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/wasm/rpc"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/wasm/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.WasmX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
