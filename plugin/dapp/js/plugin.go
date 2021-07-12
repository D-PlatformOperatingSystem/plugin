package js

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/js/executor"
	ptypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/js/types"

	// init auto test
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/js/autotest"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/js/command"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     ptypes.JsX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      command.JavaScriptCmd,
		RPC:      nil,
	})
}
