// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package retrieve

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/retrieve/cmd"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/retrieve/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/retrieve/rpc"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/retrieve/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.RetrieveX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      cmd.RetrieveCmd,
		RPC:      rpc.Init,
	})
}
