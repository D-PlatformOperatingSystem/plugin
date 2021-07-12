// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/client"
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/queue"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	evm "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common/crypto"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/runtime"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/state"
)

func TestVM(t *testing.T) {

	basePath := "testdata/"

	//
	genTestCase(basePath)

	t.Parallel()

	//
	runTestCase(t, basePath)

	//
	defer clearTestCase(basePath)
}

func TestTmp(t *testing.T) {
	//addr := common.StringToAddress("19i4kLkSrAr4ssvk1pLwjkFAnoXeJgvGvj")
	//fmt.Println(hex.EncodeToString(addr.Bytes()))
	tt := types.Now().Unix()
	fmt.Println(time.Unix(tt, 0).String())
}

type CaseFilter struct{}

var testCaseFilter = &CaseFilter{}

//
func (filter *CaseFilter) filter(num int) bool {
	return num >= 0
}

//
func (filter *CaseFilter) filterCaseName(name string) bool {
	//return name == "selfdestruct"
	return name != ""
}

func runTestCase(t *testing.T, basePath string) {
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if path == basePath || !info.IsDir() {
			return nil
		}
		runDir(t, path)
		return nil
	})
}

func runDir(tt *testing.T, basePath string) {
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		baseName := info.Name()

		if baseName[:5] == "data_" || baseName[:4] == "tpl_" || filepath.Ext(path) != ".json" {
			return nil
		}

		raw, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err.Error())
			tt.FailNow()
		}

		var data interface{}
		json.Unmarshal(raw, &data)

		cases := parseData(data.(map[string]interface{}))
		for _, c := range cases {
			//       ，
			tt.Run(c.name, func(t *testing.T) {
				runCase(t, c, baseName)
			})
		}

		return nil
	})
}

func runCase(tt *testing.T, c VMCase, file string) {
	tt.Logf("running test case:%s in file:%s", c.name, file)

	// 1        pre
	inst := evm.NewEVMExecutor()
	q := queue.New("channel")
	q.SetConfig(chainTestCfg)
	api, _ := client.New(q.Client(), nil)
	inst.SetAPI(api)
	inst.SetEnv(c.env.currentNumber, c.env.currentTimestamp, uint64(c.env.currentDifficulty))
	inst.CheckInit()
	statedb := inst.GetMStateDB()
	mdb := createStateDB(statedb, c)
	statedb.StateDB = mdb
	statedb.CoinsAccount = account.NewCoinsAccount(chainTestCfg)
	statedb.CoinsAccount.SetDB(statedb.StateDB)

	// 2        create
	vmcfg := inst.GetVMConfig()
	msg := buildMsg(c)
	context := inst.NewEVMContext(msg)
	context.Coinbase = common.StringToAddress(c.env.currentCoinbase)

	// 3        call
	env := runtime.NewEVM(context, statedb, *vmcfg, api.GetConfig())
	var (
		ret []byte
		//addr common.Address
		//leftGas uint64
		err error
	)

	if len(c.exec.address) > 0 {
		ret, _, _, err = env.Call(runtime.AccountRef(msg.From()), *common.StringToAddress(c.exec.address), msg.Data(), msg.GasLimit(), msg.Value())
	} else {
		addr := crypto.RandomContractAddress()
		ret, _, _, err = env.Create(runtime.AccountRef(msg.From()), *addr, msg.Data(), msg.GasLimit(), "testExecName", "", "")
	}

	if err != nil {
		//           ，        ，    ，   ，   post
		if len(c.err) > 0 && c.err == err.Error() {
			return
		}
		//          ，
		tt.Errorf("test case:%s, failed:%s", c.name, err)
		tt.Fail()
		return
	}
	// 4        post (  ，     Gas      ，         ，           )
	t := NewTester(tt)
	// 4.1
	t.assertEqualsB(ret, getBin(c.out))

	// 4.2
	for k, v := range c.post {
		addrStr := (*common.StringToAddress(k)).String()
		t.assertEqualsV(int(statedb.GetBalance(addrStr)), int(v.balance))

		t.assertEqualsB(statedb.GetCode(addrStr), getBin(v.code))

		for a, b := range v.storage {
			if len(a) < 1 || len(b) < 1 {
				continue
			}
			hashKey := common.BytesToHash(getBin(a))
			hashVal := common.BytesToHash(getBin(b))
			t.assertEqualsB(statedb.GetState(addrStr, hashKey).Bytes(), hashVal.Bytes())
		}
	}
}

//
func createStateDB(msdb *state.MemoryStateDB, c VMCase) *db.GoMemDB {
	//   statedb     ，
	mdb, _ := db.NewGoMemDB("test", "", 0)
	//
	for k, v := range c.pre {
		//  coins
		ac := &types.Account{Addr: c.exec.caller, Balance: v.balance}
		addAccount(mdb, k, ac)

		//
		addContractAccount(msdb, mdb, k, v, c.exec.caller)
	}

	//   MemoryStateDB
	msdb.ResetDatas()

	return mdb
}

//
func buildMsg(c VMCase) *common.Message {
	code, _ := hex.DecodeString(c.exec.code)
	addr1 := common.StringToAddress(c.exec.caller)
	addr2 := common.StringToAddress(c.exec.address)
	gasLimit := uint64(210000000)
	gasPrice := c.exec.gasPrice
	return common.NewMessage(*addr1, addr2, int64(1), uint64(c.exec.value), gasLimit, uint32(gasPrice), code, "", "")
}
