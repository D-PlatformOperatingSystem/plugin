// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

/*
             ，    4 ：
      dapp
  go              。
*/
package main

import (
	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	"github.com/D-PlatformOperatingSystem/dpos/util/cli"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin"
)

func main() {
	cli.RunDplatformOS("", "")
}
