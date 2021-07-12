// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package blackwhite
package blackwhite

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/blackwhite/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/blackwhite/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/blackwhite/rpc"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/blackwhite/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.BlackwhiteX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.BlackwhiteCmd,
		RPC:      rpc.Init,
	})
}
