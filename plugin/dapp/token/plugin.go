// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token   token
package token

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/autotest" // register token autotest package
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/rpc"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.TokenX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.TokenCmd,
		RPC:      rpc.Init,
	})
}
