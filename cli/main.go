// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

package main

import (
	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	"github.com/D-PlatformOperatingSystem/plugin/cli/buildflags"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin"

	"github.com/D-PlatformOperatingSystem/dpos/util/cli"
)

func main() {
	if buildflags.RPCAddr == "" {
		buildflags.RPCAddr = "http://localhost:28803"
	}
	cli.Run(buildflags.RPCAddr, buildflags.ParaName, "")
}
