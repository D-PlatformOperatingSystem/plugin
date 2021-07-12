// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collateralize

import (
	"github.com/D-PlatformOperatingSystem/dpos/pluginmgr"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/collateralize/commands"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/collateralize/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/collateralize/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.CollateralizeX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.CollateralizeCmd,
	})
}
