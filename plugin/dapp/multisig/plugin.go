package multisig

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/autotest" //register auto test
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/rpc"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/wallet" // register wallet package
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     mty.MultiSigX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.MultiSigCmd,
		RPC:      rpc.Init,
	})
}
