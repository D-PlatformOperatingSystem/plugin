// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package relay

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/relay/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/relay/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/relay/rpc"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/relay/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.RelayX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.RelayCmd,
		RPC:      rpc.Init,
	})
}
