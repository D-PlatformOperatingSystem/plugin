// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet_test

import (
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/util"

	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin"
	node "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/testnode"
)

func TestParaQuery(t *testing.T) {
	para := node.NewParaNode(nil, nil)
	paraCfg := para.Para.GetAPI().GetConfig()
	defer para.Close()

	var param types.ReqWalletImportPrivkey
	param.Label = "Importprivkey"
	param.Privkey = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
	para.Para.GetAPI().ExecWalletFunc("wallet", "WalletImportPrivkey", &param)

	var param1 types.ReqNewAccount
	param1.Label = "NewAccount"
	para.Para.GetAPI().ExecWalletFunc("wallet", "NewAccount", &param1)
	para.Para.GetAPI().ExecWalletFunc("wallet", "WalletLock", &types.ReqNil{})

	//  rpc
	tx := util.CreateTxWithExecer(paraCfg, para.Para.GetGenesisKey(), "user.p.test.none")
	para.Para.SendTxRPC(tx)
	para.Para.WaitHeight(1)
	tx = util.CreateTxWithExecer(paraCfg, para.Para.GetGenesisKey(), "user.p.test.none")
	para.Para.SendTxRPC(tx)
	para.Para.WaitHeight(2)

}
