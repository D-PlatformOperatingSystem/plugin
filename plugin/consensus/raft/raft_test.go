// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raft

import (
	"fmt"
	"os"
	"testing"
	"time"

	//      store,     plugin
	_ "github.com/D-PlatformOperatingSystem/dpos/system/dapp/init"
	_ "github.com/D-PlatformOperatingSystem/dpos/system/mempool/init"
	_ "github.com/D-PlatformOperatingSystem/dpos/system/store/init"
	"github.com/D-PlatformOperatingSystem/dpos/util"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"

	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/init"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/store/init"
)

//   ： go test -cover
func TestRaft(t *testing.T) {
	mockDOM := testnode.New("dplatformos.test.toml", nil)
	cfg := mockDOM.GetClient().GetConfig()
	defer mockDOM.Close()
	mockDOM.Listen()
	t.Log(mockDOM.GetGenesisAddress())
	time.Sleep(10 * time.Second)
	txs := util.GenNoneTxs(cfg, mockDOM.GetGenesisKey(), 10)
	for i := 0; i < len(txs); i++ {
		mockDOM.GetAPI().SendTx(txs[i])
	}
	mockDOM.WaitHeight(1)
	txs = util.GenNoneTxs(cfg, mockDOM.GetGenesisKey(), 10)
	for i := 0; i < len(txs); i++ {
		mockDOM.GetAPI().SendTx(txs[i])
	}
	mockDOM.WaitHeight(2)
	clearTestData()
}

func clearTestData() {
	err := os.RemoveAll("dplatformos_raft-1")
	if err != nil {
		fmt.Println("delete dplatformos_raft dir have a err:", err.Error())
	}
	fmt.Println("test data clear successfully!")
}
